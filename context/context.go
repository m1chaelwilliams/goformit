package context

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"goformit/dispatch"
	"goformit/prompttree"
	"goformit/serialization"
)

// used to lazily set the context after a model has been created
type FormModel interface {
	tea.Model
	SetCTX(ctx *AppContext)
}

type AppContext struct {
	Theme        Theme
	dispatcher   dispatch.PromptModelDispatcher
	rootPrompt   *prompttree.PromptNode
	activePrompt *prompttree.PromptNode
	formJSON     *serialization.FormJSON
	formResult   *serialization.FormResult
}

func (a *AppContext) FormResult() *serialization.FormResult {
	return a.formResult
}

func NewAppContext(
	formJSON *serialization.FormJSON,
	dispatcher dispatch.PromptModelDispatcher,
) (*AppContext, error) {
	rootPrompt, err := prompttree.NewPromptTree(formJSON, dispatcher)
	if err != nil {
		return nil, err
	}

	ctx := AppContext{
		Theme:        *NewTheme(),
		dispatcher:   dispatcher,
		rootPrompt:   rootPrompt,
		activePrompt: rootPrompt,
		formJSON:     formJSON,
		formResult:   serialization.NewFormResult(),
	}

	if rootModel, isType := rootPrompt.Model.(FormModel); isType {
		rootModel.SetCTX(&ctx)
		return &ctx, nil
	}

	return nil, errors.New("root model is not of type FormModel")
}

func (a *AppContext) ActiveModel() tea.Model {
	return a.activePrompt.Model
}

func (a *AppContext) NextModel(response []string) (tea.Model, error) {
	re := regexp.MustCompile(`\[\[(.*?)\]\]`)

	bindSubmit := a.formJSON.Prompts[a.activePrompt.Id].BindSubmit
	if bindSubmit != nil {
		varName := re.ReplaceAllStringFunc(*bindSubmit, func(s string) string {
			return strings.Trim(s, "[]")
		})
		a.formJSON.Vars[varName] = strings.Join(response, ",")
	}

	promptGroup := a.activePrompt.Group
	promptGroup = re.ReplaceAllStringFunc(promptGroup, func(s string) string {
		return strings.Trim(s, "[]")
	})

	if promptVar, exists := a.formJSON.Vars[promptGroup]; exists {
		promptGroup = promptVar
	}

	if promptGroup != "void" {
		promptResult := serialization.PromptResult{
			ID:             a.activePrompt.Id,
			Group:          promptGroup,
			Response:       strings.Join(response, ","),
			ResponseNumber: 1,
		}
		a.formResult.AddPromptResult(&promptResult)
	}

	fmt.Printf("Submitting response: %v\n", response)

	nextPromptID := "not_found"

outer:
	for _, qualifier := range a.activePrompt.Qualifiers {
		if qualifier.Qualifies(response) {
			nextPromptID = qualifier.ModelID()
			break outer
			// if nextPromptID != "[[end]]" {
			// 	break outer
			// }
		}
	}

	if nextPromptID == "not_found" {
		nextPromptID = a.formJSON.Prompts[a.activePrompt.Id].Next["_"]
	}

	if nextPromptID == "[[end]]" {
		return nil, errors.New(fmt.Sprintf("unable to find next prompt: %v", nextPromptID))
	}

	nextPromptJSON := a.formJSON.Prompts[nextPromptID]
	nextPrompt, err := prompttree.CreatePromptNode(nextPromptJSON, a.formJSON.Vars, a.dispatcher)
	if err != nil {
		return nil, err
	}
	a.activePrompt = nextPrompt
	if activeModel, isType := a.activePrompt.Model.(FormModel); isType {
		activeModel.SetCTX(a)
		return activeModel, nil
	}
	return nil, errors.New("unable to get next model. does not implemented SetCTX")
}

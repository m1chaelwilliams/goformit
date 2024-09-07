package context

import (
	"errors"
	"fmt"
	"goformit/internal/dispatch"
	"goformit/internal/logging"
	"goformit/internal/prompttree"
	"goformit/internal/serialization"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	PROMPT_NOT_FOUND = "not_found"
)

// used to lazily set the context after a model has been created
type FormModel interface {
	tea.Model
	SetCTX(ctx *AppContext)
}

// the global state of the application
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

// creates the global app context from form json
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

// attempts to process model's response and grab the next model
// errors if last prompt, error is handled by model to exit gracefully
func (a *AppContext) NextModel(response []string) (tea.Model, error) {
	// finds all variables (variables are wrapped in [[]])
	re := regexp.MustCompile(`\[\[(.*?)\]\]`)

	// grabs optional submit event
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

	logging.AppLogger.Log("CONTEXT", "Submitting response: %v\n", response)

	nextPromptID := PROMPT_NOT_FOUND

	// find the next prompt id
outer:
	for _, qualifier := range a.activePrompt.Qualifiers {
		if qualifier.Qualifies(response) {
			nextPromptID = qualifier.ModelID()
			break outer
		}
	}

	// go to catchall (base case) next if no qualifiers
	if nextPromptID == PROMPT_NOT_FOUND {
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

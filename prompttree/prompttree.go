package prompttree

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"goformit/dispatch"
	"goformit/logging"
	"goformit/serialization"
)

type Qualifier interface {
	Qualifies(input []string) bool
	ModelID() string
}

type BasicQualifier struct {
	Options     [][]string
	nextModelId string
}

func createQualifierOption(str string) []string {
	items := strings.Split(str, "+")
	sort.Strings(items)
	logging.AppLogger.Log("qualifier_init", "%v", items)
	return items
}

func NewBasicQualifier(rawStr string, nextModelId string) *BasicQualifier {
	optsRaw := strings.Split(rawStr, ",")
	opts := make([][]string, 0)
	for _, opt := range optsRaw {
		opts = append(opts, createQualifierOption(opt))
	}

	return &BasicQualifier{
		Options:     opts,
		nextModelId: nextModelId,
	}
}

func (b *BasicQualifier) Qualifies(input []string) bool {
	sort.Strings(input)

	qStatus := false
	for _, opt := range b.Options {
		logging.AppLogger.Log("qualifier", "checking %v against %v", input, opt)

		if opt[0] == "_" {
			return true
		}

		if len(input) != len(opt) {
			continue
		}

		passes := true
		for i, item := range input {
			if item != opt[i] {
				passes = false
			}
		}

		if passes {
			return true
		}
	}
	return qStatus
}

func (b *BasicQualifier) ModelID() string {
	return b.nextModelId
}

type PromptNode struct {
	Id         string
	Group      string
	Model      tea.Model
	Next       map[string]*PromptNode
	Qualifiers []Qualifier
}

func CreatePromptNode(
	promptJSON *serialization.PromptJSON,
	programVars map[string]string,
	dispatcher dispatch.PromptModelDispatcher,
) (*PromptNode, error) {
	qualifiers := make([]Qualifier, 0)
	nextMap := make(map[string]*PromptNode)
	for input, nextModelId := range promptJSON.Next {
		qualifiers = append(qualifiers, NewBasicQualifier(input, nextModelId))
		nextMap[nextModelId] = nil
	}

	model := dispatcher(promptJSON)
	if model == nil {
		return nil, errors.New(
			fmt.Sprintf("unsupported prompt type provided for model: %s", promptJSON.Id),
		)
	}

	fmt.Printf("creating prompt node: %s\n", promptJSON.Id)

	return &PromptNode{
		Id:         promptJSON.Id,
		Group:      promptJSON.Group,
		Model:      model,
		Qualifiers: qualifiers,
		Next:       nextMap,
	}, nil
}

func NewPromptTree(
	jsonContents *serialization.FormJSON,
	dispatcher dispatch.PromptModelDispatcher,
) (*PromptNode, error) {
	firstPrompt, exists := jsonContents.Prompts[jsonContents.First]
	if !exists {
		return nil, errors.New("could not find first prompt")
	}

	root, err := CreatePromptNode(firstPrompt, jsonContents.Vars, dispatcher)
	if err != nil {
		return nil, err
	}

	// createPromptTree(root, dispatcher, jsonContents)

	return root, nil
}

package prompttree

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"goformit/dispatch"
	"goformit/logging"
	"goformit/serialization"
)

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

	logging.AppLogger.Log("prompttree", "creating prompt node: %s\n", promptJSON.Id)

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

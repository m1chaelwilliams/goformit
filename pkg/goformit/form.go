package goformit

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"

	"goformit/internal/context"
	"goformit/internal/models"
	"goformit/internal/serialization"
)

type Form struct {
	// form data deserialized from the provided JSON file
	formJSON *serialization.FormJSON
	// global context of the application
	// manages submits, branching, and theming
	ctx *context.AppContext
}

// starts and manages the application mainloop
func (f *Form) Run() error {
	program := tea.NewProgram(f.ctx.ActiveModel())
	_, err := program.Run()
	return err
}

// the result of the form being executed
func (f *Form) Result() (string, error) {
	contents, err := json.Marshal(f.ctx.FormResult())
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func newFormFromJSONInternal(formJSON *serialization.FormJSON) (*Form, error) {
	ctx, err := context.NewAppContext(
		formJSON,
		func(promptJSON *serialization.PromptJSON) tea.Model {
			switch promptJSON.Type {
			case "input":
				return models.NewInputModelFromJSON(promptJSON)
			case "selection":
				return models.NewSelectionModelfromJSON(promptJSON)
			case "checkbox":
				return models.NewMultiSelectModelFromJSON(promptJSON)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &Form{
		formJSON: formJSON,
		ctx:      ctx,
	}, nil
}

// creates a new form from a json content blob
func NewFormFromJSONContent(content []byte) (*Form, error) {
	formJSON, err := serialization.NewFormJSONContent(content)
	if err != nil {
		return nil, err
	}

	return newFormFromJSONInternal(formJSON)
}

// creates a new form from a json filepath
func NewFormFromJSONFilepath(filePathWithExtension string) (*Form, error) {
	formJSON, err := serialization.NewFormJSON(filePathWithExtension)
	if err != nil {
		return nil, err
	}

	return newFormFromJSONInternal(formJSON)
}

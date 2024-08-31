package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"goformit/context"
	"goformit/logging"
	"goformit/serialization"
)

type InputModel struct {
	style       lipgloss.Style
	Title       string
	Desc        *string
	Placeholder *string
	textInput   textinput.Model
	ctx         *context.AppContext
}

// Init implements tea.Model.
func (i *InputModel) Init() tea.Cmd {
	i.textInput.Focus()
	return nil
}

func (i *InputModel) SetCTX(ctx *context.AppContext) {
	i.ctx = ctx
}

// Update implements tea.Model.
func (i *InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgStr := msg.String()
		if msgStr == "esc" {
			return i, tea.Quit
		} else if msgStr == "enter" {
			m, err := i.ctx.NextModel([]string{
				i.textInput.Value(),
			})
			if err == nil {
				return m, nil
			} else {
				logging.AppLogger.Log("input", "%s", err.Error())
				return i, tea.Quit
			}
		}
	}

	i.textInput, cmd = i.textInput.Update(msg)

	return i, cmd
}

// View implements tea.Model.
func (i *InputModel) View() string {
	entries := []string{
		i.Title,
	}
	if i.Desc != nil {
		entries = append(entries, *i.Desc)
	}
	entries = append(entries, i.style.Render(i.textInput.View()))

	return i.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			entries...,
		),
	)
}

func NewInputModelFromJSON(
	promptJSON *serialization.PromptJSON,
) *InputModel {
	textInput := textinput.New()
	textInput.Focus()
	if promptJSON.Placeholder != nil {
		textInput.Placeholder = *promptJSON.Placeholder
	}

	return &InputModel{
		Title:       promptJSON.Title,
		Desc:        promptJSON.Desc,
		style:       lipgloss.NewStyle().Margin(2, 2),
		Placeholder: promptJSON.Placeholder,
		textInput:   textInput,
	}
}

var _ tea.Model = (*InputModel)(nil)

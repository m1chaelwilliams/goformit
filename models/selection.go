package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"goformit/context"
	"goformit/logging"
	"goformit/serialization"
)

type SelectionModel struct {
	listGroup list.Model
	ctx       *context.AppContext
}

// Init implements tea.Model.
func (s *SelectionModel) Init() tea.Cmd {
	return nil
}

func (s *SelectionModel) SetCTX(ctx *context.AppContext) {
	s.ctx = ctx
}

// Update implements tea.Model.
func (s *SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgStr := msg.String()
		if msgStr == "esc" {
			return s, tea.Quit
		} else if msgStr == "enter" {
			selection := s.listGroup.SelectedItem().FilterValue()

			m, err := s.ctx.NextModel([]string{
				selection,
			})
			if err == nil {
				return m, nil
			} else {
				logging.AppLogger.Log("selection", err.Error())
				return s, tea.Quit
			}
		}
	}

	s.listGroup, cmd = s.listGroup.Update(msg)

	return s, cmd
}

// View implements tea.Model.
func (s *SelectionModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.listGroup.View(),
	)
}

func NewSelectionModelfromJSON(
	promptJSON *serialization.PromptJSON,
) *SelectionModel {
	items := []list.Item{}

	for _, choice := range promptJSON.Choices {
		items = append(items, NewListItem(choice))
	}

	// list.DefaultDelegate

	listModel := list.New(items, NewListItemDelegate(), 80, len(items)*4)
	listModel.Title = promptJSON.Title
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)

	return &SelectionModel{
		listGroup: listModel,
	}
}

var _ tea.Model = (*SelectionModel)(nil)

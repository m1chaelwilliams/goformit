package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"goformit/context"
	"goformit/logging"
	"goformit/serialization"
)

type MultiSelectModel struct {
	listGroup list.Model
	Selected  []int
	ctx       *context.AppContext
}

func NewMultiSelectModelFromJSON(
	promptJSON *serialization.PromptJSON,
) *MultiSelectModel {
	items := make([]list.Item, 0)

	for _, choice := range promptJSON.Choices {
		items = append(items, NewCheckboxItem(choice))
	}

	listGroup := list.New(items, NewMultiSelectItemDelegate(), 80, len(items)*4)
	listGroup.Title = promptJSON.Title

	return &MultiSelectModel{
		listGroup: listGroup,
		Selected:  make([]int, 0),
	}
}

func (m *MultiSelectModel) SetCTX(ctx *context.AppContext) {
	m.ctx = ctx
}

// Init implements tea.Model.
func (m *MultiSelectModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *MultiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgStr := msg.String()
		if msgStr == "esc" {
			return m, tea.Quit
		} else if msgStr == " " {
			m.listGroup.Items()[m.listGroup.Index()].(*CheckboxItem).Checked = !m.listGroup.Items()[m.listGroup.Index()].(*CheckboxItem).Checked
		} else if msgStr == "enter" {
			selections := make([]string, 0)
			for _, item := range m.listGroup.Items() {
				if item, isType := item.(*CheckboxItem); isType {
					if item.Checked {
						selections = append(selections, item.Content)
					}
				}
			}

			selStr := fmt.Sprintf("Sending: %v\n", selections)
			logging.AppLogger.Log("multiselect", "%v", selStr)

			nextModel, err := m.ctx.NextModel(selections)
			if err == nil {
				return nextModel, nil
			} else {
				fmt.Println(err)
				return m, tea.Quit
			}
		}
	}

	m.listGroup, cmd = m.listGroup.Update(msg)

	return m, cmd
}

// View implements tea.Model.
func (m *MultiSelectModel) View() string {
	return m.listGroup.View()
}

var _ tea.Model = (*MultiSelectModel)(nil)

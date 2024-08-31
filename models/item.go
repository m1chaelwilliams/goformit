package models

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CheckboxItem struct {
	Content string
	Checked bool
}

func NewCheckboxItem(content string) *CheckboxItem {
	return &CheckboxItem{
		Content: content,
		Checked: false,
	}
}

// FilterValue implements list.Item.
func (l *CheckboxItem) FilterValue() string {
	return l.Content
}

var _ list.Item = (*CheckboxItem)(nil)

type ListItem struct {
	Content string
}

func NewListItem(content string) *ListItem {
	return &ListItem{
		Content: content,
	}
}

// FilterValue implements list.Item.
func (l *ListItem) FilterValue() string {
	return l.Content
}

var _ list.Item = (*ListItem)(nil)

var (
	listItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedlistItemStyle = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("170"))
)

// delegate used for multi-select checkboxes
type MultiSelectItemDelegate struct {
	height, spacing int
}

func NewMultiSelectItemDelegate() *MultiSelectItemDelegate {
	return &MultiSelectItemDelegate{
		height:  1,
		spacing: 0,
	}
}

// Height implements list.ItemDelegate.
func (m *MultiSelectItemDelegate) Height() int {
	return m.height
}

// Render implements list.ItemDelegate.
func (*MultiSelectItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if item, isCheckbox := item.(*CheckboxItem); isCheckbox {
		if item.Checked {
			if index == m.Index() {
				renderContent := selectedlistItemStyle.Render(fmt.Sprintf("[x] %s", item.Content))
				fmt.Fprintf(w, renderContent)
			} else {
				renderContent := listItemStyle.Render(fmt.Sprintf("[x] %s", item.Content))
				fmt.Fprintf(w, renderContent)
			}
		} else {
			if index == m.Index() {
				renderContent := selectedlistItemStyle.Render(fmt.Sprintf("[ ] %s", item.Content))
				fmt.Fprintf(w, renderContent)
			} else {
				renderContent := listItemStyle.Render(fmt.Sprintf("[ ] %s", item.Content))
				fmt.Fprintf(w, renderContent)
			}
		}
	}
}

// Spacing implements list.ItemDelegate.
func (m *MultiSelectItemDelegate) Spacing() int {
	return m.spacing
}

// Update implements list.ItemDelegate.
func (*MultiSelectItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

var _ list.ItemDelegate = (*MultiSelectItemDelegate)(nil)

// delegate - handles drawing, sizing, etc. of items
type SingleSelectItemDelegate struct {
	height, spacing int
}

func NewListItemDelegate() *SingleSelectItemDelegate {
	return &SingleSelectItemDelegate{
		height:  1,
		spacing: 0,
	}
}

// Height implements list.ItemDelegate.
func (l *SingleSelectItemDelegate) Height() int {
	return l.height
}

// Render implements list.ItemDelegate.
func (l *SingleSelectItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	renderedContent := ""
	if index == m.Index() {
		renderedContent += selectedlistItemStyle.Render(fmt.Sprintf("> %s", item.FilterValue()))
	} else {
		renderedContent += listItemStyle.Render(fmt.Sprintf("- %s", item.FilterValue()))
	}
	fmt.Fprint(w, renderedContent)
}

// Spacing implements list.ItemDelegate.
func (l *SingleSelectItemDelegate) Spacing() int {
	return l.spacing
}

// Update implements list.ItemDelegate.
func (l *SingleSelectItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

var _ list.ItemDelegate = (*SingleSelectItemDelegate)(nil)

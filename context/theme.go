package context

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	ListStyle  lipgloss.Style
	TitleStyle lipgloss.Style
}

func NewTheme() *Theme {
	return &Theme{
		ListStyle:  lipgloss.NewStyle().PaddingLeft(2),
		TitleStyle: lipgloss.NewStyle().PaddingLeft(2),
	}
}

package dispatch

import (
	tea "github.com/charmbracelet/bubbletea"

	"goformit/serialization"
)

type PromptModelDispatcher func(promptJSON *serialization.PromptJSON) tea.Model

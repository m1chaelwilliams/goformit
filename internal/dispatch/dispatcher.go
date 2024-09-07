package dispatch

import (
	"goformit/internal/serialization"

	tea "github.com/charmbracelet/bubbletea"
)

type PromptModelDispatcher func(promptJSON *serialization.PromptJSON) tea.Model

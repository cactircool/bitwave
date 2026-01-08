package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pane interface {
	Update(w, h int, msg tea.Msg) tea.Cmd
	Style() lipgloss.Style
	SetStyle(lipgloss.Style)
	Gap() int
	Weight() float32
	Model() tea.Model
	View() string
	Resize(w, h int)
}

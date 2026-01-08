package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// SpacerModel is a flexible space that expands to fill available space
// Use it with weights to create flexible layouts (e.g., centering content)
type SpacerModel struct{}

func NewSpacerModel() *SpacerModel {
	return &SpacerModel{}
}

func (m *SpacerModel) Init() tea.Cmd {
	return nil
}

func (m *SpacerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *SpacerModel) View() string {
	return ""
}

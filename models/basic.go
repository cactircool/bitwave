package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BasicModel struct {
	Text string
}

func NewBasicModel(text string) *BasicModel {
	return &BasicModel{
		Text: text,
	}
}

func (m *BasicModel) Init() tea.Cmd {
	return nil
}

func (m *BasicModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *BasicModel) View() string {
	return m.Text
}

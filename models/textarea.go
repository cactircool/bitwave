package models

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextareaModel struct {
	textarea textarea.Model
	Style lipgloss.Style
}

func NewTextareaModel(placeholder string, showLineNums bool, style lipgloss.Style) *TextareaModel {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.ShowLineNumbers = showLineNums

	return &TextareaModel{
		textarea: ta,
		Style: style,
	}
}

func (m *TextareaModel) Init() tea.Cmd {
	return nil
}

func (m *TextareaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case EnterPane:
			m.textarea.Focus()
			return m, nil
		case EscapePane:
			m.textarea.Blur()
			UnregisterModel()
			return m, nil
		}
	}

	textarea, cmd := m.textarea.Update(msg)
	m.textarea = textarea
	return m, cmd
}

func (m *TextareaModel) View() string {
	return m.Style.Render(m.textarea.View())
}

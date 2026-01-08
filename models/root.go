package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	inner *NestedModel
}

func NewRootModel(horizontal bool, panes []Pane) *RootModel {
	root := &RootModel{
		inner: NewNestedModel(horizontal, panes), // Use the constructor!
	}

	// Initialize the first pane as selected
	if len(root.inner.Panes) > 0 {
		root.inner.Panes[0].SetStyle(root.inner.defaultStyles[0].Border(lipgloss.ThickBorder()))
	}

	RegisterModel(root.inner) // Register the inner model, not root
	return root
}

func (m *RootModel) Init() tea.Cmd {
	return m.inner.Init()
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.inner.Width, m.inner.Height = msg.Width, msg.Height
		m.inner.resizePanes()
	}
	return m.inner.Update(msg)
}

func (m *RootModel) View() string {
	return m.inner.View()
}

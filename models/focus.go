package models

import tea "github.com/charmbracelet/bubbletea"

var modelStack []tea.Model

func RegisterModel(model tea.Model) {
	modelStack = append(modelStack, model)
}

func UnregisterModel() {
	if len(modelStack) <= 1 { return }
	modelStack = modelStack[:len(modelStack)-1]
}

func IsCurrentModel(m tea.Model) bool {
	return len(modelStack) > 0 && modelStack[len(modelStack)-1] == m
}

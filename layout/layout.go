package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FocusState represents whether a model can receive focus and how
type FocusState int

const (
	// NotFocusable - model cannot receive focus (e.g., static text)
	NotFocusable FocusState = iota
	// Focusable - model can be highlighted but doesn't capture input (e.g., container)
	Focusable
	// Interactive - model captures keyboard input (e.g., textarea)
	Interactive
)

// SizedModel is the base interface all layout children must implement
type SizedModel interface {
	tea.Model
	SetSize(width, height int)
	GetFocusState() FocusState
	// OnFocus is called when this model receives focus
	// Returns the style to use and any command to run
	OnFocus(baseStyle lipgloss.Style) (lipgloss.Style, tea.Cmd)
	// OnBlur is called when this model loses focus
	OnBlur()
}

// Layout interface for containers that can hold children
type Layout interface {
	Add(model SizedModel, weight float64, style lipgloss.Style, gap int)
	AddStatic(model SizedModel, size int, style lipgloss.Style, gap int)
}

// LayoutModel combines both interfaces
type LayoutModel interface {
	SizedModel
	Layout
}

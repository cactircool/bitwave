package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextLayout struct {
	text   string
	width  int
	height int
}

func NewTextView(text string) *TextLayout {
	return &TextLayout{text: text}
}

func (t *TextLayout) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *TextLayout) GetFocusState() FocusState {
	return NotFocusable
}

func (t *TextLayout) OnFocus(baseStyle lipgloss.Style) (lipgloss.Style, tea.Cmd) {
	// Text doesn't show focus
	return baseStyle, nil
}

func (t *TextLayout) OnBlur() {}

func (t *TextLayout) Init() tea.Cmd {
	return nil
}

func (t *TextLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *TextLayout) View() string {
	style := lipgloss.NewStyle().
		Width(t.width).
		Height(t.height)
	return style.Render(t.text)
}

package layout

import (
	"github.com/cactircool/bitwave/bindings"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextareaLayout struct {
	textarea textarea.Model
	width    int
	height   int
	isActive bool // Whether we're actively editing (entered)
}

func NewTextareaLayout(ta textarea.Model) *TextareaLayout {
	return &TextareaLayout{
		textarea: ta,
		isActive: false,
	}
}

func DefaultTextareaLayout() *TextareaLayout {
	ta := textarea.New()
	ta.Placeholder = "Type something..."
	ta.Blur() // Start blurred
	return NewTextareaLayout(ta)
}

func (t *TextareaLayout) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.textarea.SetWidth(width)
	t.textarea.SetHeight(height)
}

func (t *TextareaLayout) GetFocusState() FocusState {
	return Interactive
}

func (t *TextareaLayout) OnFocus(baseStyle lipgloss.Style) (lipgloss.Style, tea.Cmd) {
	// When focused (highlighted), show thick border but don't activate editing yet
	return baseStyle.Border(lipgloss.ThickBorder()), nil
}

func (t *TextareaLayout) OnBlur() {
	// When we lose focus, deactivate editing
	t.isActive = false
	t.textarea.Blur()
}

func (t *TextareaLayout) Init() tea.Cmd {
	return textarea.Blink
}

func (t *TextareaLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle enter/escape for activating/deactivating editing
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case bindings.CycleEnter:
			if !t.isActive {
				// Activate editing mode
				t.isActive = true
				// Don't push to focus stack - we're already focused
				return t, t.textarea.Focus()
			}

		case bindings.CycleEscape:
			if t.isActive {
				// Deactivate editing mode
				t.isActive = false
				t.textarea.Blur()
				return t, nil
			}
			// If not active, let parent handle escape
		}
	}

	// Only forward input to textarea if we're actively editing
	if t.isActive {
		var cmd tea.Cmd
		t.textarea, cmd = t.textarea.Update(msg)
		return t, cmd
	}

	return t, nil
}

func (t *TextareaLayout) View() string {
	view := t.textarea.View()

	// Ensure the view respects our allocated size
	style := lipgloss.NewStyle().
		Width(t.width).
		Height(t.height)

	return style.Render(view)
}

// func (t *TextareaLayout) View() string {
// 	return t.textarea.View()
// }

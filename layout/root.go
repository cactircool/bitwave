package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootLayout struct {
	inner *GenericLayout
}

func NewRootLayout(direction Direction) *RootLayout {
	layout := NewLayout(direction)
	// Initialize the focus stack with the root layout
	pushFocus(layout)
	return &RootLayout{inner: layout}
}

func (r *RootLayout) Add(model SizedModel, weight float64, style lipgloss.Style, gap int) {
	r.inner.Add(model, weight, style, gap)
}

func (r *RootLayout) AddStatic(model SizedModel, size int, style lipgloss.Style, gap int) {
	r.inner.AddStatic(model, size, style, gap)
}

func (r *RootLayout) SetSize(width, height int) {
	r.inner.SetSize(width, height)
}

func (r *RootLayout) Init() tea.Cmd {
	// Initialize children first
	cmds := []tea.Cmd{r.inner.Init()}

	// Root's inner layout is already on the focus stack,
	// so directly focus its first focusable child
	if focusCmd := r.inner.focusFirst(); focusCmd != nil {
		cmds = append(cmds, focusCmd)
	}

	return tea.Batch(cmds...)
}

func (r *RootLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if v, ok := msg.(tea.WindowSizeMsg); ok {
		r.SetSize(v.Width, v.Height)
	}

	model, cmd := r.inner.Update(msg)
	r.inner = model.(*GenericLayout)
	return r, cmd
}

func (r *RootLayout) View() string {
	return r.inner.View()
}

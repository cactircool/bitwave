package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NestedModel struct {
	Width, Height int
	Horizontal    bool // true = columns, false = rows
	Panes         []Pane
	defaultStyles []lipgloss.Style
	selectedPane int
}

func NewNestedModel(horizontal bool, panes []Pane) *NestedModel {
	ds := make([]lipgloss.Style, len(panes))
	for i, p := range panes {
		ds[i] = p.Style()
	}
	return &NestedModel{
		Horizontal: horizontal,
		Panes:      panes,
		defaultStyles: ds,
	}
}

func (m *NestedModel) Init() tea.Cmd { return nil }

func (m *NestedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case NextPane:
			if len(m.Panes) > 0 && IsCurrentModel(m) {
				// deselect old
				m.Panes[m.selectedPane].SetStyle(m.defaultStyles[m.selectedPane])

				// increment selection
				m.selectedPane = (m.selectedPane + 1) % len(m.Panes)

				// select new
				m.Panes[m.selectedPane].SetStyle(m.defaultStyles[m.selectedPane].Border(lipgloss.ThickBorder()))

				// stop propogation
				msg = nil
			}
		case EnterPane:
			RegisterModel(m.Panes[m.selectedPane].Model())
		case EscapePane:
			UnregisterModel()
		case Exit:
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	if len(m.Panes) > 0 && msg != nil {
        if cmd := m.Panes[m.selectedPane].Update(m.Width, m.Height, msg); cmd != nil {
            cmds = append(cmds, cmd)
        }
    }

	// for i := range m.Panes {
	// 	if cmd := m.Panes[i].Update(m.Width, m.Height, msg); cmd != nil {
	// 		cmds = append(cmds, cmd)
	// 	}
	// }
	return m, tea.Batch(cmds...)
}

func (m *NestedModel) recalculate(w, h int, style lipgloss.Style) {
	chromeW, chromeH := chromeSize(style)
	m.Width, m.Height = w-chromeW, h-chromeH
	m.resizePanes()
}

func (m *NestedModel) resizePanes() {
	if len(m.Panes) == 0 {
		return
	}

	// Calculate total gaps
	totalGap := 0
	for _, pane := range m.Panes {
		// TODO: make more robust
		if p, ok := pane.(*GenericPane); ok {
			totalGap += p.Gap()
		}
	}

	// Calculate available space after gaps
	var availableSpace int
	if m.Horizontal {
		availableSpace = m.Width - totalGap
	} else {
		availableSpace = m.Height - totalGap
	}

	if availableSpace < 0 {
		availableSpace = 0
	}

	// Calculate total weight
	total := float32(0)
	for _, p := range m.Panes {
		total += p.Weight()
	}

	// Distribute space based on weights
	for i := range m.Panes {
		p := m.Panes[i]
		ratio := p.Weight() / total

		if m.Horizontal {
			p.Resize(int(float32(availableSpace)*ratio), m.Height)
		} else {
			p.Resize(m.Width, int(float32(availableSpace)*ratio))
		}
	}
}

func (m *NestedModel) View() string {
	if len(m.Panes) == 0 {
		return ""
	}

	views := make([]string, 0, len(m.Panes)*2-1)

	for i, p := range m.Panes {
		views = append(views, p.View())

		// Add gap after each pane except the last
		if i < len(m.Panes)-1 && p.Gap() > 0 {
			var gapStr string
			if m.Horizontal {
				// For horizontal layout, create a gap with width
				gapStr = lipgloss.NewStyle().Width(p.Gap()).Render("")
			} else {
				// For vertical layout, create newlines for the gap
				gapStr = ""
				for range p.Gap() {
					gapStr += "\n"
				}
			}
			views = append(views, gapStr)
		}
	}

	if m.Horizontal {
		return lipgloss.JoinHorizontal(lipgloss.Top, views...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

// func (m *NestedModel) selectPane(paneIndex int) {
// 	i := paneIndex
// 	if i >= len(m.Panes) || i >= len(m.defaultStyles) {
// 		return
// 	}
// 	m.Panes[i].Style = m.defaultStyles[i].Border(lipgloss.ThickBorder())
// }

// func (m *NestedModel) deselectPane(paneIndex int) {
// 	i := paneIndex
// 	if i >= len(m.Panes) || i >= len(m.defaultStyles) {
// 		return
// 	}
// 	m.Panes[i].Style = m.defaultStyles[i]
// }

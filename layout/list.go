package layout

import (
	"fmt"
	"strings"

	"github.com/cactircool/bitwave/bindings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	Value    string
	Data     interface{} // Optional data payload
	Selected bool
}

type ListLayout struct {
	items          []ListItem
	width          int
	height         int
	cursor         int
	scrollOffset   int
	maxSelections  int // 0 = unlimited, 1 = single select, n = max n selections

	normalStyle    lipgloss.Style
	cursorStyle    lipgloss.Style
	selectedStyle  lipgloss.Style
	selectedCursorStyle lipgloss.Style
	titleStyle     lipgloss.Style

	title          string
	showHelp       bool
}

func NewListLayout(title string, maxSelections int) *ListLayout {
	return &ListLayout{
		items:         []ListItem{},
		cursor:        0,
		scrollOffset:  0,
		maxSelections: maxSelections,
		title:         title,
		showHelp:      true,
		normalStyle:   lipgloss.NewStyle().Padding(0, 2),
		cursorStyle:   lipgloss.NewStyle().
			Background(lipgloss.Color("12")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1),
		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Padding(0, 1),
		selectedCursorStyle: lipgloss.NewStyle().
			Background(lipgloss.Color("10")).
			Foreground(lipgloss.Color("0")).
			Bold(true).
			Padding(0, 1),
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("13")).
			Padding(0, 2).
			MarginBottom(1),
	}
}

func (l *ListLayout) AddItem(value string, data interface{}) {
	l.items = append(l.items, ListItem{
		Value:    value,
		Data:     data,
		Selected: false,
	})
}

func (l *ListLayout) AddItems(values []string) {
	for _, v := range values {
		l.AddItem(v, nil)
	}
}

func (l *ListLayout) GetSelectedItems() []ListItem {
	selected := []ListItem{}
	for _, item := range l.items {
		if item.Selected {
			selected = append(selected, item)
		}
	}
	return selected
}

func (l *ListLayout) GetSelectedValues() []string {
	values := []string{}
	for _, item := range l.items {
		if item.Selected {
			values = append(values, item.Value)
		}
	}
	return values
}

func (l *ListLayout) ClearSelections() {
	for i := range l.items {
		l.items[i].Selected = false
	}
}

func (l *ListLayout) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l *ListLayout) GetFocusState() FocusState {
	return Interactive
}

func (l *ListLayout) OnFocus(baseStyle lipgloss.Style) (lipgloss.Style, tea.Cmd) {
	return baseStyle.Border(lipgloss.ThickBorder()), nil
}

func (l *ListLayout) OnBlur() {
	// Nothing special needed
}

func (l *ListLayout) Init() tea.Cmd {
	return nil
}

func (l *ListLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
				l.adjustScroll()
			}
			return l, nil

		case "down", "j":
			if l.cursor < len(l.items)-1 {
				l.cursor++
				l.adjustScroll()
			}
			return l, nil

		case "g":
			// Go to top
			l.cursor = 0
			l.adjustScroll()
			return l, nil

		case "G":
			// Go to bottom
			l.cursor = len(l.items) - 1
			l.adjustScroll()
			return l, nil

		case " ", bindings.CycleEnter:
			// Toggle selection
			if l.cursor >= 0 && l.cursor < len(l.items) {
				l.toggleSelection(l.cursor)
			}
			return l, nil

		case "a":
			// Select all (if unlimited or within limit)
			if l.maxSelections == 0 || l.maxSelections >= len(l.items) {
				for i := range l.items {
					l.items[i].Selected = true
				}
			}
			return l, nil

		case "A":
			// Deselect all
			l.ClearSelections()
			return l, nil

		case "?":
			// Toggle help
			l.showHelp = !l.showHelp
			return l, nil
		}
	}

	return l, nil
}

func (l *ListLayout) toggleSelection(index int) {
	if index < 0 || index >= len(l.items) {
		return
	}

	item := &l.items[index]

	if item.Selected {
		// Deselect
		item.Selected = false
	} else {
		// Check if we can select more
		selectedCount := 0
		for _, it := range l.items {
			if it.Selected {
				selectedCount++
			}
		}

		if l.maxSelections == 1 {
			// Single select - clear others
			l.ClearSelections()
			item.Selected = true
		} else if l.maxSelections == 0 || selectedCount < l.maxSelections {
			// Can select more
			item.Selected = true
		}
		// Otherwise do nothing (max reached)
	}
}

func (l *ListLayout) adjustScroll() {
	visibleItems := l.height - 2 // Account for title and help
	if l.showHelp {
		visibleItems -= 2
	}
	if visibleItems < 1 {
		visibleItems = 1
	}

	if l.cursor < l.scrollOffset {
		l.scrollOffset = l.cursor
	}
	if l.cursor >= l.scrollOffset+visibleItems {
		l.scrollOffset = l.cursor - visibleItems + 1
	}

	if l.scrollOffset < 0 {
		l.scrollOffset = 0
	}
}

// func (l *ListLayout) View() string {
// 	var b strings.Builder

// 	// Title
// 	if l.title != "" {
// 		b.WriteString(l.titleStyle.Render(l.title))
// 		b.WriteString("\n")
// 	}

// 	// Calculate visible items
// 	visibleItems := l.height - 2
// 	if l.title != "" {
// 		visibleItems--
// 	}
// 	if l.showHelp {
// 		visibleItems -= 2
// 	}
// 	if visibleItems < 1 {
// 		visibleItems = 1
// 	}

// 	// Render items
// 	endIdx := min(l.scrollOffset+visibleItems, len(l.items))
// 	for i := l.scrollOffset; i < endIdx; i++ {
// 		item := l.items[i]

// 		// Determine styling
// 		var style lipgloss.Style
// 		var prefix string

// 		if i == l.cursor && item.Selected {
// 			style = l.selectedCursorStyle
// 			prefix = "▶ ✓ "
// 		} else if i == l.cursor {
// 			style = l.cursorStyle
// 			prefix = "▶ "
// 		} else if item.Selected {
// 			style = l.selectedStyle
// 			prefix = "  ✓ "
// 		} else {
// 			style = l.normalStyle
// 			prefix = "    "
// 		}

// 		// Truncate if needed
// 		maxWidth := l.width - len(prefix) - 4
// 		if maxWidth < 0 {
// 			maxWidth = 0
// 		}
// 		text := item.Value
// 		if maxWidth > 3 && len(text) > maxWidth {
// 			text = text[:maxWidth-3] + "..."
// 		} else if len(text) > maxWidth && maxWidth > 0 {
// 			text = text[:maxWidth]
// 		} else if maxWidth == 0 {
// 			text = ""
// 		}

// 		b.WriteString(style.Render(prefix + text))
// 		b.WriteString("\n")
// 	}

// 	// Scroll indicator
// 	if l.scrollOffset > 0 {
// 		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("  ▲ More above"))
// 		b.WriteString("\n")
// 	}
// 	if endIdx < len(l.items) {
// 		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("  ▼ More below"))
// 		b.WriteString("\n")
// 	}

// 	// Help text
// 	if l.showHelp {
// 		selectedCount := 0
// 		for _, item := range l.items {
// 			if item.Selected {
// 				selectedCount++
// 			}
// 		}

// 		helpStyle := lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("240")).
// 			Padding(0, 2).
// 			MarginTop(1)

// 		var helpText string
// 		if l.maxSelections == 1 {
// 			helpText = fmt.Sprintf("↑↓=Move Space/Enter=Select ?=Help | Selected: %d", selectedCount)
// 		} else if l.maxSelections > 0 {
// 			helpText = fmt.Sprintf("↑↓=Move Space=Toggle a/A=All/None ?=Help | %d/%d selected", selectedCount, l.maxSelections)
// 		} else {
// 			helpText = fmt.Sprintf("↑↓=Move Space=Toggle a/A=All/None g/G=Top/Bottom ?=Help | %d selected", selectedCount)
// 		}

// 		b.WriteString(helpStyle.Render(helpText))
// 	}

// 	return b.String()
// }

func (l *ListLayout) View() string {
	var b strings.Builder

	// Title
	titleHeight := 0
	if l.title != "" {
		b.WriteString(l.titleStyle.Render(l.title))
		b.WriteString("\n")
		titleHeight = 2 // title + newline
	}

	// Calculate visible items
	helpHeight := 0
	if l.showHelp {
		helpHeight = 2
	}

	visibleItems := l.height - titleHeight - helpHeight
	if visibleItems < 1 {
		visibleItems = 1
	}

	// Render items
	endIdx := min(l.scrollOffset+visibleItems, len(l.items))
	renderedLines := 0

	for i := l.scrollOffset; i < endIdx && renderedLines < visibleItems; i++ {
		item := l.items[i]

		// Determine styling
		var style lipgloss.Style
		var prefix string

		if i == l.cursor && item.Selected {
			style = l.selectedCursorStyle
			prefix = "▶ ✓ "
		} else if i == l.cursor {
			style = l.cursorStyle
			prefix = "▶ "
		} else if item.Selected {
			style = l.selectedStyle
			prefix = "  ✓ "
		} else {
			style = l.normalStyle
			prefix = "    "
		}

		// Truncate if needed
		// maxWidth := l.width - len(prefix) - 4
		// if maxWidth < 0 {
		// 	maxWidth = 0
		// }
		// text := item.Value
		// if maxWidth > 3 && len(text) > maxWidth {
		// 	text = text[:maxWidth-3] + "..."
		// } else if len(text) > maxWidth && maxWidth > 0 {
		// 	text = text[:maxWidth]
		// } else if maxWidth == 0 {
		// 	text = ""
		// }

		// b.WriteString(style.Render(prefix + text))
		// Calculate available width (account for style padding)
		stylePadding := style.GetHorizontalFrameSize()
		maxWidth := l.width - len(prefix) - stylePadding
		if maxWidth < 0 {
			maxWidth = 0
		}

		text := item.Value
		if maxWidth > 3 && len(text) > maxWidth {
			text = text[:maxWidth-3] + "..."
		} else if len(text) > maxWidth && maxWidth > 0 {
			text = text[:maxWidth]
		} else if maxWidth == 0 {
			text = ""
		}

		// Render with explicit width to fill the space
		rendered := style.Width(l.width - stylePadding).Render(prefix + text)
		b.WriteString(rendered)
		b.WriteString("\n")
		renderedLines++
	}

	// Fill remaining space with blank lines
	for renderedLines < visibleItems {
		b.WriteString("\n")
		renderedLines++
	}

	// Help text at bottom
	if l.showHelp {
		selectedCount := 0
		for _, item := range l.items {
			if item.Selected {
				selectedCount++
			}
		}

		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 2)

		var helpText string
		if l.maxSelections == 1 {
			helpText = fmt.Sprintf("↑↓=Move Space/Enter=Select ?=Help | Selected: %d", selectedCount)
		} else if l.maxSelections > 0 {
			helpText = fmt.Sprintf("↑↓=Move Space=Toggle a/A=All/None ?=Help | %d/%d selected", selectedCount, l.maxSelections)
		} else {
			helpText = fmt.Sprintf("↑↓=Move Space=Toggle a/A=All/None g/G=Top/Bottom ?=Help | %d selected", selectedCount)
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render(helpText))
	}

	return b.String()
}

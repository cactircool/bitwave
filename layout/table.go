package layout

import (
	"fmt"
	"strings"

	"github.com/cactircool/bitwave/bindings"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableCell struct {
	Value    string
	Editable bool
}

type TableLayout struct {
	headers      []TableCell
	rows         [][]TableCell
	width        int
	height       int
	colWidths    []int
	selectedRow  int
	selectedCol  int
	editMode     bool
	editor       textarea.Model
	editingCell  [2]int // [row, col] - row=-1 means header
	scrollOffset int
	allowAddRows bool

	headerStyle       lipgloss.Style
	cellStyle         lipgloss.Style
	selectedStyle     lipgloss.Style
	editStyle         lipgloss.Style
	uneditableStyle   lipgloss.Style
}

func NewTableLayout(headers []string, editableHeaders bool) *TableLayout {
	headerCells := make([]TableCell, len(headers))
	for i, h := range headers {
		headerCells[i] = TableCell{Value: h, Editable: editableHeaders}
	}

	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = max(len(h), 10) // Minimum width of 10
	}

	editor := textarea.New()
	editor.SetHeight(1)
	editor.ShowLineNumbers = false
	editor.CharLimit = 500
	editor.Blur()

	return &TableLayout{
		headers:         headerCells,
		rows:            [][]TableCell{},
		colWidths:       colWidths,
		selectedRow:     0,
		selectedCol:     0,
		editMode:        false,
		editor:          editor,
		editingCell:     [2]int{-1, -1},
		scrollOffset:    0,
		allowAddRows:    false,
		headerStyle:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(0, 1),
		cellStyle:       lipgloss.NewStyle().Padding(0, 1),
		selectedStyle:   lipgloss.NewStyle().Background(lipgloss.Color("240")).Padding(0, 1),
		editStyle:       lipgloss.NewStyle().Background(lipgloss.Color("17")).Foreground(lipgloss.Color("15")).Padding(0, 1),
		uneditableStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1),
	}
}

func (t *TableLayout) AddRow(cells []string, editable []bool) {
	if len(cells) != len(t.headers) {
		return // Invalid row
	}

	row := make([]TableCell, len(cells))
	for i, c := range cells {
		ed := true
		if i < len(editable) {
			ed = editable[i]
		}
		row[i] = TableCell{Value: c, Editable: ed}

		// Update column width if needed
		if len(c) > t.colWidths[i] {
			t.colWidths[i] = len(c)
		}
	}

	t.rows = append(t.rows, row)
}

func (t *TableLayout) SetAllowAddRows(allow bool) {
	t.allowAddRows = allow
}

func (t *TableLayout) addNewRow() {
	// Create a new row with empty, editable cells
	row := make([]TableCell, len(t.headers))
	for i := range row {
		row[i] = TableCell{Value: "", Editable: true}
	}
	t.rows = append(t.rows, row)

	// Move to the new row
	t.selectedRow = len(t.rows) - 1
	t.adjustScroll()
}

func (t *TableLayout) SetSize(width, height int) {
	t.width = width
	t.height = height

	if len(t.colWidths) == 0 {
		return
	}

	// Calculate total borders and separators (fixed overhead)
	// Left border (1) + right border (1) + separators between columns
	totalBorderChars := 2 + (len(t.colWidths) - 1)

	// Each column also has padding from the style (2 chars per column from Padding(0, 1))
	totalPadding := len(t.colWidths) * 2

	// Available width for actual content
	availableWidth := width - totalBorderChars - totalPadding
	if availableWidth < len(t.colWidths) {
		availableWidth = len(t.colWidths) // At least 1 char per column
	}

	// Distribute width equally among columns
	baseWidth := availableWidth / len(t.colWidths)
	remainder := availableWidth % len(t.colWidths)

	for i := range t.colWidths {
		t.colWidths[i] = baseWidth
		if i < remainder {
			t.colWidths[i]++ // Distribute remainder
		}
		if t.colWidths[i] < 1 {
			t.colWidths[i] = 1
		}
	}
}

func (t *TableLayout) GetFocusState() FocusState {
	return Interactive
}

func (t *TableLayout) OnFocus(baseStyle lipgloss.Style) (lipgloss.Style, tea.Cmd) {
	return baseStyle.Border(lipgloss.ThickBorder()), nil
}

func (t *TableLayout) OnBlur() {
	t.editMode = false
	t.editor.Blur()
}

func (t *TableLayout) Init() tea.Cmd {
	return textarea.Blink
}

func (t *TableLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if t.editMode {
			switch key.String() {
			case bindings.CycleEscape:
				// Cancel editing
				t.editMode = false
				t.editor.Blur()
				t.editingCell = [2]int{-1, -1}
				return t, nil
			case bindings.CycleEnter:
				// Save edit
				value := strings.TrimSpace(t.editor.Value())
				if t.editingCell[0] == -1 {
					// Editing header
					t.headers[t.editingCell[1]].Value = value
				} else {
					// Editing cell
					t.rows[t.editingCell[0]][t.editingCell[1]].Value = value
				}

				// Update column width
				if len(value) > t.colWidths[t.editingCell[1]] {
					t.colWidths[t.editingCell[1]] = len(value)
				}

				t.editMode = false
				t.editor.Blur()
				t.editingCell = [2]int{-1, -1}
				return t, nil
			default:
				// Forward to editor
				var cmd tea.Cmd
				t.editor, cmd = t.editor.Update(msg)
				return t, cmd
			}
		}

		// Navigation mode
		switch key.String() {
		case "up", "k":
			if t.selectedRow > -1 {
				t.selectedRow--
				t.adjustScroll()
			}
			return t, nil
		case "down", "j":
			if t.selectedRow < len(t.rows)-1 {
				t.selectedRow++
				t.adjustScroll()
			}
			return t, nil
		case "left", "h":
			if t.selectedCol > 0 {
				t.selectedCol--
			}
			return t, nil
		case "right", "l":
			if t.selectedCol < len(t.headers)-1 {
				t.selectedCol++
			}
			return t, nil
		case "n":
			// Add new row (only if allowed)
			if t.allowAddRows {
				t.addNewRow()
			}
			return t, nil
		case "d":
			// Delete current row (only if allowed and not editing)
			if t.allowAddRows && t.selectedRow >= 0 && t.selectedRow < len(t.rows) {
				t.rows = append(t.rows[:t.selectedRow], t.rows[t.selectedRow+1:]...)
				if t.selectedRow >= len(t.rows) && t.selectedRow > 0 {
					t.selectedRow--
				}
				if len(t.rows) == 0 {
					t.selectedRow = 0
				}
				t.adjustScroll()
			}
			return t, nil
		case bindings.CycleEnter, "e", " ":
			// Enter edit mode
			var cell *TableCell
			var value string

			if t.selectedRow == -1 {
				// Editing header
				cell = &t.headers[t.selectedCol]
				value = cell.Value
				t.editingCell = [2]int{-1, t.selectedCol}
			} else {
				// Editing cell
				cell = &t.rows[t.selectedRow][t.selectedCol]
				value = cell.Value
				t.editingCell = [2]int{t.selectedRow, t.selectedCol}
			}

			if cell.Editable {
				t.editMode = true
				t.editor.SetValue(value)
				t.editor.SetWidth(t.colWidths[t.selectedCol])
				return t, t.editor.Focus()
			}
			return t, nil
		}
	}

	return t, nil
}

func (t *TableLayout) adjustScroll() {
	visibleRows := t.height - 3 // Account for header and borders
	if visibleRows < 1 {
		visibleRows = 1
	}

	if t.selectedRow < t.scrollOffset {
		t.scrollOffset = t.selectedRow
	}
	if t.selectedRow >= t.scrollOffset+visibleRows {
		t.scrollOffset = t.selectedRow - visibleRows + 1
	}

	if t.scrollOffset < 0 {
		t.scrollOffset = 0
	}
}

// func (t *TableLayout) View() string {
// 	var b strings.Builder

// 	// Render header
// 	b.WriteString("┌")
// 	for i, w := range t.colWidths {
// 		b.WriteString(strings.Repeat("─", w+2))
// 		if i < len(t.colWidths)-1 {
// 			b.WriteString("┬")
// 		}
// 	}
// 	b.WriteString("┐\n")

// 	// Header cells
// 	b.WriteString("│")
// 	for i, h := range t.headers {
// 		style := t.headerStyle
// 		if t.selectedRow == -1 && t.selectedCol == i && !t.editMode {
// 			style = t.selectedStyle
// 		}
// 		if t.editMode && t.editingCell[0] == -1 && t.editingCell[1] == i {
// 			b.WriteString(t.editStyle.Width(t.colWidths[i]).Render(t.editor.View()))
// 		} else {
// 			content := truncate(h.Value, t.colWidths[i])
// 			b.WriteString(style.Width(t.colWidths[i]).Render(content))
// 		}
// 		b.WriteString("│")
// 	}
// 	b.WriteString("\n")

// 	// Separator
// 	b.WriteString("├")
// 	for i, w := range t.colWidths {
// 		b.WriteString(strings.Repeat("─", w+2))
// 		if i < len(t.colWidths)-1 {
// 			b.WriteString("┼")
// 		}
// 	}
// 	b.WriteString("┤\n")

// 	// Render visible rows
// 	visibleRows := t.height - 4
// 	if visibleRows < 1 {
// 		visibleRows = 1
// 	}

// 	endRow := min(t.scrollOffset+visibleRows, len(t.rows))
// 	for r := t.scrollOffset; r < endRow; r++ {
// 		row := t.rows[r]
// 		b.WriteString("│")
// 		for c, cell := range row {
// 			var style lipgloss.Style

// 			if t.editMode && t.editingCell[0] == r && t.editingCell[1] == c {
// 				b.WriteString(t.editStyle.Width(t.colWidths[c]).Render(t.editor.View()))
// 			} else {
// 				if !cell.Editable {
// 					style = t.uneditableStyle
// 				} else if r == t.selectedRow && c == t.selectedCol {
// 					style = t.selectedStyle
// 				} else {
// 					style = t.cellStyle
// 				}
// 				content := truncate(cell.Value, t.colWidths[c])
// 				b.WriteString(style.Width(t.colWidths[c]).Render(content))
// 			}
// 			b.WriteString("│")
// 		}
// 		b.WriteString("\n")
// 	}

// 	// Bottom border
// 	b.WriteString("└")
// 	for i, w := range t.colWidths {
// 		b.WriteString(strings.Repeat("─", w+2))
// 		if i < len(t.colWidths)-1 {
// 			b.WriteString("┴")
// 		}
// 	}
// 	b.WriteString("┘")

// 	// Status line
// 	if t.editMode {
// 		b.WriteString("\n[EDIT MODE] Enter=Save Esc=Cancel")
// 	} else {
// 		helpText := fmt.Sprintf("↑↓←→=Navigate Enter/e/Space=Edit | Row %d/%d", t.selectedRow+1, len(t.rows))
// 		if t.allowAddRows {
// 			helpText += " | n=New d=Delete"
// 		}
// 		b.WriteString("\n" + helpText)
// 	}

// 	return b.String()
// }

func (t *TableLayout) View() string {
	var b strings.Builder

	// Render header
	b.WriteString("┌")
	for i, w := range t.colWidths {
		b.WriteString(strings.Repeat("─", w+2))
		if i < len(t.colWidths)-1 {
			b.WriteString("┬")
		}
	}
	b.WriteString("┐\n")

	// Header cells
	b.WriteString("│")
	for i, h := range t.headers {
		style := t.headerStyle
		if t.selectedRow == -1 && t.selectedCol == i && !t.editMode {
			style = t.selectedStyle
		}

		if t.editMode && t.editingCell[0] == -1 && t.editingCell[1] == i {
			// Editing a header cell
			content := truncate(t.editor.Value(), t.colWidths[i])
			b.WriteString(t.editStyle.Render(content))
		} else {
			content := truncate(h.Value, t.colWidths[i])
			b.WriteString(style.Render(content))
		}
		b.WriteString("│")
	}
	b.WriteString("\n")

	// Separator between header and body
	b.WriteString("├")
	for i, w := range t.colWidths {
		b.WriteString(strings.Repeat("─", w+2))
		if i < len(t.colWidths)-1 {
			b.WriteString("┼")
		}
	}
	b.WriteString("┤\n")

	// Render visible rows
	visibleRows := t.height - 4
	if visibleRows < 1 {
		visibleRows = 1
	}

	endRow := min(t.scrollOffset+visibleRows, len(t.rows))
	for r := t.scrollOffset; r < endRow; r++ {
		row := t.rows[r]
		b.WriteString("│")

		for c, cell := range row {
			var style lipgloss.Style

			if t.editMode && t.editingCell[0] == r && t.editingCell[1] == c {
				// Editing this cell
				content := truncate(t.editor.Value(), t.colWidths[c])
				b.WriteString(t.editStyle.Render(content))
			} else {
				// Normal cell rendering
				if !cell.Editable {
					style = t.uneditableStyle
				} else if r == t.selectedRow && c == t.selectedCol {
					style = t.selectedStyle
				} else {
					style = t.cellStyle
				}
				content := truncate(cell.Value, t.colWidths[c])
				b.WriteString(style.Render(content))
			}
			b.WriteString("│")
		}
		b.WriteString("\n")
	}

	// Bottom border
	b.WriteString("└")
	for i, w := range t.colWidths {
		b.WriteString(strings.Repeat("─", w+2))
		if i < len(t.colWidths)-1 {
			b.WriteString("┴")
		}
	}
	b.WriteString("┘")

	// Fill remaining height with empty rows if needed
	currentHeight := 4 + (endRow - t.scrollOffset) + 1 // borders + rows + status
	for currentHeight < t.height {
		b.WriteString("\n")
		currentHeight++
	}

	// Status line
	if t.editMode {
		b.WriteString("\n[EDIT MODE] Enter=Save Esc=Cancel")
	} else {
		helpText := fmt.Sprintf("↑↓←→=Navigate Enter/e/Space=Edit | Row %d/%d", t.selectedRow+1, len(t.rows))
		if t.allowAddRows {
			helpText += " | n=New d=Delete"
		}
		b.WriteString("\n" + helpText)
	}

	return b.String()
}

func truncate(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width < 3 {
		return s[:width]
	}
	return s[:width-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

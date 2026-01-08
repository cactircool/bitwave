package models

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Cell struct {
	Content  string
	Editable bool
	Style    lipgloss.Style
}

type TableModel struct {
	// Data
	Headers []string
	Rows    [][]Cell

	// Viewport state
	selectedRow int
	selectedCol int
	scrollOffsetRow int
	scrollOffsetCol int

	// Edit mode
	editing bool
	editBuffer string
	editPos int

	// Display settings
	columnWidths []int
	rowHeight int
	visibleRows int
	visibleCols int

	// Styles
	headerStyle lipgloss.Style
	selectedStyle lipgloss.Style
	normalStyle lipgloss.Style
	editingStyle lipgloss.Style
	uneditableStyle lipgloss.Style
}

type TablePane struct {
	model    *TableModel
	Viewport viewport.Model
	style    lipgloss.Style
	weight   float32
	gap      int
}

func NewTableModel(headers []string, rows int) *TableModel {
	cols := len(headers)
	// Initialize cells
	cellRows := make([][]Cell, rows)
	for i := range cellRows {
		cellRows[i] = make([]Cell, cols)
		for j := range cellRows[i] {
			cellRows[i][j] = Cell{
				Content:  "",
				Editable: true,
			}
		}
	}

	// Default column widths
	columnWidths := make([]int, cols)
	for i := range columnWidths {
		columnWidths[i] = 15 // default width
	}

	return &TableModel{
		Headers: headers,
		Rows: cellRows,
		columnWidths: columnWidths,
		rowHeight: 1,
		selectedRow: 0,
		selectedCol: 0,

		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(TableHeaderForeground)).
			Background(lipgloss.Color(TableHeaderBackground)).
			Padding(0, 1),

		selectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(TableSelectedForeground)).
			Background(lipgloss.Color(TableSelectedBackground)).
			Padding(0, 1),

		normalStyle: lipgloss.NewStyle().
			Padding(0, 1),

		editingStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(TableEditingForeground)).
			Background(lipgloss.Color(TableEditingBackground)).
			Padding(0, 1),

		uneditableStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(TableUneditableForeground)).
			Padding(0, 1),
	}
}

func NewTablePane(headers []string, rows int, weight float32, style lipgloss.Style) Pane {
	return &TablePane{
		model:    NewTableModel(headers, rows),
		Viewport: viewport.New(0, 0),
		style:    style,
		weight:   weight,
		gap:      0,
	}
}

func NewTablePaneWithGap(headers []string, rows int, cols int, weight float32, style lipgloss.Style, gap int) Pane {
	return &TablePane{
		model:    NewTableModel(headers, rows),
		Viewport: viewport.New(0, 0),
		style:    style,
		weight:   weight,
		gap:      gap,
	}
}

// Pane interface implementation
func (p *TablePane) Style() lipgloss.Style {
	return p.style
}

func (p *TablePane) SetStyle(style lipgloss.Style) {
	p.style = style
}

func (p *TablePane) Model() tea.Model {
	return p.model
}

func (p *TablePane) Gap() int {
	return p.gap
}

func (p *TablePane) Weight() float32 {
	return p.weight
}

func (p *TablePane) Resize(outerW, outerH int) {
	chromeW, chromeH := chromeSize(p.style)

	contentW := outerW - chromeW
	contentH := outerH - chromeH

	if contentW < 0 {
		contentW = 0
	}
	if contentH < 0 {
		contentH = 0
	}

	p.Viewport.Width = contentW
	p.Viewport.Height = contentH

	// Calculate visible area for the table based on new dimensions
	p.model.calculateVisibleArea(contentW, contentH)
}

func (p *TablePane) Update(w, h int, msg tea.Msg) tea.Cmd {
	// Update the table model with the message
	cmd := p.model.update(w, h, msg)

	// Update viewport content
	p.Viewport.SetContent(p.model.View())

	return cmd
}

func (p *TablePane) View() string {
	return p.style.Render(p.Viewport.View())
}

// TableModel methods

func (m *TableModel) Init() tea.Cmd { return nil }

func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// This is for tea.Model interface compatibility
	// Actual updates go through the update method below
	return m, nil
}

func (m *TableModel) update(_, _ int, msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		if m.editing {
			return m.handleEditMode(key)
		}
		return m.handleNavigationMode(key)
	}

	return nil
}

func (m *TableModel) handleNavigationMode(key tea.KeyMsg) tea.Cmd {
	keyFound := func(keys []string) bool {
		return slices.Contains(keys, key.String())
	}

	switch {
	case keyFound(TableUp):
		if m.selectedRow > 0 {
			m.selectedRow--
			m.adjustScrollOffset()
		}
	case keyFound(TableDown):
		if m.selectedRow < len(m.Rows)-1 {
			m.selectedRow++
			m.adjustScrollOffset()
		}
	case keyFound(TableLeft):
		if m.selectedCol > 0 {
			m.selectedCol--
			m.adjustScrollOffset()
		}
	case keyFound(TableRight):
		if m.selectedCol < len(m.columnWidths)-1 {
			m.selectedCol++
			m.adjustScrollOffset()
		}
	case keyFound(TableCellEnter):
		// Enter edit mode if cell is editable
		if m.selectedRow < len(m.Rows) && m.selectedCol < len(m.Rows[m.selectedRow]) {
			if m.Rows[m.selectedRow][m.selectedCol].Editable {
				m.editing = true
				m.editBuffer = m.Rows[m.selectedRow][m.selectedCol].Content
			}
		}
	case keyFound(TableCellDelete):
		// Delete cell content
		if m.selectedRow < len(m.Rows) && m.selectedCol < len(m.Rows[m.selectedRow]) {
			if m.Rows[m.selectedRow][m.selectedCol].Editable {
				m.Rows[m.selectedRow][m.selectedCol].Content = ""
			}
		}
	}
	return nil
}

func (m *TableModel) handleEditMode(key tea.KeyMsg) tea.Cmd {
	keyFound := func(keys []string) bool {
		return slices.Contains(keys, key.String())
	}

	switch {
	case keyFound(TableCellEscapeEdit):
		// Save the edited content
		if m.selectedRow < len(m.Rows) && m.selectedCol < len(m.Rows[m.selectedRow]) {
			m.Rows[m.selectedRow][m.selectedCol].Content = m.editBuffer
		}
		m.editing = false
		m.editBuffer = ""
		m.editPos = 0
	case key.String() == "left":
		m.editPos = max(0, m.editPos - 1)
	case key.String() == "right":
		m.editPos = min(len(m.editBuffer) - 1, m.editPos + 1)
	case key.String() == "up" || key.String() == "down":
		// nothing
	case key.String() == "backspace":
		if len(m.editBuffer) > 0 && m.editPos >= 0 {
			// panic(fmt.Sprintf("panicing editBuffer len = %d, editPos = %d", len(m.editBuffer), m.editPos))
			if m.editPos >= len(m.editBuffer) {
				m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
			} else {
				m.editBuffer = m.editBuffer[:m.editPos] + m.editBuffer[m.editPos+1:]
			}
			m.editPos = max(0, m.editPos - 1)
		}
	default:
		// Add character to buffer
		if len(key.String()) == 1 {
			m.editBuffer += key.String()
			m.editPos++
		}
	}
	return nil
}

func (m *TableModel) calculateVisibleArea(width, height int) {
	// Calculate how many rows and columns can fit
	totalWidth := 0
	m.visibleCols = 0
	for i := m.scrollOffsetCol; i < len(m.columnWidths); i++ {
		if totalWidth + m.columnWidths[i] + 3 > width {
			break
		}
		totalWidth += m.columnWidths[i] + 3
		m.visibleCols++
	}

	// Reserve 2 lines for header and status
	availableHeight := max(0, height - 3)
	m.visibleRows = min(len(m.Rows), availableHeight / (m.rowHeight + 1))
}

func (m *TableModel) adjustScrollOffset() {
	// Adjust row scroll
	if m.selectedRow < m.scrollOffsetRow {
		m.scrollOffsetRow = m.selectedRow
	} else if m.selectedRow >= m.scrollOffsetRow + m.visibleRows {
		m.scrollOffsetRow = m.selectedRow - m.visibleRows + 1
	}

	// Adjust column scroll
	if m.selectedCol < m.scrollOffsetCol {
		m.scrollOffsetCol = m.selectedCol
	} else if m.selectedCol >= m.scrollOffsetCol + m.visibleCols {
		m.scrollOffsetCol = m.selectedCol - m.visibleCols + 1
	}
}

func (m *TableModel) View() string {
	var b strings.Builder

	// Render headers
	headerRow := m.renderHeaderRow()
	b.WriteString(headerRow)
	b.WriteString("\n")

	// Render separator
	b.WriteString(m.renderSeparator())
	b.WriteString("\n")

	// Render visible rows
	endRow := min(len(m.Rows), m.scrollOffsetRow + m.visibleRows)

	for i := m.scrollOffsetRow; i < endRow; i++ {
		b.WriteString(m.renderRow(i))
		b.WriteString("\n")
	}

	// Render status line
	b.WriteString(m.renderStatusLine())

	return b.String()
}

func (m *TableModel) renderHeaderRow() string {
	var cells []string
	endCol := min(len(m.Headers), m.scrollOffsetCol + m.visibleCols)

	for i := m.scrollOffsetCol; i < endCol; i++ {
		content := m.Headers[i]
		if len(content) > m.columnWidths[i] {
			content = content[:m.columnWidths[i]-3] + "..."
		}
		cell := m.headerStyle.Width(m.columnWidths[i]).Render(content)
		cells = append(cells, cell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

func (m *TableModel) renderSeparator() string {
	var parts []string

	endCol := min(len(m.columnWidths), m.scrollOffsetCol + m.visibleCols)

	for i := m.scrollOffsetCol; i < endCol; i++ {
		parts = append(parts, strings.Repeat("─", m.columnWidths[i]+2))
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Join(parts, "┼"))
}

func (m *TableModel) renderRow(rowIdx int) string {
	var cells []string

	endCol := min(len(m.Rows[rowIdx]), m.scrollOffsetCol + m.visibleCols)

	for colIdx := m.scrollOffsetCol; colIdx < endCol; colIdx++ {
		cell := m.Rows[rowIdx][colIdx]

		// Determine content to display
		content := cell.Content
		if m.editing && rowIdx == m.selectedRow && colIdx == m.selectedCol {
			if m.editPos < len(m.editBuffer) {
				content = m.editBuffer[:m.editPos+1] + "█" + m.editBuffer[m.editPos+1:]
			} else {
				// Cursor at end
				content = m.editBuffer + "█"
			}
		}

		// Truncate if needed
		if !m.editing && len(content) > m.columnWidths[colIdx] {
			content = content[:m.columnWidths[colIdx]-3] + "..."
		}

		// Apply style
		var style lipgloss.Style
		if m.editing && rowIdx == m.selectedRow && colIdx == m.selectedCol {
			style = m.editingStyle
		} else if rowIdx == m.selectedRow && colIdx == m.selectedCol {
			style = m.selectedStyle
		} else if !cell.Editable {
			style = m.uneditableStyle
		} else {
			style = m.normalStyle
		}

		cellStr := style.Width(m.columnWidths[colIdx]).Render(content)
		cells = append(cells, cellStr)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

func (m *TableModel) renderStatusLine() string {
	mode := "NORMAL"
	if m.editing {
		mode = "EDIT"
	}

	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(
			"[" + mode + "] " +
			"Row: " + strconv.Itoa(m.selectedRow+1) + "/" + string(rune(len(m.Rows))) + " " +
			"Col: " + strconv.Itoa(m.selectedCol+1) + "/" + string(rune(len(m.columnWidths))) +
			fmt.Sprintf(" | ←↓↑→:navigate %s:edit %s:delete", strings.Join(TableCellEnter, ","), strings.Join(TableCellDelete, ",")),
		)

	return status
}

// Helper methods for TablePane

func (p *TablePane) SetCell(row, col int, content string, editable bool) {
	if row < len(p.model.Rows) && col < len(p.model.Rows[row]) {
		p.model.Rows[row][col].Content = content
		p.model.Rows[row][col].Editable = editable
	}
}

func (p *TablePane) SetColumnWidth(col int, width int) {
	if col < len(p.model.columnWidths) {
		p.model.columnWidths[col] = width
	}
}

func (p *TablePane) GetCell(row, col int) (string, bool) {
	if row < len(p.model.Rows) && col < len(p.model.Rows[row]) {
		return p.model.Rows[row][col].Content, true
	}
	return "", false
}

func (p *TablePane) AddRow() {
	newRow := make([]Cell, len(p.model.columnWidths))
	for i := range newRow {
		newRow[i] = Cell{
			Content:  "",
			Editable: true,
		}
	}
	p.model.Rows = append(p.model.Rows, newRow)
}

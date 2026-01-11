package app

import (
	"github.com/cactircool/bitwave/layout"
	"github.com/charmbracelet/lipgloss"
)

func ConstructRoot() *layout.RootLayout {
	root := layout.NewRootLayout(layout.Vertical)

	constructHeader(root)
	constructMain(root)
	// constructFooter(root)

	return root
}

func constructHeader(root *layout.RootLayout) {
	header := layout.NewTextView("bitwave")
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")).
		Padding(0, 1)
	root.AddStatic(header, 1, headerStyle, 0)
}

func constructMain(root *layout.RootLayout) {
	main := layout.NewLayout(layout.Horizontal)
	mainStyle := lipgloss.NewStyle()

	main.Add(layout.NewTextView("yuh"), 1, lipgloss.NewStyle().Border(lipgloss.ASCIIBorder()), 0)

	// box := layout.NewTextareaLayout(textarea.New())
	// main.Add(box, 1, lipgloss.NewStyle(), 1)

	// Left column - Todo List
	// todoList := layout.NewListLayout("Todo List", 0)
	// todoList.AddItems([]string{
	// 	"Review PRs",
	// 	"Write docs",
	// 	"Fix bug #123",
	// })
	// main.Add(todoList, 1, lipgloss.NewStyle().Border(lipgloss.RoundedBorder()), 0)

	// Center column - Table with add/delete
	// userTable := layout.NewTableLayout([]string{"Name", "Role"}, false)
	// userTable.SetAllowAddRows(true) // Enable row management
	// userTable.AddRow(
	// 	[]string{"Alice", "Engineer"},
	// 	[]bool{true, true},
	// )
	// userTable.AddRow(
	// 	[]string{"Bob", "Designer"},
	// 	[]bool{true, true},
	// )
	// main.Add(userTable, 1, lipgloss.NewStyle().Border(lipgloss.RoundedBorder()), 1)

	root.Add(main, 1, mainStyle, 0)
}

func constructFooter(root *layout.RootLayout) {
	footer := layout.NewTextView("Tab/Shift+Tab: Navigate | Enter: Focus/Edit | Esc: Exit | Ctrl+C: Quit")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)
	root.AddStatic(footer, 1, footerStyle, 0)
}

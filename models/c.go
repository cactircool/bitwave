package models

var (
	NextPane = "tab"
	EscapePane = "esc"
	EnterPane = "enter"
	Exit = "ctrl+c"

	TableUp = []string{"up", "k"}
	TableDown = []string{"down", "j"}
	TableLeft = []string{"left", "h"}
	TableRight = []string{"right", "l"}
	TableCellEnter = []string{"enter", "i", "e"}
	TableCellDelete = []string{"d"}

	TableCellEscapeEdit = []string{"enter"}

	TableHeaderForeground = "6"
	TableHeaderBackground = "235"

	TableSelectedForeground = "15"
	TableSelectedBackground = "238"

	TableEditingForeground = "15"
	TableEditingBackground = "22"

	TableUneditableForeground = "240"
)

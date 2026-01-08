// package main

// import (
// 	"log"

// 	"github.com/cactircool/bitwave/models"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// func main() {
// 	root := models.NewRootModel(
// 		true,
// 		[]models.Pane{
// 			// Left column with 3 stacked items and gaps between them
// 			models.NewPaneWithGap(
// 				models.NewNestedModel(
// 					false,
// 					[]models.Pane{
// 						models.NewPaneWithGap(
// 							models.NewBasicModel("top"),
// 							1,
// 							lipgloss.NewStyle().
// 								Border(lipgloss.RoundedBorder()),
// 							1, // 1 line gap after this pane
// 						),
// 						models.NewPaneWithGap(
// 							models.NewBasicModel("center"),
// 							1,
// 							lipgloss.NewStyle().
// 								Border(lipgloss.RoundedBorder()),
// 							1, // 1 line gap after this pane
// 						),
// 						models.NewPane(
// 							models.NewBasicModel("bottom"),
// 							1,
// 							lipgloss.NewStyle().
// 								Border(lipgloss.RoundedBorder()),
// 						),
// 					},
// 				),
// 				1,
// 				lipgloss.NewStyle(),
// 				2, // 2 character gap after this column
// 			),
// 			// Middle pane
// 			models.NewPaneWithGap(
// 				models.NewBasicModel("hello from the middle!"),
// 				1,
// 				lipgloss.NewStyle().
// 					Border(lipgloss.RoundedBorder()).
// 					Padding(3),
// 				2, // 2 character gap after this pane
// 			),
// 			// Right column with spacers for vertical centering
// 			models.NewPane(
// 				models.NewNestedModel(
// 					false,
// 					[]models.Pane{
// 						// Top spacer - takes up 1 unit of space
// 						models.NewPane(
// 							models.NewSpacerModel(),
// 							1,
// 							lipgloss.NewStyle(),
// 						),
// 						// Centered content - takes up 0.5 units of space
// 						models.NewPane(
// 							models.NewBasicModel("I'm centered vertically!"),
// 							0.5,
// 							lipgloss.NewStyle().
// 								Border(lipgloss.ThickBorder()).
// 								Padding(2),
// 						),
// 						// Bottom spacer - takes up 1 unit of space
// 						models.NewPane(
// 							models.NewSpacerModel(),
// 							1,
// 							lipgloss.NewStyle(),
// 						),
// 					},
// 				),
// 				1,
// 				lipgloss.NewStyle(),
// 			),
// 		},
// 	)

// 	p := tea.NewProgram(root, tea.WithAltScreen())
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"log"

	"github.com/cactircool/bitwave/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	root := models.NewRootModel(
		true,
		[]models.Pane{
			models.NewTablePane(
				[]string{
					"Name",
					"Age",
					"Occupation",
				},
				10,
				1.0,
				lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
			),
			// models.NewPane(
			// 	models.NewNestedModel(
			// 		false,
			// 		[]models.Pane{
			// 			models.NewPane(
			// 				models.NewBasicModel("top"),
			// 				1,
			// 				lipgloss.NewStyle().
			// 					Border(lipgloss.RoundedBorder()),
			// 			),
			// 			models.NewPane(
			// 				models.NewBasicModel("center"),
			// 				1,
			// 				lipgloss.NewStyle().
			// 					Border(lipgloss.RoundedBorder()),
			// 			),
			// 			models.NewPane(
			// 				models.NewBasicModel("bottom"),
			// 				1,
			// 				lipgloss.NewStyle().
			// 					Border(lipgloss.RoundedBorder()),
			// 			),
			// 			models.NewPane(
			// 				models.NewSpacerModel(),
			// 				2,
			// 				lipgloss.NewStyle(),
			// 			),
			// 		},
			// 	),
			// 	1,
			// 	lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1),
			// ),
			// models.NewPane(
			// 	models.NewTextareaModel(
			// 		"placeholder",
			// 		false,
			// 		lipgloss.NewStyle().Border(lipgloss.ASCIIBorder()),
			// 	),
			// 	1,
			// 	lipgloss.NewStyle().
			// 		Border(lipgloss.RoundedBorder()),
			// ),
			// models.NewPane(
			// 	models.NewBasicModel("hello from the right!"),
			// 	2,
			// 	lipgloss.NewStyle().
			// 		Border(lipgloss.RoundedBorder()).
			// 		Padding(0),
			// ),
		},
	)

	p := tea.NewProgram(root, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

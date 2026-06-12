package main

import (
	"log"

	"github.com/cactircool/bitwave/app"
	"github.com/cactircool/bitwave/cmd"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cmd.Execute()

	root := app.ConstructRoot()
	p := tea.NewProgram(root, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

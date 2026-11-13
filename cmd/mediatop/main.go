package main

import (
	"log"

	"github.com/Harxhit/mediatop/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	p := tea.NewProgram(ui.Model{})

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}

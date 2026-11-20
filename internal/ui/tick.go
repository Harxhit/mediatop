package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

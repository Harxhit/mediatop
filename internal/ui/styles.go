package ui

import "github.com/charmbracelet/lipgloss"

var (
	cGreen   = lipgloss.Color("2")
	cYellow  = lipgloss.Color("3")
	cRed     = lipgloss.Color("1")
	cCyan    = lipgloss.Color("6")
	cMagenta = lipgloss.Color("5")
	cBlue    = lipgloss.Color("12")
	cDim     = lipgloss.Color("8")

	labelStyle  = lipgloss.NewStyle().Bold(true).Foreground(cCyan)
	valueStyle  = lipgloss.NewStyle().Bold(true)
	dimStyle    = lipgloss.NewStyle().Foreground(cDim)
	errStyle    = lipgloss.NewStyle().Bold(true).Foreground(cRed)
	rxStyle     = lipgloss.NewStyle().Foreground(cGreen)
	txStyle     = lipgloss.NewStyle().Foreground(cMagenta)
	barGreen    = lipgloss.NewStyle().Foreground(cGreen)
	barYellow   = lipgloss.NewStyle().Foreground(cYellow)
	barRed      = lipgloss.NewStyle().Foreground(cRed)
	titleBar    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("4"))
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(cGreen)
	headerSort  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(cCyan)
	selStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(cCyan)
	fKeyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("0"))
	fLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(cCyan)
)

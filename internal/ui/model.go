package ui

import (
	"fmt"
	"strings"

	"github.com/Harxhit/mediatop/internal/system"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cpuUsage    float64
	cpuError    string
	width       int
	height      int
	memoryUsage float64
	memoryError string
}

type cpuUsageMsg struct {
	usage float64
	err   error
}

type memoryUsageMsg struct {
	usage float64
	err   error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7DD3FC"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ADE80")).
			Bold(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#475569")).
			Padding(1, 2)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))

	valueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8FAFC"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64748B"))

	barActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#38BDF8"))

	barInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#334155"))
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchCPU(), fetchMemory())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case cpuUsageMsg:
		if msg.err != nil {
			m.cpuError = msg.err.Error()
			return m, fetchCPU()
		}

		m.cpuError = ""
		m.cpuUsage = msg.usage
		return m, fetchCPU()

	case memoryUsageMsg:
		if msg.err != nil {
			m.memoryError = msg.err.Error()
			return m, fetchMemory()
		}

		m.memoryError = ""
		m.memoryUsage = msg.usage
		return m, fetchMemory()
	}

	return m, nil
}

func (m Model) View() string {
	width := m.width - 4
	if width < 40 {
		width = 40
	}
	if width > 90 {
		width = 90
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		titleStyle.Render("MEDIATOP"),
		"    ",
		statusStyle.Render("● LIVE"),
	)

	barWidth := width - 18
	if barWidth < 10 {
		barWidth = 10
	}
	if barWidth > 55 {
		barWidth = 55
	}

	usage := m.cpuUsage
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}

	filled := int(float64(barWidth) * usage / 100)
	empty := barWidth - filled

	bar := barActive.Render(strings.Repeat("█", filled)) +
		barInactive.Render(strings.Repeat("░", empty))

	cpuPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("CPU UTILIZATION"),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			valueStyle.Render(fmt.Sprintf("%.2f%%", m.cpuUsage)),
			"   ",
			bar,
		),
		"",
		labelStyle.Render("Aggregate CPU  •  All logical processors"),
	)

	memoryPanel := lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render("MEMORY UTILIZATION"),
		"",
		valueStyle.Render(fmt.Sprintf("%.2f%%", m.memoryUsage)),
		"",
		labelStyle.Render("RAM  •  Used / Available"),
	)

	panel := lipgloss.JoinVertical(
		lipgloss.Left,
		panelStyle.Width(width-6).Render(cpuPanel),
		panelStyle.Width(width-6).Render(memoryPanel),
	)

	footer := footerStyle.Render("q: quit    •    Refresh: 1s")

	if m.cpuError != "" {
		footer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F87171")).
			Render("CPU ERROR: " + m.cpuError)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		panel,
		"",
		footer,
	)
}

func fetchCPU() tea.Cmd {
	return func() tea.Msg {
		cpuStats, err := system.ReadCPUStats()
		if err != nil {
			return cpuUsageMsg{err: err}
		}

		delta := system.CalculateTotalDelta(cpuStats)
		usage := system.CalculateCPUutilization(cpuStats.Idle, delta)

		return cpuUsageMsg{
			usage: usage,
		}
	}
}

func fetchMemory() tea.Cmd {
	return func() tea.Msg {
		stats, err := system.ReadMemoryStats()

		if err != nil {
			return memoryUsageMsg{err: err}
		}

		usage := system.CalculateMemoryUsage(stats)

		return memoryUsageMsg{
			usage: usage,
		}
	}
}

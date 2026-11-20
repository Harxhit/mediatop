package ui

import (
	"github.com/Harxhit/mediatop/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cpuUsage     float64
	cpuError     string
	width        int
	height       int
	memoryUsage  float64
	memoryError  string
	networkRates []system.NetworkRate
	networkError string

	selected int
	sortBy   sortKey
	paused   bool

	cpuHist []float64
	memHist []float64
	rxHist  []float64
	txHist  []float64
}

type cpuUsageMsg struct {
	usage float64
	err   error
}

type memoryUsageMsg struct {
	usage float64
	err   error
}

type networkUsageMsg struct {
	rates []system.NetworkRate
	err   error
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchCPU(), fetchMemory(), fetchNetwork(), Tick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "f10":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.networkRates)-1 {
				m.selected++
			}
		case "s":
			name := SelectedName(m.networkRates, m.selected)
			m.sortBy = (m.sortBy + 1) % sortKeyCount
			SortRates(m.networkRates, m.sortBy)
			m.selected = IndexOf(m.networkRates, name, m.selected)
		case "p", " ":
			m.paused = !m.paused
		}
		return m, nil

	case tickMsg:
		if m.paused {
			return m, Tick()
		}
		return m, tea.Batch(fetchCPU(), fetchMemory(), fetchNetwork(), Tick())

	case cpuUsageMsg:
		if msg.err != nil {
			m.cpuError = msg.err.Error()
			return m, nil
		}
		m.cpuError = ""
		m.cpuUsage = msg.usage
		m.cpuHist = AppendHist(m.cpuHist, msg.usage)
		return m, nil

	case memoryUsageMsg:
		if msg.err != nil {
			m.memoryError = msg.err.Error()
			return m, nil
		}
		m.memoryError = ""
		m.memoryUsage = msg.usage
		m.memHist = AppendHist(m.memHist, msg.usage)
		return m, nil

	case networkUsageMsg:
		if msg.err != nil {
			m.networkError = msg.err.Error()
			return m, nil
		}
		m.networkError = ""

		name := SelectedName(m.networkRates, m.selected)
		SortRates(msg.rates, m.sortBy)
		m.networkRates = msg.rates
		m.selected = IndexOf(m.networkRates, name, m.selected)

		var rx, tx float64
		for _, r := range m.networkRates {
			rx += float64(r.RXBytesPerSec)
			tx += float64(r.TXBytesPerSec)
		}
		m.rxHist = AppendHist(m.rxHist, rx)
		m.txHist = AppendHist(m.txHist, tx)
		return m, nil
	}

	return m, nil
}

func fetchCPU() tea.Cmd {
	return func() tea.Msg {
		cpuStats, err := system.ReadCPUStats()
		if err != nil {
			return cpuUsageMsg{err: err}
		}

		delta := system.CalculateTotalDelta(cpuStats)
		usage := system.CalculateCPUutilization(cpuStats.Idle, delta)

		return cpuUsageMsg{usage: usage}
	}
}

func fetchMemory() tea.Cmd {
	return func() tea.Msg {
		stats, err := system.ReadMemoryStats()
		if err != nil {
			return memoryUsageMsg{err: err}
		}

		return memoryUsageMsg{usage: system.CalculateMemoryUsage(stats)}
	}
}

func fetchNetwork() tea.Cmd {
	return func() tea.Msg {
		rates, err := system.ReadNetworkRates()
		if err != nil {
			return networkUsageMsg{err: err}
		}

		return networkUsageMsg{rates: rates}
	}
}

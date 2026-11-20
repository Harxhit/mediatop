package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	w := m.width
	if w < 64 {
		w = 64
	}
	colL := (w - 2) / 2
	colR := w - 2 - colL

	// title bar
	status := "● LIVE"
	if m.paused {
		status = "❚❚ PAUSED"
	}
	left := " MEDIATOP"
	right := status + "   " + time.Now().Format("15:04:05") + " "
	gap := w - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	title := titleBar.Render(left + strings.Repeat(" ", gap) + right)

	// meters + history
	meters := Meter("CPU", m.cpuUsage, fmt.Sprintf("%.1f%%", m.cpuUsage), colL) + "  " +
		Meter("Mem", m.memoryUsage, fmt.Sprintf("%.1f%%", m.memoryUsage), colR)

	indent := "    "
	sparks := indent + barGreen.Render(Spark(m.cpuHist, colL-5, 100)) + " " + "  " +
		indent + lipgloss.NewStyle().Foreground(cBlue).Render(Spark(m.memHist, colR-5, 100)) + " "

	netRow := NetCell("RX", LastOf(m.rxHist), m.rxHist, colL, cGreen) + "  " +
		NetCell("TX", LastOf(m.txHist), m.txHist, colR, cMagenta)

	// table
	barW := (w - 12*3 - 4) / 2
	if barW < 6 {
		barW = 6
	}
	if barW > 40 {
		barW = 40
	}

	hc := func(text string, active bool) string {
		if active {
			return headerSort.Render(text)
		}
		return headerStyle.Render(text)
	}
	sp := headerStyle.Render(" ")
	header := hc(fmt.Sprintf("%-12s", "IFACE ▲"), m.sortBy == sortName) + sp +
		hc(fmt.Sprintf("%12s", "RX/s ▼"), m.sortBy == sortRX) + sp +
		headerStyle.Render(strings.Repeat(" ", barW)) + sp +
		hc(fmt.Sprintf("%12s", "TX/s ▼"), m.sortBy == sortTX) + sp +
		headerStyle.Render(strings.Repeat(" ", barW))
	if rest := w - lipgloss.Width(header); rest > 0 {
		header += headerStyle.Render(strings.Repeat(" ", rest))
	}

	maxRate := 1.0
	for _, r := range m.networkRates {
		if v := float64(r.RXBytesPerSec); v > maxRate {
			maxRate = v
		}
		if v := float64(r.TXBytesPerSec); v > maxRate {
			maxRate = v
		}
	}

	// fixed lines: title, meters, sparks, net, blank, header, footer
	avail := len(m.networkRates)
	if m.height > 0 {
		avail = m.height - 7
	}
	if avail < 1 {
		avail = 1
	}
	start := 0
	if m.selected >= avail {
		start = m.selected - avail + 1
	}
	end := start + avail
	if end > len(m.networkRates) {
		end = len(m.networkRates)
	}

	lines := []string{title, meters, sparks, netRow, "", header}

	if len(m.networkRates) == 0 {
		lines = append(lines, dimStyle.Render("  no network interfaces reported"))
	}

	for i := start; i < end; i++ {
		r := m.networkRates[i]
		rx, tx := float64(r.RXBytesPerSec), float64(r.TXBytesPerSec)
		name := fmt.Sprintf("%-12s", Truncate(r.Name, 12))

		if i == m.selected {
			line := fmt.Sprintf("%s %12s %s %12s %s",
				name,
				HumanRate(rx), LoadBar(rx, maxRate, barW, rxStyle, false),
				HumanRate(tx), LoadBar(tx, maxRate, barW, txStyle, false),
			)
			lines = append(lines, selStyle.Width(w).Render(line))
			continue
		}

		rxText := fmt.Sprintf("%12s", HumanRate(rx))
		txText := fmt.Sprintf("%12s", HumanRate(tx))
		if rx == 0 {
			rxText = dimStyle.Render(rxText)
		} else {
			rxText = rxStyle.Render(rxText)
		}
		if tx == 0 {
			txText = dimStyle.Render(txText)
		} else {
			txText = txStyle.Render(txText)
		}

		lines = append(lines,
			valueStyle.Render(name)+" "+
				rxText+" "+LoadBar(rx, maxRate, barW, rxStyle, true)+" "+
				txText+" "+LoadBar(tx, maxRate, barW, txStyle, true),
		)
	}

	upper := strings.Join(lines, "\n")
	if pad := m.height - lipgloss.Height(upper) - 1; pad > 0 {
		upper += strings.Repeat("\n", pad)
	}

	pauseLabel := " Pause "
	if m.paused {
		pauseLabel = " Resume "
	}
	footer := Fkeys(w, [][2]string{
		{"↑↓", "Select "},
		{"S", "Sort:" + m.sortBy.String() + " "},
		{"P", pauseLabel},
		{"Q", "Quit "},
	})
	switch {
	case m.cpuError != "":
		footer = errStyle.Render("CPU ERROR: " + m.cpuError)
	case m.memoryError != "":
		footer = errStyle.Render("MEM ERROR: " + m.memoryError)
	case m.networkError != "":
		footer = errStyle.Render("NET ERROR: " + m.networkError)
	}

	return upper + "\n" + footer
}

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var sparkBlocks = []rune("▁▂▃▄▅▆▇█")

func Spark(vals []float64, width int, ceil float64) string {
	if width < 1 {
		return ""
	}
	if len(vals) > width {
		vals = vals[len(vals)-width:]
	}
	if ceil <= 0 {
		ceil = 1
		for _, v := range vals {
			if v > ceil {
				ceil = v
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(strings.Repeat(" ", width-len(vals)))
	for _, v := range vals {
		idx := int(v / ceil * float64(len(sparkBlocks)-1))
		if idx < 0 {
			idx = 0
		}
		if idx > len(sparkBlocks)-1 {
			idx = len(sparkBlocks) - 1
		}
		sb.WriteRune(sparkBlocks[idx])
	}
	return sb.String()
}

func Meter(label string, pct float64, text string, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	const textW = 6
	inner := width - lipgloss.Width(label) - 2
	if inner < 14 {
		inner = 14
	}
	barLen := inner - textW - 1
	filled := int(float64(barLen)*pct/100 + 0.5)

	var sb strings.Builder
	for i := 0; i < filled; i++ {
		pos := float64(i) / float64(barLen) * 100
		switch {
		case pos >= 80:
			sb.WriteString(barRed.Render("|"))
		case pos >= 50:
			sb.WriteString(barYellow.Render("|"))
		default:
			sb.WriteString(barGreen.Render("|"))
		}
	}
	sb.WriteString(strings.Repeat(" ", barLen-filled))
	sb.WriteString(" ")
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%*s", textW, text)))

	return labelStyle.Render(label) + dimStyle.Render("[") + sb.String() + dimStyle.Render("]")
}

func NetCell(label string, rate float64, hist []float64, width int, col lipgloss.Color) string {
	sparkW := width - 15
	if sparkW < 4 {
		sparkW = 4
	}
	return labelStyle.Render(fmt.Sprintf("%-3s", label)) +
		valueStyle.Render(fmt.Sprintf("%-11s", HumanRate(rate))) + " " +
		lipgloss.NewStyle().Foreground(col).Render(Spark(hist, sparkW, 0))
}

func LoadBar(v, max float64, width int, style lipgloss.Style, colored bool) string {
	n := int(v/max*float64(width) + 0.5)
	if v > 0 && n == 0 {
		n = 1
	}
	if n > width {
		n = width
	}
	on := strings.Repeat("█", n)
	off := strings.Repeat("░", width-n)
	if !colored {
		return on + off
	}
	return style.Render(on) + dimStyle.Render(off)
}

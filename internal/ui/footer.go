package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Fkeys(w int, keys [][2]string) string {
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fKeyStyle.Render(k[0]) + fLabelStyle.Render(k[1]))
	}
	if pad := w - lipgloss.Width(sb.String()); pad > 0 {
		sb.WriteString(fLabelStyle.Render(strings.Repeat(" ", pad)))
	}
	return sb.String()
}

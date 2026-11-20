package ui

import "fmt"

func formatBytesPerSecond(bytesPerSec float64) string {
	units := []string{"B/s", "KB/s", "MB/s", "GB/s"}

	i := 0
	for bytesPerSec >= 1024 && i < len(units)-1 {
		bytesPerSec /= 1024
		i++
	}

	return fmt.Sprintf("%.2f %s", bytesPerSec, units[i])
}

func HumanRate(b float64) string {
	units := []string{"B/s", "KB/s", "MB/s", "GB/s"}
	i := 0
	for b >= 1024 && i < len(units)-1 {
		b /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", b, units[i])
	}
	return fmt.Sprintf("%.1f %s", b, units[i])
}

func Truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

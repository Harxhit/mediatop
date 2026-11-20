package ui

import (
	"sort"

	"github.com/Harxhit/mediatop/internal/system"
)

type sortKey int

const (
	sortName sortKey = iota
	sortRX
	sortTX
	sortKeyCount
)

func (s sortKey) String() string {
	switch s {
	case sortRX:
		return "RX"
	case sortTX:
		return "TX"
	default:
		return "Name"
	}
}

func SelectedName(rates []system.NetworkRate, sel int) string {
	if sel >= 0 && sel < len(rates) {
		return rates[sel].Name
	}
	return ""
}

func IndexOf(rates []system.NetworkRate, name string, fallback int) int {
	for i, r := range rates {
		if r.Name == name {
			return i
		}
	}
	if fallback >= len(rates) {
		fallback = len(rates) - 1
	}
	if fallback < 0 {
		fallback = 0
	}
	return fallback
}

func SortRates(rates []system.NetworkRate, by sortKey) {
	sort.SliceStable(rates, func(i, j int) bool {
		a, b := rates[i], rates[j]
		switch by {
		case sortRX:
			if a.RXBytesPerSec != b.RXBytesPerSec {
				return float64(a.RXBytesPerSec) > float64(b.RXBytesPerSec)
			}
		case sortTX:
			if a.TXBytesPerSec != b.TXBytesPerSec {
				return float64(a.TXBytesPerSec) > float64(b.TXBytesPerSec)
			}
		}
		return a.Name < b.Name
	})
}

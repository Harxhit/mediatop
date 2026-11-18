package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MemoryStats struct {
	Total     uint64
	Free      uint64
	Available uint64
	Cached    uint64
}

func ReadMemoryStats() (MemoryStats, error) {

	var stats MemoryStats
	data, err := os.ReadFile("/proc/meminfo")

	if err != nil {
		return stats, fmt.Errorf("reading /proc/meminfo: %w", err)
	}

	parts := strings.Split(string(data), "\n")
	if len(parts) == 0 {
		return stats, fmt.Errorf("empty /proc/stat")
	}

	for _, part := range parts {

		parts := strings.Fields(part)

		if len(parts) < 2 {
			continue
		}

		value, err := strconv.ParseUint(parts[1], 10, 64)

		if err != nil {
			return stats, fmt.Errorf("MemTotal /proc/meminfo: %w", err)
		}
		switch parts[0] {
		case "MemTotal:":
			stats.Total = value
		case "MemFree:":
			stats.Free = value
		case "MemAvailable:":
			stats.Available = value
		case "Cached:":
			stats.Cached = value
		}
	}
	return stats, nil
}

func CalculateMemoryUsage(stats MemoryStats) float64 {
	if stats.Total == 0 || stats.Available > stats.Available {
		return 0
	}

	usage := (float64(stats.Total-stats.Available) / float64(stats.Total)) * 100
	return usage
}

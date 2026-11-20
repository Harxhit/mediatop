package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type InterfaceStats struct {
	Name      string
	RXBytes   uint64
	RXPackets uint64
	RXErrors  uint64
	RXDrops   uint64
	TXBytes   uint64
	TXPackets uint64
	TXErrors  uint64
	TXDrops   uint64
}

type NetworkRate struct {
	Name          string
	RXBytesPerSec float64
	TXBytesPerSec float64
}

func ReadNetworkStats() ([]InterfaceStats, error) {
	var netStats []InterfaceStats

	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, fmt.Errorf("read network stats: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines[2:] {
		iface, rawStats, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		iface = strings.TrimSpace(iface)
		fields := strings.Fields(rawStats)

		if len(fields) < 12 {
			return nil, fmt.Errorf(
				"invalid network stats for interface %q: expected at least 12 counters, got %d",
				iface, len(fields),
			)
		}

		var parsedStats []uint64

		for i, field := range fields {
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"parse counter %d for interface %q: %w",
					i, iface, err,
				)
			}

			parsedStats = append(parsedStats, value)
		}

		ifaceStats := InterfaceStats{
			Name:      iface,
			RXBytes:   parsedStats[0],
			RXPackets: parsedStats[1],
			RXErrors:  parsedStats[2],
			RXDrops:   parsedStats[3],
			TXBytes:   parsedStats[8],
			TXPackets: parsedStats[9],
			TXErrors:  parsedStats[10],
			TXDrops:   parsedStats[11],
		}

		netStats = append(netStats, ifaceStats)
	}

	return netStats, nil
}

func CalculateNetworkRate(previous, current InterfaceStats, elapsed time.Duration) NetworkRate {
	if elapsed <= 0 ||
		current.RXBytes < previous.RXBytes ||
		current.TXBytes < previous.TXBytes ||
		current.Name != previous.Name {
		return NetworkRate{Name: current.Name}
	}

	return NetworkRate{
		Name:          current.Name,
		RXBytesPerSec: float64(current.RXBytes-previous.RXBytes) / elapsed.Seconds(),
		TXBytesPerSec: float64(current.TXBytes-previous.TXBytes) / elapsed.Seconds(),
	}
}

func ReadNetworkRates() ([]NetworkRate, error) {
	var rates []NetworkRate

	previousStats, err := ReadNetworkStats()
	if err != nil {
		return nil, fmt.Errorf("read previous network snapshot: %w", err)
	}

	previousTime := time.Now()

	time.Sleep(time.Second)

	currentStats, err := ReadNetworkStats()
	if err != nil {
		return nil, fmt.Errorf("read current network snapshot: %w", err)
	}

	elapsed := time.Since(previousTime)

	for _, current := range currentStats {
		for _, previous := range previousStats {
			if current.Name == previous.Name {
				rate := CalculateNetworkRate(previous, current, elapsed)
				rates = append(rates, rate)
				break
			}
		}
	}

	return rates, nil
}

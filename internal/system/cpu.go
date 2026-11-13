package system

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUStats struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	IOWait    uint64
	IRQ       uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

func parseCPUStats() (CPUStats, error) {
	var stats CPUStats

	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return stats, fmt.Errorf("reading /proc/stat: %w", err)
	}

	parts := strings.SplitN(string(data), "\n", 2)
	if len(parts) == 0 {
		return stats, fmt.Errorf("empty /proc/stat")
	}

	fields := strings.Fields(parts[0])

	if len(fields) < 11 {
		return stats, fmt.Errorf("invalid CPU statistics: expected 10 numeric fields, got %d", len(fields)-1)
	}

	if fields[0] != "cpu" {
		return stats, fmt.Errorf("expected aggregate CPU line, got %q", fields[0])
	}

	values := make([]uint64, 10)

	for i := 0; i < 10; i++ {
		value, err := strconv.ParseUint(fields[i+1], 10, 64)
		if err != nil {
			return stats, fmt.Errorf("parsing CPU field %d: %w", i, err)
		}
		values[i] = value
	}

	stats = CPUStats{
		User:      values[0],
		Nice:      values[1],
		System:    values[2],
		Idle:      values[3],
		IOWait:    values[4],
		IRQ:       values[5],
		SoftIRQ:   values[6],
		Steal:     values[7],
		Guest:     values[8],
		GuestNice: values[9],
	}
	return stats, nil
}

func ReadCPUStats() (CPUStats, error) {
	var stats CPUStats
	firstSnapShot, err := parseCPUStats()
	if err != nil {
		return stats, fmt.Errorf("failed to read first CPU snapshot: %w", err)
	}
	time.Sleep(1 * time.Second)
	secondSnapShot, err := parseCPUStats()
	if err != nil {
		return stats, fmt.Errorf("failed to read second CPU snapshot: %w", err)
	}

	stats = CPUStats{
		User:      secondSnapShot.User - firstSnapShot.User,
		Nice:      secondSnapShot.Nice - firstSnapShot.Nice,
		System:    secondSnapShot.System - firstSnapShot.System,
		Idle:      secondSnapShot.Idle - firstSnapShot.Idle,
		IOWait:    secondSnapShot.IOWait - firstSnapShot.IOWait,
		IRQ:       secondSnapShot.IRQ - firstSnapShot.IRQ,
		SoftIRQ:   secondSnapShot.SoftIRQ - firstSnapShot.SoftIRQ,
		Steal:     secondSnapShot.Steal - firstSnapShot.Steal,
		Guest:     secondSnapShot.Guest - firstSnapShot.Guest,
		GuestNice: secondSnapShot.GuestNice - firstSnapShot.GuestNice,
	}

	return stats, nil
}

func CalculateTotalDelta(stats CPUStats) uint64 {
	delta := stats.User +
		stats.Nice +
		stats.System +
		stats.Idle +
		stats.IOWait +
		stats.IRQ +
		stats.SoftIRQ +
		stats.Steal

	return delta
}

func CalculateCPUutilization(idleDelta uint64, delta uint64) float64 {
	usage := float64(delta-idleDelta) / float64(delta) * 100
	return usage
}

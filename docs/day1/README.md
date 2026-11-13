# mediatop

A terminal-based observability tool for monitoring WebRTC calls and inspecting real-time media statistics.

Built in Go, inspired by `htop`.

## Day 1 — Project Foundation & CPU Monitoring

**Date:** November 13, 2026

### Today's Objectives

- Set up the Go project.
- Understand Linux CPU statistics.
- Read and parse CPU metrics from `/proc/stat`.
- Calculate CPU utilization.
- Build the initial terminal UI using Bubble Tea.

### 1. Project Structure

Initial project structure:

```text
mediatop/
├── cmd/
│   └── mediatop/
│       └── main.go
├── internal/
│   ├── system/
│   │   ├── cpu.go
│   │   └── memory.go
│   ├── network/
│   ├── media/
│   └── ui/
│       └── model.go
├── go.mod
└── README.md
```

### 2. Linux CPU Monitoring

Implemented CPU statistics collection using:

```text
/proc/stat
```

Learned about the aggregate CPU counters:

- User
- Nice
- System
- Idle
- IOWait
- IRQ
- SoftIRQ
- Steal
- Guest
- GuestNice

Implemented parsing using:

- `os.ReadFile()`
- `strings.Fields()`
- `strconv.ParseUint()`

Created a `CPUStats` struct to represent CPU counters.

### 3. CPU Utilization Calculation

Implemented a two-snapshot approach:

1. Read the first CPU snapshot.
2. Wait one second.
3. Read the second CPU snapshot.
4. Calculate differences between counters.
5. Calculate total CPU time.
6. Calculate CPU utilization.

CPU utilization formula:

```text
CPU Usage = (Total Delta - Idle Delta) / Total Delta × 100
```

Implemented:

```go
ReadCPUStats()
CalculateTotalDelta()
CalculateCPUutilization()
```

Important Linux detail: Guest and GuestNice counters are already included in User and Nice respectively, so they must not be counted again in total CPU time.

### 4. Continuous CPU Monitoring

Implemented continuous CPU monitoring with a Go `for` loop.

Successfully tested CPU utilization updates at approximately one-second intervals.

Example:

```text
CPU Usage: 3.58%
CPU Usage: 5.32%
CPU Usage: 5.19%
CPU Usage: 3.35%
CPU Usage: 5.06%
```

### 5. Terminal UI Foundation

Installed and started working with:

- **Bubble Tea** — terminal application framework.
- **Lip Gloss** — terminal styling library.

Learned the Bubble Tea architecture:

| Component | Responsibility |
|---|---|
| Model | Stores application state |
| Init | Initializes the application |
| Update | Handles messages and updates state |
| View
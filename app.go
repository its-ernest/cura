package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/its-ernest/cura/pkg/memory"

	"github.com/BurntSushi/toml"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemStats defines the data structure sent to the React frontend
type SystemStats struct {
	CPUUsage     float64 `json:"cpuUsage"`
	RAMUsage     float64 `json:"ramUsage"`
	TotalRAM     uint64  `json:"totalRam"`
	ProcessCount int     `json:"processCount"`
}

// fixed:
// EnforcementConfig separate instead of anonymous nested structs
type EnforcementConfig struct {
	IsEnforced bool    `toml:"is_enforced" json:"is_enforced"`
	MemoryCap  float64 `toml:"memory_cap" json:"memory_cap"`
	CPUCeiling float64 `toml:"cpu_ceiling" json:"cpu_ceiling"`
}

// Config to match settings.toml structure
type Config struct {
	Version     string            `toml:"version" json:"version"`
	Enforcement EnforcementConfig `toml:"enforcement" json:"enforcement"`
}

type App struct {
	ctx           context.Context
	memoryManager *memory.Manager
	config        Config
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{
		// default to 80% usage cap (20% reserve)
		memoryManager: memory.NewManager(80.0),
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// automatically load settings on startup
	cfg, err := a.LoadSettings()
	if err == nil && cfg.Enforcement.IsEnforced {
		a.memoryManager.StartEnforcer(ctx)
		fmt.Println("Auto-enforcer started")
	}
}

// SetMemoryCap bridges the React slider to the Go memory manager
func (a *App) SetMemoryCap(percent float64) {
	a.memoryManager.SetCap(percent)
}

// GetLiveStats provides the pulse for the React Dashboard rings
func (a *App) GetLiveStats() (SystemStats, error) {
	// get CPU aggregate usage
	c, _ := cpu.Percent(0, false)
	cpuVal := 0.0
	if len(c) > 0 {
		cpuVal = c[0]
	}

	// get Memory Stats
	v, _ := mem.VirtualMemory()

	// logic: calculate actual usage (total - available)
	actualUsed := float64(v.Total - v.Available)
	actualUsagePercent := (actualUsed / float64(v.Total)) * 100

	// round total RAM to nearest GiB (prevents the "7GB" display on 8GB sticks)
	totalGiB := math.Round(float64(v.Total) / 1073741824.0)

	// get current process count
	pids, _ := process.Pids()

	return SystemStats{
		CPUUsage:     cpuVal,
		RAMUsage:     actualUsagePercent,
		TotalRAM:     uint64(totalGiB),
		ProcessCount: len(pids),
	}, nil
}

// StartEnforcement is called from React when the toggle is turned ON
func (a *App) StartEnforcement() {
	a.memoryManager.StartEnforcer(a.ctx)
}

// StopEnforcement is called from React when the toggle is turned OFF
func (a *App) StopEnforcement() {
	a.memoryManager.StopEnforcer()
}

// LoadSettings reads the TOML file and returns it to React
func (a *App) LoadSettings() (Config, error) {
	_, err := toml.DecodeFile("settings.toml", &a.config)
	if err != nil {
		return Config{}, err
	}
	// sync the memory manager with loaded settings
	a.memoryManager.SetCap(a.config.Enforcement.MemoryCap)
	return a.config, nil
}

// SaveSettings writes the current config state to the TOML file
func (a *App) SaveSettings(cfg Config) error {
	a.config = cfg
	f, err := os.Create("settings.toml")
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}
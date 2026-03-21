package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/its-ernest/cura/pkg/logging"
	"github.com/its-ernest/cura/pkg/memory"
	"github.com/its-ernest/cura/pkg/whitelist"

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
	Version     string                      `toml:"version" json:"version"`
	Enforcement EnforcementConfig           `toml:"enforcement" json:"enforcement"`
	Apps        map[string]memory.AppStatus `toml:"apps" json:"apps"`
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

var l *logging.Logger = logging.NewLogger("cura.log")

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// automatically load settings on startup
	cfg, err := a.LoadSettings()

	// get installed apps list
	a.RefreshApps()

	// start enforcer if previously enabled
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
	// get cpu aggregate usage
	c, _ := cpu.Percent(0, false)
	cpuVal := 0.0
	if len(c) > 0 {
		cpuVal = c[0]
	}

	// get memory stats
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
	if a.config.Apps != nil {
		a.memoryManager.AppMap = a.config.Apps
	}
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

func (a *App) GetLogs(limit int) (string, error) {

	content, err := l.ReadLogs(limit)
	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return "Waiting for system events...", nil
	}

	start := 0
	if len(lines) > limit {
		start = len(lines) - limit
	}

	finalLogs := strings.Join(lines[start:], "\n")
	fmt.Printf("DEBUG: Sent %d lines to UI.\n", len(lines[start:]))

	return finalLogs, nil
}

// Whitelist exposed methods
// GetAppMap returns the current memory map to the UI
func (a *App) GetAppMap() map[string]memory.AppStatus {
	return a.memoryManager.AppMap
}

func (a *App) RefreshApps() {
	rawApps, err := whitelist.GetWindowsApps()
	if err != nil {
		l.Write(fmt.Sprintf("ERROR: Registry scan failed: %v", err))
		return
	}

	newCount := 0
	for _, app := range rawApps {
		// Crucial: Only add if it doesn't exist to avoid overwriting user 'IsExempt' settings
		if _, exists := a.memoryManager.AppMap[app.Name]; !exists {
			a.memoryManager.AppMap[app.Name] = memory.AppStatus{
				Directory: app.ExePath,
				IsExempt:  false,
			}
			newCount++
		}
	}

	// 3. Sync the local config and save if new apps were found
	if newCount > 0 {
		a.config.Apps = a.memoryManager.AppMap
		a.SaveSettings(a.config)
		l.Write(fmt.Sprintf("SYSTEM: Found %d new apps. Total registered: %d", newCount, len(a.memoryManager.AppMap)))
	}
}

// ToggleExemption flips the is_exempt status and persists
func (a *App) ToggleExemption(name string) {
	if app, ok := a.memoryManager.AppMap[name]; ok {
		// toggle the state in the manager
		app.IsExempt = !app.IsExempt
		a.memoryManager.AppMap[name] = app

		// persist the change to the toml file via the config
		a.config.Apps = a.memoryManager.AppMap
		a.SaveSettings(a.config)

		status := "MONITORED"
		if app.IsExempt {
			status = "EXEMPT"
		}
		l.Write(fmt.Sprintf("CONFIG: %s is now %s", name, status))
	}
}

// RemoveApp deletes an entry from the map
func (a *App) RemoveApp(name string) {
	delete(a.memoryManager.AppMap, name)
	l.Write(fmt.Sprintf("CONFIG: Removed %s from Registry", name))
}

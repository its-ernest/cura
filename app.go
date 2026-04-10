package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/its-ernest/cura/pkg/logging"
	"github.com/its-ernest/cura/pkg/memory"
	"github.com/its-ernest/cura/pkg/routine"
	"github.com/its-ernest/cura/pkg/updater"
	"github.com/wailsapp/wails/v2/pkg/runtime"

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
	AutoUpdate bool    `toml:"auto_update" json:"auto_update"`
}

// Config to match settings.toml structure
type Config struct {
	Version     string                      `toml:"version" json:"version"`
	Enforcement EnforcementConfig           `toml:"enforcement" json:"enforcement"`
	Apps        map[string]memory.AppStatus `toml:"apps" json:"apps"`
}

type App struct {
	ctx            context.Context
	path           string
	memoryManager  *memory.Manager
	routineManager *routine.Manager
	config         Config
	configMu       sync.Mutex
}

// NewApp creates a new App instance
func NewApp() *App {
	execPath, _ := os.Executable()
	dir := filepath.Dir(execPath)
	mm := memory.NewManager(80.0)
	return &App{
		path: filepath.Join(dir, "settings.toml"),
		// default to 80% usage cap (20% reserve)
		memoryManager:  mm,
		routineManager: routine.NewManager(mm),
	}
}

var l *logging.Logger = logging.NewLogger("cura.log")

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	routines, err := a.GetRoutines()
	if err == nil {
		a.routineManager.Routines = routines
	} else {
		l.Write(fmt.Sprintf("SYSTEM: No routines found or folder missing: %v", err))
	}

	go a.routineManager.Run()

	// automatically load settings on startup
	cfg, err := a.LoadSettings()

	// get installed apps list
	a.RefreshApps()

	// update check
	go func() {
		time.Sleep(2 * time.Second)
		release, hasUpdate, err := updater.CheckForUpdates(a.config.Version)
		if err == nil && hasUpdate {
			if a.config.Enforcement.AutoUpdate {
				l.Write("UPDATE: Auto-update triggered for " + release.TagName)
			} else {
				runtime.EventsEmit(a.ctx, "update_available", release)
			}
		}

		// cleanup .old update files
		a.cleanupLegacyFiles()
	}()

	// start enforcer if previously enabled
	if err == nil && cfg.Enforcement.IsEnforced {
		go func() {
			time.Sleep(500 * time.Millisecond)
			a.memoryManager.StartEnforcer(a.ctx)
			l.Write("SYSTEM: Auto-enforcer resumed from settings.")
		}()
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
	a.configMu.Lock()
	defer a.configMu.Unlock()
	_, err := toml.DecodeFile(a.path, &a.config)
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
	a.configMu.Lock()
	defer a.configMu.Unlock()
	a.config = cfg
	f, err := os.Create(a.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// Updates logic
func (a *App) cleanupLegacyFiles() {
	execPath, _ := os.Executable()
	oldFile := execPath + ".old"

	if _, err := os.Stat(oldFile); err == nil {
		err := os.Remove(oldFile)
		if err == nil {
			l.Write("SYSTEM: Cleaned up legacy update files.")
		}
	}
}

func (a *App) TriggerUpdate(release updater.GitHubRelease) string {
	err := updater.DownloadAndInstall(&release)
	if err != nil {
		l.Write(fmt.Sprintf("UPDATE ERROR: %v", err))
		return "Update failed: " + err.Error()
	}

	a.configMu.Lock()
	a.config.Version = release.TagName
	a.configMu.Unlock()

	err = a.SaveSettings(a.config)
	if err != nil {
		l.Write("ERROR: Could not sync new version to TOML: " + err.Error())
	}

	l.Write("SYSTEM: Update installed. Binary swapped. Restarting in 1s...")

	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	return "Success"
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
	//fmt.Printf("DEBUG: Sent %d lines to UI.\n", len(lines[start:]))

	return finalLogs, nil
}

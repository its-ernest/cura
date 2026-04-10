package main

import (
	"fmt"
	"path/filepath"

	"github.com/its-ernest/cura/pkg/whitelist"
	"github.com/its-ernest/cura/pkg/memory"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Whitelist exposed methods
// GetAppMap returns the current memory map to the UI
func (a *App) GetAppMap() map[string]memory.AppStatus {
	return a.memoryManager.AppMap
}

// ToggleExemption flips the is_exempt status and persists
func (a *App) ToggleExemption(name string) {
	a.configMu.Lock()

	if app, ok := a.memoryManager.AppMap[name]; ok {
		// toggle the state in the manager
		app.IsExempt = !app.IsExempt
		a.memoryManager.AppMap[name] = app

		// persist the change to the toml file via the config
		a.config.Apps = a.memoryManager.AppMap
		a.configMu.Unlock()
		a.SaveSettings(a.config)

		status := "MONITORED"
		if app.IsExempt {
			status = "EXEMPT"
		}
		l.Write(fmt.Sprintf("CONFIG: %s is now %s", name, status))
	} else {
		a.configMu.Unlock()
	}
}

// RemoveApp deletes an entry from the map
func (a *App) RemoveApp(name string) {
	a.configMu.Lock()
	delete(a.memoryManager.AppMap, name)
	a.config.Apps = a.memoryManager.AppMap
	a.configMu.Unlock()

	a.SaveSettings(a.config)
	l.Write(fmt.Sprintf("CONFIG: Removed %s", name))
}

// SelectProcess opens a native dialog and adds the result to the whitelist
func (a *App) SelectProcess() string {
	//define options
	options := runtime.OpenDialogOptions{
		Title: "Select Process or Folder to Exempt",
		Filters: []runtime.FileFilter{
			{DisplayName: "Executables (*.exe)", Pattern: "*.exe"},
		},
	}

	// open file dialog
	path, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil || path == "" {
		return ""
	}

	// extract the name (e.g., "chrome.exe")
	name := filepath.Base(path)

	a.configMu.Lock()
	// create the AppStatus struct for the process, then assign it
	a.memoryManager.AppMap[name] = memory.AppStatus{
		Directory: path,
		IsExempt:  true,
	}
	a.config.Apps = a.memoryManager.AppMap // Sync back to config
	a.configMu.Unlock()

	a.SaveSettings(a.config)
	return name
}

func (a *App) RefreshApps() {
	a.configMu.Lock()
	rawApps, err := whitelist.GetWindowsApps()
	if err != nil {
		l.Write(fmt.Sprintf("ERROR: Registry scan failed: %v", err))
		return
	}

	newCount := 0
	for _, app := range rawApps {
		// crucial: only add if it doesn't exist to avoid overwriting user 'IsExempt' settings
		if _, exists := a.memoryManager.AppMap[app.Name]; !exists {
			a.memoryManager.AppMap[app.Name] = memory.AppStatus{
				Directory: app.ExePath,
				IsExempt:  false,
			}
			newCount++
		}
	}

	// sync the local config and save if new apps were found
	a.configMu.Unlock()
	if newCount > 0 {
		a.config.Apps = a.memoryManager.AppMap
		a.SaveSettings(a.config)
		l.Write(fmt.Sprintf("SYSTEM: Found %d new apps. Total registered: %d", newCount, len(a.memoryManager.AppMap)))
	}
}

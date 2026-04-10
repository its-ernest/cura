package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/its-ernest/cura/pkg/routine"
	"gopkg.in/yaml.v3"
)

// GetRoutines scans the routines folder and returns a list of available pipelines
func (a *App) GetRoutines() ([]*routine.Routine, error) {
	routineDir := "./routines"
	files, err := os.ReadDir(routineDir)
	if err != nil {
		return nil, err
	}

	var list []*routine.Routine
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".yaml" || ext == ".yml" {
			r, err := routine.LoadRoutine(filepath.Join(routineDir, file.Name()))
			if err == nil {
				r.Path = filepath.Join(routineDir, file.Name())
				list = append(list, r)
			}
		}
	}
	return list, nil
}

// ToggleRoutine allows the UI to enable/disable a specific routine file
func (a *App) ToggleRoutine(name string, enabled bool) {
	fmt.Printf("DEBUG: Go received ToggleRoutine(%s, %v)\n", name, enabled)

	for _, r := range a.routineManager.Routines {
		if r.Name == name {
			r.Enabled = enabled

			// 1. Persist to Disk
			if r.Path != "" {
				data, err := yaml.Marshal(r)
				if err != nil {
					l.Write(fmt.Sprintf("ERROR: Failed to marshal YAML for %s: %v", name, err))
				} else {
					err = os.WriteFile(r.Path, data, 0644)
					if err != nil {
						l.Write(fmt.Sprintf("ERROR: Failed to write YAML for %s: %v", name, err))
					}
				}
			}

			// 2. Cleanup if deactivating
			if !enabled && r.IsActive {
				a.routineManager.Deactivate(r)
			}
			break
		}
	}
	l.Write(fmt.Sprintf("ROUTINE: Toggle %s to %v (Persisted)", name, enabled))
}

// CreateRoutine saves a new routine to a YAML file and adds it to the manager
func (a *App) CreateRoutine(r routine.Routine) error {
	// 1. Sanitize name for filename
	filename := fmt.Sprintf("%s.yaml", strings.ToLower(strings.ReplaceAll(r.Name, " ", "_")))
	filepath := filepath.Join("./routines", filename)

	r.Enabled = true // Default to enabled
	r.Path = filepath

	// 2. Marshal to YAML
	data, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("failed to encode routine: %v", err)
	}

	// 3. Save to disk
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	// 4. Update memory state so the background loop sees it immediately
	a.routineManager.Routines = append(a.routineManager.Routines, &r)

	l.Write(fmt.Sprintf("ROUTINE: Created new pipeline '%s'", r.Name))
	return nil
}

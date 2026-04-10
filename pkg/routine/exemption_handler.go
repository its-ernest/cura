package routine

import (
	"github.com/its-ernest/cura/pkg/memory"
)

type ExemptionHandler struct {
	// Store the previous state: map[appName]wasExempt
	PreviousStates map[string]bool
}

func (h *ExemptionHandler) Apply(mm *memory.Manager, apps []interface{}) {
	if h.PreviousStates == nil {
		h.PreviousStates = make(map[string]bool)
	}

	for _, val := range apps {
		appName, ok := val.(string)
		if !ok {
			continue
		}

		// 1. Record the current state before touching it
		if app, exists := mm.AppMap[appName]; exists {
			h.PreviousStates[appName] = app.IsExempt

			// 2. Force it to Exempt for the routine
			app.IsExempt = true
			mm.AppMap[appName] = app
		} else {
			// 3. If it's not in the map, add it as a "Temporary" entry
			// mark it for deletion in Cleanup since it wasn't there originally
			mm.AppMap[appName] = memory.AppStatus{
				Directory: "ROUTINE_TEMP",
				IsExempt:  true,
			}
			h.PreviousStates[appName] = false // It didn't exist, so logically it wasn't exempt
		}
	}
}

func (h *ExemptionHandler) Cleanup(mm *memory.Manager) {
	for name, wasExempt := range h.PreviousStates {
		if app, exists := mm.AppMap[name]; exists {
			if app.Directory == "ROUTINE_TEMP" {
				// If we added it just for this routine, now we can delete it
				delete(mm.AppMap, name)
			} else {
				// If it was an existing app, just revert its status
				app.IsExempt = wasExempt
				mm.AppMap[name] = app
			}
		}
	}
	h.PreviousStates = nil // Clear the snapshot
}

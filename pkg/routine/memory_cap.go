package routine

import (
	"fmt"

	"github.com/its-ernest/cura/pkg/memory"
)

// MemoryCapHandler manages the transition between global and routine-specific caps.
type MemoryCapHandler struct {
	OriginalCap  float64
	IsOverridden bool
}

// ApplyCap takes the value from the YAML (percentage) and updates the enforcer.
func (h *MemoryCapHandler) ApplyCap(mm *memory.Manager, newValue interface{}) error {
	// 1. Convert the interface{} value from YAML to a float
	val, ok := newValue.(int)
	if !ok {
		// fallback check for float64 if the YAML parser treated it as float
		fVal, ok := newValue.(float64)
		if !ok {
			return fmt.Errorf("invalid memory_cap value type")
		}
		val = int(fVal)
	}

	// 2. Save the original cap if this is the first override in the session
	if !h.IsOverridden {
		h.OriginalCap = mm.CapPercentage
		h.IsOverridden = true
		l.Write(fmt.Sprintf("ROUTINE: Original cap saved: %.0f%%", h.OriginalCap))
	}

	// 3. Routine cap ALWAYS takes priority — lock out preset/slider changes
	mm.RoutineOverride = true
	mm.CapPercentage = float64(val)
	l.Write(fmt.Sprintf("ROUTINE: Cap overridden to %.0f%% (preset: %.0f%%)", float64(val), mm.PresetCap))

	return nil
}

// Restore restores the preset cap and unlocks routine override
func (h *MemoryCapHandler) Restore(mm *memory.Manager) {
	if h.IsOverridden {
		mm.RoutineOverride = false
		// use the latest preset (not stale OriginalCap) in case user changed it mid-routine
		mm.CapPercentage = mm.PresetCap
		h.IsOverridden = false
		l.Write(fmt.Sprintf("ROUTINE: Cap restored to preset %.0f%%", mm.PresetCap))
	}
}

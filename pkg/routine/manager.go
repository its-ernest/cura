package routine

import (
	"fmt"
	"time"

	"github.com/its-ernest/cura/pkg/logging"
	"github.com/its-ernest/cura/pkg/memory"
)

var l *logging.Logger = logging.NewLogger("cura.log")

type Manager struct {
	Routines      []*Routine
	Ticker        *time.Ticker
	MemoryManager *memory.Manager
	CapHandler    *MemoryCapHandler // State restorer
}

func NewManager(mm *memory.Manager) *Manager {
	return &Manager{
		Routines:      make([]*Routine, 0),
		Ticker:        time.NewTicker(2 * time.Second),
		MemoryManager: mm,
		CapHandler:    &MemoryCapHandler{},
	}
}

func (m *Manager) Run() {
	l.Write("ROUTINE: Observer engine started.")
	for range m.Ticker.C {
		for _, r := range m.Routines {
			if !r.Enabled {
				continue
			}

			isRunning := CheckProcess(r.Trigger.Target)

			if isRunning && !r.IsActive {
				l.Write(fmt.Sprintf("ROUTINE: Trigger met for '%s'. Activating...", r.Name))
				m.Activate(r)
			} else if !isRunning && r.IsActive {
				l.Write(fmt.Sprintf("ROUTINE: Stop condition met for '%s'. Reverting...", r.Name))
				m.Deactivate(r)
			}
		}
	}
}

func (m *Manager) Activate(r *Routine) {
	r.IsActive = true

	for _, action := range r.Actions {
		switch action.Type {
		case "set_memory_cap":
			// Call the ApplyCap logic from memory_cap.go
			err := m.CapHandler.ApplyCap(m.MemoryManager, action.Value)
			if err != nil {
				l.Write(fmt.Sprintf("ERROR: Routine '%s' failed: %v", r.Name, err))
			}
		case "boost_priority":
			l.Write(fmt.Sprintf("ROUTINE: Priority boost requested for %s", action.Target))
		}
	}
}

func (m *Manager) Deactivate(r *Routine) {
	r.IsActive = false
	// Restore cap config
	m.CapHandler.Restore(m.MemoryManager)
}

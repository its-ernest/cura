package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type Manager struct {
	CapPercentage float64
	IsActive      bool
	cancelFunc    context.CancelFunc
}

func NewManager(initialCap float64) *Manager {
	return &Manager{
		CapPercentage: initialCap,
		IsActive:      false,
	}
}

func (m *Manager) SetCap(percent float64) {
	m.CapPercentage = percent
}

// StartEnforcer uses adaptive frequency to handle rapid RAM spikes.
func (m *Manager) StartEnforcer(ctx context.Context) {
	if m.IsActive {
		return
	}

	m.IsActive = true
	enforceCtx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	go func() {
		l.Write("SYSTEM: Adaptive Enforcer routine initialized.")

		for {
			select {
			case <-enforceCtx.Done():
				m.IsActive = false
				l.Write("SYSTEM: Enforcer routine terminated.")
				return
			default:
				v, _ := mem.VirtualMemory()
				actualUsage := (float64(v.Total-v.Available) / float64(v.Total)) * 100

				if actualUsage > m.CapPercentage {
					l.Write(fmt.Sprintf("ALERT: Pressure %.1f%% exceeds Cap %.1f%%. Enforcing...", actualUsage, m.CapPercentage))
					m.enforce()
					
					// BURST MODE: Re-check quickly to catch rapid spikes
					time.Sleep(500 * time.Millisecond)
				} else {
					// IDLE MODE: System healthy
					time.Sleep(3 * time.Second)
				}
			}
		}
	}()
}

func (m *Manager) StopEnforcer() {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}
}
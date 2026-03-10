package memory

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
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

// StartEnforcer launches the background routine that stays active until the RAM is under the cap
func (m *Manager) StartEnforcer(ctx context.Context) {
	if m.IsActive {
		return
	}

	m.IsActive = true
	enforceCtx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	go func() {
		ticker := time.NewTicker(3 * time.Second) // check frequency
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				v, _ := mem.VirtualMemory()
				actualUsage := (float64(v.Total-v.Available) / float64(v.Total)) * 100

				if actualUsage > m.CapPercentage {
					fmt.Printf("[CURA] Usage %.1f%% exceeds Cap %.1f%%. Enforcing...\n", actualUsage, m.CapPercentage)
					m.enforce()
				}
			case <-enforceCtx.Done():
				m.IsActive = false
				return
			}
		}
	}()
}

func (m *Manager) StopEnforcer() {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}
}

// enforce finds the best candidates to close to free up RAM
func (m *Manager) enforce() {
	processes, err := process.Processes()
	if err != nil {
		return
	}

	// critical Windows processes that should be NEVER touched
	protected := map[string]bool{
		"explorer.exe": true,
		"System":       true,
		"svchost.exe":  true,
		"lsass.exe":    true,
		"wininit.exe":  true,
		"csrss.exe":    true,
		"services.exe": true,
		"Wails.exe":    true, "cura.exe": true, "cura-dev.exe": true, // don't kill self!
	}

	type candidate struct {
		proc *process.Process
		mem  uint64
		name string
	}

	var candidates []candidate

	for _, p := range processes {
		name, _ := p.Name()
		if protected[name] || p.Pid <= 4 {
			continue
		}

		// Windows specific: Idle detection
		// if cpu is near zero, it a prime candidate for background cleanup
		cpu, _ := p.CPUPercent()
		if cpu > 0.5 { 
			continue
		}

		memInfo, _ := p.MemoryInfo()
		if memInfo == nil || memInfo.RSS == 0 {
			continue
		}

		candidates = append(candidates, candidate{
			proc: p,
			mem:  memInfo.RSS,
			name: name,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].mem > candidates[j].mem
	})

	if len(candidates) > 0 {
		target := candidates[0]
		// equivalent of clicking "End Task" in Task Manager.
		fmt.Printf("[CURA] Enforcing: Closing %s (%d MB)\n", target.name, target.mem/1024/1024)
		target.proc.Terminate()
	}
}
package memory

import (
	"context"
	"fmt"
	"os"
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

// logToFile appends a timestamped message to cura.log
func (m *Manager) logToFile(message string) {
	f, err := os.OpenFile("cura.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
}

func (m *Manager) SetCap(percent float64) {
	m.CapPercentage = percent
}

func (m *Manager) StartEnforcer(ctx context.Context) {
	if m.IsActive {
		return
	}

	m.IsActive = true
	enforceCtx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				v, _ := mem.VirtualMemory()
				actualUsage := (float64(v.Total-v.Available) / float64(v.Total)) * 100

				if actualUsage > m.CapPercentage {
					m.logToFile(fmt.Sprintf("ALERT: Usage %.1f%% exceeds Cap %.1f%%. Enforcing...", actualUsage, m.CapPercentage))
					m.enforce()
				}
			case <-enforceCtx.Done():
				m.IsActive = false
				m.logToFile("SYSTEM: Enforcer stopped.")
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

func (m *Manager) enforce() {
	processes, err := process.Processes()
	if err != nil {
		return
	}

	// Critical Windows processes + Development tools
    protected := map[string]bool{
        "explorer.exe": true, "System": true, "svchost.exe": true,
        "lsass.exe": true, "wininit.exe": true, "csrss.exe": true,
        "services.exe": true, "ShellHost.exe": true, "sihost.exe": true,
        "ShellExperienceHost.exe": true, "dllhost.exe": true, "ctfmon.exe": true,
		
		// dev tools 
		"WindowsTerminal.exe": true, "OpenConsole.exe": true, "powershell.exe": true,

		// critical for app to run
		"msedgewebview2.exe": true, "cura.exe": true, "cura-dev.exe": true, 

        // temp, needed for testing
        "Code.exe": true, "node.exe": true, "taskhostw.exe": true, "wails.exe": true, 
        "MSPCManagerService.exe": true, "esbuild.exe": true,
    }

	type candidate struct {
		proc *process.Process
		mem  uint64
		name string
		cpu  float64
	}

	var candidates []candidate

	for _, p := range processes {
		name, _ := p.Name()
		if protected[name] || p.Pid <= 4 {
			continue
		}

		memInfo, _ := p.MemoryInfo()
		if memInfo == nil || memInfo.RSS == 0 {
			continue
		}

		cpuVal, _ := p.CPUPercent()

		candidates = append(candidates, candidate{
			proc: p,
			mem:  memInfo.RSS,
			name: name,
			cpu:  cpuVal,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].mem > candidates[j].mem
	})

	for _, target := range candidates {
		v, _ := mem.VirtualMemory()
		currentUsage := (float64(v.Total-v.Available) / float64(v.Total)) * 100

		if currentUsage <= m.CapPercentage {
			break
		}

		if target.cpu < 2.0 || currentUsage > (m.CapPercentage + 5.0) {
			m.logToFile(fmt.Sprintf("ACTION: Terminating %s (%d MB) | System Pressure: %.2f%%",
				target.name, target.mem/1024/1024, currentUsage))

			target.proc.Kill()
			time.Sleep(100 * time.Millisecond)
		}
	}
}
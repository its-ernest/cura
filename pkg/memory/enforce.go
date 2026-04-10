package memory

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// global backoff map to track processes that auto-restart
var killBackoff = make(map[string]time.Time)

func (m *Manager) enforce() {
	processes, err := process.Processes()
	if err != nil {
		return
	}

	protected := GetProtectedProcesses()
	type candidate struct {
		proc *process.Process
		mem  uint64
		name string
		cpu  float64
	}

	var candidates []candidate
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			// name is often unavailable for kernel tasks even as admin
			//l.Write(fmt.Sprintf("SILENT ERROR: Couldn't obtain process info for PID %d", p.Pid))
			continue
		}

		// filter: protected or system
		if protected[name] || p.Pid <= 4 {
			//l.Write(fmt.Sprintf("ACTION: Skipping %s", name))
			continue
		}

		// dynamic filter: protects hardware drivers and system32 sub-binaries
		if IsSystemOrDriver(p) {
			//l.Write(fmt.Sprintf("SKIP: Skipping system %s", name))
			continue
		}

		// filter: backoff (stop killing the same immortal process repeatedly)
		if lastAttempt, exists := killBackoff[name]; exists && time.Since(lastAttempt) < 40*time.Second {
			continue
		}

		memInfo, err := p.MemoryInfo()
		if err != nil {
			l.Write(fmt.Sprintf("DEBUG: Access denied for %s (PID %d)", name, p.Pid))
			continue
		}
		if memInfo == nil || memInfo.RSS == 0 {
			continue
		}

		// filter stale/foreground
		if !m.IsStale(p, 10*time.Minute) {
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

	// sort: highest RAM usage first
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].mem > candidates[j].mem
	})

	for _, target := range candidates {
		v, _ := mem.VirtualMemory()
		currentUsage := (float64(v.Total-v.Available) / float64(v.Total)) * 100

		if currentUsage <= m.CapPercentage {
			break
		}

		// idle hogs vs critical pressure
		if target.cpu < 2.0 || currentUsage > (m.CapPercentage+5.0) {
			l.Write(fmt.Sprintf("ACTION: Terminating %s (%d MB) | System Pressure: %.2f%%",
				target.name, bytesToMB(target.mem), currentUsage))

			target.proc.Kill()

			// mark as killed to trigger backoff timer
			killBackoff[target.name] = time.Now()

			// small delay to allow system to update memory tables
			time.Sleep(100 * time.Millisecond)
		}
	}
}

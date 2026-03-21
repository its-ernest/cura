package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/go-foreground"
	"github.com/shirou/gopsutil/v3/process"
)

// IsStale checks if a process is old enough and not currently in use
func (m *Manager) IsStale(p *process.Process, threshold time.Duration) bool {
	pid := p.Pid
	name, err := p.Name()
	if err != nil {
		l.Write(fmt.Sprintf("SILENT ERROR: Couldn't obtain process info name for PID %d", p.Pid))
		return false
	}
	exePath, err := p.Exe()
	if err != nil {
		l.Write(fmt.Sprintf("SILENT ERROR: Couldn't obtain process info exePath for PID %d", p.Pid))
		return false
	}

	// check if it is the foreground window
	fgPID, _ := foreground.GetForegroundPID()
	if int32(fgPID) == pid {
		// if so, update the last-seen timestamp right now
		m.LastForegroundMap[pid] = time.Now()
		return false // never stale if the user is looking at it
	}

	// if not foreground, check grace period
	if lastSeen, exists := m.LastForegroundMap[pid]; exists {
		gracePeriod := 2 * time.Minute // tolerate 2 mins of being minimized/backgrounded
		if time.Since(lastSeen) < gracePeriod {
			return false // too soon to kill
		}
	}

	// if not too soon, verify if it breaches threshold
	createTime, err := p.CreateTime() // returns ms since epoch
	if err != nil {
		return false
	}

	created := time.Unix(0, createTime*int64(time.Millisecond))

	// if it was created less than the-threshold ago, it's not stale
	if time.Since(created) < threshold {
		return false
	}

	// slight performance fix: moved whitelist looping into IsStale
	// ensures certain system protected processes doesn't have to be looped
	// whitelist loop checker
	for appName, status := range m.AppMap {
		if status.IsExempt {
			// if the names match or the process lives in an exempted directory
			if strings.EqualFold(name, appName) || strings.HasPrefix(exePath, status.Directory) {
				return false // shield this process
			}
		}
	}

	return true
}

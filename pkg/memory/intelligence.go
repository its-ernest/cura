package memory

import (
	"time"
	"github.com/LightningDev1/go-foreground"
	"github.com/shirou/gopsutil/v3/process"
)

// IsStale checks if a process is old enough and not currently in use
func IsStale(p *process.Process, threshold time.Duration) bool {
	createTime, err := p.CreateTime() // returns ms since epoch
	if err != nil {
		return false
	}

	created := time.Unix(0, createTime*int64(time.Millisecond))
	
	// if it was created less than 'threshold' ago, it's not stale
	if time.Since(created) < threshold {
		return false
	}

	// check if it is the foreground window
	fgPID, _ := foreground.GetForegroundPID()
	if int32(fgPID) == p.Pid {
		return false // never stale if the user is looking at it
	}

	return true
}
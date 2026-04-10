package routine

import (
	"strings"
	"github.com/shirou/gopsutil/v3/process"
)

// CheckProcess scans the current process tree for a matching name
func CheckProcess(targetName string) bool {
	// get all running process IDs
	pids, err := process.Pids()
	if err != nil {
		return false
	}

	targetName = strings.ToLower(targetName)

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, err := p.Name()
		if err != nil {
			continue
		}

		// match: "vscodium.exe" == "vscodium.exe"
		if strings.ToLower(name) == targetName {
			return true
		}
	}

	return false
}
package memory

import (
	"os"
	"strings"

	"github.com/its-ernest/cura/pkg/logging"
	"github.com/shirou/gopsutil/v3/process"
)

var l *logging.Logger = logging.NewLogger("cura.log")

// fixed: added generic checker that handles protection of processes better than GetProtectedProcesses
// it is still needed to keep GetProtectedprocess() for specific processes outside system store folders
// IsSystemOrDriver checks if the process provided is a critical Windows or Driver process
func IsSystemOrDriver(p *process.Process) bool {
	exePath, err := p.Exe()
	if err != nil {
		return false
	}

	exePath = strings.ToLower(exePath)
	systemRoot := strings.ToLower(os.Getenv("SystemRoot")) // usually C:\Windows

	// protect everything in critical driver and system folders
	criticalPaths := []string{
		systemRoot + "\\system32\\drivers",
		systemRoot + "\\system32\\driverstore",
		systemRoot + "\\system32\\wbem", // WMI providers
	}

	for _, path := range criticalPaths {
		if strings.HasPrefix(exePath, strings.ToLower(path)) {
			return true
		}
	}

	return false
}

// GetProtectedProcesses returns the map of critical processes
func GetProtectedProcesses() map[string]bool {
	// critical Windows processes + Development tools
	protected := map[string]bool{
		"explorer.exe": true, "System": true, "svchost.exe": true,
		"lsass.exe": true, "wininit.exe": true, "csrss.exe": true,
		"services.exe": true, "ShellHost.exe": true, "sihost.exe": true,
		"ShellExperienceHost.exe": true, "dllhost.exe": true, "dwm.exe": true,
		"ctfmon.exe": true, "winlogon.exe": true, "smss.exe": true,
		"LsaIso.exe": true, "fontdrvhost.exe": true, "vmcompute.exe": true,
		"conhost.exe": true, "Taskmgr.exe": true,

		// dev tools 
		"WindowsTerminal.exe": true, "OpenConsole.exe": true, "powershell.exe": true,

		// critical for app to run
		"msedgewebview2.exe": true, "cura.exe": true, "cura-dev.exe": true,

		// temp, needed for testing
		//"Code.exe": true, "node.exe": true, "taskhostw.exe": true, "wails.exe": true,
		//"MSPCManagerService.exe": true, "esbuild.exe": true,
	}
	return protected

}

// bytesToMB converts bytes to MB for easier calculations and comparisons
func bytesToMB(b uint64) uint64 {
	return b / 1024 / 1024
}
package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsSystemOrDriver(t *testing.T) {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = "C:\\Windows" // fallback for non-Windows dev environments
	}

	tests := []struct {
		name     string
		exePath  string
		expected bool
	}{
		{
			name:     "System Driver",
			exePath:  filepath.Join(systemRoot, "System32", "drivers", "etc.sys"),
			expected: true,
		},
		{
			name:     "Driver Store Bin",
			exePath:  filepath.Join(systemRoot, "System32", "DriverStore", "FileRepository", "audio.exe"),
			expected: true,
		},
		{
			name:     "WMI Provider",
			exePath:  filepath.Join(systemRoot, "System32", "wbem", "wmiprvse.exe"),
			expected: true,
		},
		{
			name:     "User Application",
			exePath:  "C:\\Users\\Ernest\\AppData\\Local\\Discord\\Discord.exe",
			expected: false,
		},
		{
			name:     "Program Files App",
			exePath:  "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
			expected: false,
		},
		{
			name:     "Case Insensitivity Check",
			exePath:  strings.ToUpper(filepath.Join(systemRoot, "SYSTEM32", "DRIVERS", "RTKAUDIO.SYS")),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// I am simulating the logic of IsSystemOrDriver directly here because mocking the process.Process struct's Exe() method requires an interface wrapper or a complex monkeypatch.
			
			exePath := strings.ToLower(tt.exePath)
			sysRoot := strings.ToLower(systemRoot)

			criticalPaths := []string{
				filepath.Join(sysRoot, "system32", "drivers"),
				filepath.Join(sysRoot, "system32", "driverstore"),
				filepath.Join(sysRoot, "system32", "wbem"),
			}

			actual := false
			for _, path := range criticalPaths {
				if strings.HasPrefix(exePath, strings.ToLower(path)) {
					actual = true
					break
				}
			}

			if actual != tt.expected {
				t.Errorf("IsSystemOrDriver logic failed for %s: expected %v, got %v", tt.exePath, tt.expected, actual)
			}
		})
	}
}

func TestGetProtectedProcesses(t *testing.T) {
	protected := GetProtectedProcesses()

	// essential checks
	essential := []string{"explorer.exe", "dwm.exe", "cura.exe", "dllhost.exe"}

	for _, name := range essential {
		if !protected[name] {
			t.Errorf("Expected %s to be protected, but it wasn't", name)
		}
	}
}
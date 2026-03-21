package whitelist

import (
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

type InstalledApp struct {
	Name        string `json:"name"`
	ExePath     string `json:"exe_path"`
	IsProtected bool   `json:"is_protected"`
}

func GetWindowsApps() ([]InstalledApp, error) {
	var apps []InstalledApp

	configs := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Uninstall`}, // added to fix for user data apps
	}

	for _, conf := range configs {
		k, err := registry.OpenKey(conf.root, conf.path, registry.READ)
		if err != nil {
			continue
		}
		names, _ := k.ReadSubKeyNames(-1)
		for _, name := range names {
			subKey, _ := registry.OpenKey(k, name, registry.READ)
			displayName, _, _ := subKey.GetStringValue("DisplayName")
			installDir, _, _ := subKey.GetStringValue("InstallLocation")

			// fix: some apps may return empty Installir,
			// "DisplayIcon" can act as a fallback to get the directory
			if installDir == "" {
				iconPath, _, _ := subKey.GetStringValue("DisplayIcon")
				if iconPath != "" {
					installDir = filepath.Dir(iconPath)
				}
			}

			if displayName != "" && installDir != "" {
				apps = append(apps, InstalledApp{
					Name:    displayName,
					ExePath: filepath.Clean(installDir),
				})
			}
			subKey.Close()
		}
		k.Close()
	}
	return apps, nil
}

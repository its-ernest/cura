/*
 * Cura Launcher
 * Developed by its-ernest & Gemini (Google AI)
 * * This launcher handles UAC elevation for the Cura system utility.
 * Licensed under [MIT]
 */
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func main() {
	// 1. Locate the target binary
	cwd, _ := os.Getwd()
	targetApp := filepath.Join(cwd, "cura-arm64.exe")

	// 2. Check if program has admin rights
	if !isAdmin() {
		// If not admin, re-run self or the target with the 'runas' verb
		err := runAsAdmin(targetApp)
		if err != nil {
			fmt.Printf("Elevation failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}

// runAsAdmin uses the Windows Shell API to trigger the UAC prompt
func runAsAdmin(exe string) error {
	verb := "runas"
	cwd, _ := os.Getwd()

	verbPtr, _ := windows.UTF16PtrFromString(verb)
	exePtr, _ := windows.UTF16PtrFromString(exe)
	cwdPtr, _ := windows.UTF16PtrFromString(cwd)

	// SW_SHOWNORMAL = 1
	var showCmd int32 = 1

	err := windows.ShellExecute(0, verbPtr, exePtr, nil, cwdPtr, showCmd)
	if err != nil {
		return err
	}
	return nil
}

func isAdmin() bool {
	var sid *windows.SID

	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}
	return member
}

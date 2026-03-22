package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func DownloadAndInstall(release *GitHubRelease) error {
	// determine architecture suffix
	suffix := "amd64.zip"
	if runtime.GOARCH == "arm64" {
		suffix = "arm64.zip"
	}

	// find asset URL
	var downloadURL string
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, suffix) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no suitable asset found for %s", runtime.GOARCH)
	}

	// download the zip to a temp file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpZip, err := os.CreateTemp("", "cura-update-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpZip.Name())

	if _, err := io.Copy(tmpZip, resp.Body); err != nil {
		return err
	}
	tmpZip.Close()

	// get current executable path and base directory
	executablePath, _ := os.Executable()
	baseDir := filepath.Dir(executablePath)

	// open the zip and extract all contents
	r, err := zip.OpenReader(tmpZip.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// determine the target path relative to the app directory
		targetPath := filepath.Join(baseDir, f.Name)
		if strings.HasSuffix(f.Name, ".toml") || strings.HasSuffix(f.Name, ".log") {
			continue
		}
		fmt.Println("UPDATING: Captured file: " + f.Name)

		// handle directory creation
		if f.FileInfo().IsDir() {
			os.MkdirAll(targetPath, f.Mode())
			continue
		}

		// ensure parent directories exist for new nested files
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// windows os swap:
		// if this file is the currently running executable, rename it first
		if strings.EqualFold(targetPath, executablePath) {
			oldPath := executablePath + ".old"
			os.Remove(oldPath) // clean up any previous failed updates

			err := os.Rename(executablePath, oldPath)
			if err != nil {
				return fmt.Errorf("failed to rename current binary: %v", err)
			}
		}

		// extract the file from the archive
		rc, err := f.Open()
		if err != nil {
			return err
		}

		newFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		if _, err := io.Copy(newFile, rc); err != nil {
			newFile.Close()
			rc.Close()
			return err
		}

		newFile.Close()
		rc.Close()
	}

	return nil
}

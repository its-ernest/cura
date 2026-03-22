package updater

import (
	"encoding/json"
	"net/http"

	//"path/filepath"
	"time"
)

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

func CheckForUpdates(currentVersion string) (*GitHubRelease, bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/its-ernest/cura/releases/latest")
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, err
	}

	if release.TagName != currentVersion {
		return &release, true, nil
	}

	return &release, false, nil
}

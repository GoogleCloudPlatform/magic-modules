package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// PRFile represents a file modified in a pull request.
type PRFile struct {
	Path string `json:"path"`
}

// PR represents a GitHub Pull Request with its number and modified files.
type PR struct {
	Number int      `json:"number"`
	Files  []PRFile `json:"files"`
}

// NormalizePath strips the "mmv1/" prefix and cleans the path so that paths
// match across public magic-modules and private EAP overrides layout.
func NormalizePath(p string) string {
	p = filepath.Clean(p)
	p = strings.TrimPrefix(p, "mmv1/")
	return p
}

// GetFilesTouchedByOpenPRs queries the GitHub CLI for open pull requests
// updated in the last N days and returns a map of normalized paths to PR numbers.
func GetFilesTouchedByOpenPRs(days int) (map[string][]int, error) {
	sinceDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	searchQuery := fmt.Sprintf("state:open updated:>=%s", sinceDate)

	cmd := exec.Command("gh", "pr", "list",
		"-R", "GoogleCloudPlatform/magic-modules",
		"--limit", "1000",
		"--search", searchQuery,
		"--json", "number,files",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run gh command: %w (stderr: %s)", err, stderr.String())
	}

	var prs []PR
	if err := json.Unmarshal(stdout.Bytes(), &prs); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from gh output: %w", err)
	}

	touchedFiles := make(map[string][]int)
	for _, pr := range prs {
		for _, f := range pr.Files {
			norm := NormalizePath(f.Path)
			touchedFiles[norm] = append(touchedFiles[norm], pr.Number)
		}
	}

	return touchedFiles, nil
}

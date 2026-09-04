package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const disallowLegacyApiaryDesc = "Check git diff for forbidden legacy Apiary Compute client imports and invocations."

var (
	legacyComputeImportRegex = regexp.MustCompile(`^\+[[:space:]]*(?:import[[:space:]]+)?(?:[a-zA-Z0-9_.]+[[:space:]]+)?"google\.golang\.org/api/compute/`)
	legacyComputeCallRegex   = regexp.MustCompile(`^\+[[:space:]]*.*(?:DEPRECATED_LegacyApiaryClient|tpgcompute\.(?:NewClient|DEPRECATED_LegacyApiaryClient)|compute_tpg\.(?:NewClient|DEPRECATED_LegacyApiaryClient))`)

	checkedPathspecs = []string{
		"mmv1/third_party/terraform/**.go",
		"mmv1/third_party/terraform/**.go.tmpl",
		"mmv1/third_party/cai2hcl/**.go",
		"mmv1/third_party/tgc*/**.go*",
	}
)

type DisallowedMatch struct {
	File    string
	Line    string
	Pattern string
}

type disallowLegacyApiaryOptions struct {
	rootOptions *rootOptions
	stdout      io.Writer
	stderr      io.Writer
	baseRef     string
	repoDir     string
	diffFile    string
}

func newDisallowLegacyApiaryCmd(rootOptions *rootOptions) *cobra.Command {
	o := &disallowLegacyApiaryOptions{
		rootOptions: rootOptions,
		stdout:      os.Stdout,
		stderr:      os.Stderr,
	}

	command := &cobra.Command{
		Use:   "disallow-legacy-apiary",
		Short: disallowLegacyApiaryDesc,
		Long:  disallowLegacyApiaryDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}

	command.Flags().StringVar(&o.baseRef, "base-ref", "", "base git ref to diff against (e.g. origin/main)")
	command.Flags().StringVar(&o.repoDir, "repo-dir", "", "path to repository root (defaults to git repository root)")
	command.Flags().StringVar(&o.diffFile, "diff-file", "", "optional path to a unified diff file to inspect instead of running git diff")

	return command
}

func (o *disallowLegacyApiaryOptions) run() error {
	var diffReader io.Reader

	if o.diffFile != "" {
		if o.diffFile == "-" {
			diffReader = os.Stdin
		} else {
			f, err := os.Open(o.diffFile)
			if err != nil {
				return fmt.Errorf("failed to open diff file %s: %w", o.diffFile, err)
			}
			defer f.Close()
			diffReader = f
		}
	} else {
		repoDir := o.repoDir
		if repoDir == "" {
			var err error
			repoDir, err = findRepoRoot(".")
			if err != nil {
				return err
			}
		}

		baseRef := o.baseRef
		if baseRef == "" || baseRef == "origin/" {
			var err error
			baseRef, err = resolveDefaultBaseRef(repoDir)
			if err != nil {
				return err
			}
		}

		diffOutput, err := runGitDiff(repoDir, baseRef, checkedPathspecs)
		if err != nil {
			return err
		}
		diffReader = bytes.NewReader(diffOutput)
	}

	matches, err := CheckLegacyApiaryDiff(diffReader)
	if err != nil {
		return err
	}

	if len(matches) > 0 {
		if os.Getenv("GITHUB_ACTIONS") == "true" {
			fmt.Fprintln(o.stderr, "::error::New calls or imports of the legacy Apiary Compute client are strictly forbidden:")
		} else {
			fmt.Fprintln(o.stderr, "New calls or imports of the legacy Apiary Compute client are strictly forbidden:")
		}

		for _, m := range matches {
			if m.File != "" {
				fmt.Fprintf(o.stderr, "  %s: %s\n", m.File, strings.TrimSpace(m.Line))
			} else {
				fmt.Fprintf(o.stderr, "  %s\n", strings.TrimSpace(m.Line))
			}
		}

		fmt.Fprintln(o.stderr, "Please use transport_tpg.SendRequest instead. See PR #16847 for details.")
		return fmt.Errorf("found %d forbidden legacy Apiary Compute client reference(s)", len(matches))
	}

	return nil
}

// CheckLegacyApiaryDiff scans a unified diff stream and identifies newly added lines
// that import or call the legacy Apiary Compute client.
func CheckLegacyApiaryDiff(r io.Reader) ([]DisallowedMatch, error) {
	var matches []DisallowedMatch
	var currentFile string

	scanner := bufio.NewScanner(r)
	// Allow scanning lines up to 10MB to handle large diff hunks
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "diff --git ") {
			parts := strings.Fields(line)
			if len(parts) >= 4 && strings.HasPrefix(parts[3], "b/") {
				currentFile = strings.TrimPrefix(parts[3], "b/")
			}
			continue
		}

		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			continue
		}

		// Ignore diff headers and non-added lines
		if strings.HasPrefix(line, "+++") || !strings.HasPrefix(line, "+") {
			continue
		}

		if legacyComputeImportRegex.MatchString(line) {
			matches = append(matches, DisallowedMatch{
				File:    currentFile,
				Line:    strings.TrimPrefix(line, "+"),
				Pattern: "legacy compute import",
			})
		} else if legacyComputeCallRegex.MatchString(line) {
			matches = append(matches, DisallowedMatch{
				File:    currentFile,
				Line:    strings.TrimPrefix(line, "+"),
				Pattern: "legacy compute invocation",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading diff: %w", err)
	}

	return matches, nil
}

func findRepoRoot(startDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = startDir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to find git repository root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func resolveDefaultBaseRef(repoDir string) (string, error) {
	candidates := []string{"origin/main", "upstream/main", "main"}
	for _, c := range candidates {
		cmd := exec.Command("git", "rev-parse", "--verify", c)
		cmd.Dir = repoDir
		if err := cmd.Run(); err == nil {
			return c, nil
		}
	}
	return "HEAD~1", nil
}

func runGitDiff(repoDir, baseRef string, pathspecs []string) ([]byte, error) {
	args := []string{"diff", "-U0", baseRef, "--"}
	args = append(args, pathspecs...)

	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git diff failed: %s (stderr: %s)", exitErr, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	return out, nil
}

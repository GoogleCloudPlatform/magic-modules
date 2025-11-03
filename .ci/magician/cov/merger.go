package cov

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Merger stores coverage related folders and files with a command line runner.
// File structure is:
// workdir/unit-test-cov/
// workdir/vcr-test-cov/
// workdir/merged-test-cov/
// workdir/cov.txt
// workdir/cov.html
type Merger struct {
	rnr            ExecRunner
	workDir        string
	UnitTestCovDir string
	VcrTestCovDir  string
	MergedDir      string
}

func NewTestCovMerger(rnr ExecRunner, workDir string) (*Merger, error) {
	unitTestDir := filepath.Join(workDir, "unit-test-cov")
	vcrTestDir := filepath.Join(workDir, "vcr-test-cov")
	mergedDir := filepath.Join(workDir, "merged-test-cov")
	for _, dir := range []string{unitTestDir, vcrTestDir, mergedDir} {
		if err := rnr.Mkdir(dir); err != nil {
			return nil, fmt.Errorf("failed to create dir %s: %w", dir, err)
		}
	}
	return &Merger{
		workDir:        workDir,
		rnr:            rnr,
		UnitTestCovDir: unitTestDir,
		VcrTestCovDir:  vcrTestDir,
		MergedDir:      mergedDir,
	}, nil
}

func (m *Merger) HTMLCovPath() string {
	return filepath.Join(m.workDir, "cov.html")
}

func (m *Merger) Merge() error {
	if isFolderEmpty(m.UnitTestCovDir) && isFolderEmpty(m.VcrTestCovDir) {
		return fmt.Errorf("no coverage data found in provided folders")
	}

	covTxtPath := filepath.Join(m.workDir, "cov.txt")
	covHTMLPath := m.HTMLCovPath()

	if _, err := m.rnr.Run(
		"go",
		[]string{
			"tool",
			"covdata",
			"merge",
			fmt.Sprintf("-i=%s,%s", m.UnitTestCovDir, m.VcrTestCovDir),
			"-o=" + m.MergedDir,
		},
		nil,
	); err != nil {
		return fmt.Errorf("failed to merge coverage data: %s", err)
	}

	if _, err := m.rnr.Run(
		"go",
		[]string{
			"tool",
			"covdata",
			"textfmt",
			"-i=" + m.MergedDir,
			"-o=" + covTxtPath,
		},
		nil,
	); err != nil {
		return fmt.Errorf("failed to convert coverage data to text format: %s", err)
	}

	if _, err := m.rnr.Run(
		"go",
		[]string{
			"tool",
			"cover",
			"-html=" + covTxtPath,
			"-o=" + covHTMLPath,
		},
		nil,
	); err != nil {
		return fmt.Errorf("failed to convert coverage data to text format: %s", err)
	}
	return nil
}

func (m *Merger) UploadToGCS(gcsPrefix string, buildID string) (string, error) {
	bucketName := strings.TrimPrefix(gcsPrefix, "gs://")
	gcsPath := fmt.Sprintf("gs://%s/%s/", bucketName, buildID)
	fmt.Printf("Uploading coverage result to %s\n", gcsPath)
	args := []string{"-m", "cp", m.HTMLCovPath(), gcsPath}
	if _, err := m.rnr.Run("gsutil", args, nil); err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error upload cov html: %w", err)
	}
	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s/cov.html", bucketName, buildID), nil
}

func (m *Merger) PackageCovComment() (string, error) {
	out, err := m.rnr.Run(
		"go",
		[]string{
			"tool",
			"covdata",
			"percent",
			"-i=" + m.MergedDir,
		},
		nil,
	)
	if err != nil {
		return "", err
	}

	covList := strings.Split(strings.TrimSpace(out), "\n")

	commentTemplate := `
<details>
<summary>Click here to see Test Coverage Metrics </summary>
<blockquote>
<ul>
{{range .}}
<li>{{. -}}</li>
{{end}}
</ul>
</blockquote>
</details>
	`

	// Create a new template and parse the letter into it.
	sb := new(strings.Builder)
	t := template.Must(template.New("commentTemplate").Parse(commentTemplate))
	err = t.Execute(sb, covList)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(sb.String()), nil
}

func isFolderEmpty(dirPath string) bool {
	file, err := os.Open(dirPath)
	if err != nil {
		return true
	}
	defer file.Close()

	_, err = file.Readdirnames(1)
	return err != nil
}

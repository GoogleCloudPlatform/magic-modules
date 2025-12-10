package cov

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type MockRunner interface {
	ExecRunner
}

type ParameterList []any

type mockRunner struct {
	dirs       []string
	commands   []ParameterList
	cmdResults map[string]string
}

func (r *mockRunner) Mkdir(path string) error {
	r.dirs = append(r.dirs, path)
	return nil
}

func (r *mockRunner) Run(name string, args []string, env map[string]string) (string, error) {
	r.commands = append(r.commands, ParameterList{name, args})
	cmd := fmt.Sprintf("%s %v", name, args)
	if result, ok := r.cmdResults[cmd]; ok {
		return result, nil
	}
	return "", nil
}

func TestNewTestCovMerger(t *testing.T) {
	rnr := &mockRunner{}
	_, err := NewTestCovMerger(rnr, "/tmp")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{"/tmp/unit-test-cov", "/tmp/vcr-test-cov", "/tmp/merged-test-cov"}, rnr.dirs); diff != "" {
		t.Errorf("NewTestCovMerger did not create expected folders: (-want, +got) = %s", diff)
	}
}

func TestMergeFail(t *testing.T) {
	workdir := t.TempDir()

	rnr := &mockRunner{}
	merger, err := NewTestCovMerger(rnr, workdir)
	if err != nil {
		t.Fatal(err)
	}
	err = merger.Merge()
	if err == nil {
		t.Fatal("expect failure since folders are empty, but got nil err")
	}

	if !strings.Contains(err.Error(), "no coverage data found") {
		t.Errorf("Merge() got unexpected error: %s", err)
	}

}

func TestMerge(t *testing.T) {
	workdir := t.TempDir()
	err := os.MkdirAll(filepath.Join(workdir, "unit-test-cov"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(filepath.Join(workdir, "unit-test-cov", "cov.data"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	rnr := &mockRunner{}
	merger := Merger{
		UnitTestCovDir: filepath.Join(workdir, "unit-test-cov"),
		VcrTestCovDir:  filepath.Join(workdir, "vcr-test-cov"),
		MergedDir:      filepath.Join(workdir, "merged-test-cov"),
		workDir:        workdir,
		rnr:            rnr,
	}

	if err := merger.Merge(); err != nil {
		t.Fatal(err)
	}

	want := []ParameterList{
		{"go", []string{"tool", "covdata", "merge", "-i=" + workdir + "/unit-test-cov," + workdir + "/vcr-test-cov", "-o=" + workdir + "/merged-test-cov"}},
		{"go", []string{"tool", "covdata", "textfmt", "-i=" + workdir + "/merged-test-cov", "-o=" + workdir + "/cov.txt"}},
		{"go", []string{"tool", "cover", "-html=" + workdir + "/cov.txt", "-o=" + workdir + "/cov.html"}},
	}
	if diff := cmp.Diff(want, rnr.commands); diff != "" {
		t.Errorf("Merge got different commands: %s", diff)
	}
}

func TestUploadToGCS(t *testing.T) {
	workdir := os.TempDir()
	rnr := &mockRunner{}
	merger := Merger{
		rnr:     rnr,
		workDir: workdir,
	}

	got, err := merger.UploadToGCS("gs://bucket/path", "12345")
	if err != nil {
		t.Fatal(err)
	}

	want := "https://storage.cloud.google.com/bucket/path/12345/cov.html"
	if got != want {
		t.Errorf("UploadToGCS got = %s, want = %s", got, want)
	}

	wantCommands := []ParameterList{
		{"gsutil", []string{"-m", "cp", filepath.Join(workdir, "cov.html"), "gs://bucket/path/12345/"}},
	}
	if diff := cmp.Diff(wantCommands, rnr.commands); diff != "" {
		t.Errorf("UploadToGCS got different commands: %s", diff)
	}
}

func TestPackageCovComment(t *testing.T) {
	workdir := os.TempDir()
	rnr := &mockRunner{
		cmdResults: map[string]string{
			"go [tool covdata percent -i=" + filepath.Join(workdir, "merged-test-cov]"): "pkg1 10%\npkg2 20%\n\n",
		},
	}
	merger := Merger{
		rnr:       rnr,
		workDir:   workdir,
		MergedDir: filepath.Join(workdir, "merged-test-cov"),
	}

	got, err := merger.PackageCovComment()
	if err != nil {
		t.Fatal(err)
	}

	want := strings.TrimSpace(`
<details>
<summary>Click here to see Test Coverage Metrics </summary>
<blockquote>
<ul>

<li>pkg1 10%</li>

<li>pkg2 20%</li>

</ul>
</blockquote>
</details>
	`)
	if got != want {
		t.Errorf("PackageCovComment got = %s, want = %s", got, want)
	}
}

package vcr

import (
	"fmt"
	"magician/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const numDirs = 11 // number of directories in the provider to run tests in, -1 for all

type Tester interface {
	CloneProvider(goPath, githubUsername, githubToken string, version Version) error
	FetchCassettes(version Version) error
	Run(mode Mode, version Version) (*Result, error)
	Cleanup() error
}

type Result struct {
	FailedTests []string
}

type Version int

const (
	GA Version = iota
	Beta
)

const numVersions = 2

func (v Version) String() string {
	switch v {
	case GA:
		return "ga"
	case Beta:
		return "beta"
	}
	return "unknown"
}

func (v Version) BucketPath() string {
	if v == GA {
		return ""
	}
	return v.String() + "/"
}

// TODO: move this into magician/github
func (v Version) RepoName() string {
	switch v {
	case GA:
		return "terraform-provider-google"
	case Beta:
		return "terraform-provider-google-beta"
	}
	return "unknown"
}

// TODO: move this into magician/github
func (v Version) URL(githubUsername, githubToken string) string {
	return fmt.Sprintf("https://%s:%s@github.com/%s/%s", githubUsername, githubToken, githubUsername, v.RepoName())
}

type Mode int

const (
	Replaying Mode = iota
	Recording
)

const numModes = 2

func (m Mode) Lower() string {
	switch m {
	case Replaying:
		return "replaying"
	case Recording:
		return "recording"
	}
	return "unknown"
}

func (m Mode) Upper() string {
	return strings.ToUpper(m.Lower())
}

type logKey struct {
	mode    Mode
	version Version
}

type vcrTester struct {
	r             vcrRunner          // for running commands and manipulating files
	version       Version            // either "ga" or "beta", defaults to "ga"
	mode          Mode               // either "REPLAYING" or "RECORDING", defaults to "REPLAYING"
	baseDir       string             // the directory in which this tester was created
	saKeyPath     string             // where sa_key.json is relative to baseDir
	cassettePaths map[Version]string // where cassettes are relative to baseDir by version
	logPaths      map[logKey]string  // where logs are relative to baseDir by version and mode
	repoPaths     map[Version]string // relative paths of already cloned repos by version
}

type vcrRunner interface {
	GetCWD() string
	PushDir(path string) error
	PopDir() error
	Mkdir(path string) error
	WriteFile(name, data string) error
	ReadFile(name string) (string, error)
	RemoveAll(path string) error
	Run(name string, args, env []string) (string, error)
}

const accTestParalellism = 32

const replayingTimeout = "240m"

var failedTestsExpression = regexp.MustCompile(`(?m:^--- FAIL: TestAcc(\w+))`)

// Create a new tester in the current working directory and write the service account key file.
func NewTester(goPath, saKey string) (Tester, error) {
	r, err := exec.NewRunner()
	if err != nil {
		return nil, err
	}
	saKeyPath := "sa_key.json"
	if err := r.WriteFile(saKeyPath, saKey); err != nil {
		return nil, err
	}
	return &vcrTester{
		r:         r,
		baseDir:   r.GetCWD(),
		saKeyPath: saKeyPath,
		cassettePaths: map[Version]string{
			// here for local testing purposes
			Beta: "cassettes/beta",
		},
		logPaths: make(map[logKey]string, numVersions*numModes),
		repoPaths: map[Version]string{
			// here for local testing purposes
			GA:   "/usr/local/google/home/thomasrodgers/go/src/github.com/hashicorp/terraform-provider-google",
			Beta: "/usr/local/google/home/thomasrodgers/go/src/github.com/hashicorp/terraform-provider-google-beta",
		},
	}, nil
}

// Clone the provider with the given url to a local path under the go path and store it by version in the paths map.
func (vt *vcrTester) CloneProvider(goPath, githubUsername, githubToken string, version Version) error {
	if _, ok := vt.repoPaths[version]; ok {
		// Skip cloning an already cloned repo.
		return nil
	}
	repoPath := filepath.Join(goPath, "src", "github.com", githubUsername, version.RepoName())
	vt.repoPaths[version] = repoPath
	if _, err := vt.r.Run("git", []string{"clone", version.URL(githubUsername, githubToken), repoPath}, nil); err != nil {
		return err
	}
	return nil
}

// Fetch the cassettes for the current version.
// Should be run from the base dir.
func (vt *vcrTester) FetchCassettes(version Version) error {
	cassettePath, ok := vt.cassettePaths[version]
	if ok {
		return nil
	}
	cassettePath = filepath.Join("cassettes", version.String())
	vt.r.Mkdir(cassettePath)
	bucketPath := fmt.Sprintf("gs://ci-vcr-cassettes/%sfixtures/*", version.BucketPath())
	// Fetch the cassettes.
	args := []string{"-m", "-q", "cp", bucketPath, cassettePath}
	fmt.Println("Fetching cassettes:\n", "gsutil", strings.Join(args, " "))
	if _, err := vt.r.Run("gsutil", args, nil); err != nil {
		return err
	}
	return nil
}

// Run the vcr tests in the given mode and provider version and return the result.
// This will overwrite any existing logs for the current version and the given mode.
func (vt *vcrTester) Run(mode Mode, version Version) (*Result, error) {
	logPath, ok := vt.logPaths[logKey{mode, version}]
	if !ok {
		// We've never run this mode and version.
		logPath = filepath.Join("testlogs", mode.Lower(), version.String())
		if err := vt.r.Mkdir(logPath); err != nil {
			return nil, err
		}
	}

	repoPath, ok := vt.repoPaths[version]
	if !ok {
		return nil, fmt.Errorf("no repo cloned for version %s in %v", version, vt.repoPaths)
	}
	if err := vt.r.PushDir(repoPath); err != nil {
		return nil, err
	}
	testDirs, err := vt.googleTestDirectory()
	if err != nil {
		return nil, err
	}

	if numDirs > -1 {
		testDirs = testDirs[:numDirs]
	}

	args := []string{"test"}
	args = append(args, testDirs...)
	args = append(args,
		"-parallel",
		strconv.Itoa(accTestParalellism),
		"-v",
		"-run=TestAcc",
		"-timeout",
		replayingTimeout,
		`-ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc`,
	)
	env := []string{
		"GOOGLE_REGION=us-central1",
		"GOOGLE_ZONE=us-central1-a",
		fmt.Sprintf("VCR_PATH=%s", filepath.Join(vt.baseDir, vt.cassettePaths[version])),
		"VCR_MODE=" + mode.Upper(),
		fmt.Sprintf("ACCTEST_PARALLELISM=%d", accTestParalellism),
		"GOOGLE_CREDENTIALS=" + vt.saKeyPath,
		fmt.Sprintf("GOOGLE_APPLICATION_CREDENTIALS=%s/sa_key.json", vt.baseDir),
		"GOOGLE_TEST_DIRECTORY=" + strings.Join(testDirs, " "),
		"TF_LOG=DEBUG",
		"TF_LOG_SDK_FRAMEWORK=INFO",
		fmt.Sprintf("TF_LOG_PATH_MASK=%s", filepath.Join(vt.baseDir, logPath, "%s.log")),
		"TF_ACC=1",
		"TF_SCHEMA_PANIC_ON_ERROR=1",
	}
	fmt.Printf(`Running go:
	env:
%s
	args:
%s
`, strings.Join(env, "\n"), strings.Join(args, "\n"))
	output, err := vt.r.Run("go", args, env)
	if err != nil {
		// Use error as output for log.
		output = fmt.Sprintf("Error replaying tests:\n%v", err)
	}
	// Leave repo directory.
	if err := vt.r.PopDir(); err != nil {
		return nil, err
	}

	logFileName := filepath.Join(logPath, "all_tests.log")
	// Write output (or error) to test log.
	if err := vt.r.WriteFile(logFileName, output); err != nil {
		return nil, fmt.Errorf("error writing replaying log: %v, test output: %v", err, output)
	}
	// Collect failed tests from output.
	failed, err := vt.matchingTests(output, failedTestsExpression)
	if err != nil {
		return nil, err
	}
	return &Result{FailedTests: failed}, nil
}

// Deletes the service account key and the repos.
func (vt *vcrTester) Cleanup() error {
	if err := vt.r.RemoveAll(vt.saKeyPath); err != nil {
		return err
	}

	/* skip for testing
	for _, path := range vt.repoPaths {
		if err := vt.r.RemoveAll(path); err != nil {
			return err
		}
	}*/
	return nil
}

// Returns a list of all directories to run tests in.
// Must be called after chaging into the provider dir.
func (vt *vcrTester) googleTestDirectory() ([]string, error) {
	var testDirs []string
	if allPackages, err := vt.r.Run("go", []string{"list", "./..."}, nil); err != nil {
		return nil, err
	} else {
		for _, dir := range strings.Split(allPackages, "\n") {
			if !strings.Contains(dir, "github.com/hashicorp/terraform-provider-google-beta/scripts") {
				testDirs = append(testDirs, dir)
			}
		}
	}
	return testDirs, nil
}

// Reads the log at the given path and returns a slice of tests matching the expression.
func (vt *vcrTester) matchingTests(output string, exp *regexp.Regexp) ([]string, error) {
	matches := exp.FindAllStringSubmatch(output, -1)
	var tests []string
	for _, submatches := range matches {
		if len(submatches) != 2 {
			return nil, fmt.Errorf("unexpected submatches for failed tests expression: %v", submatches)
		}
		tests = append(tests, submatches[1])
	}
	return tests, nil
}

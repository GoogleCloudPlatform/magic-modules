package vcr

import (
	"fmt"
	"io/fs"
	"magician/exec"
	"magician/provider"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Result struct {
	PassedTests  []string
	SkippedTests []string
	FailedTests  []string
	Panics       []string
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
	version provider.Version
}

type Tester struct {
	env           map[string]string           // shared environment variables for running tests
	rnr           exec.ExecRunner             // for running commands and manipulating files
	baseDir       string                      // the directory in which this tester was created
	saKeyPath     string                      // where sa_key.json is relative to baseDir
	cassettePaths map[provider.Version]string // where cassettes are relative to baseDir by version
	logPaths      map[logKey]string           // where logs are relative to baseDir by version and mode
	repoPaths     map[provider.Version]string // relative paths of already cloned repos by version
}

const accTestParallelism = 32
const parallelJobs = 16

const replayingTimeout = "240m"

var testResultsExpression = regexp.MustCompile(`(?m:^--- (PASS|FAIL|SKIP): (TestAcc\w+))`)

var testPanicExpression = regexp.MustCompile(`^panic: .*`)

// Create a new tester in the current working directory and write the service account key file.
func NewTester(env map[string]string, rnr exec.ExecRunner) (*Tester, error) {
	saKeyPath := "sa_key.json"
	if err := rnr.WriteFile(saKeyPath, env["SA_KEY"]); err != nil {
		return nil, err
	}
	return &Tester{
		env:           env,
		rnr:           rnr,
		baseDir:       rnr.GetCWD(),
		saKeyPath:     saKeyPath,
		cassettePaths: make(map[provider.Version]string, provider.NumVersions),
		logPaths:      make(map[logKey]string, provider.NumVersions*numModes),
		repoPaths:     make(map[provider.Version]string, provider.NumVersions),
	}, nil
}

func (vt *Tester) SetRepoPath(version provider.Version, repoPath string) {
	vt.repoPaths[version] = repoPath
}

// Fetch the cassettes for the current version if not already fetched.
// Should be run from the base dir.
func (vt *Tester) FetchCassettes(version provider.Version, baseBranch, prNumber string) error {
	_, ok := vt.cassettePaths[version]
	if ok {
		return nil
	}
	cassettePath := filepath.Join(vt.baseDir, "cassettes", version.String())
	vt.rnr.Mkdir(cassettePath)
	if baseBranch != "FEATURE-BRANCH-major-release-6.0.0" {
		// pull main cassettes (major release uses branch specific casssettes as primary ones)
		bucketPath := fmt.Sprintf("gs://ci-vcr-cassettes/%sfixtures/*", version.BucketPath())
		if err := vt.fetchBucketPath(bucketPath, cassettePath); err != nil {
			fmt.Println("Error fetching cassettes: ", err)
		}
	}
	if baseBranch != "main" {
		bucketPath := fmt.Sprintf("gs://ci-vcr-cassettes/%srefs/branches/%s/fixtures/*", version.BucketPath(), baseBranch)
		if err := vt.fetchBucketPath(bucketPath, cassettePath); err != nil {
			fmt.Println("Error fetching cassettes: ", err)
		}
	}
	if prNumber != "" {
		bucketPath := fmt.Sprintf("gs://ci-vcr-cassettes/%srefs/heads/auto-pr-%s/fixtures/*", version.BucketPath(), prNumber)
		if err := vt.fetchBucketPath(bucketPath, cassettePath); err != nil {
			fmt.Println("Error fetching cassettes: ", err)
		}
	}
	vt.cassettePaths[version] = cassettePath
	return nil
}

func (vt *Tester) fetchBucketPath(bucketPath, cassettePath string) error {
	// Fetch the cassettes.
	args := []string{"-m", "-q", "cp", bucketPath, cassettePath}
	fmt.Println("Fetching cassettes:\n", "gsutil", strings.Join(args, " "))
	if _, err := vt.rnr.Run("gsutil", args, nil); err != nil {
		return err
	}
	return nil
}

// Run the vcr tests in the given mode and provider version and return the result.
// This will overwrite any existing logs for the given mode and version.
func (vt *Tester) Run(mode Mode, version provider.Version, testDirs []string) (*Result, error) {
	logPath, err := vt.getLogPath(mode, version)
	if err != nil {
		return nil, err
	}

	repoPath, ok := vt.repoPaths[version]
	if !ok {
		return nil, fmt.Errorf("no repo cloned for version %s in %v", version, vt.repoPaths)
	}
	if err := vt.rnr.PushDir(repoPath); err != nil {
		return nil, err
	}
	if len(testDirs) == 0 {
		var err error
		testDirs, err = vt.googleTestDirectory()
		if err != nil {
			return nil, err
		}
	}

	cassettePath := filepath.Join(vt.baseDir, "cassettes", version.String())
	switch mode {
	case Replaying:
		cassettePath, ok = vt.cassettePaths[version]
		if !ok {
			return nil, fmt.Errorf("cassettes not fetched for version %s", version)
		}
	case Recording:
		if err := vt.rnr.RemoveAll(cassettePath); err != nil {
			return nil, fmt.Errorf("error removing cassettes: %v", err)
		}
		if err := vt.rnr.Mkdir(cassettePath); err != nil {
			return nil, fmt.Errorf("error creating cassette dir: %v", err)
		}
		vt.cassettePaths[version] = cassettePath
	}

	args := []string{"test"}
	args = append(args, testDirs...)
	args = append(args,
		"-parallel",
		strconv.Itoa(accTestParallelism),
		"-v",
		"-run=TestAcc",
		"-timeout",
		replayingTimeout,
		"-ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc",
		"-vet=off",
	)
	env := map[string]string{
		"VCR_PATH":                       cassettePath,
		"VCR_MODE":                       mode.Upper(),
		"ACCTEST_PARALLELISM":            strconv.Itoa(accTestParallelism),
		"GOOGLE_CREDENTIALS":             vt.env["SA_KEY"],
		"GOOGLE_APPLICATION_CREDENTIALS": filepath.Join(vt.baseDir, vt.saKeyPath),
		"GOOGLE_TEST_DIRECTORY":          strings.Join(testDirs, " "),
		"TF_LOG":                         "DEBUG",
		"TF_LOG_SDK_FRAMEWORK":           "INFO",
		"TF_LOG_PATH_MASK":               filepath.Join(logPath, "%s.log"),
		"TF_ACC":                         "1",
		"TF_SCHEMA_PANIC_ON_ERROR":       "1",
	}
	for ev, val := range vt.env {
		env[ev] = val
	}
	var printedEnv string
	for ev, val := range env {
		if ev == "SA_KEY" || strings.HasPrefix(ev, "GITHUB_TOKEN") {
			val = "{hidden}"
		}
		printedEnv += fmt.Sprintf("%s=%s\n", ev, val)
	}
	fmt.Printf(`Running go:
	env:
%v
	args:
%s
`, printedEnv, strings.Join(args, " "))
	output, testErr := vt.rnr.Run("go", args, env)
	if testErr != nil {
		// Use error as output for log.
		output = fmt.Sprintf("Error %s tests:\n%v", mode.Lower(), testErr)
	}
	// Leave repo directory.
	if err := vt.rnr.PopDir(); err != nil {
		return nil, err
	}

	logFileName := filepath.Join(vt.baseDir, "testlogs", fmt.Sprintf("%s_test.log", mode.Lower()))
	// Write output (or error) to test log.
	// Append to existing log file.
	allOutput, _ := vt.rnr.ReadFile(logFileName)
	if allOutput != "" {
		allOutput += "\n"
	}
	allOutput += output
	if err := vt.rnr.WriteFile(logFileName, allOutput); err != nil {
		return nil, fmt.Errorf("error writing log: %v, test output: %v", err, allOutput)
	}
	return collectResult(output), testErr
}

func (vt *Tester) RunParallel(mode Mode, version provider.Version, testDirs, tests []string) (*Result, error) {
	logPath, err := vt.getLogPath(mode, version)
	if err != nil {
		return nil, err
	}
	if err := vt.rnr.Mkdir(filepath.Join(vt.baseDir, "testlogs", mode.Lower()+"_build")); err != nil {
		return nil, err
	}
	repoPath, ok := vt.repoPaths[version]
	if !ok {
		return nil, fmt.Errorf("no repo cloned for version %s in %v", version, vt.repoPaths)
	}
	if err := vt.rnr.PushDir(repoPath); err != nil {
		return nil, err
	}
	if len(testDirs) == 0 {
		var err error
		testDirs, err = vt.googleTestDirectory()
		if err != nil {
			return nil, err
		}
	}

	cassettePath := filepath.Join(vt.baseDir, "cassettes", version.String())
	switch mode {
	case Replaying:
		cassettePath, ok = vt.cassettePaths[version]
		if !ok {
			return nil, fmt.Errorf("cassettes not fetched for version %s", version)
		}
	case Recording:
		if err := vt.rnr.RemoveAll(cassettePath); err != nil {
			return nil, fmt.Errorf("error removing cassettes: %v", err)
		}
		if err := vt.rnr.Mkdir(cassettePath); err != nil {
			return nil, fmt.Errorf("error creating cassette dir: %v", err)
		}
		vt.cassettePaths[version] = cassettePath
	}

	running := make(chan struct{}, parallelJobs)
	outputs := make(chan string, len(testDirs)*len(tests))
	wg := &sync.WaitGroup{}
	wg.Add(len(testDirs) * len(tests))
	errs := make(chan error, len(testDirs)*len(tests)*2)
	for _, testDir := range testDirs {
		for _, test := range tests {
			running <- struct{}{}
			go vt.runInParallel(mode, version, testDir, test, logPath, cassettePath, running, wg, outputs, errs)
		}
	}

	wg.Wait()

	close(outputs)
	close(errs)

	// Leave repo directory.
	if err := vt.rnr.PopDir(); err != nil {
		return nil, err
	}
	var output string
	for otpt := range outputs {
		output += otpt
	}
	logFileName := filepath.Join(vt.baseDir, "testlogs", fmt.Sprintf("%s_test.log", mode.Lower()))
	if err := vt.rnr.WriteFile(logFileName, output); err != nil {
		return nil, err
	}
	var testErr error
	for err := range errs {
		if err != nil {
			testErr = err
			break
		}
	}
	return collectResult(output), testErr
}

func (vt *Tester) runInParallel(mode Mode, version provider.Version, testDir, test, logPath, cassettePath string, running <-chan struct{}, wg *sync.WaitGroup, outputs chan<- string, errs chan<- error) {
	args := []string{
		"test",
		testDir,
		"-parallel",
		"1",
		"-v",
		"-run=" + test + "$",
		"-timeout",
		replayingTimeout,
		"-ldflags=-X=github.com/hashicorp/terraform-provider-google-beta/version.ProviderVersion=acc",
		"-vet=off",
	}
	env := map[string]string{
		"VCR_PATH":                       cassettePath,
		"VCR_MODE":                       mode.Upper(),
		"ACCTEST_PARALLELISM":            "1",
		"GOOGLE_CREDENTIALS":             vt.env["SA_KEY"],
		"GOOGLE_APPLICATION_CREDENTIALS": filepath.Join(vt.baseDir, vt.saKeyPath),
		"GOOGLE_TEST_DIRECTORY":          testDir,
		"TF_LOG":                         "DEBUG",
		"TF_LOG_SDK_FRAMEWORK":           "INFO",
		"TF_LOG_PATH_MASK":               filepath.Join(logPath, "%s.log"),
		"TF_ACC":                         "1",
		"TF_SCHEMA_PANIC_ON_ERROR":       "1",
	}
	for ev, val := range vt.env {
		env[ev] = val
	}
	output, testErr := vt.rnr.Run("go", args, env)
	outputs <- output
	if testErr != nil {
		// Use error as output for log.
		output = fmt.Sprintf("Error %s tests:\n%v", mode.Lower(), testErr)
		errs <- testErr
	}
	logFileName := filepath.Join(vt.baseDir, "testlogs", mode.Lower()+"_build", fmt.Sprintf("%s_%s_test.log", test, mode.Lower()))
	// Write output (or error) to test log.
	// Append to existing log file.
	previousLog, _ := vt.rnr.ReadFile(logFileName)
	if previousLog != "" {
		output = previousLog + "\n" + output
	}
	if err := vt.rnr.WriteFile(logFileName, output); err != nil {
		errs <- fmt.Errorf("error writing log: %v, test output: %v", err, output)
	}
	<-running
	wg.Done()
}

func (vt *Tester) getLogPath(mode Mode, version provider.Version) (string, error) {
	lgky := logKey{mode, version}
	logPath, ok := vt.logPaths[lgky]
	if !ok {
		// We've never run this mode and version.
		logPath = filepath.Join(vt.baseDir, "testlogs", mode.Lower(), version.String())
		if err := vt.rnr.Mkdir(logPath); err != nil {
			return "", err
		}
		vt.logPaths[lgky] = logPath
	}
	return logPath, nil
}

func (vt *Tester) UploadLogs(logBucket, prNumber, buildID string, parallel, afterRecording bool, mode Mode, version provider.Version) error {
	bucketPath := fmt.Sprintf("gs://%s/%s/", logBucket, version)
	if prNumber != "" {
		bucketPath += fmt.Sprintf("refs/heads/auto-pr-%s/", prNumber)
	}
	if buildID != "" {
		bucketPath += fmt.Sprintf("artifacts/%s/", buildID)
	}
	lgky := logKey{mode, version}
	logPath, ok := vt.logPaths[lgky]
	if !ok {
		return fmt.Errorf("no log path found for mode %s and version %s", mode.Lower(), version)
	}
	args := []string{"-h", "Content-Type:text/plain", "-q", "cp", "-r", filepath.Join(vt.baseDir, "testlogs", fmt.Sprintf("%s_test.log", mode.Lower())), bucketPath + "build-log/"}
	fmt.Println("Uploading build log:\n", "gsutil", strings.Join(args, " "))
	if _, err := vt.rnr.Run("gsutil", args, nil); err != nil {
		fmt.Println("Error uploading build log: ", err)
	}
	var suffix string
	if afterRecording {
		suffix = "_after_recording"
	}
	if parallel {
		args := []string{"-h", "Content-Type:text/plain", "-m", "-q", "cp", "-r", filepath.Join(vt.baseDir, "testlogs", mode.Lower()+"_build", "*"), fmt.Sprintf("%sbuild-log/%s_build%s/", bucketPath, mode.Lower(), suffix)}
		fmt.Println("Uploading build logs:\n", "gsutil", strings.Join(args, " "))
		if _, err := vt.rnr.Run("gsutil", args, nil); err != nil {
			fmt.Println("Error uploading build logs: ", err)
		}
	}
	args = []string{"-h", "Content-Type:text/plain", "-m", "-q", "cp", "-r", filepath.Join(logPath, "*"), fmt.Sprintf("%s%s%s/", bucketPath, mode.Lower(), suffix)}
	fmt.Println("Uploading logs:\n", "gsutil", strings.Join(args, " "))
	if _, err := vt.rnr.Run("gsutil", args, nil); err != nil {
		fmt.Println("Error uploading logs: ", err)
	}
	return nil
}

func (vt *Tester) UploadCassettes(logBucket, prNumber string, version provider.Version) error {
	cassettePath, ok := vt.cassettePaths[version]
	if !ok {
		return fmt.Errorf("no cassettes found for version %s", version)
	}
	args := []string{"-m", "-q", "cp", filepath.Join(cassettePath, "*"), fmt.Sprintf("gs://%s/%s/refs/heads/auto-pr-%s/fixtures/", logBucket, version, prNumber)}
	fmt.Println("Uploading cassettes:\n", "gsutil", strings.Join(args, " "))
	if _, err := vt.rnr.Run("gsutil", args, nil); err != nil {
		fmt.Println("Error uploading cassettes: ", err)
	}
	return nil
}

// Deletes the service account key.
func (vt *Tester) Cleanup() error {
	if err := vt.rnr.RemoveAll(vt.saKeyPath); err != nil {
		return err
	}
	return nil
}

// Returns a list of all directories to run tests in.
// Must be called after changing into the provider dir.
func (vt *Tester) googleTestDirectory() ([]string, error) {
	var testDirs []string
	if allPackages, err := vt.rnr.Run("go", []string{"list", "./..."}, nil); err != nil {
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

// Print all log file names and contents, except for all_tests.log.
// Must be called after running tests.
func (vt *Tester) printLogs(logPath string) {
	vt.rnr.Walk(logPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Name() == "all_tests.log" {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println("======= ", info.Name(), " =======")
		if logContent, err := vt.rnr.ReadFile(path); err == nil {
			fmt.Println(logContent)
		}
		return nil
	})
}

func collectResult(output string) *Result {
	matches := testResultsExpression.FindAllStringSubmatch(output, -1)
	resultSets := make(map[string]map[string]struct{}, 4)
	for _, submatches := range matches {
		if len(submatches) != 3 {
			fmt.Printf("Warning: unexpected regex match found in test output: %v", submatches)
			continue
		}
		if _, ok := resultSets[submatches[1]]; !ok {
			resultSets[submatches[1]] = make(map[string]struct{})
		}
		resultSets[submatches[1]][submatches[2]] = struct{}{}
	}
	results := make(map[string][]string, 4)
	results["PANIC"] = testPanicExpression.FindAllString(output, -1)
	sort.Strings(results["PANIC"])
	for _, kind := range []string{"FAIL", "PASS", "SKIP"} {
		for test := range resultSets[kind] {
			results[kind] = append(results[kind], test)
		}
		sort.Strings(results[kind])
	}
	return &Result{
		FailedTests:  results["FAIL"],
		PassedTests:  results["PASS"],
		SkippedTests: results["SKIP"],
		Panics:       results["PANIC"],
	}
}

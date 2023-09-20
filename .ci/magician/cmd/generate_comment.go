package cmd

import (
	"fmt"
	"magician/github"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

type gcGithub interface {
	GetPullRequestLabelIDs(prNumber string) (map[int]struct{}, error)
	PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error
	PostComment(prNumber, comment string) error
}

type runner interface {
	Getwd() (string, error)
	Copy(src, dest string) error
	RemoveAll(path string) error
	Run(path, name string, args, env []string) (string, string, error)
}

type actualRunner struct{}

func (ar *actualRunner) Getwd() (string, error) {
	return os.Getwd()
}

func (ar *actualRunner) Copy(src, dest string) error {
	return cp.Copy(src, dest)
}

func (ar *actualRunner) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (ar actualRunner) Run(path, name string, args, env []string) (string, string, error) {
	if path != "" {
		os.Chdir(path)
	}
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.Output()
	if err != nil {
		exitErr := err.(*exec.ExitError)
		return string(out), string(exitErr.Stderr), err
	}
	return string(out), "", nil
}

var generateCommentCmd = &cobra.Command{
	Use:   "generate-comment",
	Short: "Run presubmit generate comment",
	Long: `This command processes pull requests and performs various validations and actions based on the PR's metadata and author.

	The following PR details are expected as environment variables:
	1. BUILD_ID
	2. PROJECT_ID
	3. BUILD_STEP
	4. COMMIT_SHA
	5. PR_NUMBER
	6. GITHUB_TOKEN

	The command performs the following steps:
	1. Clone the tpg, tpgb, tfc, and tfoics repos from modular-magician.
	2. Compute the diffs between auto-pr-# and auto-pr-#-old branches.
	3. Run the diff processor to detect breaking changes.
	4. Run the missing test detector to detect missing tests for fields changed.
	5. Report the results in a PR comment.
	6. Run unit tests for the missing test detector.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		buildID := os.Getenv("BUILD_ID")
		fmt.Println("Build ID: ", buildID)

		projectID := os.Getenv("PROJECT_ID")
		fmt.Println("Project ID: ", projectID)

		buildStep := os.Getenv("BUILD_STEP")
		fmt.Println("Build Step: ", buildStep)

		commit := os.Getenv("COMMIT_SHA")
		fmt.Println("Commit SHA: ", commit)

		pr := os.Getenv("PR_NUMBER")
		fmt.Println("PR Number: ", pr)

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			fmt.Println("Did not provide GITHUB_TOKEN environment variable")
			os.Exit(1)
		}

		gh := github.NewGithubService()
		execGenerateComment(buildID, projectID, buildStep, commit, pr, githubToken, gh, &actualRunner{})
	},
}

func execGenerateComment(buildID, projectID, buildStep, commit, pr, githubToken string, gh gcGithub, r runner) {
	newBranch := "auto-pr-" + pr
	oldBranch := "auto-pr-" + pr + "-old"
	wd, err := r.Getwd()
	if err != nil {
		fmt.Println("Failed to get current working directory: ", err)
		os.Exit(1)
	}
	mmLocalPath := filepath.Join(wd, "..", "..")
	tpgRepoName := "terraform-provider-google"
	tpgLocalPath := filepath.Join(mmLocalPath, "..", "tpg")
	tpgbRepoName := "terraform-provider-google-beta"
	tpgbLocalPath := filepath.Join(mmLocalPath, "..", "tpgb")
	tfoicsRepoName := "docs-examples"
	tfoicsLocalPath := filepath.Join(mmLocalPath, "..", "tfoics")
	// For backwards compatibility until at least Nov 15 2021
	tfcRepoName := "terraform-google-conversion"
	tfcLocalPath := filepath.Join(mmLocalPath, "..", "tfc")

	// TPG/TPGB difference
	diffs, err := cloneAndDiff(tpgRepoName, tpgLocalPath, oldBranch, newBranch, "Terraform GA", githubToken, r)
	if err != nil {
		fmt.Printf("Error cloning and diffing tpg repo: %v\n", err)
		os.Exit(1)
	}

	tpgbDiffs, err := cloneAndDiff(tpgbRepoName, tpgbLocalPath, oldBranch, newBranch, "Terraform Beta", githubToken, r)
	if err != nil {
		fmt.Printf("Error cloning and diffing tpgb repo: %v\n", err)
		os.Exit(1)
	}
	diffs += tpgbDiffs

	// Breaking change setup and execution
	diffProcessorPath := filepath.Join(mmLocalPath, "tools", "diff-processor")
	for _, path := range []string{"old", "new"} {
		if err := r.Copy(tpgLocalPath, filepath.Join(diffProcessorPath, path)); err != nil {
			fmt.Printf("Error copying files: %v\n", err)
			os.Exit(1)
		}
	}
	var tpgBreaking, tpgbBreaking, breakingChanges string
	var diffProccessorErr error
	if stdout, stderr, err := r.Run(diffProcessorPath, "make", []string{"build"}, []string{"OLD_REF=" + oldBranch, "NEW_REF=" + newBranch}); err != nil {
		fmt.Printf("Error running make build in %s: %v\nstdout:\n%s\nstderr:\n%s\n", diffProcessorPath, err, stdout, stderr)
		diffProccessorErr = err
	} else {
		tpgBreaking, stderr, err = r.Run(diffProcessorPath, "bin/diff-processor", []string{"breaking-changes"}, nil)
		if err != nil {
			fmt.Printf("Error running diff-processor: %v\nstdout:\n%s\nstderr:\n%s\n", err, tpgBreaking, stderr)
			diffProccessorErr = err
		}
	}
	for _, path := range []string{"old", "new", "bin"} {
		if err := r.RemoveAll(filepath.Join(diffProcessorPath, path)); err != nil {
			fmt.Printf("error removing files: %v\n", err)
			os.Exit(1)
		}
	}
	for _, path := range []string{"old", "new"} {
		if err := r.Copy(tpgbLocalPath, filepath.Join(diffProcessorPath, path)); err != nil {
			fmt.Printf("Error copying files: %v\n", err)
			os.Exit(1)
		}
	}

	if diffProccessorErr != nil {
		fmt.Println("Breaking changes failed")
		breakingChanges = `## Breaking Change Detection Failed
The breaking change detector crashed during execution. This is usually due to the downstream provider(s) failing to compile. Please investigate or follow up with your reviewer.`
	} else {
		fmt.Println("Breaking changes succeeded")
		breakingChanges = compareBreakingChanges(tpgBreaking, tpgbBreaking)
	}

	tpgbLocalPathOld := tpgbLocalPath + "old"
	if err := r.Copy(tpgbLocalPath, tpgbLocalPathOld); err != nil {
		fmt.Printf("Error copying files from %s to %s: %v\n", tpgbLocalPath, tpgbLocalPathOld, err)
		os.Exit(1)
	}
	if stdout, stderr, err := r.Run(tpgbLocalPathOld, "git", []string{"checkout", "origin/" + oldBranch}, nil); err != nil {
		fmt.Printf("Error checking out %s in %s: %v\nstdout:\n%s\nstderr:\n%s\n", oldBranch, tpgbLocalPathOld, err, stdout, stderr)
		os.Exit(1)
	}

	if err := updatePackageName("old", tpgbLocalPathOld, r); err != nil {
		fmt.Printf("Error updating package name in %s: %v\n", tpgbLocalPathOld, err)
		os.Exit(1)
	}
	if err := updatePackageName("new", tpgbLocalPath, r); err != nil {
		fmt.Printf("Error updating package name in %s: %v\n", tpgbLocalPath, err)
		os.Exit(1)
	}

	missingTestDetectorPath := filepath.Join(mmLocalPath, "tools", "missing-test-detector")
	if stdout, stderr, err := r.Run(missingTestDetectorPath, "go", []string{"mod", "edit", "-replace", fmt.Sprintf("google/provider/%s=%s", "new", tpgbLocalPath)}, nil); err != nil {
		fmt.Printf("Error running go mod edit: %v\nstdout:\n%s\nstderr:\n%s\n", err, stdout, stderr)
	}
	if stdout, stderr, err := r.Run(missingTestDetectorPath, "go", []string{"mod", "edit", "-replace", fmt.Sprintf("google/provider/%s=%s", "old", tpgbLocalPathOld)}, nil); err != nil {
		fmt.Printf("Error running go mod edit: %v\nstdout:\n%s\nstderr:\n%s\n", err, stdout, stderr)
	}
	if stdout, stderr, err := r.Run(missingTestDetectorPath, "go", []string{"mod", "tidy"}, nil); err != nil {
		fmt.Printf("Error running go mod tidy: %v\nstdout:\n%s\nstderr:\n%s\n", err, stdout, stderr)
	}
	missingTests, stderr, err := r.Run(missingTestDetectorPath, "go", []string{"run", ".", fmt.Sprintf("-services-dir=%s/google-beta/services", tpgbLocalPath)}, nil)
	if err != nil {
		fmt.Printf("Error running missing test detector: %v\nstdout:\n%s\nstderr:\n%s\n", err, missingTests, stderr)
		missingTests = ""
	} else {
		fmt.Printf("Successfully ran missing test detector:\n%s\n", missingTests)
	}

	// TF Conversion - for compatibility until at least Nov 15 2021
	// allow this to fail for compatibility during tfv/tgc transition phase
	tfcDiffs, err := cloneAndDiff(tfcRepoName, tfcLocalPath, oldBranch, newBranch, "TF Conversion", githubToken, r)
	if err != nil {
		fmt.Printf("Error getting tfc diffs: %v\n", err)
	}
	diffs += tfcDiffs

	tfoicsDiffs, err := cloneAndDiff(tfoicsRepoName, tfoicsLocalPath, oldBranch, newBranch, "TF OiCS", githubToken, r)
	if err != nil {
		fmt.Printf("Error getting tf oics diffs: %v\n", err)
		os.Exit(1)
	}
	diffs += tfoicsDiffs

	message := "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n"
	breakingState := "success"
	if breakingChanges != "" {
		message += breakingChanges + "\n\n"

		labels, err := gh.GetPullRequestLabelIDs(pr)
		if err != nil {
			fmt.Printf("Error getting pull request labels: %v\n", err)
			os.Exit(1)
		}
		if _, ok := labels[4598495472]; !ok {
			breakingState = "failure"
		}
	}

	if diffs == "" {
		message += "## Diff report\nYour PR hasn't generated any diffs, but I'll let you know if a future commit does."
	} else {
		message += "## Diff report\nYour PR generated some diffs in downstreams - here they are.\n" + diffs
		if missingTests != "" {
			message += "\n" + missingTests + "\n"
		}
	}

	if err := gh.PostComment(pr, message); err != nil {
		fmt.Printf("Error posting comment to PR %s: %v\n", pr, err)
	}

	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(pr, "terraform-provider-breaking-change-test", breakingState, targetURL, commit); err != nil {
		fmt.Printf("Error posting build status for pr %s commit %s: %v\n", pr, commit, err)
		os.Exit(1)
	}

	if diffs, stderr, err := r.Run(mmLocalPath, "git", []string{"diff", "HEAD", "origin/main", "tools/missing-test-detector"}, nil); err != nil {
		fmt.Printf("Error diffing with origin/main: %v\nstdout:\n%s\nstderr:\n%s\n", err, diffs, stderr)
		os.Exit(1)
	} else if diffs != "" {
		fmt.Printf("Found diffs in missing test detector:\n%s\nRunning tests.\n", diffs)
		if err := testTools(mmLocalPath, tpgbLocalPath, pr, commit, buildID, buildStep, projectID, gh, r); err != nil {
			fmt.Printf("Error testing tools in %s: %v\n", mmLocalPath, err)
			os.Exit(1)
		}
	}
}

func cloneAndDiff(repoName, path, oldBranch, newBranch, diffTitle, githubToken string, r runner) (string, error) {
	url := fmt.Sprintf("https://modular-magician:%s@github.com/modular-magician/%s", githubToken, repoName)
	if stdout, stderr, err := r.Run("", "git", []string{"clone", "-b", newBranch, url, path}, nil); err != nil {
		return "", fmt.Errorf("error cloning %s: %v\nstdout:\n%s\nstderr:\n%s", repoName, err, stdout, stderr)
	}
	if stdout, stderr, err := r.Run(path, "git", []string{"fetch", "origin", oldBranch}, nil); err != nil {
		return "", fmt.Errorf("error fetching branch %s in repo %s: %v\nstdout:\n%s\nstderr:\n%s", oldBranch, repoName, err, stdout, stderr)
	}

	if summary, stderr, err := r.Run(path, "git", []string{"diff", "origin/" + oldBranch, "origin/" + newBranch, "--shortstat"}, nil); err != nil {
		return "", fmt.Errorf("error diffing %s and %s: %v\nstdout:\n%s\nstderr:\n%s", oldBranch, newBranch, err, summary, stderr)
	} else if summary != "" {
		return fmt.Sprintf("\n%s: [Diff](https://github.com/modular-magician/%s/compare/%s..%s) (%s)", diffTitle, repoName, oldBranch, newBranch, strings.TrimSuffix(summary, "\n")), nil
	}
	return "", nil
}

// Get the breaking change message including the unique tpg messages and all tpgb messages.
func compareBreakingChanges(tpgBreaking, tpgbBreaking string) string {
	var allMessages []string
	if tpgBreaking == "" {
		if tpgbBreaking == "" {
			return ""
		}
		allMessages = strings.Split(tpgbBreaking, "\n")
	} else if tpgbBreaking == "" {
		allMessages = strings.Split(tpgBreaking, "\n")
	} else {
		dashExp := regexp.MustCompile("-.*")
		tpgMessages := strings.Split(tpgBreaking, "\n")
		tpgbMessages := strings.Split(tpgbBreaking, "\n")
		tpgbSet := make(map[string]struct{}, len(tpgbMessages))
		var tpgUnique []string
		for _, message := range tpgbMessages {
			simple := dashExp.ReplaceAllString(message, "")
			tpgbSet[simple] = struct{}{}
		}
		for _, message := range tpgMessages {
			simple := dashExp.ReplaceAllString(message, "")
			if _, ok := tpgbSet[simple]; !ok {
				tpgUnique = append(tpgUnique, message)
			}
		}
		allMessages = append(tpgUnique, tpgbMessages...)
	}
	if len(allMessages) > 0 {
		return `Breaking Change(s) Detected
The following breaking change(s) were detected within your pull request.

* ` + strings.Join(allMessages, "\n* ") + `

If you believe this detection to be incorrect please raise the concern with your reviewer.
If you intend to make this change you will need to wait for a [major release](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#example-major-number-increments) window.
An ` + "`override-breaking-change`" + `label can be added to allow merging.
`
	}
	return ""
}

// Update the provider package name to the given name in the given path.
// name should be either "old" or "new".
func updatePackageName(name, path string, r runner) error {
	oldPackageName := "github.com/hashicorp/terraform-provider-google-beta"
	newPackageName := "google/provider/" + name
	fmt.Printf("Updating package name in %s from %s to %s\n", path, oldPackageName, newPackageName)
	if stdout, stderr, err := r.Run(path, "find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", fmt.Sprintf("s~%s~%s~g", oldPackageName, newPackageName), "{}", "+"}, nil); err != nil {
		return fmt.Errorf("error running find: %v\nstdout:\n%s\nstderr:\n%s\n", err, stdout, stderr)
	}
	if stdout, stderr, err := r.Run(path, "sed", []string{"-i.bak", fmt.Sprintf("s|%s|%s|g", oldPackageName, newPackageName), "go.mod"}, nil); err != nil {
		return fmt.Errorf("error running sed: %v\nstdout:\n%s\nstderr:\n%s\n", err, stdout, stderr)
	}
	if stdout, stderr, err := r.Run(path, "sed", []string{"-i.bak", fmt.Sprintf("s|%s|%s|g", oldPackageName, newPackageName), "go.sum"}, nil); err != nil {
		return fmt.Errorf("error running sed: %v\nstdout:\n%s\nstderr:\n%s\n", err, stdout, stderr)
	}
	return nil
}

// Run unit tests for the missing test detector and diff processor.
// Report results using Github API.
func testTools(mmLocalPath, tpgbLocalPath, pr, commit, buildID, buildStep, projectID string, gh gcGithub, r runner) error {
	missingTestDetectorPath := filepath.Join(mmLocalPath, "tools", "missing-test-detector")
	if stdout, stderr, err := r.Run(missingTestDetectorPath, "go", []string{"mod", "tidy"}, nil); err != nil {
		fmt.Printf("error running go mod tidy in %s: %v\nstdout:\n%s\nstderr:\n%s\n", missingTestDetectorPath, err, stdout, stderr)
	}
	servicesDir := filepath.Join(tpgbLocalPath, "google-beta", "services")
	state := "success"
	if stdout, stderr, err := r.Run(missingTestDetectorPath, "go", []string{"test"}, []string{"SERVICES_DIR=" + servicesDir}); err != nil {
		fmt.Printf("error from running go test in %s: %v\nstdout:\n%s\nstderr:\n%s\n", missingTestDetectorPath, err, stdout, stderr)
		state = "failure"
	}
	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	return gh.PostBuildStatus(pr, "unit-tests-missing-test-detector", state, targetURL, commit)
}

func init() {
	rootCmd.AddCommand(generateCommentCmd)
}

package cmd

import (
	"fmt"
	"magician/exec"
	"magician/github"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const allowBreakingChangesLabel = 4598495472

type gcGithub interface {
	GetPullRequestLabelIDs(prNumber string) (map[int]struct{}, error)
	PostBuildStatus(prNumber, title, state, targetURL, commitSha string) error
	PostComment(prNumber, comment string) error
}

type gcRunner interface {
	Getwd() (string, error)
	Copy(src, dest string) error
	RemoveAll(path string) error
	Chdir(path string)
	Run(name string, args, env []string) (string, error)
	MustRun(name string, args, env []string) string
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
		execGenerateComment(buildID, projectID, buildStep, commit, pr, githubToken, gh, exec.NewRunner())
	},
}

func execGenerateComment(buildID, projectID, buildStep, commit, pr, githubToken string, gh gcGithub, r gcRunner) {
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

	var diffs string
	for _, repo := range []struct {
		name    string
		title   string
		path    string
		canFail bool
	}{
		{
			name:  tpgRepoName,
			title: "Terraform GA",
			path:  tpgLocalPath,
		},
		{
			name:  tpgbRepoName,
			title: "Terraform Beta",
			path:  tpgbLocalPath,
		},
		{
			name:    tfcRepoName,
			title:   "TF Conversion",
			path:    tfcLocalPath,
			canFail: true,
		},
		{
			name:  tfoicsRepoName,
			title: "TF OiCS",
			path:  tfoicsLocalPath,
		},
	} {
		// TPG/TPGB difference
		repoDiffs, err := cloneAndDiff(repo.name, repo.path, oldBranch, newBranch, repo.title, githubToken, r)
		if err != nil {
			fmt.Printf("Error cloning and diffing tpg repo: %v\n", err)
			if !repo.canFail {
				os.Exit(1)
			}
		}
		diffs += repoDiffs
	}

	breakingChanges, err := detectBreakingChanges(mmLocalPath, tpgLocalPath, tpgbLocalPath, oldBranch, newBranch, r)
	if err != nil {
		fmt.Println("Error setting up breaking change detector: ", err)
		os.Exit(1)
	}

	missingTests, err := detectMissingTests(mmLocalPath, tpgbLocalPath, oldBranch, r)
	if err != nil {
		fmt.Println("Error setting up missing test detector: ", err)
		os.Exit(1)
	}

	message := "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n"
	breakingState := "success"
	if breakingChanges != "" {
		message += breakingChanges + "\n\n"

		labels, err := gh.GetPullRequestLabelIDs(pr)
		if err != nil {
			fmt.Printf("Error getting pull request labels: %v\n", err)
			os.Exit(1)
		}
		if _, ok := labels[allowBreakingChangesLabel]; !ok {
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

	r.Chdir(mmLocalPath)
	if diffs := r.MustRun("git", []string{"diff", "HEAD", "origin/main", "tools/missing-test-detector"}, nil); diffs != "" {
		fmt.Printf("Found diffs in missing test detector:\n%s\nRunning tests.\n", diffs)
		if err := testTools(mmLocalPath, tpgbLocalPath, pr, commit, buildID, buildStep, projectID, gh, r); err != nil {
			fmt.Printf("Error testing tools in %s: %v\n", mmLocalPath, err)
			os.Exit(1)
		}
	}
}

func cloneAndDiff(repoName, path, oldBranch, newBranch, diffTitle, githubToken string, r gcRunner) (string, error) {
	url := fmt.Sprintf("https://modular-magician:%s@github.com/modular-magician/%s", githubToken, repoName)
	if _, err := r.Run("git", []string{"clone", "-b", newBranch, url, path}, nil); err != nil {
		return "", fmt.Errorf("error cloning %s: %v\n", repoName, err)
	}
	r.Chdir(path)
	if _, err := r.Run("git", []string{"fetch", "origin", oldBranch}, nil); err != nil {
		return "", fmt.Errorf("error fetching branch %s in repo %s: %v\n", oldBranch, repoName, err)
	}

	if summary, err := r.Run("git", []string{"diff", "origin/" + oldBranch, "origin/" + newBranch, "--shortstat"}, nil); err != nil {
		return "", fmt.Errorf("error diffing %s and %s: %v\n", oldBranch, newBranch, err)
	} else if summary != "" {
		return fmt.Sprintf("\n%s: [Diff](https://github.com/modular-magician/%s/compare/%s..%s) (%s)", diffTitle, repoName, oldBranch, newBranch, strings.TrimSuffix(summary, "\n")), nil
	}
	return "", nil
}

// Run the breaking change detector and return the results.
// Returns an empty string unless there are breaking changes or the detector failed.
// Error will be nil unless an error occurs manipulating files.
func detectBreakingChanges(mmLocalPath, tpgLocalPath, tpgbLocalPath, oldBranch, newBranch string, r gcRunner) (string, error) {
	// Breaking change setup and execution
	diffProcessorPath := filepath.Join(mmLocalPath, "tools", "diff-processor")
	for _, path := range []string{"old", "new"} {
		if err := r.Copy(tpgLocalPath, filepath.Join(diffProcessorPath, path)); err != nil {
			return "", err
		}
	}
	var tpgBreaking, tpgbBreaking, breakingChanges string
	var diffProccessorErr error
	r.Chdir(diffProcessorPath)
	if _, err := r.Run("make", []string{"build"}, []string{"OLD_REF=" + oldBranch, "NEW_REF=" + newBranch}); err != nil {
		fmt.Printf("Error running make build in %s: %v\n", diffProcessorPath, err)
		diffProccessorErr = err
	} else {
		tpgBreaking, err = r.Run("bin/diff-processor", []string{"breaking-changes"}, nil)
		if err != nil {
			fmt.Println("Diff processor error: ", err)
			diffProccessorErr = err
		}
	}
	for _, path := range []string{"old", "new", "bin"} {
		if err := r.RemoveAll(filepath.Join(diffProcessorPath, path)); err != nil {
			return "", err
		}
	}
	for _, path := range []string{"old", "new"} {
		if err := r.Copy(tpgbLocalPath, filepath.Join(diffProcessorPath, path)); err != nil {
			return "", err
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
	return breakingChanges, nil
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

// Run the missing test detector and return the results.
// Returns an empty string unless there are missing tests.
// Error will be nil unless an error occurs during setup.
func detectMissingTests(mmLocalPath, tpgbLocalPath, oldBranch string, r gcRunner) (string, error) {
	tpgbLocalPathOld := tpgbLocalPath + "old"
	if err := r.Copy(tpgbLocalPath, tpgbLocalPathOld); err != nil {
		return "", err
	}
	oldDir, err := r.Getwd()
	if err != nil {
		return "", err
	}
	defer r.Chdir(oldDir)
	r.Chdir(tpgbLocalPathOld)
	if _, err := r.Run("git", []string{"checkout", "origin/" + oldBranch}, nil); err != nil {
		return "", err
	}

	if err := updatePackageName("old", tpgbLocalPathOld, r); err != nil {
		return "", err
	}
	if err := updatePackageName("new", tpgbLocalPath, r); err != nil {
		return "", err
	}

	missingTestDetectorPath := filepath.Join(mmLocalPath, "tools", "missing-test-detector")
	r.Chdir(missingTestDetectorPath)
	if _, err := r.Run("go", []string{"mod", "edit", "-replace", fmt.Sprintf("google/provider/%s=%s", "new", tpgbLocalPath)}, nil); err != nil {
		fmt.Printf("Error running go mod edit: %v\n", err)
	}
	if _, err := r.Run("go", []string{"mod", "edit", "-replace", fmt.Sprintf("google/provider/%s=%s", "old", tpgbLocalPathOld)}, nil); err != nil {
		fmt.Printf("Error running go mod edit: %v\n", err)
	}
	if _, err := r.Run("go", []string{"mod", "tidy"}, nil); err != nil {
		fmt.Printf("Error running go mod tidy: %v\n", err)
	}
	missingTests, err := r.Run("go", []string{"run", ".", fmt.Sprintf("-services-dir=%s/google-beta/services", tpgbLocalPath)}, nil)
	if err != nil {
		fmt.Printf("Error running missing test detector: %v\n", err)
		missingTests = ""
	} else {
		fmt.Printf("Successfully ran missing test detector:\n%s\n", missingTests)
	}
	return missingTests, nil
}

// Update the provider package name to the given name in the given path.
// name should be either "old" or "new".
func updatePackageName(name, path string, r gcRunner) error {
	oldPackageName := "github.com/hashicorp/terraform-provider-google-beta"
	newPackageName := "google/provider/" + name
	fmt.Printf("Updating package name in %s from %s to %s\n", path, oldPackageName, newPackageName)
	oldDir, err := r.Getwd()
	if err != nil {
		return err
	}
	defer r.Chdir(oldDir)
	r.Chdir(path)
	if _, err := r.Run("find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", fmt.Sprintf("s~%s~%s~g", oldPackageName, newPackageName), "{}", "+"}, nil); err != nil {
		return fmt.Errorf("error running find: %v\n", err)
	}
	if _, err := r.Run("sed", []string{"-i.bak", fmt.Sprintf("s|%s|%s|g", oldPackageName, newPackageName), "go.mod"}, nil); err != nil {
		return fmt.Errorf("error running sed: %v\n", err)
	}
	if _, err := r.Run("sed", []string{"-i.bak", fmt.Sprintf("s|%s|%s|g", oldPackageName, newPackageName), "go.sum"}, nil); err != nil {
		return fmt.Errorf("error running sed: %v\n", err)
	}
	return nil
}

// Run unit tests for the missing test detector and diff processor.
// Report results using Github API.
func testTools(mmLocalPath, tpgbLocalPath, pr, commit, buildID, buildStep, projectID string, gh gcGithub, r gcRunner) error {
	missingTestDetectorPath := filepath.Join(mmLocalPath, "tools", "missing-test-detector")
	oldDir, err := r.Getwd()
	if err != nil {
		return err
	}
	defer r.Chdir(oldDir)
	r.Chdir(missingTestDetectorPath)
	if _, err := r.Run("go", []string{"mod", "tidy"}, nil); err != nil {
		fmt.Printf("error running go mod tidy in %s: %v\n", missingTestDetectorPath, err)
	}
	servicesDir := filepath.Join(tpgbLocalPath, "google-beta", "services")
	state := "success"
	if _, err := r.Run("go", []string{"test"}, []string{"SERVICES_DIR=" + servicesDir}); err != nil {
		fmt.Printf("error from running go test in %s: %v\n", missingTestDetectorPath, err)
		state = "failure"
	}
	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	return gh.PostBuildStatus(pr, "unit-tests-missing-test-detector", state, targetURL, commit)
}

func init() {
	rootCmd.AddCommand(generateCommentCmd)
}

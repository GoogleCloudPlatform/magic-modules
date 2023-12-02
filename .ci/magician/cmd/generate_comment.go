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

const allowBreakingChangesLabel = "override-breaking-change"

type ProviderVersion string

type Repository struct {
	Name        string // Name in GitHub (e.g. magic-modules)
	Title       string // Title for display (e.g. Magic Modules)
	Path        string // local Path once checked out, including Name
	Version     ProviderVersion
	DiffCanFail bool // whether to allow the command to continue if cloning or diffing the repo fails
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

		prNumber := os.Getenv("PR_NUMBER")
		fmt.Println("PR Number: ", prNumber)

		githubToken, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			fmt.Println("Did not provide GITHUB_TOKEN environment variable")
			os.Exit(1)
		}

		gh := github.NewClient()
		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating a runner: ", err)
			os.Exit(1)
		}
		execGenerateComment(buildID, projectID, buildStep, commit, prNumber, githubToken, gh, rnr)
	},
}

func execGenerateComment(buildID, projectID, buildStep, commit, prNumber, githubToken string, gh GithubClient, r ExecRunner) {
	newBranch := "auto-pr-" + prNumber
	oldBranch := "auto-pr-" + prNumber + "-old"
	wd := r.GetCWD()
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
	for _, repo := range []Repository{
		{
			Name:  tpgRepoName,
			Title: "Terraform GA",
			Path:  tpgLocalPath,
		},
		{
			Name:  tpgbRepoName,
			Title: "Terraform Beta",
			Path:  tpgbLocalPath,
		},
		{
			Name:        tfcRepoName,
			Title:       "TF Conversion",
			Path:        tfcLocalPath,
			DiffCanFail: true,
		},
		{
			Name:  tfoicsRepoName,
			Title: "TF OiCS",
			Path:  tfoicsLocalPath,
		},
	} {
		// TPG/TPGB difference
		repoDiffs, err := cloneAndDiff(repo, oldBranch, newBranch, githubToken, r)
		if err != nil {
			fmt.Printf("Error cloning and diffing tpg repo: %v\n", err)
			if !repo.DiffCanFail {
				os.Exit(1)
			}
		}
		diffs += "\n" + repoDiffs
	}

	var showBreakingChangesFailed bool
	var err error
	diffProcessorPath := filepath.Join(mmLocalPath, "tools", "diff-processor")
	// versionedBreakingChanges is a map of breaking change output by provider version.
	versionedBreakingChanges := make(map[ProviderVersion]string, 2)

	for _, repo := range []Repository{
		{
			Title:   "TPG",
			Path:    tpgLocalPath,
			Version: "ga",
		},
		{
			Title:   "TPGB",
			Path:    tpgbLocalPath,
			Version: "beta",
		},
	} {
		// TPG diff processor
		err = buildDiffProcessor(diffProcessorPath, repo.Path, oldBranch, newBranch, r)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output, err := computeBreakingChanges(diffProcessorPath, r)
		if err != nil {
			fmt.Println("Error computing TPG breaking changes: ", err)
			showBreakingChangesFailed = true
		}
		versionedBreakingChanges[repo.Version] = strings.TrimSuffix(output, "\n")
		err = addLabels(diffProcessorPath, githubToken, prNumber, r)
		if err != nil {
			fmt.Println("Error adding TPG labels to PR: ", err)
		}
		err = cleanDiffProcessor(diffProcessorPath, r)
		if err != nil {
			fmt.Println("Error cleaning up diff processor: ", err)
			os.Exit(1)
		}
	}

	var breakingChanges string
	if showBreakingChangesFailed {
		breakingChanges = `## Breaking Change Detection Failed
The breaking change detector crashed during execution. This is usually due to the downstream provider(s) failing to compile. Please investigate or follow up with your reviewer.`
	} else {
		breakingChanges = combineBreakingChanges(versionedBreakingChanges["ga"], versionedBreakingChanges["beta"])
	}

	// Missing test detector
	missingTests, err := detectMissingTests(mmLocalPath, tpgbLocalPath, oldBranch, r)
	if err != nil {
		fmt.Println("Error setting up missing test detector: ", err)
		os.Exit(1)
	}

	message := "Hi there, I'm the Modular magician. I've detected the following information about your changes:\n\n"
	breakingState := "success"
	if breakingChanges != "" {
		message += breakingChanges + "\n\n"

		pullRequest, err := gh.GetPullRequest(prNumber)
		if err != nil {
			fmt.Printf("Error getting pull request: %v\n", err)
			os.Exit(1)
		}

		breakingChangesAllowed := false
		for _, label := range pullRequest.Labels {
			if label.Name == allowBreakingChangesLabel {
				breakingChangesAllowed = true
				break
			}
		}
		if !breakingChangesAllowed {
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

	if err := gh.PostComment(prNumber, message); err != nil {
		fmt.Printf("Error posting comment to PR %s: %v\n", prNumber, err)
	}

	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(prNumber, "terraform-provider-breaking-change-test", breakingState, targetURL, commit); err != nil {
		fmt.Printf("Error posting build status for pr %s commit %s: %v\n", prNumber, commit, err)
		os.Exit(1)
	}

	if err := r.PushDir(mmLocalPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if diffs := r.MustRun("git", []string{"diff", "HEAD", "origin/main", "tools/missing-test-detector"}, nil); diffs != "" {
		fmt.Printf("Found diffs in missing test detector:\n%s\nRunning tests.\n", diffs)
		if err := testTools(mmLocalPath, tpgbLocalPath, prNumber, commit, buildID, buildStep, projectID, gh, r); err != nil {
			fmt.Printf("Error testing tools in %s: %v\n", mmLocalPath, err)
			os.Exit(1)
		}
	}
	if err := r.PopDir(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cloneAndDiff(repo Repository, oldBranch, newBranch, githubToken string, r ExecRunner) (string, error) {
	// Clone the repo to the desired repo.Path.
	url := fmt.Sprintf("https://modular-magician:%s@github.com/modular-magician/%s", githubToken, repo.Name)
	if _, err := r.Run("git", []string{"clone", "-b", newBranch, url, repo.Path}, nil); err != nil {
		return "", fmt.Errorf("error cloning %s: %v\n", repo.Name, err)
	}

	// Push dir to the newly cloned repo.
	if err := r.PushDir(repo.Path); err != nil {
		return "", err
	}
	if _, err := r.Run("git", []string{"fetch", "origin", oldBranch}, nil); err != nil {
		return "", fmt.Errorf("error fetching branch %s in repo %s: %v\n", oldBranch, repo.Name, err)
	}

	// Return summary, if any, and return to original directory.
	if summary, err := r.Run("git", []string{"diff", "origin/" + oldBranch, "origin/" + newBranch, "--shortstat"}, nil); err != nil {
		return "", fmt.Errorf("error diffing %s and %s: %v\n", oldBranch, newBranch, err)
	} else if summary != "" {
		summary = strings.TrimSuffix(summary, "\n")
		return fmt.Sprintf("%s: [Diff](https://github.com/modular-magician/%s/compare/%s..%s) (%s)", repo.Title, repo.Name, oldBranch, newBranch, summary), r.PopDir()
	}
	return "", r.PopDir()
}

// Build the diff processor for tpg or tpgb
func buildDiffProcessor(diffProcessorPath, providerLocalPath, oldBranch, newBranch string, r ExecRunner) error {
	if err := r.PushDir(diffProcessorPath); err != nil {
		return err
	}
	for _, path := range []string{"old", "new"} {
		if err := r.Copy(providerLocalPath, filepath.Join(diffProcessorPath, path)); err != nil {
			return err
		}
	}
	if _, err := r.Run("make", []string{"build"}, []string{"OLD_REF=" + oldBranch, "NEW_REF=" + newBranch}); err != nil {
		return fmt.Errorf("Error running make build in %s: %v\n", diffProcessorPath, err)
	}
	return r.PopDir()
}

func computeBreakingChanges(diffProcessorPath string, r ExecRunner) (string, error) {
	if err := r.PushDir(diffProcessorPath); err != nil {
		return "", err
	}
	breakingChanges, err := r.Run("bin/diff-processor", []string{"breaking-changes"}, nil)
	if err != nil {
		return "", err
	}
	return breakingChanges, r.PopDir()
}

func addLabels(diffProcessorPath, githubToken, prNumber string, r ExecRunner) error {
	if err := r.PushDir(diffProcessorPath); err != nil {
		return err
	}
	output, err := r.Run("bin/diff-processor", []string{"add-labels", prNumber}, []string{fmt.Sprintf("GITHUB_TOKEN=%s", githubToken)})
	fmt.Println(output)
	if err != nil {
		return err
	}
	return r.PopDir()
}

func cleanDiffProcessor(diffProcessorPath string, r ExecRunner) error {
	for _, path := range []string{"old", "new", "bin"} {
		if err := r.RemoveAll(filepath.Join(diffProcessorPath, path)); err != nil {
			return err
		}
	}
	return nil
}

// Get the breaking change message including the unique tpg messages and all tpgb messages.
func combineBreakingChanges(tpgBreaking, tpgbBreaking string) string {
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
		return `## Breaking Change(s) Detected
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
func detectMissingTests(mmLocalPath, tpgbLocalPath, oldBranch string, r ExecRunner) (string, error) {
	tpgbLocalPathOld := tpgbLocalPath + "old"

	if err := r.Copy(tpgbLocalPath, tpgbLocalPathOld); err != nil {
		return "", err
	}

	if err := r.PushDir(tpgbLocalPathOld); err != nil {
		return "", err
	}
	if _, err := r.Run("git", []string{"checkout", "origin/" + oldBranch}, nil); err != nil {
		return "", err
	}

	if err := updatePackageName("old", tpgbLocalPathOld, r); err != nil {
		return "", err
	}
	if err := updatePackageName("new", tpgbLocalPath, r); err != nil {
		return "", err
	}
	if err := r.PopDir(); err != nil {
		return "", err
	}

	missingTestDetectorPath := filepath.Join(mmLocalPath, "tools", "missing-test-detector")
	if err := r.PushDir(missingTestDetectorPath); err != nil {
		return "", err
	}
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
	return missingTests, r.PopDir()
}

// Update the provider package name to the given name in the given path.
// name should be either "old" or "new".
func updatePackageName(name, path string, r ExecRunner) error {
	oldPackageName := "github.com/hashicorp/terraform-provider-google-beta"
	newPackageName := "google/provider/" + name
	fmt.Printf("Updating package name in %s from %s to %s\n", path, oldPackageName, newPackageName)
	if err := r.PushDir(path); err != nil {
		return err
	}
	if _, err := r.Run("find", []string{".", "-type", "f", "-name", "*.go", "-exec", "sed", "-i.bak", fmt.Sprintf("s~%s~%s~g", oldPackageName, newPackageName), "{}", "+"}, nil); err != nil {
		return fmt.Errorf("error running find: %v\n", err)
	}
	if _, err := r.Run("sed", []string{"-i.bak", fmt.Sprintf("s|%s|%s|g", oldPackageName, newPackageName), "go.mod"}, nil); err != nil {
		return fmt.Errorf("error running sed: %v\n", err)
	}
	if _, err := r.Run("sed", []string{"-i.bak", fmt.Sprintf("s|%s|%s|g", oldPackageName, newPackageName), "go.sum"}, nil); err != nil {
		return fmt.Errorf("error running sed: %v\n", err)
	}
	return r.PopDir()
}

// Run unit tests for the missing test detector and diff processor.
// Report results using Github API.
func testTools(mmLocalPath, tpgbLocalPath, prNumber, commit, buildID, buildStep, projectID string, gh GithubClient, r ExecRunner) error {
	missingTestDetectorPath := filepath.Join(mmLocalPath, "tools", "missing-test-detector")
	r.PushDir(missingTestDetectorPath)
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
	if err := gh.PostBuildStatus(prNumber, "unit-tests-missing-test-detector", state, targetURL, commit); err != nil {
		return err
	}
	return r.PopDir()
}

func init() {
	rootCmd.AddCommand(generateCommentCmd)
}

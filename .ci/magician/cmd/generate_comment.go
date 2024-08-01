/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"magician/exec"
	"magician/github"
	"magician/provider"
	"magician/source"

	"github.com/GoogleCloudPlatform/magic-modules/tools/issue-labeler/labeler"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	_ "embed"
)

var (
	//go:embed DIFF_COMMENT.md
	diffComment string
)

type Diff struct {
	Title     string
	Repo      string
	ShortStat string
}

type BreakingChange struct {
	Message                string
	DocumentationReference string
}

type MissingTestInfo struct {
	SuggestedTest string
	Tests         []string
}

type Errors struct {
	Title  string
	Errors []string
}

type diffCommentData struct {
	PrNumber        int
	Diffs           []Diff
	BreakingChanges []BreakingChange
	MissingTests    map[string]*MissingTestInfo
	Errors          []Errors
}

const allowBreakingChangesLabel = "override-breaking-change"

var gcEnvironmentVariables = [...]string{
	"BUILD_ID",
	"BUILD_STEP",
	"COMMIT_SHA",
	"GOPATH",
	"HOME",
	"PATH",
	"PR_NUMBER",
	"PROJECT_ID",
}

var generateCommentCmd = &cobra.Command{
	Use:   "generate-comment",
	Short: "Run presubmit generate comment",
	Long: `This command processes pull requests and performs various validations and actions based on the PR's metadata and author.

	The following PR details are expected as environment variables:
` + listGCEnvironmentVariables() + `

	The command performs the following steps:
	1. Clone the tpg, tpgb, tfc, and tfoics repos from modular-magician.
	2. Compute the diffs between auto-pr-# and auto-pr-#-old branches.
	3. Run the diff processor to detect breaking changes.
	4. Run the missing test detector to detect missing tests for fields changed.
	5. Report the results in a PR comment.
	6. Run unit tests for the missing test detector.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string, len(gcEnvironmentVariables))
		for _, ev := range gcEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}

		for _, tokenName := range []string{"GITHUB_TOKEN_DOWNSTREAMS", "GITHUB_TOKEN_MAGIC_MODULES"} {
			val, ok := lookupGithubTokenOrFallback(tokenName)
			if !ok {
				return fmt.Errorf("did not provide %s or GITHUB_TOKEN environment variable", tokenName)
			}
			env[tokenName] = val
		}
		gh := github.NewClient(env["GITHUB_TOKEN_MAGIC_MODULES"])
		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating a runner: %w", err)
		}
		ctlr := source.NewController(filepath.Join("workspace", "go"), "modular-magician", env["GITHUB_TOKEN_DOWNSTREAMS"], rnr)
		prNumber, err := strconv.Atoi(env["PR_NUMBER"])
		if err != nil {
			return fmt.Errorf("error parsing PR_NUMBER: %w", err)
		}
		return execGenerateComment(
			prNumber,
			env["GITHUB_TOKEN_MAGIC_MODULES"],
			env["BUILD_ID"],
			env["BUILD_STEP"],
			env["PROJECT_ID"],
			env["COMMIT_SHA"],
			gh,
			rnr,
			ctlr,
		)
	},
}

func listGCEnvironmentVariables() string {
	var result string
	for i, ev := range gcEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execGenerateComment(prNumber int, ghTokenMagicModules, buildId, buildStep, projectId, commitSha string, gh GithubClient, rnr ExecRunner, ctlr *source.Controller) error {
	errors := map[string][]string{"Other": []string{}}

	pullRequest, err := gh.GetPullRequest(strconv.Itoa(prNumber))
	if err != nil {
		fmt.Printf("Error getting pull request: %v\n", err)
		errors["Other"] = append(errors["Other"], "Failed to fetch PR data")
	}

	newBranch := fmt.Sprintf("auto-pr-%d", prNumber)
	oldBranch := fmt.Sprintf("auto-pr-%d-old", prNumber)
	wd := rnr.GetCWD()
	mmLocalPath := filepath.Join(wd, "..", "..")

	tpgRepo := source.Repo{
		Name:    "terraform-provider-google",
		Title:   "`google` provider",
		Path:    filepath.Join(mmLocalPath, "..", "tpg"),
		Version: provider.GA,
	}
	tpgbRepo := source.Repo{
		Name:    "terraform-provider-google-beta",
		Title:   "`google-beta` provider",
		Path:    filepath.Join(mmLocalPath, "..", "tpgb"),
		Version: provider.Beta,
	}
	tgcRepo := source.Repo{
		Name:    "terraform-google-conversion",
		Title:   "`terraform-google-conversion`",
		Path:    filepath.Join(mmLocalPath, "..", "tgc"),
		Version: provider.Beta,
	}
	tfoicsRepo := source.Repo{
		Name:  "docs-examples",
		Title: "Open in Cloud Shell",
		Path:  filepath.Join(mmLocalPath, "..", "tfoics"),
	}

	// Initialize repos
	data := diffCommentData{
		PrNumber: prNumber,
	}
	for _, repo := range []*source.Repo{&tpgRepo, &tpgbRepo, &tgcRepo, &tfoicsRepo} {
		errors[repo.Title] = []string{}
		repo.Branch = newBranch
		repo.Cloned = true
		if err := ctlr.Clone(repo); err != nil {
			fmt.Println("Failed to clone repo at new branch: ", err)
			errors[repo.Title] = append(errors[repo.Title], "Failed to clone repo at new branch")
			repo.Cloned = false
		}
		if err := ctlr.Fetch(repo, oldBranch); err != nil {
			fmt.Println("Failed to fetch old branch: ", err)
			errors[repo.Title] = append(errors[repo.Title], "Failed to clone repo at old branch")
			repo.Cloned = false
			continue
		}
		if repo.Name == "terraform-provider-google-beta" || repo.Name == "terraform-provider-google" {
			if err := ctlr.Checkout(repo, oldBranch); err != nil {
				errors[repo.Title] = append(errors[repo.Title], fmt.Sprintf("Failed to checkout branch %s", oldBranch))
				repo.Cloned = false
				continue
			}
			rnr.PushDir(repo.Path)
			if _, err := rnr.Run("make", []string{"build"}, nil); err != nil {
				errors[repo.Title] = append(errors[repo.Title], fmt.Sprintf("Failed to build branch %s", oldBranch))
				repo.Cloned = false
			}
			rnr.PopDir()
			ctlr.Checkout(repo, newBranch)
		}
	}

	diffs := []Diff{}
	for _, repo := range []*source.Repo{&tpgRepo, &tpgbRepo, &tgcRepo, &tfoicsRepo} {
		if !repo.Cloned {
			fmt.Println("Skipping diff; repo failed to clone: ", repo.Name)
			continue
		}
		shortStat, err := ctlr.DiffShortStat(repo, oldBranch, newBranch)
		if err != nil {
			fmt.Println("Failed to compute repo diff --shortstat: ", err)
			errors[repo.Title] = append(errors[repo.Title], "Failed to compute repo diff shortstats")
		}
		if shortStat != "" {
			diffs = append(diffs, Diff{
				Title:     repo.Title,
				Repo:      repo.Name,
				ShortStat: shortStat,
			})
			repo.ChangedFiles, err = ctlr.DiffNameOnly(repo, oldBranch, newBranch)
			if err != nil {
				fmt.Println("Failed to compute repo diff --name-only: ", err)
				errors[repo.Title] = append(errors[repo.Title], "Failed to compute repo changed filenames")
			}
		}
	}
	data.Diffs = diffs

	// The breaking changes are unique across both provider versions
	uniqueAffectedResources := map[string]struct{}{}
	uniqueBreakingChanges := map[string]BreakingChange{}
	diffProcessorPath := filepath.Join(mmLocalPath, "tools", "diff-processor")
	diffProcessorEnv := map[string]string{
		"OLD_REF": oldBranch,
		"NEW_REF": newBranch,
		// Passthrough vars required for a valid build environment.
		"PATH":   os.Getenv("PATH"),
		"GOPATH": os.Getenv("GOPATH"),
		"HOME":   os.Getenv("HOME"),
	}
	for _, repo := range []source.Repo{tpgRepo, tpgbRepo} {
		if !repo.Cloned {
			fmt.Println("Skipping diff processor; repo failed to clone: ", repo.Name)
			continue
		}
		if len(repo.ChangedFiles) == 0 {
			fmt.Println("Skipping diff processor; no diff: ", repo.Name)
			continue
		}
		err = buildDiffProcessor(diffProcessorPath, repo.Path, diffProcessorEnv, rnr)
		if err != nil {
			fmt.Println("building diff processor: ", err)
			errors[repo.Title] = append(errors[repo.Title], "The diff processor failed to build. This is usually due to the downstream provider failing to compile.")
			continue
		}

		breakingChanges, err := computeBreakingChanges(diffProcessorPath, rnr)
		if err != nil {
			fmt.Println("computing breaking changes: ", err)
			errors[repo.Title] = append(errors[repo.Title], "The diff processor crashed while computing breaking changes. This is usually due to the downstream provider failing to compile.")
		}
		for _, breakingChange := range breakingChanges {
			uniqueBreakingChanges[breakingChange.Message] = breakingChange
		}

		if repo.Name == "terraform-provider-google-beta" {
			// Run missing test detector (currently only for beta)
			missingTests, err := detectMissingTests(diffProcessorPath, repo.Path, rnr)
			if err != nil {
				fmt.Println("Error running missing test detector: ", err)
				errors[repo.Title] = append(errors[repo.Title], "The missing test detector failed to run.")
			}
			data.MissingTests = missingTests
		}

		affectedResources, err := changedSchemaResources(diffProcessorPath, rnr)
		if err != nil {
			fmt.Println("computing changed resource schemas: ", err)
			errors[repo.Title] = append(errors[repo.Title], "The diff processor crashed while computing changed resource schemas.")
		}
		for _, resource := range affectedResources {
			uniqueAffectedResources[resource] = struct{}{}
		}
	}
	breakingChangesSlice := maps.Values(uniqueBreakingChanges)
	sort.Slice(breakingChangesSlice, func(i, j int) bool {
		return breakingChangesSlice[i].Message < breakingChangesSlice[j].Message
	})
	data.BreakingChanges = breakingChangesSlice

	// Compute affected resources based on changed files
	changedFilesAffectedResources := map[string]struct{}{}
	for _, repo := range []source.Repo{tpgRepo, tpgbRepo} {
		if !repo.Cloned {
			fmt.Println("Skipping changed file service labels; repo failed to clone: ", repo.Name)
			continue
		}
		for _, path := range repo.ChangedFiles {
			if r := fileToResource(path); r != "" {
				uniqueAffectedResources[r] = struct{}{}
				changedFilesAffectedResources[r] = struct{}{}
			}
		}
	}
	fmt.Printf("affected resources based on changed files: %v\n", maps.Keys(changedFilesAffectedResources))

	// Compute service labels based on affected resources
	uniqueServiceLabels := map[string]struct{}{}
	regexpLabels, err := labeler.BuildRegexLabels(labeler.EnrolledTeamsYaml)
	if err != nil {
		fmt.Println("error building regexp labels: ", err)
		errors["Other"] = append(errors["Other"], "Failed to parse service label mapping")
	}
	if len(regexpLabels) > 0 {
		for _, label := range labeler.ComputeLabels(maps.Keys(uniqueAffectedResources), regexpLabels) {
			uniqueServiceLabels[label] = struct{}{}
		}
	}

	// Add service labels to PR if it doesn't already have service labels
	if len(uniqueServiceLabels) > 0 {
		// short-circuit if service labels have already been added to the PR
		hasServiceLabels := false
		for _, label := range pullRequest.Labels {
			if strings.HasPrefix(label.Name, "service/") {
				hasServiceLabels = true
			}
		}
		if !hasServiceLabels {
			serviceLabelsSlice := maps.Keys(uniqueServiceLabels)
			sort.Strings(serviceLabelsSlice)
			if len(serviceLabelsSlice) > 3 {
				// Treat this as a cross-provider change
				serviceLabelsSlice = []string{"service/terraform"}
			}
			if err = gh.AddLabels(strconv.Itoa(prNumber), serviceLabelsSlice); err != nil {
				fmt.Printf("Error posting new service labels %q: %s", serviceLabelsSlice, err)
				errors["Other"] = append(errors["Other"], "Failed to update service labels")
			}
		}
	}

	// Update breaking changes status on PR
	breakingState := "success"
	if len(uniqueBreakingChanges) > 0 {
		breakingState = "failure"
		// If fetching the PR failed, Labels will be empty
		for _, label := range pullRequest.Labels {
			if label.Name == allowBreakingChangesLabel {
				breakingState = "success"
				break
			}
		}
	}
	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildId, buildStep, projectId)
	if err = gh.PostBuildStatus(strconv.Itoa(prNumber), "terraform-provider-breaking-change-test", breakingState, targetURL, commitSha); err != nil {
		fmt.Printf("Error posting build status for pr %d commit %s: %v\n", prNumber, commitSha, err)
		errors["Other"] = append(errors["Other"], "Failed to update breaking-change status check with state: "+breakingState)
	}

	// Add errors to data as an ordered list
	errorsList := []Errors{}
	for _, repo := range []source.Repo{tpgRepo, tpgbRepo, tgcRepo, tfoicsRepo} {
		if len(errors[repo.Title]) > 0 {
			errorsList = append(errorsList, Errors{
				Title:  repo.Title,
				Errors: errors[repo.Title],
			})
		}
	}
	if len(errors["Other"]) > 0 {
		errorsList = append(errorsList, Errors{
			Title:  "Other",
			Errors: errors["Other"],
		})
	}
	data.Errors = errorsList

	// Post diff comment
	message, err := formatDiffComment(data)
	if err != nil {
		fmt.Printf("Data: %v\n", data)
		return fmt.Errorf("error formatting message: %w", err)
	}
	if err := gh.PostComment(strconv.Itoa(prNumber), message); err != nil {
		fmt.Println("Comment: ", message)
		return fmt.Errorf("error posting comment to PR %d: %w", prNumber, err)
	}
	return nil
}

// Build the diff processor for tpg or tpgb
func buildDiffProcessor(diffProcessorPath, providerLocalPath string, env map[string]string, rnr ExecRunner) error {
	for _, path := range []string{"old", "new", "bin"} {
		if err := rnr.RemoveAll(filepath.Join(diffProcessorPath, path)); err != nil {
			return err
		}
	}
	if err := rnr.PushDir(diffProcessorPath); err != nil {
		return err
	}
	for _, path := range []string{"old", "new"} {
		if err := rnr.Copy(providerLocalPath, filepath.Join(diffProcessorPath, path)); err != nil {
			return err
		}
	}
	if _, err := rnr.Run("make", []string{"build"}, env); err != nil {
		return fmt.Errorf("Error running make build in %s: %v\n", diffProcessorPath, err)
	}
	return rnr.PopDir()
}

func computeBreakingChanges(diffProcessorPath string, rnr ExecRunner) ([]BreakingChange, error) {
	if err := rnr.PushDir(diffProcessorPath); err != nil {
		return nil, err
	}
	output, err := rnr.Run("bin/diff-processor", []string{"breaking-changes"}, nil)
	if err != nil {
		return nil, err
	}

	if output == "" {
		return nil, nil
	}

	var changes []BreakingChange
	if err = json.Unmarshal([]byte(output), &changes); err != nil {
		return nil, err
	}
	return changes, rnr.PopDir()
}

func changedSchemaResources(diffProcessorPath string, rnr ExecRunner) ([]string, error) {
	if err := rnr.PushDir(diffProcessorPath); err != nil {
		return nil, err
	}

	output, err := rnr.Run("bin/diff-processor", []string{"changed-schema-resources"}, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Resources with changed schemas: " + output)

	var labels []string
	if err = json.Unmarshal([]byte(output), &labels); err != nil {
		return nil, err
	}

	if err = rnr.PopDir(); err != nil {
		return nil, err
	}
	return labels, nil
}

// Run the missing test detector and return the results.
// Returns an empty string unless there are missing tests.
// Error will be nil unless an error occurs during setup.
func detectMissingTests(diffProcessorPath, tpgbLocalPath string, rnr ExecRunner) (map[string]*MissingTestInfo, error) {
	if err := rnr.PushDir(diffProcessorPath); err != nil {
		return nil, err
	}

	output, err := rnr.Run("bin/diff-processor", []string{"detect-missing-tests", fmt.Sprintf("%s/google-beta/services", tpgbLocalPath)}, nil)
	if err != nil {
		return nil, err
	}

	var missingTests map[string]*MissingTestInfo
	if err = json.Unmarshal([]byte(output), &missingTests); err != nil {
		return nil, err
	}
	return missingTests, rnr.PopDir()
}

func formatDiffComment(data diffCommentData) (string, error) {
	tmpl, err := template.New("DIFF_COMMENT.md").Parse(diffComment)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse DIFF_COMMENT.md: %s", err))
	}
	sb := new(strings.Builder)
	err = tmpl.Execute(sb, data)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

var resourceFileRegexp = regexp.MustCompile(`^.*/services/[^/]+/(?:data_source_|resource_|iam_)(.*?)(?:_test|_sweeper|_iam_test|_generated_test|_internal_test)?.go`)
var resourceDocsRegexp = regexp.MustCompile(`^.*website/docs/(?:r|d)/(.*).html.markdown`)

func fileToResource(path string) string {
	var submatches []string
	if strings.HasSuffix(path, ".go") {
		submatches = resourceFileRegexp.FindStringSubmatch(path)
	} else if strings.HasSuffix(path, ".html.markdown") {
		submatches = resourceDocsRegexp.FindStringSubmatch(path)
	}

	if len(submatches) == 0 {
		return ""
	}

	// The regexes will each return the resource name as the first
	// submatch, stripping any prefixes or suffixes.
	resource := submatches[1]

	if !strings.HasPrefix(resource, "google_") {
		resource = "google_" + resource
	}
	return resource
}

func pathChanged(path string, changedFiles []string) bool {
	for _, f := range changedFiles {
		if strings.HasPrefix(f, path) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(generateCommentCmd)
}

package cmd

import (
	"fmt"
	"magician/exec"
	"magician/github"
	"magician/source"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var testTGCIntegrationCmd = &cobra.Command{
	Use:   "test-tgc-integration",
	Short: "Run tgc integration tests via workflow dispatch",
	Long: `This command runs tgc unit tests via workflow dispatch

	The following PR details are expected as environment variables:
	1. GOPATH
	2. GITHUB_TOKEN_MAGIC_MODULES
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		goPath, ok := os.LookupEnv("GOPATH")
		if !ok {
			return fmt.Errorf("did not provide GOPATH environment variable")
		}

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			return fmt.Errorf("did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating runner: %w", err)
		}

		ctlr := source.NewController(goPath, "modular-magician", githubToken, rnr)

		gh := github.NewClient(githubToken)

		return execTestTGCIntegration(args[0], args[1], args[2], args[3], args[4], args[5], "modular-magician", rnr, ctlr, gh)
	},
}

func execTestTGCIntegration(prNumber, mmCommit, buildID, projectID, buildStep, ghRepo, githubUsername string, rnr exec.ExecRunner, ctlr *source.Controller, gh GithubClient) error {
	newBranch := "auto-pr-" + prNumber
	repo := &source.Repo{
		Name:   ghRepo,
		Branch: newBranch,
	}
	ctlr.SetPath(repo)
	if err := ctlr.Clone(repo); err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}
	if err := rnr.PushDir(repo.Path); err != nil {
		return fmt.Errorf("error changing to repo dir: %w", err)
	}
	diffs, err := rnr.Run("git", []string{"diff", "--name-only", "HEAD~1"}, nil)
	if err != nil {
		return fmt.Errorf("error diffing repo: %w", err)
	}
	hasGoFiles := false
	for _, diff := range strings.Split(diffs, "\n") {
		if strings.HasSuffix(diff, ".go") {
			hasGoFiles = true
			break
		}
	}
	if !hasGoFiles {
		fmt.Println("Skipping tests: No go files changed")
		return nil
	}

	fmt.Println("Running tests: Go files changed")

	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(prNumber, ghRepo+"-test-integration", "pending", targetURL, mmCommit); err != nil {
		return fmt.Errorf("error posting build status: %w", err)
	}

	if _, err := rnr.Run("go", []string{"mod", "edit", "-replace", fmt.Sprintf("github.com/hashicorp/terraform-provider-google-beta=github.com/%s/terraform-provider-google-beta@%s", githubUsername, newBranch)}, nil); err != nil {
		fmt.Println("Error running go mod edit: ", err)
	}
	if _, err := rnr.Run("go", []string{"mod", "tidy"}, nil); err != nil {
		fmt.Println("Error running go mod tidy: ", err)
	}

	if _, err := rnr.Run("make", []string{"build"}, nil); err != nil {
		fmt.Println("Error running make build: ", err)
	}
	state := "success"
	if _, err := rnr.Run("make", []string{"test-integration"}, nil); err != nil {
		fmt.Println("Error running make test-integration: ", err)
		state = "failure"
	}

	if err := gh.PostBuildStatus(prNumber, ghRepo+"-test-integration", state, targetURL, mmCommit); err != nil {
		return fmt.Errorf("error posting build status: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(testTGCIntegrationCmd)
}

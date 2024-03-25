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
	Run: func(cmd *cobra.Command, args []string) {
		goPath, ok := os.LookupEnv("GOPATH")
		if !ok {
			fmt.Println("Did not provide GOPATH environment variable")
			os.Exit(1)
		}

		githubToken, ok := lookupGithubTokenOrFallback("GITHUB_TOKEN_MAGIC_MODULES")
		if !ok {
			fmt.Println("Did not provide GITHUB_TOKEN_MAGIC_MODULES or GITHUB_TOKEN environment variables")
			os.Exit(1)
		}

		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating runner: ", err)
			os.Exit(1)
		}

		ctlr := source.NewController(goPath, "modular-magician", githubToken, rnr)

		gh := github.NewClient(githubToken)

		execTestTGCIntegration(args[0], args[1], args[2], args[3], args[4], args[5], "modular-magician", rnr, ctlr, gh)
	},
}

func execTestTGCIntegration(prNumber, mmCommit, buildID, projectID, buildStep, ghRepo, githubUsername string, rnr ExecRunner, ctlr *source.Controller, gh GithubClient) {
	newBranch := "auto-pr-" + prNumber
	repo := &source.Repo{
		Name:   ghRepo,
		Branch: newBranch,
	}
	ctlr.SetPath(repo)
	if err := ctlr.Clone(repo); err != nil {
		fmt.Println("Error cloning repo: ", err)
		os.Exit(1)
	}
	if err := rnr.PushDir(repo.Path); err != nil {
		fmt.Println("Error changing to repo dir: ", err)
		os.Exit(1)
	}
	diffs, err := rnr.Run("git", []string{"diff", "--name-only", "HEAD~1"}, nil)
	if err != nil {
		fmt.Println("Error diffing repo: ", err)
		os.Exit(1)
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
		os.Exit(0)
	}

	fmt.Println("Running tests: Go files changed")

	targetURL := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds;region=global/%s;step=%s?project=%s", buildID, buildStep, projectID)
	if err := gh.PostBuildStatus(prNumber, ghRepo+"-test-integration", "pending", targetURL, mmCommit); err != nil {
		fmt.Println("Error posting build status: ", err)
		os.Exit(1)
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
		fmt.Println("Error posting build status: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(testTGCIntegrationCmd)
}

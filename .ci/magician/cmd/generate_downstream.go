package cmd

import (
	"fmt"
	"magician/exec"
	"magician/github"
	"magician/provider"
	"magician/source"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var changelogExp = regexp.MustCompile("(?s)```release-note.*?```")

var gdEnvironmentVariables = [...]string{
	"BASE_BRANCH",
	"GOPATH",
}

var gdTokenEnvironmentVariables = [...]string{
	"GITHUB_TOKEN_CLASSIC",
	"GITHUB_TOKEN_DOWNSTREAMS",
	"GITHUB_TOKEN",
}

var generateDownstreamCmd = &cobra.Command{
	Use:   "generate-downstream",
	Short: "Run generate downstream",
	Long: `This command runs after pull requests are merged to generate corresponding changes in downstream repos.

	It expects the following arguments:
	1. Command, either head, base, or downstream
	2. Name of the downstream repo, either terraform, terraform-google-conversion, or tf-oics
	3. Version of the downstream
	4. Commit SHA of the squashed merge commit

	The following environment variables should be set:
` + listGDEnvironmentVariables(),
	RunE: func(cmd *cobra.Command, args []string) error {
		env := make(map[string]string, len(gdEnvironmentVariables))
		for _, ev := range gdEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				return fmt.Errorf("did not provide %s environment variable", ev)
			}
			env[ev] = val
		}

		var githubToken string
		for _, ev := range gdTokenEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if ok {
				env[ev] = val
				githubToken = val
				break
			}
		}

		gh := github.NewClient(githubToken)
		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating a runner: %w", err)
		}
		ctlr := source.NewController(env["GOPATH"], "modular-magician", githubToken, rnr)
		oldToken := os.Getenv("GITHUB_TOKEN")
		if err := os.Setenv("GITHUB_TOKEN", githubToken); err != nil {
			return fmt.Errorf("error setting GITHUB_TOKEN environment variable: %w", err)
		}
		defer func() {
			if err := os.Setenv("GITHUB_TOKEN", oldToken); err != nil {
				fmt.Println("Error setting GITHUB_TOKEN environment variable: ", err)
			}
		}()

		if len(args) != 4 {
			return fmt.Errorf("wrong number of arguments %d, expected 4", len(args))
		}

		return execGenerateDownstream(env["BASE_BRANCH"], args[0], args[1], args[2], args[3], gh, rnr, ctlr)
	},
}

func listGDEnvironmentVariables() string {
	var result string
	for i, ev := range gdEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execGenerateDownstream(baseBranch, command, repo, version, ref string, gh GithubClient, rnr exec.ExecRunner, ctlr *source.Controller) error {
	if baseBranch == "" {
		baseBranch = "main"
	}
	if command == "downstream" {
		var syncBranchPrefix string
		if repo == "terraform" {
			if version == "beta" {
				syncBranchPrefix = "tpgb-sync"
			} else if version == "ga" {
				syncBranchPrefix = "tpg-sync"
			}
		} else if repo == "terraform-google-conversion" {
			syncBranchPrefix = "tgc-sync"
		} else if repo == "tf-oics" {
			syncBranchPrefix = "tf-oics-sync"
		}
		syncBranch := getSyncBranch(syncBranchPrefix, baseBranch)
		if syncBranchHasCommit(ref, syncBranch, rnr) {
			fmt.Printf("Sync branch %s already has commit %s, skipping generation\n", syncBranch, ref)
			os.Exit(0)
		}
	}

	mmLocalPath := filepath.Join(rnr.GetCWD(), "..", "..")
	mmCopyPath := filepath.Join(mmLocalPath, "..", fmt.Sprintf("mm-%s-%s-%s", repo, version, command))
	if _, err := rnr.Run("cp", []string{"-rp", mmLocalPath, mmCopyPath}, nil); err != nil {
		return fmt.Errorf("error copying magic modules: %w", err)
	}
	mmRepo := &source.Repo{
		Name: "magic-modules",
		Path: mmCopyPath,
	}

	downstreamRepo, scratchRepo, commitMessage, err := cloneRepo(mmRepo, baseBranch, repo, version, command, ref, rnr, ctlr)
	if err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	if err := rnr.PushDir(mmCopyPath); err != nil {
		return fmt.Errorf("error changing directory to copied magic modules: %w", err)
	}

	if err := setGitConfig(rnr); err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}

	if err := runMake(downstreamRepo, command, rnr); err != nil {
		return fmt.Errorf("error running make: %w", err)
	}

	var pullRequest *github.PullRequest
	if command == "downstream" {
		pullRequest, err = getPullRequest(baseBranch, ref, gh)
		if err != nil {
			return fmt.Errorf("error getting pull request: %w", err)
		}
		if repo == "terraform" {
			if err := addChangelogEntry(scratchRepo, pullRequest, rnr); err != nil {
				return fmt.Errorf("error adding changelog entry: %w", err)
			}
		}
	}

	scratchCommitSha, commitErr := createCommit(scratchRepo, commitMessage, rnr)
	if commitErr != nil {
		fmt.Println("Error creating commit: ", commitErr)
		if !strings.Contains(commitErr.Error(), "nothing to commit") {
			return fmt.Errorf("error creating commit: %w", commitErr)
		}
	}

	if _, err := rnr.Run("git", []string{"push", ctlr.URL(scratchRepo), scratchRepo.Branch, "-f"}, nil); err != nil {
		return fmt.Errorf("error pushing commit: %w", err)
	}

	if commitErr == nil && command == "downstream" {
		if err := mergePullRequest(downstreamRepo, scratchRepo, scratchCommitSha, pullRequest, rnr, gh); err != nil {
			return fmt.Errorf("error merging pull request: %w", err)
		}
	}
	return nil
}

func cloneRepo(mmRepo *source.Repo, baseBranch, repo, version, command, ref string, rnr exec.ExecRunner, ctlr *source.Controller) (*source.Repo, *source.Repo, string, error) {
	downstreamRepo := &source.Repo{
		Title:  repo,
		Branch: baseBranch,
	}
	switch repo {
	case "terraform":
		if version == "ga" {
			downstreamRepo.Name = "terraform-provider-google"
			downstreamRepo.Version = provider.GA
		} else if version == "beta" {
			downstreamRepo.Name = "terraform-provider-google-beta"
			downstreamRepo.Version = provider.Beta
		} else {
			return nil, nil, "", fmt.Errorf("unrecognized version %s", version)
		}
		downstreamRepo.Owner = "hashicorp"
	case "terraform-google-conversion":
		downstreamRepo.Name = "terraform-google-conversion"
		downstreamRepo.Owner = "GoogleCloudPlatform"
	case "tf-oics":
		if downstreamRepo.Branch == "main" {
			downstreamRepo.Branch = "master"
		}
		downstreamRepo.Name = "docs-examples"
		downstreamRepo.Owner = "terraform-google-modules"
	case "tf-cloud-docs":
		fmt.Println(repo, " is no longer available.")
		return nil, nil, "", nil
	default:
		return nil, nil, "", fmt.Errorf("unrecognized repo %s", repo)
	}
	ctlr.SetPath(downstreamRepo)
	if err := ctlr.Clone(downstreamRepo); err != nil {
		return nil, nil, "", err
	}
	scratchRepo := &source.Repo{
		Name:    downstreamRepo.Name,
		Owner:   "modular-magician",
		Path:    downstreamRepo.Path,
		Version: downstreamRepo.Version,
	}
	var commitMessage string
	switch command {
	case "head":
		scratchRepo.Branch = "auto-pr-" + ref
		commitMessage = fmt.Sprintf("New generated code for MM PR %s.", ref)
	case "base":
		// In this case, there is guaranteed to be a merge commit,
		// and the *left* side of it is the old main branch.
		// the *right* side of it is the code to be merged.
		if err := ctlr.Checkout(mmRepo, "HEAD~"); err != nil {
			return nil, nil, "", err
		}
		scratchRepo.Branch = fmt.Sprintf("auto-pr-%s-old", ref)
		commitMessage = fmt.Sprintf("Old generated code for MM PR %s.", ref)
	case "downstream":
		scratchRepo.Branch = "downstream-pr-" + ref
		originalMessage, err := rnr.Run("git", []string{"log", "-1", "--pretty=%B", ref}, nil)
		if err != nil {
			return nil, nil, "", err
		}
		commitMessage = fmt.Sprintf("%s\n[upstream:%s]", originalMessage, ref)
	}
	return downstreamRepo, scratchRepo, commitMessage, nil
}

func setGitConfig(rnr exec.ExecRunner) error {
	if _, err := rnr.Run("git", []string{"config", "--local", "user.name", "Modular Magician"}, nil); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"config", "--local", "user.email", "magic-modules@google.com"}, nil); err != nil {
		return err
	}
	return nil
}

func runMake(downstreamRepo *source.Repo, command string, rnr exec.ExecRunner) error {
	switch downstreamRepo.Title {
	case "terraform-google-conversion":
		if _, err := rnr.Run("make", []string{"clean-tgc", "OUTPUT_PATH=" + downstreamRepo.Path}, nil); err != nil {
			return err
		}
		if _, err := rnr.Run("make", []string{"tgc", "OUTPUT_PATH=" + downstreamRepo.Path}, nil); err != nil {
			return err
		}
		if command == "downstream" {
			if err := rnr.PushDir(downstreamRepo.Path); err != nil {
				return err
			}
			if _, err := rnr.Run("go", []string{"get", "-d", "github.com/hashicorp/terraform-provider-google-beta@" + downstreamRepo.Branch}, nil); err != nil {
				return err
			}
			if _, err := rnr.Run("go", []string{"mod", "tidy"}, nil); err != nil {
				return err
			}
			if _, err := rnr.Run("make", []string{"build"}, nil); err != nil {
				fmt.Println("Error building tgc: ", err)
			}
			if err := rnr.PopDir(); err != nil {
				return err
			}
		}
	case "tf-oics":
		if _, err := rnr.Run("make", []string{"tf-oics", "OUTPUT_PATH=" + downstreamRepo.Path}, nil); err != nil {
			return err
		}
	case "terraform":
		if _, err := rnr.Run("make", []string{"clean-provider", "OUTPUT_PATH=" + downstreamRepo.Path}, nil); err != nil {
			return err
		}
		if _, err := rnr.Run("make", []string{"provider", "OUTPUT_PATH=" + downstreamRepo.Path, fmt.Sprintf("VERSION=%s", downstreamRepo.Version)}, nil); err != nil {
			return err
		}
	}
	return nil
}

func getPullRequest(baseBranch, ref string, gh GithubClient) (*github.PullRequest, error) {
	prs, err := gh.GetPullRequests("closed", baseBranch, "updated", "desc")
	if err != nil {
		return nil, err
	}
	for _, pr := range prs {
		if pr.MergeCommitSha == ref {
			return &pr, nil
		}
	}
	return nil, fmt.Errorf("no pr found with merge commit sha %s and base branch %s", ref, baseBranch)
}

func createCommit(scratchRepo *source.Repo, commitMessage string, rnr exec.ExecRunner) (string, error) {
	if err := rnr.PushDir(scratchRepo.Path); err != nil {
		return "", err
	}
	if err := setGitConfig(rnr); err != nil {
		return "", err
	}

	if _, err := rnr.Run("git", []string{"add", "."}, nil); err != nil {
		return "", err
	}
	if _, err := rnr.Run("git", []string{"checkout", "-b", scratchRepo.Branch}, nil); err != nil {
		return "", err
	}

	if _, err := rnr.Run("git", []string{"commit", "--signoff", "-m", commitMessage}, nil); err != nil {
		return "", err
	}

	commitSha, err := rnr.Run("git", []string{"rev-parse", "HEAD"}, nil)
	if err != nil {
		return "", fmt.Errorf("error retrieving commit sha: %w", err)
	}

	commitSha = strings.TrimSpace(commitSha)
	fmt.Printf("Commit sha on the branch is: `%s`\n", commitSha)

	return commitSha, err
}

func addChangelogEntry(downstreamRepo *source.Repo, pullRequest *github.PullRequest, rnr exec.ExecRunner) error {
	if err := rnr.PushDir(downstreamRepo.Path); err != nil {
		return err
	}
	rnr.Mkdir(".changelog")
	if err := rnr.WriteFile(filepath.Join(".changelog", fmt.Sprintf("%d.txt", pullRequest.Number)), strings.Join(changelogExp.FindAllString(pullRequest.Body, -1), "\n")); err != nil {
		return err
	}
	return rnr.PopDir()
}

func mergePullRequest(downstreamRepo, scratchRepo *source.Repo, scratchRepoSha string, pullRequest *github.PullRequest, rnr exec.ExecRunner, gh GithubClient) error {
	fmt.Printf(`Base: %s:%s
Head: %s:%s
`, downstreamRepo.Owner, downstreamRepo.Branch, scratchRepo.Owner, scratchRepo.Branch)
	newPRURL, err := rnr.Run("hub", []string{
		"pull-request",
		"-b",
		fmt.Sprintf("%s:%s",
			downstreamRepo.Owner,
			downstreamRepo.Branch),
		"-h",
		fmt.Sprintf("%s:%s",
			scratchRepo.Owner,
			scratchRepo.Branch),
		"-m",
		pullRequest.Title,
		"-m",
		pullRequest.Body,
		"-m",
		"Derived from " + pullRequest.HTMLUrl,
	}, nil)
	if err != nil {
		return err
	}
	fmt.Println("Created PR ", newPRURL)
	newPRURLParts := strings.Split(newPRURL, "/")
	newPRNumber := strings.TrimSuffix(newPRURLParts[len(newPRURLParts)-1], "\n")

	// Wait a few seconds, then merge the PR.
	time.Sleep(5 * time.Second)
	fmt.Println("Merging PR ", newPRURL)
	if err := gh.MergePullRequest(downstreamRepo.Owner, downstreamRepo.Name, newPRNumber, scratchRepoSha); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateDownstreamCmd)
}

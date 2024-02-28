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
	"GITHUB_TOKEN",
	"GOPATH",
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
	Run: func(cmd *cobra.Command, args []string) {
		env := make(map[string]string, len(gdEnvironmentVariables))
		for _, ev := range gdEnvironmentVariables {
			val, ok := os.LookupEnv(ev)
			if !ok {
				fmt.Printf("Did not provide %s environment variable\n", ev)
				os.Exit(1)
			}
			env[ev] = val
		}

		gh := github.NewClient()
		rnr, err := exec.NewRunner()
		if err != nil {
			fmt.Println("Error creating a runner: ", err)
			os.Exit(1)
		}
		ctlr := source.NewController(env["GOPATH"], "modular-magician", env["GITHUB_TOKEN"], rnr)

		if len(args) != 4 {
			fmt.Printf("Wrong number of arguments %d, expected 4\n", len(args))
			os.Exit(1)
		}

		execGenerateDownstream(env["BASE_BRANCH"], args[0], args[1], args[2], args[3], gh, rnr, ctlr)
	},
}

func listGDEnvironmentVariables() string {
	var result string
	for i, ev := range gdEnvironmentVariables {
		result += fmt.Sprintf("\t%2d. %s\n", i+1, ev)
	}
	return result
}

func execGenerateDownstream(baseBranch, command, repo, version, ref string, gh GithubClient, rnr ExecRunner, ctlr *source.Controller) {
	if baseBranch == "" {
		baseBranch = "main"
	}

	mmLocalPath := filepath.Join(rnr.GetCWD(), "..", "..")
	mmCopyPath := filepath.Join(mmLocalPath, "..", fmt.Sprintf("mm-%s-%s-%s", repo, version, command))
	if _, err := rnr.Run("cp", []string{"-rp", mmLocalPath, mmCopyPath}, nil); err != nil {
		fmt.Println("Error copying magic modules: ", err)
		os.Exit(1)
	}
	mmRepo := &source.Repo{
		Name: "magic-modules",
		Path: mmCopyPath,
	}

	downstreamRepo, scratchRepo, commitMessage, err := cloneRepo(mmRepo, baseBranch, repo, version, command, ref, rnr, ctlr)
	if err != nil {
		fmt.Println("Error cloning repo: ", err)
		os.Exit(1)
	}

	if err := rnr.PushDir(mmCopyPath); err != nil {
		fmt.Println("Error changing directory to copied magic modules: ", err)
		os.Exit(1)
	}

	if err := setGitConfig(rnr); err != nil {
		fmt.Println("Error setting config: ", err)
		os.Exit(1)
	}

	if err := runMake(downstreamRepo, command, rnr); err != nil {
		fmt.Println("Error running make: ", err)
		os.Exit(1)
	}

	commitErr := createCommit(scratchRepo, commitMessage, rnr)
	if commitErr != nil {
		fmt.Println("Error creating commit: ", commitErr)
	}

	var pullRequest *github.PullRequest
	if commitErr == nil && command == "downstream" {
		pullRequest, err = getPullRequest(baseBranch, ref, gh)
		if err != nil {
			fmt.Println("Error getting pull request: ", err)
			os.Exit(1)
		}
		if repo == "terraform" {
			if err := addChangelogEntry(pullRequest, rnr); err != nil {
				fmt.Println("Error adding changelog entry: ", err)
				os.Exit(1)
			}
		}
	}

	if _, err := rnr.Run("git", []string{"push", ctlr.URL(scratchRepo), scratchRepo.Branch, "-f"}, nil); err != nil {
		fmt.Println("Error pushing commit: ", err)
		os.Exit(1)
	}

	if commitErr == nil && command == "downstream" {
		if err := mergePullRequest(downstreamRepo, scratchRepo, pullRequest, rnr, gh); err != nil {
			fmt.Println("Error merging pull request: ", err)
			os.Exit(1)
		}
	}
}

func cloneRepo(mmRepo *source.Repo, baseBranch, repo, version, command, ref string, rnr ExecRunner, ctlr *source.Controller) (*source.Repo, *source.Repo, string, error) {
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

func setGitConfig(rnr ExecRunner) error {
	if _, err := rnr.Run("git", []string{"config", "--local", "user.name", "Modular Magician"}, nil); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"config", "--local", "user.email", "magic-modules@google.com"}, nil); err != nil {
		return err
	}
	return nil
}

func runMake(downstreamRepo *source.Repo, command string, rnr ExecRunner) error {
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

func createCommit(scratchRepo *source.Repo, commitMessage string, rnr ExecRunner) error {
	if err := rnr.PushDir(scratchRepo.Path); err != nil {
		return err
	}
	if err := setGitConfig(rnr); err != nil {
		return err
	}

	if _, err := rnr.Run("git", []string{"add", "."}, nil); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"checkout", "-b", scratchRepo.Branch}, nil); err != nil {
		return err
	}

	if _, err := rnr.Run("git", []string{"commit", "--signoff", "-m", commitMessage}, nil); err != nil {
		return err
	}

	return nil
}

func addChangelogEntry(pullRequest *github.PullRequest, rnr ExecRunner) error {
	rnr.Mkdir(".changelog")
	if err := rnr.WriteFile(filepath.Join(".changelog", fmt.Sprintf("%d.txt", pullRequest.Number)), strings.Join(changelogExp.FindAllString(pullRequest.Body, -1), "\n")); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"add", "."}, nil); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"commit", "--signoff", "--amend", "--no-edit"}, nil); err != nil {
		return err
	}
	return nil
}

func mergePullRequest(downstreamRepo, scratchRepo *source.Repo, pullRequest *github.PullRequest, rnr ExecRunner, gh GithubClient) error {
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
	if err := gh.MergePullRequest(downstreamRepo.Owner, downstreamRepo.Name, newPRNumber); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateDownstreamCmd)
}

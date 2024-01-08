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

func execGenerateDownstream(baseBranch, command, repo, version, ref string, gh GithubClient, rnr *exec.Runner, ctlr *source.Controller) {
	if baseBranch == "" {
		baseBranch = "main"
	}
	upstreamRepo, scratchRepo, err := cloneRepo(baseBranch, repo, version, ctlr)
	if err != nil {
		fmt.Println("Error cloning repo: ", err)
		os.Exit(1)
	}

	mmLocalPath := filepath.Join(rnr.GetCWD(), "..", "..")
	mmCopyPath := filepath.Join(mmLocalPath, "..", fmt.Sprintf("mm-%s-%s-%s", repo, version, command))
	if _, err := rnr.Run("cp", []string{"-rp", mmLocalPath, mmCopyPath}, nil); err != nil {
		fmt.Println("Error copying magic modules: ", err)
		os.Exit(1)
	}

	if err := rnr.PushDir(mmCopyPath); err != nil {
		fmt.Println("Error changing directory to copied magic modules: ", err)
		os.Exit(1)
	}

	if err := setConfig(rnr); err != nil {
		fmt.Println("Error setting config: ", err)
		os.Exit(1)
	}

	mmRepo := &source.Repo{
		Name: "magic-modules",
		Path: mmCopyPath,
	}
	commitMessage, err := branchAndCommitMessage(mmRepo, scratchRepo, command, ref, rnr, ctlr)
	if err != nil {
		fmt.Println("Error getting branch and commit message: ", err)
		os.Exit(1)
	}

	if err := runMake(upstreamRepo, command, rnr); err != nil {
		fmt.Println("Error running make: ", err)
		os.Exit(1)
	}

	if err := pushCommit(upstreamRepo, scratchRepo, baseBranch, command, commitMessage, ref, gh, rnr, ctlr); err != nil {
		fmt.Println("Error pushing commit: ", err)
		os.Exit(1)
	}
}

func cloneRepo(baseBranch, repo, version string, ctlr *source.Controller) (*source.Repo, *source.Repo, error) {
	upstreamRepo := &source.Repo{
		Title:  repo,
		Branch: baseBranch,
	}
	switch repo {
	case "terraform":
		if version == "ga" {
			upstreamRepo.Name = "terraform-provider-google"
			upstreamRepo.Version = provider.GA
		} else if version == "beta" {
			upstreamRepo.Name = "terraform-provider-google-beta"
			upstreamRepo.Version = provider.Beta
		} else {
			return nil, nil, fmt.Errorf("unrecognized version %s", version)
		}
		upstreamRepo.Owner = "hashicorp"
	case "terraform-google-conversion":
		upstreamRepo.Name = "terraform-google-conversion"
		upstreamRepo.Owner = "GoogleCloudPlatform"
	case "tf-oics":
		if upstreamRepo.Branch == "main" {
			upstreamRepo.Branch = "master"
		}
		upstreamRepo.Name = "docs-examples"
		upstreamRepo.Owner = "terraform-google-modules"
	case "tf-cloud-docs":
		fmt.Println(repo, " is no longer available.")
		return nil, nil, nil
	default:
		return nil, nil, fmt.Errorf("unrecognized repo %s", repo)
	}
	ctlr.SetPath(upstreamRepo)
	if err := ctlr.Clone(upstreamRepo); err != nil {
		return nil, nil, err
	}
	scratchRepo := &source.Repo{
		Name:    upstreamRepo.Name,
		Branch:  baseBranch,
		Owner:   "modular-magician",
		Path:    upstreamRepo.Path,
		Version: upstreamRepo.Version,
	}
	return upstreamRepo, scratchRepo, nil
}

func setConfig(rnr *exec.Runner) error {
	if _, err := rnr.Run("git", []string{"config", "--local", "user.name", "Modular Magician"}, nil); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"config", "--local", "user.email", "magic-modules@google.com"}, nil); err != nil {
		return err
	}
	return nil
}

func branchAndCommitMessage(mmRepo, scratchRepo *source.Repo, command, ref string, rnr *exec.Runner, ctlr *source.Controller) (string, error) {
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
			return "", err
		}
		scratchRepo.Branch = fmt.Sprintf("auto-pr-%s-old", ref)
		commitMessage = fmt.Sprintf("Old generated code for MM PR %s.", ref)
	case "downstream":
		scratchRepo.Branch = "downstream-pr-" + ref
		originalMessage, err := rnr.Run("git", []string{"log", "-1", "--pretty=%B", ref}, nil)
		if err != nil {
			return "", err
		}
		commitMessage = fmt.Sprintf("%s\n[upstream:%s]", originalMessage, ref)
	}
	return commitMessage, nil
}

func runMake(upstreamRepo *source.Repo, command string, rnr *exec.Runner) error {
	switch upstreamRepo.Title {
	case "terraform-google-conversion":
		if _, err := rnr.Run("make", []string{"clean-tgc", "OUTPUT_PATH=" + upstreamRepo.Path}, nil); err != nil {
			return err
		}
		if _, err := rnr.Run("make", []string{"tgc", "OUTPUT_PATH=" + upstreamRepo.Path}, nil); err != nil {
			return err
		}
		if command == "downstream" {
			if err := rnr.PushDir(upstreamRepo.Path); err != nil {
				return err
			}
			if _, err := rnr.Run("go", []string{"get", "-d", "github.com/hashicorp/terraform-provider-google-beta@" + upstreamRepo.Branch}, nil); err != nil {
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
		if _, err := rnr.Run("make", []string{"tf-oics", "OUTPUT_PATH=" + upstreamRepo.Path}, nil); err != nil {
			return err
		}
	case "terraform":
		if _, err := rnr.Run("make", []string{"clean-provider", "OUTPUT_PATH=" + upstreamRepo.Path}, nil); err != nil {
			return err
		}
		if _, err := rnr.Run("make", []string{"provider", "OUTPUT_PATH=" + upstreamRepo.Path, fmt.Sprintf("VERSION=%s", upstreamRepo.Version)}, nil); err != nil {
			return err
		}
	}
	return nil
}

func pushCommit(upstreamRepo, scratchRepo *source.Repo, baseBranch, command, commitMessage, ref string, gh GithubClient, rnr *exec.Runner, ctlr *source.Controller) error {
	if err := rnr.PushDir(scratchRepo.Path); err != nil {
		return err
	}
	if err := setConfig(rnr); err != nil {
		return err
	}

	if _, err := rnr.Run("git", []string{"add", "."}, nil); err != nil {
		return err
	}
	if _, err := rnr.Run("git", []string{"checkout", "-b", scratchRepo.Branch}, nil); err != nil {
		return err
	}

	if _, err := rnr.Run("git", []string{"commit", "--signoff", "-m", commitMessage}, nil); err != nil {
		if strings.Contains(err.Error(), "nothing to commit, working tree clean") {
			if _, err := rnr.Run("git", []string{"push", ctlr.URL(scratchRepo), scratchRepo.Branch, "-f"}, nil); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	prs, err := gh.GetPullRequests("closed", baseBranch, "updated", "desc")
	if err != nil {
		return err
	}
	var pullRequest github.PullRequest
	for _, pr := range prs {
		if pr.MergeCommitSha == ref {
			pullRequest = pr
			break
		}
	}

	if command == "downstream" && upstreamRepo.Title == "terraform" {
		// Add the changelog entry.
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
	}

	if _, err := rnr.Run("git", []string{"push", ctlr.URL(scratchRepo), scratchRepo.Branch, "-f"}, nil); err != nil {
		return err
	}

	if command == "downstream" {
		fmt.Printf(`Base: %s:%s
Head: %s:%s
`, upstreamRepo.Owner, upstreamRepo.Branch, scratchRepo.Owner, scratchRepo.Branch)
		newPRURL, err := rnr.Run("hub", []string{
			"pull-request",
			"-b",
			fmt.Sprintf("%s:%s",
				upstreamRepo.Owner,
				upstreamRepo.Branch),
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
		if err := gh.MergePullRequest(upstreamRepo.Owner, upstreamRepo.Name, newPRNumber); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(generateDownstreamCmd)
}

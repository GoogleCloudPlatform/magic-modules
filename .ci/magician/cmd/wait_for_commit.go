package cmd

import (
	"fmt"
	"magician/exec"
	"magician/source"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var waitForCommitCmd = &cobra.Command{
	Use:   "wait-for-commit",
	Short: "Wait for the given commit to be ready for downstream push",
	Long: `This command waits until the given commit should be the sync branch's next commit by comparing the history of the base branch and the sync branch. There could be the case when several commits are merged at the same time to the base branch, and they need to be pushed in the same sequence as in base branch to a downstream sync branch.

	The command expects the following as arguments:
	1. SYNC_BRANCH_PREFIX
	2. BASE_BRANCH
	3. SHA

	It then performs the following operations:
	1. Quit if the given sha is already in the sync branch.
	2. Loop until the given sha's parent is equal to the sync branch head.
	`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		syncBranchPrefix := args[0]
		baseBranch := args[1]
		sha := args[2]

		rnr, err := exec.NewRunner()
		if err != nil {
			return fmt.Errorf("error creating Runner: %w", err)
		}

		return execWaitForCommit(syncBranchPrefix, baseBranch, sha, rnr)
	},
}

var waitFunc = func() {
	time.Sleep(5 * time.Second)
}

func execWaitForCommit(syncBranchPrefix, baseBranch, sha string, runner source.Runner) error {
	syncBranch := getSyncBranch(syncBranchPrefix, baseBranch)
	fmt.Println("SYNC_BRANCH: ", syncBranch)

	if syncBranchHasCommit(sha, syncBranch, runner) {
		fmt.Printf("found %s in history of %s - skipping wait\n", sha, syncBranch)
		return nil
	}

	for {
		output, err := gitRevParse("origin/"+syncBranch, runner)
		if err != nil {
			return err
		}
		syncHead := strings.TrimSpace(output)

		output, err = gitRevParse(sha+"~", runner)
		if err != nil {
			return err
		}
		baseParent := strings.TrimSpace(output)
		if syncHead == baseParent {
			return nil
		}
		fmt.Println("sync branch is at: ", syncHead)
		fmt.Println("current commit is: ", sha)

		if _, err := runner.Run("git", []string{"fetch", "origin", syncBranch}, nil); err != nil {
			return err
		}
		waitFunc()
	}
}

func getSyncBranch(syncBranchPrefix, baseBranch string) string {
	if baseBranch == "main" {
		return syncBranchPrefix
	}
	return fmt.Sprintf("%s-%s", syncBranchPrefix, baseBranch)
}

func syncBranchHasCommit(sha, syncBranch string, runner source.Runner) bool {
	if _, err := runner.Run("git", []string{"merge-base", "--is-ancestor", sha, "origin/" + syncBranch}, nil); err == nil {
		return true
	}
	return false
}

func gitRevParse(target string, runner source.Runner) (string, error) {
	return runner.Run("git", []string{"rev-parse", "--short", target}, nil)
}

func init() {
	rootCmd.AddCommand(waitForCommitCmd)
}

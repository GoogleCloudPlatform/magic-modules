package main

import (
	"fmt"
	"os"
)

func main() {
	GITHUB_TOKEN, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("Did not provide GITHUB_TOKEN environment variable")
		os.Exit(1)
	}
	if len(os.Args) <= 7 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	target := os.Args[1]
	fmt.Println("Trigger Target: ", target)

	prNumber := os.Args[2]
	fmt.Println("PR Number: ", prNumber)

	commitSha := os.Args[3]
	fmt.Println("Commit SHA: ", commitSha)

	branchName := os.Args[4]
	fmt.Println("Branch Name: ", branchName)

	headRepoUrl := os.Args[5]
	fmt.Println("Head Repo URL: ", headRepoUrl)

	headBranch := os.Args[6]
	fmt.Println("Head Branch: ", headBranch)

	baseBranch := os.Args[7]
	fmt.Println("Base Branch: ", baseBranch)

	substitutions := map[string]string{
		"BRANCH_NAME":    branchName,
		"_PR_NUMBER":     prNumber,
		"_HEAD_REPO_URL": headRepoUrl,
		"_HEAD_BRANCH":   headBranch,
		"_BASE_BRANCH":   baseBranch,
	}

	author, err := getPullRequestAuthor(prNumber, GITHUB_TOKEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	if target == "auto_run_and_gcbrun" {
		err = requestReviewer(author, prNumber, GITHUB_TOKEN)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	trusted := isTrustedUser(author, GITHUB_TOKEN)

	// auto_run_and_gcbrun will be run on every commit or /gcbrun: only trigger builds for trusted users
	// gcbrun_only will be run on every /gcbrun: only trigger builds for untrusted users (because trusted users will be handled by auto_run_and_gcbrun)
	if (target == "auto_run_and_gcbrun" && trusted) || (target == "gcbrun_only" && !trusted) {
		err = triggerMMPresubmitRuns("graphite-docker-images", "magic-modules", commitSha, substitutions)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

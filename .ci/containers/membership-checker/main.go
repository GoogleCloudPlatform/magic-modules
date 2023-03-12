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
	fmt.Println("PR Number: ", target)

	prNumber := os.Args[2]
	fmt.Println("PR Number: ", prNumber)

	commitSha := os.Args[3]
	fmt.Println("Commit SHA: ", commitSha)

	branchName := os.Args[4]
	fmt.Println("Branch Name: ", branchName)

	headRepoUrl := os.Args[5]
	fmt.Println("Branch Name: ", headRepoUrl)

	headBranch := os.Args[6]
	fmt.Println("Branch Name: ", headBranch)

	baseBranch := os.Args[7]
	fmt.Println("Branch Name: ", baseBranch)

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

	if target == "check_auto_run_contributor" {
		err = reviewerAssignment(author, prNumber, GITHUB_TOKEN)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	trusted := isTrustedUser(author, GITHUB_TOKEN)

	if (target == "check_auto_run_contributor" && trusted) || (target == "check_community_contributor" && !trusted) {
		err = triggerMMPresubmitRuns("graphite-docker-images", "magic-modules", commitSha, substitutions)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

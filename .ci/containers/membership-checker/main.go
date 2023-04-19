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

	projectId := "graphite-docker-images"
	repoName := "magic-modules"

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
		os.Exit(1)
	}

	if target == "auto_run" {
		err = requestReviewer(author, prNumber, GITHUB_TOKEN)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	trusted := isTrustedUser(author, GITHUB_TOKEN)

	// auto_run(contributor-membership-checker) will be run on every commit or /gcbrun:
	// only triggers builds for trusted users

	// needs_approval(community-checker) will be run after approval:
	// 1. will be auto approved (by contributor-membership-checker) for trusted users
	// 2. needs approval from team reviewer via cloud build for untrusted users
	// 3. only triggers build for untrusted users (because trusted users will be handled by auto_run)
	if (target == "auto_run" && trusted) || (target == "needs_approval" && !trusted) {
		err = triggerMMPresubmitRuns(projectId, repoName, commitSha, substitutions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// in contributor-membership-checker job:
	// 1. auto approve community-checker run for trusted users
	// 2. add awaiting-approval label to external contributor PRs
	if target == "auto_run" {
		if trusted {
			err = approveCommunityChecker(prNumber, projectId, commitSha)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			addAwaitingApprovalLabel(prNumber, GITHUB_TOKEN)
			postAwaitingApprovalBuildLink(prNumber, GITHUB_TOKEN, projectId, commitSha)
		}
	}

	// in community-checker job:
	// remove awaiting-approval label from external contributor PRs
	if target == "needs_approval" {
		removeAwaitingApprovalLabel(prNumber, GITHUB_TOKEN)
	}
}

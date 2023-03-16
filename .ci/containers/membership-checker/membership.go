package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

var (
	// TODO: add unit tests to ensure that
	// 1. people in vacationList are also  in reviewers
	// 2. people in trustedContributors are not in reviewers

	// This is for the random-assignee rotation.
	reviewerRotation = []string{"megan07", "slevenick", "c2thorn", "rileykarson", "melinath", "ScottSuarez", "shuyama1", "SarahFrench", "roaks3", "zli82016", "trodge", "hao-nan-li"}

	// This is for new team members who are onboarding
	trustedContributors = []string{"NickElliot"}

	// This is for reviewers who are "on vacation" will not receive new review assignments but will still receive re-requests for assigned PRs.
	vacationList = []string{"zli82016"}
)

// Check if a user is team member to not request a random reviewer
func isTeamMember(author string) bool {
	return slices.Contains(reviewerRotation, author) || slices.Contains(trustedContributors, author)
}

func isTeamReviewer(reviewer string) bool {
	return slices.Contains(reviewerRotation, reviewer)
}

// Check if a user is safe to run tests automatically
func isTrustedUser(author, GITHUB_TOKEN string) bool {
	if isTeamMember(author) {
		fmt.Println("User is on the list")
		return true
	}

	if isOrgMember(author, "GoogleCloudPlatform", GITHUB_TOKEN) {
		fmt.Println("User is a GCP org member")
		return true
	}

	if isOrgMember(author, "googlers", GITHUB_TOKEN) {
		fmt.Println("User is a googlers org member")
		return true
	}

	return false
}

func isOrgMember(author, org, GITHUB_TOKEN string) bool {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, author)
	res, _ := requestCall(url, "GET", GITHUB_TOKEN, nil, nil)

	return res != 404
}

func getAvailableReviewers()

func getRandomReviewer() string {
	rand.Seed(time.Now().Unix())
	reviewer := reviewerRotation[rand.Intn(len(reviewerRotation))]
	return reviewer, nil
}

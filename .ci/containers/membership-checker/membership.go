package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

var (
	// This is for the random-assignee rotation.
	reviewerRotation = []string{
		"megan07",
		"slevenick",
		"c2thorn",
		"rileykarson",
		"melinath",
		"ScottSuarez",
		"shuyama1",
		"SarahFrench",
		"roaks3",
		"zli82016",
		"trodge",
		"hao-nan-li",
		"NickElliot",
	}

	// This is for new team members who are onboarding
	trustedContributors = []string{}

	// This is for reviewers who are "on vacation": will not receive new review assignments but will still receive re-requests for assigned PRs.
	onVacationReviewers = []string{"ScottSuarez"}
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
		fmt.Println("User is a team member")
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

func getRandomReviewer() string {
	availableReviewers := removes(reviewerRotation, onVacationReviewers)
	rand.Seed(time.Now().UnixNano())
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

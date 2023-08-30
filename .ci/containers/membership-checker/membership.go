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
	onVacationReviewers = []string{
		"zli82016",
	}
)

type userType int64

const (
	communityUserType userType = iota
	googlerUserType
	coreContributorUserType
)

func (ut userType) String() string {
	switch ut {
	case googlerUserType:
		return "Googler"
	case coreContributorUserType:
		return "Core Contributor"
	default:
		return "Community Contributor"
	}
}

// Check if a user is team member to not request a random reviewer
func isTeamMember(author string) bool {
	return slices.Contains(reviewerRotation, author) || slices.Contains(trustedContributors, author)
}

func isTeamReviewer(reviewer string) bool {
	return slices.Contains(reviewerRotation, reviewer)
}

func getUserType(user, GITHUB_TOKEN string) userType {
	if isTeamMember(user) {
		fmt.Println("User is a team member")
		return coreContributorUserType
	}

	if isOrgMember(user, "GoogleCloudPlatform", GITHUB_TOKEN) {
		fmt.Println("User is a GCP org member")
		return googlerUserType
	}

	if isOrgMember(user, "googlers", GITHUB_TOKEN) {
		fmt.Println("User is a googlers org member")
		return googlerUserType
	}

	return communityUserType
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

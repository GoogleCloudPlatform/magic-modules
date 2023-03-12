package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	// This is where you add users who do not need to have an assignee chosen for them
	noAssigneeList = []string{"megan07", "slevenick", "c2thorn", "rileykarson", "melinath", "ScottSuarez", "shuyama1", "SarahFrench", "roaks3", "zli82016", "trodge", "hao-nan-li"}

	// This is where you add people to the random-assignee rotation.
	reviewerRotationList = []string{"megan07", "slevenick", "c2thorn", "rileykarson", "melinath", "ScottSuarez", "shuyama1", "SarahFrench", "roaks3", "zli82016", "trodge", "hao-nan-li"}

	// This is where your add reviewers who will be re-requested reviews when PR authors make new commits
	// This should mostly be identical to reviewerRotationList, but if someone is temporally removed from assignee list, they can still be on this list to keep getting review alert for current PRs
	rerequestReviewerRotationList = []string{"megan07", "slevenick", "c2thorn", "rileykarson", "melinath", "ScottSuarez", "shuyama1", "SarahFrench", "roaks3", "zli82016", "trodge", "hao-nan-li"}

	// This is where you add trusted users (besides the users who are already in noAssigneeList) that do not need a '/gcbrun' comment from team to run tests
	trustedMemberList = []string{}
)

func isNoAssigneeUser(author string) bool {
	return onList(author, noAssigneeList)
}

func isTeamReviewer(reviewer string) bool {
	return onList(reviewer, rerequestReviewerRotationList)
}

// Check if a user is safe to run tests automatically
func isTrustedUser(author, GITHUB_TOKEN string) bool {
	if isTrustedMember(author) {
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

func isTrustedMember(author string) bool {
	return onList(author, noAssigneeList) || onList(author, trustedMemberList)
}

func isOrgMember(author, org, GITHUB_TOKEN string) bool {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, author)
	res, _ := requestCall(url, "GET", GITHUB_TOKEN, nil, nil)

	return res != 404
}

func getRamdomReviewer() string {
	assignee := reviewerRotationList[rand.Intn(len(reviewerRotationList))]
	rand.Seed(time.Now().Unix())
	return assignee
}

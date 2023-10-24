package github

import (
	"fmt"
	utils "magician/utility"
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
	onVacationReviewers = []string{}
)

type UserType int64

const (
	CommunityUserType UserType = iota
	GooglerUserType
	CoreContributorUserType
)

func (ut UserType) String() string {
	switch ut {
	case GooglerUserType:
		return "Googler"
	case CoreContributorUserType:
		return "Core Contributor"
	default:
		return "Community Contributor"
	}
}

func (gh *github) GetUserType(user string) UserType {
	if isTeamMember(user, gh.token) {
		fmt.Println("User is a team member")
		return CoreContributorUserType
	}

	if isOrgMember(user, "GoogleCloudPlatform", gh.token) {
		fmt.Println("User is a GCP org member")
		return GooglerUserType
	}

	if isOrgMember(user, "googlers", gh.token) {
		fmt.Println("User is a googlers org member")
		return GooglerUserType
	}

	return CommunityUserType
}

// Check if a user is team member to not request a random reviewer
func isTeamMember(author, githubToken string) bool {
	return slices.Contains(reviewerRotation, author) || slices.Contains(trustedContributors, author)
}

func IsTeamReviewer(reviewer string) bool {
	return slices.Contains(reviewerRotation, reviewer)
}

func isOrgMember(author, org, githubToken string) bool {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, author)
	res, _ := utils.RequestCall(url, "GET", githubToken, nil, nil)

	return res != 404
}

func GetRandomReviewer() string {
	availableReviewers := utils.Removes(reviewerRotation, onVacationReviewers)
	rand.Seed(time.Now().UnixNano())
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

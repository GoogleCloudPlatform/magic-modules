/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
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
		"BBBmau",
	}

	// This is for new team members who are onboarding
	trustedContributors = []string{}

	// This is for reviewers who are "on vacation": will not receive new review assignments but will still receive re-requests for assigned PRs.
	onVacationReviewers = []string{
		"zli82016",
		"NickElliot",
		"ScottSuarez",
	}
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

func (gh *Client) GetUserType(user string) UserType {
	if IsCoreContributor(user) {
		fmt.Println("User is a core contributor")
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
func IsCoreContributor(user string) bool {
	return slices.Contains(reviewerRotation, user) || slices.Contains(trustedContributors, user)
}

func IsCoreReviewer(reviewer string) bool {
	return slices.Contains(reviewerRotation, reviewer)
}

func isOrgMember(author, org, githubToken string) bool {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, author)
	err := utils.RequestCall(url, "GET", githubToken, nil, nil)

	if err != nil {
		return false
	}
	return true
}

func GetRandomReviewer() string {
	availableReviewers := AvailableReviewers()
	rand.Seed(time.Now().UnixNano())
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

func AvailableReviewers() []string {
	return utils.Removes(reviewerRotation, onVacationReviewers)
}

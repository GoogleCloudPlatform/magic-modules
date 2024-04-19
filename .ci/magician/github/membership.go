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
	// User can specify the time zone like this, and following the example below:
	pdtLoc, _           = time.LoadLocation("America/Los_Angeles")
	onVacationReviewers = []onVacationReviewer{
		// Example: taking vacation from 2024-03-28 to 2024-04-02 in pdt time zone.
		// both ends are inclusive:
		// {
		// 	id:        "xyz",
		// 	startDate: newDate(2024, 3, 28, pdtLoc),
		// 	endDate:   newDate(2024, 4, 2, pdtLoc),
		// },
		{
			id:        "hao-nan-li",
			startDate: newDate(2024, 4, 11, pdtLoc),
			endDate:   newDate(2024, 6, 14, pdtLoc),
		},
		{
	  		id: "slevenick",
	 		startDate: newDate(2024, 4, 20, pdtLoc),
	  		endDate: newDate(2024, 4, 27, pdtLoc),
		},
	}
)

type UserType int64

type date struct {
	year  int
	month int
	day   int
	loc   *time.Location
}

type onVacationReviewer struct {
	id        string
	startDate date
	endDate   date
}

func newDate(year, month, day int, loc *time.Location) date {
	return date{
		year:  year,
		month: month,
		day:   day,
		loc:   loc,
	}
}

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

	return err == nil
}

func GetRandomReviewer() string {
	availableReviewers := AvailableReviewers()
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

func AvailableReviewers() []string {
	return available(time.Now(), reviewerRotation, onVacationReviewers)
}

func available(nowTime time.Time, allReviewers []string, vacationList []onVacationReviewer) []string {
	onVacationList := onVacation(nowTime, vacationList)
	return utils.Removes(allReviewers, onVacationList)
}

func onVacation(nowTime time.Time, vacationList []onVacationReviewer) []string {
	var onVacationList []string
	for _, reviewer := range vacationList {
		start := time.Date(reviewer.startDate.year, time.Month(reviewer.startDate.month), reviewer.startDate.day, 0, 0, 0, 0, reviewer.startDate.loc)
		end := time.Date(reviewer.endDate.year, time.Month(reviewer.endDate.month), reviewer.endDate.day, 0, 0, 0, 0, reviewer.endDate.loc).AddDate(0, 0, 1).Add(-1 * time.Millisecond)
		if nowTime.Before(start) || nowTime.After(end) {
			continue
		}
		onVacationList = append(onVacationList, reviewer.id)
	}
	return onVacationList
}

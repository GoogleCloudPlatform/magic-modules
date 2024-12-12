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

	"golang.org/x/exp/maps"
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
	_, isTrustedContributor := trustedContributors[user]
	return IsCoreReviewer(user) || isTrustedContributor
}

func IsCoreReviewer(user string) bool {
	_, isCoreReviewer := reviewerRotation[user]
	return isCoreReviewer
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

// Return a random reviewer other than the old reviewer
func GetNewRandomReviewer(oldReviewer string) string {
	availableReviewers := AvailableReviewers()
	availableReviewers = utils.Removes(availableReviewers, []string{oldReviewer})
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

func AvailableReviewers() []string {
	return available(time.Now(), maps.Keys(reviewerRotation), onVacationReviewers)
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

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
	"slices"
	"time"

	"golang.org/x/exp/maps"
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

	if gh.IsTeamMember("GoogleCloudPlatform", "terraform", user) {
		fmt.Printf("Debug test --- User '%s' is an active member of the 'terraform' team in 'GoogleCloudPlatform' organization\n", user)
	} else {
		fmt.Printf("Debug test --- User '%s' is not an active member of the 'terraform' team in 'GoogleCloudPlatform' organization\n", user)
	}

	if gh.IsOrgMember(user, "GoogleCloudPlatform") {
		fmt.Println("User is a GCP org member")
		return GooglerUserType
	}

	if gh.IsOrgMember(user, "googlers") {
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

// GetRandomReviewer returns a random available reviewer (optionally excluding some people from the reviewer pool)
func GetRandomReviewer(excludedReviewers []string) string {
	availableReviewers := AvailableReviewers(excludedReviewers)
	reviewer := availableReviewers[rand.Intn(len(availableReviewers))]
	return reviewer
}

func AvailableReviewers(excludedReviewers []string) []string {
	return available(time.Now(), reviewerRotation, excludedReviewers)
}

func available(nowTime time.Time, reviewerRotation map[string]ReviewerConfig, excludedReviewers []string) []string {
	excludedReviewers = append(excludedReviewers, onVacation(nowTime, reviewerRotation)...)
	ret := utils.Removes(maps.Keys(reviewerRotation), excludedReviewers)
	slices.Sort(ret)
	return ret
}

func onVacation(nowTime time.Time, reviewerRotation map[string]ReviewerConfig) []string {
	var onVacationList []string
	for reviewer, config := range reviewerRotation {
		for _, v := range config.vacations {
			if nowTime.Before(v.GetStart(config.timezone)) || nowTime.After(v.GetEnd(config.timezone)) {
				continue
			}
			onVacationList = append(onVacationList, reviewer)
		}
	}
	return onVacationList
}

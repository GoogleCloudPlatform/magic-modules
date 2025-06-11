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
		fmt.Println("User is an active member of the 'terraform' team in 'GoogleCloudPlatform' organization")
		return GooglerUserType
	} else {
		fmt.Printf("User '%s' is not an active member of the 'terraform' team in 'GoogleCloudPlatform' organization\n", user)
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
	return reviewerRotation.isCoreReviewer(user)
}

func GetRandomReviewer(excludedReviewers []string) string {
	return reviewerRotation.getRandomReviewer(excludedReviewers)
}

func AvailableReviewers(excludedReviewers []string) []string {
	return reviewerRotation.availableReviewers(excludedReviewers)
}

// unused unless exporting hardcoded reviewer rotation to yaml
func WriteReviewerRotation() ([]byte, error) {
	return reviewerRotation.write()
}

func ReadReviewerRotation(data []byte) error {
	reviewerRotation.setStartEnd()
	return reviewerRotation.read(data)
}

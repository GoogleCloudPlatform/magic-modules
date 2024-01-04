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
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestChooseCoreReviewers(t *testing.T) {
	cases := map[string]struct {
		RequestedReviewers                               []User
		PreviousReviewers                                []User
		ExpectReviewersFromList, ExpectSpecificReviewers []string
		ExpectPrimaryReviewer                            bool
	}{
		"no previous review requests assigns new reviewer from team": {
			RequestedReviewers:      []User{},
			PreviousReviewers:       []User{},
			ExpectReviewersFromList: utils.Removes(reviewerRotation, onVacationReviewers),
			ExpectPrimaryReviewer:   true,
		},
		"requested reviewer from team means that primary reviewer was already selected": {
			RequestedReviewers:    []User{User{Login: reviewerRotation[0]}},
			PreviousReviewers:     []User{},
			ExpectPrimaryReviewer: false,
		},
		"requested off-team reviewer does not mean that primary reviewer was already selected": {
			RequestedReviewers:    []User{User{Login: "foobar"}},
			PreviousReviewers:     []User{},
			ExpectPrimaryReviewer: true,
		},
		"previously involved team member reviewers should have review requested and mean that primary reviewer was already selected": {
			RequestedReviewers:      []User{},
			PreviousReviewers:       []User{User{Login: reviewerRotation[0]}},
			ExpectSpecificReviewers: []string{reviewerRotation[0]},
			ExpectPrimaryReviewer:   false,
		},
		"previously involved reviewers that are not team members are ignored": {
			RequestedReviewers:      []User{},
			PreviousReviewers:       []User{User{Login: "foobar"}},
			ExpectReviewersFromList: utils.Removes(reviewerRotation, onVacationReviewers),
			ExpectPrimaryReviewer:   true,
		},
		"only previously involved team member reviewers will have review requested": {
			RequestedReviewers:      []User{},
			PreviousReviewers:       []User{User{Login: reviewerRotation[0]}, User{Login: "foobar"}, User{Login: reviewerRotation[1]}},
			ExpectSpecificReviewers: []string{reviewerRotation[0], reviewerRotation[1]},
			ExpectPrimaryReviewer:   false,
		},
		"primary reviewer will not have review requested even if other team members previously reviewed": {
			RequestedReviewers:      []User{User{Login: reviewerRotation[1]}},
			PreviousReviewers:       []User{User{Login: reviewerRotation[0]}},
			ExpectSpecificReviewers: []string{reviewerRotation[0]},
			ExpectPrimaryReviewer:   false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			reviewers, primaryReviewer := ChooseCoreReviewers(tc.RequestedReviewers, tc.PreviousReviewers)
			if tc.ExpectPrimaryReviewer && primaryReviewer == "" {
				t.Error("wanted primary reviewer to be returned; got none")
			}
			if !tc.ExpectPrimaryReviewer && primaryReviewer != "" {
				t.Errorf("wanted no primary reviewer; got %s", primaryReviewer)
			}
			if len(tc.ExpectReviewersFromList) > 0 {
				for _, reviewer := range reviewers {
					if !slices.Contains(tc.ExpectReviewersFromList, reviewer) {
						t.Errorf("wanted reviewer %s to be in list %v but they were not", reviewer, tc.ExpectReviewersFromList)
					}
				}
			}
			if len(tc.ExpectSpecificReviewers) > 0 {
				if !slices.Equal(reviewers, tc.ExpectSpecificReviewers) {
					t.Errorf("wanted reviewers to be %v; instead got %v", tc.ExpectSpecificReviewers, reviewers)
				}
			}
		})
	}
}

func TestFormatReviewerComment(t *testing.T) {
	cases := map[string]struct {
		Reviewer       string
		AuthorUserType UserType
		Trusted        bool
	}{
		"community contributor": {
			Reviewer:       "foobar",
			AuthorUserType: CommunityUserType,
			Trusted:        false,
		},
		"googler": {
			Reviewer:       "foobar",
			AuthorUserType: GooglerUserType,
			Trusted:        true,
		},
		"core contributor": {
			Reviewer:       "foobar",
			AuthorUserType: CoreContributorUserType,
			Trusted:        true,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			comment := FormatReviewerComment(tc.Reviewer, tc.AuthorUserType, tc.Trusted)
			t.Log(comment)
			if !strings.Contains(comment, fmt.Sprintf("@%s", tc.Reviewer)) {
				t.Errorf("wanted comment to contain @%s; does not.", tc.Reviewer)
			}
			if !strings.Contains(comment, tc.AuthorUserType.String()) {
				t.Errorf("wanted comment to contain user type (%s); does not.", tc.AuthorUserType.String())
			}
			if strings.Contains(comment, fmt.Sprintf("~%s~", tc.AuthorUserType.String())) {
				t.Errorf("wanted user type (%s) in comment to not be crossed out, but it is", tc.AuthorUserType.String())
			}
			for _, ut := range []UserType{CommunityUserType, GooglerUserType, CoreContributorUserType} {
				if ut != tc.AuthorUserType && !strings.Contains(comment, fmt.Sprintf("~%s~", ut.String())) {
					t.Errorf("wanted other user type (%s) in comment to be crossed out, but it is not", ut)
				}
			}

			if tc.Trusted && !strings.Contains(comment, "Tests will run automatically") {
				t.Errorf("wanted comment to say tests will run automatically; does not")
			}
			if !tc.Trusted && !strings.Contains(comment, "Tests will require approval") {
				t.Errorf("wanted comment to say tests will require approval; does not")
			}
		})

	}
}

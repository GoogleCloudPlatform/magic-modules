package github

import (
	"fmt"
	utils "magician/utility"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestChooseReviewers(t *testing.T) {
	cases := map[string]struct {
		FirstRequestedReviewer                           string
		PreviouslyInvolvedReviewers                      []string
		ExpectReviewersFromList, ExpectSpecificReviewers []string
		ExpectPrimaryReviewer                            bool
	}{
		"no previous review requests assigns new reviewer from team": {
			FirstRequestedReviewer:      "",
			PreviouslyInvolvedReviewers: []string{},
			ExpectReviewersFromList:     utils.Removes(reviewerRotation, onVacationReviewers),
			ExpectPrimaryReviewer:       true,
		},
		"first requested reviewer means that primary reviewer was already selected": {
			FirstRequestedReviewer:      "foobar",
			PreviouslyInvolvedReviewers: []string{},
			ExpectPrimaryReviewer:       false,
		},
		"previously involved team member reviewers should have review requested and mean that primary reviewer was already selected": {
			FirstRequestedReviewer:      "",
			PreviouslyInvolvedReviewers: []string{reviewerRotation[0]},
			ExpectSpecificReviewers:     []string{reviewerRotation[0]},
			ExpectPrimaryReviewer:       false,
		},
		"previously involved reviewers that are not team members are ignored": {
			FirstRequestedReviewer:      "",
			PreviouslyInvolvedReviewers: []string{"foobar"},
			ExpectReviewersFromList:     utils.Removes(reviewerRotation, onVacationReviewers),
			ExpectPrimaryReviewer:       true,
		},
		"only previously involved team member reviewers will have review requested": {
			FirstRequestedReviewer:      "",
			PreviouslyInvolvedReviewers: []string{reviewerRotation[0], "foobar", reviewerRotation[1]},
			ExpectSpecificReviewers:     []string{reviewerRotation[0], reviewerRotation[1]},
			ExpectPrimaryReviewer:       false,
		},
		"primary reviewer will not have review requested even if other team members previously reviewed": {
			FirstRequestedReviewer:      reviewerRotation[1],
			PreviouslyInvolvedReviewers: []string{reviewerRotation[0]},
			ExpectSpecificReviewers:     []string{reviewerRotation[0]},
			ExpectPrimaryReviewer:       false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			reviewers, primaryReviewer := ChooseReviewers(tc.FirstRequestedReviewer, tc.PreviouslyInvolvedReviewers)
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

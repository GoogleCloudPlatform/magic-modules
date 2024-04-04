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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
)

func TestTrustedContributors(t *testing.T) {
	for _, member := range trustedContributors {
		if slices.Contains(reviewerRotation, member) {
			t.Fatalf(`%v should not be on reviewerRotation list`, member)
		}
	}
}

func TestOnVacationReviewers(t *testing.T) {
	for _, member := range onVacationReviewers {
		if !slices.Contains(reviewerRotation, member.id) {
			t.Fatalf(`%v is not on reviewerRotation list`, member)
		}
	}
}

func TestAvailableReviewers(t *testing.T) {
	tests := []struct {
		name       string
		rotation   []string
		onVacation []onVacationReviewer
		timeNow    time.Time
		want       []string
	}{
		{
			name:     "reviewers on vacation are excluded",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
					endDate:   time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			timeNow: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			want:    []string{"id1"},
		},
		{
			name:     "reviewers are included after vacation ends",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
					endDate:   time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			timeNow: time.Date(2024, 4, 2, 10, 0, 0, 0, time.UTC),
			want:    []string{"id1", "id2"},
		},
		{
			name:     "reviewers are included before vacation starts",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
					endDate:   time.Date(2024, 4, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			timeNow: time.Date(2024, 3, 28, 10, 0, 0, 0, time.UTC),
			want:    []string{"id1", "id2"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			origRotation := reviewerRotation
			origOnVacation := onVacationReviewers
			reviewerRotation = test.rotation
			onVacationReviewers = test.onVacation
			defer func() {
				reviewerRotation = origRotation
				onVacationReviewers = origOnVacation
			}()

			got := AvailableReviewers(test.timeNow)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("AvailableReviewers() got diff: %s", diff)
			}
		})
	}

}

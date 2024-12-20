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
)

func TestTrustedContributors(t *testing.T) {
	for member, _ := range trustedContributors {
		if IsCoreReviewer(member) {
			t.Fatalf(`%v should not be on reviewerRotation list`, member)
		}
	}
}

func TestOnVacationReviewers(t *testing.T) {
	for _, member := range onVacationReviewers {
		if !IsCoreReviewer(member.id) {
			t.Fatalf(`%v is not on reviewerRotation list`, member)
		}
	}
}

func TestAvailable(t *testing.T) {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		rotation   []string
		onVacation []onVacationReviewer
		timeNow    time.Time
		want       []string
	}{
		{
			name:     "reviewers on vacation start date are excluded",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: newDate(2024, 3, 29, time.UTC),
					endDate:   newDate(2024, 4, 2, time.UTC),
				},
			},
			timeNow: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			want:    []string{"id1"},
		},
		{
			name:     "reviewers on vacation end date are excluded",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: newDate(2024, 3, 29, time.UTC),
					endDate:   newDate(2024, 4, 2, time.UTC),
				},
			},
			timeNow: time.Date(2024, 4, 2, 10, 0, 0, 0, time.UTC),
			want:    []string{"id1"},
		},
		{
			name:     "reviewers are included after vacation ends",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: newDate(2024, 3, 29, time.UTC),
					endDate:   newDate(2024, 4, 2, time.UTC),
				},
			},
			timeNow: time.Date(2024, 4, 3, 0, 0, 0, 0, time.UTC),
			want:    []string{"id1", "id2"},
		},
		{
			name:     "reviewers are included before vacation starts",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: newDate(2024, 3, 29, time.UTC),
					endDate:   newDate(2024, 4, 2, time.UTC),
				},
			},
			timeNow: time.Date(2024, 3, 28, 23, 0, 0, 0, time.UTC),
			want:    []string{"id1", "id2"},
		},
		{
			name:     "reviewers are excluded since vacation still not ends in the specified time zone",
			rotation: []string{"id1", "id2"},
			onVacation: []onVacationReviewer{
				{
					id:        "id2",
					startDate: newDate(2024, 3, 29, la),
					endDate:   newDate(2024, 4, 2, la),
				},
			},
			// it's still 2024-04-02 in LA time zone
			timeNow: time.Date(2024, 4, 3, 0, 0, 0, 0, newYork),
			want:    []string{"id1"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := available(test.timeNow, test.rotation, test.onVacation)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("available(%v, %v, %v) got diff: %s", test.timeNow, test.rotation, test.onVacation, diff)
			}
		})
	}

}

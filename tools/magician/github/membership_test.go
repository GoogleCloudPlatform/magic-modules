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

func TestAvailable(t *testing.T) {
	// Double-check that timezones are loadable first.
	_, err := time.LoadLocation("US/Eastern")
	if err != nil {
		t.Fatal(err)
	}
	_, err = time.LoadLocation("US/Pacific")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name              string
		rotation          map[string]ReviewerConfig
		timeNow           time.Time
		excludedReviewers []string
		want              []string
	}{
		{
			name: "reviewers on vacation start date are excluded",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					timezone: time.UTC,
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
					},
				},
			},
			timeNow: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			want:    []string{"id1"},
		},
		{
			name: "reviewers on vacation end date are excluded",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					timezone: time.UTC,
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
					},
				},
			},
			timeNow: time.Date(2024, 4, 2, 10, 0, 0, 0, time.UTC),
			want:    []string{"id1"},
		},
		{
			name: "reviewers are included after vacation ends",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					timezone: time.UTC,
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
					},
				},
			},
			timeNow: time.Date(2024, 4, 3, 0, 0, 0, 0, time.UTC),
			want:    []string{"id1", "id2"},
		},
		{
			name: "reviewers are included before vacation starts",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					timezone: time.UTC,
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
					},
				},
			},
			timeNow: time.Date(2024, 3, 28, 23, 0, 0, 0, time.UTC),
			want:    []string{"id1", "id2"},
		},
		{
			name: "reviewers are excluded if vacation has not ended in the specified time zone",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
					},
				},
			},
			// it's still 2024-04-02 in Pacific time zone
			timeNow: time.Date(2024, 4, 3, 0, 0, 0, 0, usEastern),
			want:    []string{"id1"},
		},
		{
			name: "included before vacations",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
						{
							startDate: newDate(2024, 5, 2),
							endDate:   newDate(2024, 5, 5),
						},
					},
				},
			},
			timeNow: time.Date(2024, 3, 1, 0, 0, 0, 0, usPacific),
			want:    []string{"id1", "id2"},
		},
		{
			name: "excluded during first vacation",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
						{
							startDate: newDate(2024, 5, 2),
							endDate:   newDate(2024, 5, 5),
						},
					},
				},
			},
			timeNow: time.Date(2024, 4, 1, 0, 0, 0, 0, usPacific),
			want:    []string{"id1"},
		},
		{
			name: "included between vacations",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
						{
							startDate: newDate(2024, 5, 2),
							endDate:   newDate(2024, 5, 5),
						},
					},
				},
			},
			timeNow: time.Date(2024, 4, 4, 0, 0, 0, 0, usPacific),
			want:    []string{"id1", "id2"},
		},
		{
			name: "excluded during second vacation",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
						{
							startDate: newDate(2024, 5, 2),
							endDate:   newDate(2024, 5, 5),
						},
					},
				},
			},
			timeNow: time.Date(2024, 5, 3, 0, 0, 0, 0, usPacific),
			want:    []string{"id1"},
		},
		{
			name: "included after vacations",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {
					vacations: []Vacation{
						{
							startDate: newDate(2024, 3, 29),
							endDate:   newDate(2024, 4, 2),
						},
						{
							startDate: newDate(2024, 5, 2),
							endDate:   newDate(2024, 5, 5),
						},
					},
				},
			},
			timeNow: time.Date(2024, 6, 1, 0, 0, 0, 0, usPacific),
			want:    []string{"id1", "id2"},
		},
		{
			name: "explicitly excluded reviewers",
			rotation: map[string]ReviewerConfig{
				"id1": {vacations: []Vacation{}},
				"id2": {vacations: []Vacation{}},
			},
			excludedReviewers: []string{"id2"},
			want:              []string{"id1"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := available(test.timeNow, test.rotation, test.excludedReviewers)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("available(%v, %v, %v) got diff: %s", test.timeNow, test.rotation, test.excludedReviewers, diff)
			}
		})
	}

}

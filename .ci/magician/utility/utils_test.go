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
package utility

import (
	"reflect"
	"testing"
	"time"
)

func TestRemovesList(t *testing.T) {
	cases := map[string]struct {
		Original, Removal, Expected []string
	}{
		"Remove list": {
			Original: []string{"a", "b", "c"},
			Removal:  []string{"b"},
			Expected: []string{"a", "c"},
		},
		"Remove case sensitive elements": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"b", "c", "A"},
			Expected: []string{"a", "B"},
		},
		"Remove nonexistent elements": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"a", "A", "d"},
			Expected: []string{"b", "c", "B"},
		},
		"Remove none": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{},
			Expected: []string{"a", "b", "c", "A", "B"},
		},
		"Remove all": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"a", "b", "c", "A", "B"},
			Expected: []string{},
		},
		"Remove all and extra nonexistent elements": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"a", "b", "c", "A", "B", "D"},
			Expected: []string{},
		},
	}
	for tn, tc := range cases {
		result := Removes(tc.Original, tc.Removal)
		if !reflect.DeepEqual(result, tc.Expected) {
			t.Errorf("bad: %s, '%s' removes '%s' expect result: %s, but got: %s", tn, tc.Original, tc.Removal, tc.Expected, result)
		}
	}
}

// Test the shouldRetry function
func TestShouldRetry(t *testing.T) {
	config := defaultRetryConfig()

	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"Should retry on 500", 500, true},
		{"Should retry on 503", 503, true},
		{"Should not retry on 200", 200, false},
		{"Should not retry on 400", 400, false},
		{"Should not retry on 404", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetry(tt.statusCode, config); got != tt.want {
				t.Errorf("shouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test the calculateBackoff function
func TestCalculateBackoff(t *testing.T) {
	config := retryConfig{
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		BackoffFactor:  2.0,
	}

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{"First attempt", 0, 100 * time.Millisecond},
		{"Second attempt", 1, 200 * time.Millisecond},
		{"Third attempt", 2, 400 * time.Millisecond},
		{"Fourth attempt", 3, 800 * time.Millisecond},
		{"Fifth attempt (capped)", 4, 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateBackoff(tt.attempt, config); got != tt.want {
				t.Errorf("calculateBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}

package main

import (
	"regexp"
	"testing"
)


func Test_listDifference(t *testing.T) {
	testCases := map[string]struct {
		a          []string
		b          []string
		expectDiff bool
		errorRegex *regexp.Regexp
	}{
		"detects when lists match": {
			a: []string{"a", "c", "b"},
			b: []string{"a", "b", "c"},
		},
		"detects when items from list A is missing items present in list B - 1 missing": {
			a:          []string{"a", "b"},
			b:          []string{"a", "c", "b"},
			expectDiff: true,
			errorRegex: regexp.MustCompile("[c]"),
		},
		"detects when items from list A is missing items present in list B - 2 missing": {
			a:          []string{"b"},
			b:          []string{"a", "c", "b"},
			expectDiff: true,
			errorRegex: regexp.MustCompile("[a, c]"),
		},
		"doesn't detect differences if list A is a superset of list B": {
			a:          []string{"a", "b", "c"},
			b:          []string{"a", "c"},
			expectDiff: false,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			err := listDifference(tc.a, tc.b)
			if !tc.expectDiff && (err != nil) {
				t.Fatalf("saw an unexpected diff error: %s", err)
			}
			if tc.expectDiff && (err == nil) {
				t.Fatalf("expected a diff error but saw none")
			}
			if !tc.expectDiff && err == nil {
				// Stop assertions in no error cases
				return
			}
			if !tc.errorRegex.MatchString(err.Error()) {
				t.Fatalf("expected diff error to contain a match for regex %s, error string: %s", tc.errorRegex.String(), err)
			}
		})
	}
}

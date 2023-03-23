package main

import (
	"reflect"
	"testing"

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
		if !slices.Contains(reviewerRotation, member) {
			t.Fatalf(`%v is not on reviewerRotation list`, member)
		}
	}
}

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
		result := removes(tc.Original, tc.Removal)
		if !reflect.DeepEqual(result, tc.Expected) {
			t.Errorf("bad: %s, '%s' removes '%s' expect result: %s, but got: %s", tn, tc.Original, tc.Removal, tc.Expected, result)
		}
	}
}

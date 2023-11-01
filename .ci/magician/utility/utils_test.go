package utility

import (
	"reflect"
	"testing"
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

package diff

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func removeZeroPadding(s string) string {
	return strings.ReplaceAll(s, ".0", "")
}

func TestRemoveZeroPadding(t *testing.T) {
	for _, tc := range []struct {
		name       string
		zeroPadded string
		expected   string
	}{
		{
			name:       "no zeroes",
			zeroPadded: "a.b.c",
			expected:   "a.b.c",
		},
		{
			name:       "one zero",
			zeroPadded: "a.0.b.c",
			expected:   "a.b.c",
		},
		{
			name:       "two zeroes",
			zeroPadded: "a.0.b.0.c",
			expected:   "a.b.c",
		},
	} {
		if got := removeZeroPadding(tc.zeroPadded); got != tc.expected {
			t.Errorf("unexpected result for test case %s: %s (expected %s)", tc.name, got, tc.expected)
		}
	}
}

func TestIsSubsetOf(t *testing.T) {
	for _, tc := range []struct {
		name     string
		set1     FieldSet
		set2     FieldSet
		isSubset bool
	}{
		{
			name:     "empty set is subset of any set",
			set1:     FieldSet{},
			set2:     FieldSet{"a": {}},
			isSubset: true,
		},
		{
			name:     "set is subset of itself",
			set1:     FieldSet{"a": {}},
			set2:     FieldSet{"a": {}},
			isSubset: true,
		},
		{
			name:     "set is not subset of smaller set",
			set1:     FieldSet{"a": {}, "b": {}},
			set2:     FieldSet{"a": {}},
			isSubset: false,
		},
		{
			name:     "set is subset of larger set",
			set1:     FieldSet{"a": {}},
			set2:     FieldSet{"a": {}, "b": {}},
			isSubset: true,
		},
	} {
		if got := tc.set1.IsSubsetOf(tc.set2); got != tc.isSubset {
			t.Errorf("unexpected result for test case %s: %t (expected %t)", tc.name, got, tc.isSubset)
		}
	}
}

func TestDifference(t *testing.T) {
	for _, tc := range []struct {
		name     string
		set1     FieldSet
		set2     FieldSet
		expected FieldSet
	}{
		{
			name:     "empty set difference is empty set",
			set1:     FieldSet{},
			set2:     FieldSet{"a": {}},
			expected: FieldSet{},
		},
		{
			name:     "set difference with itself is empty set",
			set1:     FieldSet{"a": {}},
			set2:     FieldSet{"a": {}},
			expected: FieldSet{},
		},
		{
			name:     "set difference with subset is diff",
			set1:     FieldSet{"a": {}, "b": {}},
			set2:     FieldSet{"a": {}},
			expected: FieldSet{"b": {}},
		},
		{
			name:     "set difference with superset is empty set",
			set1:     FieldSet{"a": {}},
			set2:     FieldSet{"a": {}, "b": {}},
			expected: FieldSet{},
		},
	} {
		got := tc.set1.Difference(tc.set2)
		gotKeys := setToSortedSlice(got)
		expectedKeys := setToSortedSlice(tc.expected)
		if !cmp.Equal(gotKeys, expectedKeys) {
			t.Errorf("unexpected result for test case %s: %v (expected %v)", tc.name, gotKeys, expectedKeys)
		}
	}

}

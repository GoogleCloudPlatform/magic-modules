package diff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestFieldConflictSetsMerge(t *testing.T) {
	cases := []struct {
		name     string
		fcs      *FieldConflictSets
		other    *FieldConflictSets
		expected *FieldConflictSets
	}{
		{
			name:     "merging with nil",
			fcs:      makeFieldConflictSetsFromKey("field1///"),
			other:    nil,
			expected: makeFieldConflictSetsFromKey("field1///"),
		},
		{
			name:     "merging non-empty sets",
			fcs:      makeFieldConflictSetsFromKey("field1/field2,field3//"),
			other:    makeFieldConflictSetsFromKey("//field4/field5,field6"),
			expected: makeFieldConflictSetsFromKey("field1/field2,field3/field4/field5,field6"),
		},
		{
			name:     "merging overlapping sets",
			fcs:      makeFieldConflictSetsFromKey("field1;field2///"),
			other:    makeFieldConflictSetsFromKey("field2;field3///"),
			expected: makeFieldConflictSetsFromKey("field1;field2;field3///"),
		},
		{
			name:     "merging from other first",
			fcs:      makeFieldConflictSetsFromKey("///field2;field3;field4"),
			other:    makeFieldConflictSetsFromKey("///field1;field3"),
			expected: makeFieldConflictSetsFromKey("///field1;field2;field3;field4"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fcs.Merge(tc.other)
			if diff := cmp.Diff(tc.expected, tc.fcs); diff != "" {
				t.Errorf("merged FieldConflictSets not equal (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestRemoveZeroPadding(t *testing.T) {
	cases := map[string]string{
		"field1.0.field2":   "field1.field2",
		"field1.0.field2.0": "field1.field2",
		"field1.0":          "field1",
		"0.field1.0.field2": "field1.field2",
		"field1":            "field1",
		"":                  "",
		"0":                 "",
	}
	for input, expected := range cases {
		assert.Equal(t, expected, removeZeroPadding(input))
	}
}

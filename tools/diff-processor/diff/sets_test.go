package diff

import (
	"testing"
)

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

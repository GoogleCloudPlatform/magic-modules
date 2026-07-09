package gkehub2

import (
	"testing"
)

func TestRolloutSequenceDurationDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"different values, same units": {
			Old:                "60s",
			New:                "65s",
			ExpectDiffSuppress: false,
		},
		"different values, different units": {
			Old:                "65s",
			New:                "1d",
			ExpectDiffSuppress: false,
		},
		"same values, same units": {
			Old:                "60s",
			New:                "60s",
			ExpectDiffSuppress: true,
		},
		"same values, different units": {
			Old:                "60s",
			New:                "1m",
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		if rolloutSequenceDurationDiffSuppress("duration", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

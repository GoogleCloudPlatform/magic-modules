func TestSecondaryIpDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"empty strings": {
			Old:                "",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"auto range": {
			Old:                "",
			New:                "auto",
			ExpectDiffSuppress: false,
		},
		"auto on already applied range": {
			Old:                "10.0.0.0/28",
			New:                "auto",
			ExpectDiffSuppress: true,
		},
		"same ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.0.0.0/28",
			ExpectDiffSuppress: true,
		},
    "different ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.1.2.3/28",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if secondaryIpDiffSuppress("whatever", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

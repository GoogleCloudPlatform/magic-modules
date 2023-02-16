package google

import "testing"

func TestDocumentAIProcessorDefaultVersionVersionDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"API returns project id": {
			Old:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/my-version",
			ExpectDiffSuppress: false,
		},
		"API returns weird format": {
			Old:                "what-is-happening",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/my-version",
			ExpectDiffSuppress: false,
		},
		"User provides weird format": {
			Old:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "something-is-off-here",
			ExpectDiffSuppress: false,
		},
		"Only difference is project number vs id": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/my-version",
			ExpectDiffSuppress: true,
		},
		"different location": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/ca/processors/abcdef01234/processorVersions/my-version",
			ExpectDiffSuppress: false,
		},
		"different location with channel": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/ca/processors/abcdef01234/processorVersions/stable",
			ExpectDiffSuppress: false,
		},
		"different processor": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/qwerty01234/processorVersions/my-version",
			ExpectDiffSuppress: false,
		},
		"different processor with channel": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/qwerty01234/processorVersions/stable",
			ExpectDiffSuppress: false,
		},
		"different processor version": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/qwerty01234/processorVersions/my-version-2",
			ExpectDiffSuppress: false,
		},
		"different processor version with channel": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/qwerty01234/processorVersions/stable",
			ExpectDiffSuppress: false,
		},
		"stable channel": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/stable",
			ExpectDiffSuppress: true,
		},
		"rc channel": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/rc",
			ExpectDiffSuppress: true,
		},
		"pretrained channel (legacy)": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/pretrained",
			ExpectDiffSuppress: true,
		},
		"pretrained-next channel (legacy)": {
			Old:                "projects/1234567890/locations/us/processors/abcdef01234/processorVersions/my-version",
			New:                "projects/my-project/locations/us/processors/abcdef01234/processorVersions/pretrained-next",
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if documentAIProcessorDefaultVersionVersionDiffSuppress("version", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", tc.Old, tc.New, tc.ExpectDiffSuppress)
			}
		})
	}
}

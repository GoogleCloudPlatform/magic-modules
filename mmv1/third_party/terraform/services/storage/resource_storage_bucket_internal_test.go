package storage

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestLabelDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		K, Old, New        string
		ExpectDiffSuppress bool
	}{
		"missing goog-dataplex-asset-id": {
			K:                  "labels.goog-dataplex-asset-id",
			Old:                "test-bucket",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-asset-id": {
			K:                  "labels.goog-dataplex-asset-id",
			Old:                "test-bucket",
			New:                "test-bucket-1",
			ExpectDiffSuppress: false,
		},
		"missing goog-dataplex-lake-id": {
			K:                  "labels.goog-dataplex-lake-id",
			Old:                "test-lake",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-lake-id": {
			K:                  "labels.goog-dataplex-lake-id",
			Old:                "test-lake",
			New:                "test-lake-1",
			ExpectDiffSuppress: false,
		},
		"missing goog-dataplex-project-id": {
			K:                  "labels.goog-dataplex-project-id",
			Old:                "test-project-12345",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-project-id": {
			K:                  "labels.goog-dataplex-project-id",
			Old:                "test-project-12345",
			New:                "test-project-12345-1",
			ExpectDiffSuppress: false,
		},
		"missing goog-dataplex-zone-id": {
			K:                  "labels.goog-dataplex-zone-id",
			Old:                "test-zone1",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit goog-dataplex-zone-id": {
			K:                  "labels.goog-dataplex-zone-id",
			Old:                "test-zone1",
			New:                "test-zone1-1",
			ExpectDiffSuppress: false,
		},
		"labels.%": {
			K:                  "labels.%",
			Old:                "5",
			New:                "1",
			ExpectDiffSuppress: true,
		},
		"deleted custom key": {
			K:                  "labels.my-label",
			Old:                "my-value",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"added custom key": {
			K:                  "labels.my-label",
			Old:                "",
			New:                "my-value",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		if resourceDataplexLabelDiffSuppress(tc.K, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q: %q => %q expect DiffSuppress to return %t", tn, tc.K, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestIsIgnorableStorageObjectDeleteError(t *testing.T) {
	cases := map[string]struct {
		err      error
		expected bool
	}{
		// Callers check err != nil before calling, so nil is not an ignorable error.
		"nil": {
			err:      nil,
			expected: false,
		},
		"404 no such object": {
			err:      &googleapi.Error{Code: 404, Message: "No such object: example-bucket/path/to/object"},
			expected: true,
		},
		"410 gone": {
			err:      &googleapi.Error{Code: 410, Message: "Gone"},
			expected: true,
		},
		"wrapped 404": {
			err:      fmt.Errorf("deleting object: %w", &googleapi.Error{Code: 404, Message: "No such object: example-bucket/path/to/object"}),
			expected: true,
		},
		"403 forbidden": {
			err:      &googleapi.Error{Code: 403, Message: "Forbidden"},
			expected: false,
		},
		"409 conflict": {
			err:      &googleapi.Error{Code: 409, Message: "Conflict"},
			expected: false,
		},
		"wrapped 403": {
			err:      fmt.Errorf("deleting object: %w", &googleapi.Error{Code: 403, Message: "Forbidden"}),
			expected: false,
		},
		"non-googleapi error": {
			err:      errors.New("network timeout"),
			expected: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := isIgnorableStorageObjectDeleteError(tc.err); got != tc.expected {
				t.Fatalf("isIgnorableStorageObjectDeleteError(%v) = %v, want %v", tc.err, got, tc.expected)
			}
		})
	}
}

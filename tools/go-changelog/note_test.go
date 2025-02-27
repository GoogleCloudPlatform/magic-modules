// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package changelog

import (
	"reflect"
	"testing"
)

func TestValidateNote(t *testing.T) {
	cases := []struct {
		name          string
		changelogNote Note
		expectedError *EntryValidationError
	}{
		{
			name: "invalid type",
			changelogNote: Note{
				Type: "feature",
				Body: "this is to add a feature",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorUnknownTypes,
			},
		},
		{
			name: "newline after changelog content",
			changelogNote: Note{
				Type: "note",
				Body: "test only change\n",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorMultipleLines,
			},
		},
		{
			name: "valid new resource changelog format",
			changelogNote: Note{
				Type: "new-resource",
				Body: "`google_new_resource`",
			},
			expectedError: nil,
		},
		{
			name: "invalid new resource/datasource changelog format: missing backticks",
			changelogNote: Note{
				Type: "new-resource",
				Body: "google_new_resource",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorInvalidNewReourceOrDatasourceFormat,
			},
		},
		{
			name: "invalid new resource/datasource changelog format: missing google prefix",
			changelogNote: Note{
				Type: "new-datasource",
				Body: "`new_datasource`",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorInvalidNewReourceOrDatasourceFormat,
			},
		},
		{
			name: "invalid new resource/datasource changelog format: including spaces",
			changelogNote: Note{
				Type: "new-datasource",
				Body: "`google new datasource`",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorInvalidNewReourceOrDatasourceFormat,
			},
		},
		{
			name: "valid enhancement/bug fix changelog format",
			changelogNote: Note{
				Type: "enhancement",
				Body: "compute: added a new field to google_resource resource",
			},
			expectedError: nil,
		},
		{
			name: "valid enhancement/bug: allow underscore in product name",
			changelogNote: Note{
				Type: "enhancement",
				Body: "backup_dr: added a new field to google_resource resource",
			},
			expectedError: nil,
		},
		{
			name: "invalid enhancement/bug fix changelog format: missing product",
			changelogNote: Note{
				Type: "enhancement",
				Body: "added a new field to google_resource resource",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorInvalidEnhancementOrBugFixFormat,
			},
		},
		{
			name: "invalid enhancement/bug fix changelog format: incorrect product name",
			changelogNote: Note{
				Type: "enhancement",
				Body: "compute engine: added a new field to google_resource resource",
			},
			expectedError: &EntryValidationError{
				Code: EntryErrorInvalidEnhancementOrBugFixFormat,
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actualError := tc.changelogNote.Validate()
			if actualError != nil && tc.expectedError != nil {
				if !reflect.DeepEqual(actualError.Code, tc.expectedError.Code) {
					t.Errorf("want %v; got %v", tc.expectedError.Code, actualError.Code)
				}
			} else if actualError != nil {
				t.Errorf("want no error; got %v", actualError)
			} else if tc.expectedError != nil {
				t.Errorf("want %v; got no error", tc.expectedError)
			}

		})
	}
}

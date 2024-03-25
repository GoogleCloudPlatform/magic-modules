// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package changelog

import (
	"fmt"
	"reflect"
	"testing"
)

func TestValidateNote(t *testing.T) {
	cases := map[string]struct {
		changelogNote Note
		expectedError *EntryValidationError
	}{
		"invalid type": {
			changelogNote: Note{
				Type: "feature",
				Body: "this is to add a feature",
			},
			expectedError: &EntryValidationError{
				message: fmt.Sprintf("unknown changelog types %v: please use only the configured changelog entry types: %v", "feature", "this is to add a feature"),
				Code:    EntryErrorUnknownTypes,
				Details: map[string]interface{}{
					"type": "feature",
					"note": "this is to add a feature",
				},
			},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			fmt.Println(tc.changelogNote)
			error := tc.changelogNote.Validate()
			if !reflect.DeepEqual(error, tc.expectedError) {
				t.Errorf("want %v; got %v", tc.expectedError, error)
			}
		})
	}
}

package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestDetectMissingTests(t *testing.T) {
	allTests, err := readAllTests("testdata")
	if err != nil {
		t.Errorf("error reading tests before testing detect missing tests: %v", err)
	}
	for _, test := range []struct {
		name                   string
		changedFields          map[string]ResourceChanges
		expectedUntestedFields []string
	}{
		{
			name: "covered-resource",
			changedFields: map[string]ResourceChanges{
				"covered_resource": {
					"field_one": &Field{Added: true},
					"field_two": ResourceChanges{
						"field_three": &Field{Changed: true},
					},
					"field_four": ResourceChanges{
						"field_five": ResourceChanges{
							"field_six": &Field{Added: true},
						},
					},
				},
			},
		},
		{
			name: "uncovered-resource",
			changedFields: map[string]ResourceChanges{
				"uncovered_resource": {
					"field_one": &Field{Changed: true},
					"field_two": ResourceChanges{
						"field_three": &Field{Added: true},
					},
					"field_four": ResourceChanges{
						"field_five": ResourceChanges{
							"field_six": &Field{Changed: true},
						},
					},
				},
			},
			expectedUntestedFields: []string{"field_four.field_five.field_six", "field_one"},
		},
		{
			name: "config-variable-resource",
			changedFields: map[string]ResourceChanges{
				"config_variable": {
					"field_one": &Field{Added: true},
				},
			},
		},
		{
			name: "no-test-resource",
			changedFields: map[string]ResourceChanges{
				"no_test": {
					"field_one": &Field{Added: true},
				},
			},
			expectedUntestedFields: []string{"field_one"},
		},
	} {
		missingTests, err := detectMissingTests(test.changedFields, allTests)
		if err != nil {
			t.Errorf("error detecting missing tests for %s: %s", test.name, err)
		}
		if len(test.expectedUntestedFields) == 0 {
			if len(missingTests) > 0 {
				for resourceName, missingTest := range missingTests {
					t.Errorf("found unexpected untested fields in %s for resource %s: %v", test.name, resourceName, missingTest.UntestedFields)
				}
			}
		} else {
			if len(missingTests) == 1 {
				for _, missingTest := range missingTests {
					sort.Strings(missingTest.UntestedFields)
					if !reflect.DeepEqual(missingTest.UntestedFields, test.expectedUntestedFields) {
						t.Errorf(
							"did not find expected untested fields in %s, found %v, expected %v",
							test.name, missingTest.UntestedFields, test.expectedUntestedFields)
					}
				}
			} else {
				t.Errorf("found unexpected number of missing tests in %s: %d", test.name, len(missingTests))
			}
		}
	}
}

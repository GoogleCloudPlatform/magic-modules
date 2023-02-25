package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestDetectMissingTests(t *testing.T) {
	t.Fatal("this test should fail")
	allTests, err := readAllTests("testdata")
	if err != nil {
		t.Errorf("error reading tests before testing detect missing tests: %v", err)
	}
	for _, test := range []struct {
		name                   string
		changedFields          map[string]FieldCoverage
		expectedUntestedFields []string
	}{
		{
			name: "covered-resource",
			changedFields: map[string]FieldCoverage{
				"covered_resource": {
					"field_one": false,
					"field_two": FieldCoverage{
						"field_three": false,
					},
					"field_four": FieldCoverage{
						"field_five": FieldCoverage{
							"field_six": false,
						},
					},
				},
			},
		},
		{
			name: "uncovered-resource",
			changedFields: map[string]FieldCoverage{
				"uncovered_resource": {
					"field_one": false,
					"field_two": FieldCoverage{
						"field_three": false,
					},
					"field_four": FieldCoverage{
						"field_five": FieldCoverage{
							"field_six": false,
						},
					},
				},
			},
			expectedUntestedFields: []string{"field_four.field_five.field_six", "field_one"},
		},
		{
			name: "config-variable-resource",
			changedFields: map[string]FieldCoverage{
				"config_variable": {
					"field_one": false,
				},
			},
		},
		{
			name: "no-test-resource",
			changedFields: map[string]FieldCoverage{
				"no_test": {
					"field_one": false,
				},
			},
			expectedUntestedFields: []string{"field_one"},
		},
	} {
		missingTests := detectMissingTests(test.changedFields, allTests)
		if len(test.expectedUntestedFields) == 0 {
			if len(missingTests) > 0 {
				for resourceName, missingTest := range missingTests {
					t.Errorf("found unexpected untested fields for resource %s: %v", resourceName, missingTest.UntestedFields)
				}
			}
		} else {
			if len(missingTests) == 1 {
				for _, missingTest := range missingTests {
					sort.Strings(missingTest.UntestedFields)
					if !reflect.DeepEqual(missingTest.UntestedFields, test.expectedUntestedFields) {
						t.Errorf("did not find expected untested fields, found %v, expected %v", missingTest.UntestedFields, test.expectedUntestedFields)
					}
				}
			} else {
				t.Errorf("found unexpected number of missing tests: %d", len(missingTests))
			}
		}
	}
}

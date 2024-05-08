package detector

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/reader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetChangedFieldsFromSchemaDiff(t *testing.T) {
	for _, test := range []struct {
		name          string
		schemaDiff    diff.SchemaDiff
		changedFields map[string]ResourceChanges
	}{
		{
			name: "covered-resource",
			schemaDiff: diff.SchemaDiff{
				"covered_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"field_one": {
							New: &schema.Schema{},
						},
						"field_two.field_three": {
							New: &schema.Schema{},
							Old: &schema.Schema{},
						},
						"field_four": {
							New: &schema.Schema{
								Elem: &schema.Resource{},
							},
						},
						"field_four.field_five.field_six": {
							New: &schema.Schema{},
						},
						"field_seven": {
							New: &schema.Schema{Computed: true},
						},
						"project": {
							New: &schema.Schema{},
						},
					},
				},
			},
			changedFields: map[string]ResourceChanges{
				"covered_resource": {
					"field_one":                       &Field{Added: true},
					"field_two.field_three":           &Field{Changed: true},
					"field_four.field_five.field_six": &Field{Added: true},
				},
			},
		},
	} {
		if changedFields := getChangedFieldsFromSchemaDiff(test.schemaDiff); !reflect.DeepEqual(changedFields, test.changedFields) {
			t.Errorf("got unexpected changed fields: %v, expected %v", changedFields, test.changedFields)
		}
	}

}

func TestGetMissingTestsForChanges(t *testing.T) {
	allTests, errs := reader.ReadAllTests("../reader/testdata")
	if len(errs) > 0 {
		t.Errorf("errors reading tests before testing detect missing tests: %v", errs)
	}
	for _, test := range []struct {
		name                 string
		changedFields        map[string]ResourceChanges
		expectedMissingTests map[string]MissingTestInfo
	}{
		{
			name: "covered-resource",
			changedFields: map[string]ResourceChanges{
				"covered_resource": {
					"field_one":                       &Field{Added: true},
					"field_two.field_three":           &Field{Changed: true},
					"field_four.field_five.field_six": &Field{Added: true},
				},
			},
		},
		{
			name: "uncovered-resource",
			changedFields: map[string]ResourceChanges{
				"uncovered_resource": {
					"field_one":                       &Field{Changed: true},
					"field_two.field_three":           &Field{Added: true},
					"field_four.field_five.field_six": &Field{Changed: true},
				},
			},
			expectedMissingTests: map[string]MissingTestInfo{
				"uncovered_resource": {
					UntestedFields: []string{"field_four.field_five.field_six", "field_one"},
					SuggestedTest: `resource "uncovered_resource" "primary" {
  field_four {
    field_five {
      field_six = # value needed
    }
  }
  field_one = # value needed
}
`,
				},
			},
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
			expectedMissingTests: map[string]MissingTestInfo{
				"no_test": {
					UntestedFields: []string{"field_one"},
					SuggestedTest: `resource "no_test" "primary" {
  field_one = # value needed
}
`,
				},
			},
		},
		{
			name: "multiple-resources-missing-tests",
			changedFields: map[string]ResourceChanges{
				"no_test": {
					"field_one": &Field{Added: true},
				},
				"uncovered_resource": {
					"field_one":                       &Field{Changed: true},
					"field_two.field_three":           &Field{Added: true},
					"field_four.field_five.field_six": &Field{Changed: true},
				},
			},
			expectedMissingTests: map[string]MissingTestInfo{
				"no_test": {
					UntestedFields: []string{"field_one"},
					SuggestedTest: `resource "no_test" "primary" {
  field_one = # value needed
}
`,
				},
				"uncovered_resource": {
					UntestedFields: []string{"field_four.field_five.field_six", "field_one"},
					SuggestedTest: `resource "uncovered_resource" "primary" {
  field_four {
    field_five {
      field_six = # value needed
    }
  }
  field_one = # value needed
}
`,
				},
			},
		},
	} {
		missingTests, err := getMissingTestsForChanges(test.changedFields, allTests)
		if err != nil {
			t.Errorf("error detecting missing tests for %s: %s", test.name, err)
		}
		if len(test.expectedMissingTests) == 0 {
			if len(missingTests) > 0 {
				for resourceName, missingTest := range missingTests {
					t.Errorf("found unexpected untested fields in %s for resource %s: %v", test.name, resourceName, missingTest.UntestedFields)
				}
			}
		} else {
			if len(missingTests) == len(test.expectedMissingTests) {
				for resourceName, missingTest := range missingTests {
					expectedMissingTest := test.expectedMissingTests[resourceName]
					if !reflect.DeepEqual(missingTest.UntestedFields, expectedMissingTest.UntestedFields) {
						t.Errorf(
							"did not find expected untested fields in %s, found %v, expected %v",
							test.name, missingTest.UntestedFields, expectedMissingTest.UntestedFields)
					}
					if missingTest.SuggestedTest != expectedMissingTest.SuggestedTest {
						t.Errorf("did not find expected suggested test in %s, found %s, expected %s",
							test.name, missingTest.SuggestedTest, expectedMissingTest.SuggestedTest)
					}
				}
			} else {
				t.Errorf("found unexpected number of missing tests in %s: %d", test.name, len(missingTests))
			}
		}
	}
}

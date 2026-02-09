package detector

import (
	"reflect"
	"sort"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetChangedFieldsFromSchemaDiff(t *testing.T) {
	cases := []struct {
		name       string
		schemaDiff diff.SchemaDiff
		want       map[string]ResourceChanges
	}{
		{
			name: "covered-resource",
			schemaDiff: diff.SchemaDiff{
				"google_compute_instance": diff.ResourceDiff{
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
			want: map[string]ResourceChanges{
				"google_compute_instance": {
					"field_one":                       &Field{Added: true},
					"field_two.field_three":           &Field{Changed: true},
					"field_four.field_five.field_six": &Field{Added: true},
				},
			},
		},
		{
			name: "iam-resource",
			schemaDiff: diff.SchemaDiff{
				"google_resource_iam_member": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"condition": {
							New: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"expression": {
											Type: schema.TypeString,
										},
										"title": {
											Type: schema.TypeString,
										},
										"description": {
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"condition.description": {
							New: &schema.Schema{},
						},
						"condition.expression": {
							New: &schema.Schema{},
						},
						"condition.title": {
							New: &schema.Schema{},
						},
					},
				},
				"google_resource_iam_binding": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"condition": {
							New: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"expression": {
											Type: schema.TypeString,
										},
										"title": {
											Type: schema.TypeString,
										},
										"description": {
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"condition.description": {
							New: &schema.Schema{},
						},
						"condition.expression": {
							New: &schema.Schema{},
						},
						"condition.title": {
							New: &schema.Schema{},
						},
					},
				},
				"google_resource_iam_policy": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"condition": {
							New: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"expression": {
											Type: schema.TypeString,
										},
										"title": {
											Type: schema.TypeString,
										},
										"description": {
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"condition.description": {
							New: &schema.Schema{},
						},
						"condition.expression": {
							New: &schema.Schema{},
						},
						"condition.title": {
							New: &schema.Schema{},
						},
					},
				},
			},
			want: map[string]ResourceChanges{},
		},
	}

	for _, tc := range cases {
		got := getChangedFieldsFromSchemaDiff(tc.schemaDiff)
		if diff := cmp.Diff(got, tc.want); diff != "" {
			t.Errorf("getChangedFieldsFromSchemaDiff(%s) returned unexpected diff. +got, -want:\n%s", tc.name, diff)
		}
	}

}

func TestGetMissingTestsForChanges(t *testing.T) {
	allTests, errs := reader.ReadAllTests("../../test-reader/reader/testdata")
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

func TestDetectMissingDocs(t *testing.T) {
	// If repo is not temp dir, then the doc file points to tools/diff-processor/testdata/website/docs/r/a_resource.html.markdown.
	for _, test := range []struct {
		name       string
		schemaDiff diff.SchemaDiff
		repo       string
		want       map[string]MissingDocDetails
	}{
		{
			name: "doc file not exist",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"field_one": {
							New: &schema.Schema{},
						},
						"field_one.a": {
							New: &schema.Schema{},
						},
						"field_one.b": {
							New: &schema.Schema{},
						},
						"field_two.a": {
							New: &schema.Schema{},
							Old: &schema.Schema{},
						},
						"field_two.b": {
							New: &schema.Schema{},
						},
						"field_three": {
							New: &schema.Schema{
								Computed: true,
								Optional: true,
							},
						},
						"field_four": {
							New: &schema.Schema{
								Computed: true,
							},
						},
					},
				},
			},
			repo: t.TempDir(),
			want: map[string]MissingDocDetails{
				"a_resource": {
					Name:     "a_resource",
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields:   []string{"field_one", "field_one.a", "field_one.b", "field_two.b", "field_three", "field_four"},
				},
			},
		},
		{
			name: "doc file exist",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"field_one": {
							New: &schema.Schema{},
						},
						"field_one.a": {
							New: &schema.Schema{},
						},
						"field_one.b": {
							New: &schema.Schema{},
						},
						"field_two.a": {
							New: &schema.Schema{},
							Old: &schema.Schema{},
						},
						"field_two.b": {
							New: &schema.Schema{},
						},
						"field_three": {
							New: &schema.Schema{
								Computed: true,
								Optional: true,
							},
						},
						"field_four": {
							New: &schema.Schema{
								Computed: true,
							},
						},
					},
				},
			},
			repo: "../testdata",
			want: map[string]MissingDocDetails{
				"a_resource": {
					Name:     "a_resource",
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields:   []string{"field_one.b", "field_two.b", "field_three", "field_four"},
				},
			},
		},
		{
			name: "nested new field missing doc",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"field_one.c": {
							New: &schema.Schema{},
						},
					},
				},
			},
			repo: "../testdata",
			want: map[string]MissingDocDetails{
				"a_resource": {
					Name:     "a_resource",
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields:   []string{"field_one.c"},
				},
			},
		},
		{
			name: "member and members is member/members in doc",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"member": {
							New: &schema.Schema{},
						},
						"members": {
							New: &schema.Schema{},
						},
					},
				},
			},
			repo: "../testdata",
			want: map[string]MissingDocDetails{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DetectMissingDocs(test.schemaDiff, test.repo)
			if err != nil {
				t.Fatalf("DetectMissingDocs = %v, want = nil", err)
			}
			for r := range test.want {
				sort.Strings(test.want[r].Fields)
			}
			for r := range got {
				sort.Strings(got[r].Fields)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("DetectMissingDocs =  %v, want = %v", got, test.want)
			}
		})
	}
}

func TestDetectMissingDocsForDatasource(t *testing.T) {
	for _, test := range []struct {
		name       string
		schemaDiff diff.SchemaDiff
		repo       string
		want       map[string]MissingDocDetails
	}{
		{
			name: "doc file not exist",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
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
								Computed: true,
								Optional: true,
							},
						},
						"field_five": {
							New: &schema.Schema{
								Computed: true,
							},
						},
					},
				},
			},
			repo: t.TempDir(),
			want: map[string]MissingDocDetails{
				"a_resource": {
					Name:     "a_resource",
					FilePath: "/website/docs/d/a_resource.html.markdown",
					Fields: []string{
						"field_five",
						"field_one",
						"field_four",
					},
				},
			},
		},
		{
			name: "doc file exist",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
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
								Computed: true,
								Optional: true,
							},
						},
						"field_five": {
							New: &schema.Schema{
								Computed: true,
							},
						},
					},
				},
			},
			repo: "../testdata",
			want: map[string]MissingDocDetails{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DetectMissingDocsForDatasource(test.schemaDiff, test.repo)
			if err != nil {
				t.Fatalf("DetectMissingDocsForDatasource = %v, want = nil", err)
			}
			for r := range test.want {
				sort.Strings(test.want[r].Fields)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("got unexpected added fields: %v, expected %v", got, test.want)
			}
		})
	}
}

package detector

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	newProvider "google/provider/new/google/provider"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExistingMissingDoc(t *testing.T) {
	changes := diff.ComputeSchemaDiff(map[string]*schema.Resource{}, newProvider.ResourceMap())
	ret, err := DetectMissingDocs(changes, "/Users/ciris/gows/magic-modules/tools/diff-processor/new", newProvider.ResourceMap())
	if err != nil {
		t.Fatal(err)
	}
	for resource, info := range ret {
		if !(strings.Contains(resource, "_iam_policy") || strings.Contains(resource, "_iam_member") || strings.Contains(resource, "_iam_binding")) {
			t.Log(resource)
			t.Logf("info = %+v", info)
		}
	}
	if len(ret) > 0 {
		t.Fatalf("found %d resource has fields missing docs", len(ret))
	}
}

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
				"iam_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"condition": {
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
	// top level field_one is argument, field_two is attribute.
	resourceSchema := map[string]*schema.Resource{
		"a_resource": {
			Schema: map[string]*schema.Schema{
				"field_one": {
					Computed: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"a": {
								Computed: true,
								Optional: true,
							},
							"b": {
								Computed: true,
								Optional: false,
							},
							"c": {
								Computed: true,
								Optional: false,
							},
						},
					},
				},
				"field_two": {
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"a": {
								Computed: true,
								Optional: false,
							},
							"b": {
								Computed: true,
								Optional: false,
							},
							"c": {
								Computed: true,
								Optional: false,
							},
						},
					},
				},
				"field_three": {
					Computed: true,
					Optional: true,
				},
				"field_four": {
					Computed: true,
				},
			},
		},
	}

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
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields: []MissingDocField{
						{
							Field:   "field_one",
							Section: "Argument Reference",
						},
						{
							Field:   "field_one.a",
							Section: "Argument Reference",
						},
						{
							Field:   "field_one.b",
							Section: "Argument Reference",
						},
						{
							Field:   "field_two.b",
							Section: "Attributes Reference",
						},
						{
							Field:   "field_three",
							Section: "Argument Reference",
						},
						{
							Field:   "field_four",
							Section: "Attributes Reference",
						},
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
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields: []MissingDocField{
						{
							Field:   "field_one.b",
							Section: "Argument Reference",
						},
						{
							Field:   "field_two.b",
							Section: "Attributes Reference",
						},
						{
							Field:   "field_three",
							Section: "Argument Reference",
						},
						{
							Field:   "field_four",
							Section: "Attributes Reference",
						},
					},
				},
			},
		},
		{
			name: "argument's nested new field missing doc",
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
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields: []MissingDocField{
						{
							Field:   "field_one.c",
							Section: "Argument Reference",
						},
					},
				},
			},
		},
		{
			name: "attribute's nested new field missing doc",
			schemaDiff: diff.SchemaDiff{
				"a_resource": diff.ResourceDiff{
					Fields: map[string]diff.FieldDiff{
						"field_two.c": {
							New: &schema.Schema{},
						},
					},
				},
			},
			repo: "../testdata",
			want: map[string]MissingDocDetails{
				"a_resource": {
					FilePath: "/website/docs/r/a_resource.html.markdown",
					Fields: []MissingDocField{
						{
							Field:   "field_two.c",
							Section: "Attributes Reference",
						},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := DetectMissingDocs(test.schemaDiff, test.repo, resourceSchema)
			if err != nil {
				t.Fatalf("DetectMissingDocs = %v, want = nil", err)
			}
			for r := range test.want {
				sort.Slice(test.want[r].Fields, func(i, j int) bool {
					return test.want[r].Fields[i].Field < test.want[r].Fields[j].Field
				})
			}
			for r := range got {
				sort.Slice(got[r].Fields, func(i, j int) bool {
					return got[r].Fields[i].Field < got[r].Fields[j].Field
				})
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("got unexpected added fields: %v, expected %v", got, test.want)
			}
		})
	}
}

func TestIsAttribute(t *testing.T) {
	resourceSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"one": {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"a": {
							Computed: false,
							Optional: true,
						},
					},
				},
			},
			"two": {
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"a": {
							Computed: true,
							Optional: false,
						},
					},
				},
			},
		},
	}
	for _, test := range []struct {
		name  string
		field string
		want  bool
	}{
		{
			name:  "field itself is top level attribute",
			field: "one",
			want:  true,
		},
		{
			name:  "field's top level is attribute",
			field: "one.a",
			want:  true,
		},
		{
			name:  "field itself is not top level attribute",
			field: "two",
			want:  false,
		},
		{
			name:  "field's top level is not attribute",
			field: "two.a",
			want:  false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := isAttribute(test.field, resourceSchema)
			if test.want != got {
				t.Fatalf("isAttribute = %v, want = %v", got, test.want)
			}

		})
	}
}

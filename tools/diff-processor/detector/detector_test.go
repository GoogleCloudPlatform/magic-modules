package detector

import (
	"os"
	"path/filepath"
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
			// block_types_test.go declares block_types_resource as a managed
			// resource, a data source, an ephemeral resource and a list
			// resource. Only field_one, on the managed resource, is coverage
			// for the managed resource's schema; the fields of the same-named
			// blocks of other kinds belong to different schemas.
			name: "other-block-kinds-are-not-managed-resource-coverage",
			changedFields: map[string]ResourceChanges{
				"block_types_resource": {
					"field_one":   &Field{Added: true},
					"field_two":   &Field{Added: true},
					"field_three": &Field{Added: true},
				},
			},
			expectedMissingTests: map[string]MissingTestInfo{
				"block_types_resource": {
					UntestedFields: []string{"field_three", "field_two"},
					SuggestedTest: `resource "block_types_resource" "primary" {
  field_three = # value needed
  field_two   = # value needed
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
		{
			name: "expect-error-and-non-applying-steps-are-not-coverage",
			changedFields: map[string]ResourceChanges{
				"expect_error_resource": {
					"field_one":   &Field{Added: true},
					"field_two":   &Field{Added: true},
					"field_three": &Field{Added: true},
				},
			},
			expectedMissingTests: map[string]MissingTestInfo{
				"expect_error_resource": {
					UntestedFields: []string{"field_three", "field_two"},
					SuggestedTest: `resource "expect_error_resource" "primary" {
  field_three = # value needed
  field_two   = # value needed
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

func TestDetectMissingUpdateTests(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "resource_update_coverage_test.go")
	src := `package service_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccOneStepOnly(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: ` + "`" + `
resource "google_one_step" "primary" {
  mutable_field   = "only-create"
  immutable_field = "create-only"
}
` + "`" + `,
			},
			&resource.TestStep{
				Config: ` + "`" + `
resource "google_one_step" "primary" {
  mutable_field   = "rejected-update"
  immutable_field = "create-only"
  error_only      = "invalid"
}
` + "`" + `,
				ExpectError: regexp.MustCompile("invalid"),
			},
		},
	})
}

func TestAccTwoStepMutations(t *testing.T) {
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"static_ctx":    "unchanged",
		"mutated_ctx":   "before",
		"flag_ctx":      true,
	}
	context_1 := map[string]interface{}{
		"random_suffix": context["random_suffix"],
		"static_ctx":    "unchanged",
		"mutated_ctx":   "after",
		"flag_ctx":      false,
	}

	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccTwoStepConfigStep1(context),
			},
			{
				Config: testAccTwoStepConfigStep2(context_1),
			},
		},
	})
}

func testAccTwoStepConfigStep1(context map[string]interface{}) string {
	return acctest.Nprintf(` + "`" + `
resource "google_two_step" "primary" {
  name             = "tf-test-%{random_suffix}"
  identical_lit    = "same"
  static_ctx_field = "%{static_ctx}"
  mutated_ctx      = "%{mutated_ctx}"
  mutated_bool     = %{flag_ctx}
  cleared_in_step2 = "present"
}
` + "`" + `, context)
}

func testAccTwoStepConfigStep2(context map[string]interface{}) string {
	return acctest.Nprintf(` + "`" + `
resource "google_two_step" "primary" {
  name             = "tf-test-%{random_suffix}"
  identical_lit    = "same"
  static_ctx_field = "%{static_ctx}"
  mutated_ctx      = "%{mutated_ctx}"
  mutated_bool     = %{flag_ctx}
  added_in_step2   = "added"
}
` + "`" + `, context)
}
`
	if err := os.WriteFile(testFile, []byte(src), 0o600); err != nil {
		t.Fatalf("failed writing temp test file: %v", err)
	}
	parsedTests, errs := reader.ReadTestFiles([]string{testFile})
	if len(errs) > 0 {
		t.Fatalf("unexpected reader errors: %v", errs)
	}

	// Verify ExpectError pointer step does not satisfy positive coverage in DetectMissingTests.
	missingPositive, err := DetectMissingTests(diff.SchemaDiff{
		"google_one_step": {
			Fields: map[string]diff.FieldDiff{
				"mutable_field": {New: &schema.Schema{Optional: true}},
				"error_only":    {New: &schema.Schema{Optional: true}},
			},
		},
	}, parsedTests)
	if err != nil {
		t.Fatalf("DetectMissingTests returned error: %v", err)
	}
	if got := missingPositive["google_one_step"]; got == nil || !reflect.DeepEqual(got.UntestedFields, []string{"error_only"}) {
		t.Fatalf("expected error_only to be flagged by DetectMissingTests, got %+v", got)
	}

	schemaDiff := diff.SchemaDiff{
		"google_one_step": {
			Fields: map[string]diff.FieldDiff{
				"mutable_field":     {New: &schema.Schema{Optional: true}},
				"immutable_field":   {New: &schema.Schema{Optional: true, ForceNew: true}},
				"computed_only":     {New: &schema.Schema{Computed: true}},
				"newly_updatable":   {Old: &schema.Schema{Optional: true, ForceNew: true}, New: &schema.Schema{Optional: true, ForceNew: false}},
				"existing_modified": {Old: &schema.Schema{Optional: true}, New: &schema.Schema{Optional: true, Description: "updated docs"}},
			},
		},
		"google_two_step": {
			Fields: map[string]diff.FieldDiff{
				"name":             {New: &schema.Schema{Required: true, ForceNew: true}},
				"identical_lit":    {New: &schema.Schema{Optional: true}},
				"static_ctx_field": {New: &schema.Schema{Optional: true}},
				"mutated_ctx":      {New: &schema.Schema{Optional: true}},
				"mutated_bool":     {New: &schema.Schema{Optional: true}},
				"cleared_in_step2": {New: &schema.Schema{Optional: true}},
				"added_in_step2":   {New: &schema.Schema{Optional: true}},
			},
		},
	}

	gotUpdateMissing, err := DetectMissingUpdateTests(schemaDiff, parsedTests)
	if err != nil {
		t.Fatalf("DetectMissingUpdateTests returned error: %v", err)
	}

	wantUpdateMissing := map[string]*MissingUpdateTestInfo{
		"google_one_step": {
			UntestedUpdateFields: []string{"mutable_field", "newly_updatable"},
			Tests:                []string{"TestAccOneStepOnly"},
		},
		"google_two_step": {
			UntestedUpdateFields: []string{"identical_lit", "static_ctx_field"},
			Tests:                []string{"TestAccTwoStepMutations"},
		},
	}
	if d := cmp.Diff(wantUpdateMissing, gotUpdateMissing); d != "" {
		t.Errorf("DetectMissingUpdateTests mismatch (-want +got):\n%s", d)
	}
}

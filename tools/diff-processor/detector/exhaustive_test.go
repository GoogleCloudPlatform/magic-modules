package detector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func TestDetectExhaustiveMissingTests_CoreRecursiveLifecycle(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"description": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"nested": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     &schema.Resource{},
					},
				},
				"nested.opt_a": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"nested.opt_b": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
			},
		},
	}

	// Step 1: set all fields
	// Step 2: modify mutable fields (description, nested, nested.opt_a, nested.opt_b)
	// Step 3: unset nested.opt_b while nested remains present with nested.opt_a
	// Step 4: unset nested.opt_a while nested remains present as an empty block ("{}")
	// Step 5: unset nested and description
	tests := []*reader.Test{
		{
			Name: "TestAccExample_exhaustive",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":         "\"test-name\"",
								"description":  "\"desc-1\"",
								"nested.opt_a": "\"a-1\"",
								"nested.opt_b": "\"b-1\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":         "\"test-name\"",
								"description":  "\"desc-2\"",
								"nested.opt_a": "\"a-2\"",
								"nested.opt_b": "\"b-2\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":         "\"test-name\"",
								"description":  "\"desc-2\"",
								"nested.opt_a": "\"a-2\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":        "\"test-name\"",
								"description": "\"desc-2\"",
								"nested":      "{}",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name": "\"test-name\"",
							},
						},
					},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, tests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected all fields to be exhaustively covered, got incomplete: %s", FormatExhaustiveMarkdown(got))
	}
}

func TestDetectExhaustiveMissingTests_AntiShortcutMultiChildUnset(t *testing.T) {
	// If nested has two optional children (opt_a and opt_b) both present in Step 2,
	// and Step 3 unsets the entire nested block at once without ever unsetting opt_a or opt_b
	// while nested remained present, opt_a and opt_b must be flagged as missing UNSET.
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"nested": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     &schema.Resource{},
					},
				},
				"nested.opt_a": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"nested.opt_b": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
			},
		},
	}

	tests := []*reader.Test{
		{
			Name: "TestAccExample_shortcut",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"nested.opt_a": "\"a-1\"",
								"nested.opt_b": "\"b-1\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"nested.opt_a": "\"a-2\"",
								"nested.opt_b": "\"b-2\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {},
						},
					},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, tests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info, ok := got["google_example_resource"]
	if !ok {
		t.Fatalf("expected google_example_resource to have missing tests")
	}
	wantIncomplete := []string{"nested.opt_a", "nested.opt_b"}
	if diffStr := cmp.Diff(wantIncomplete, info.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}
	if info.Fields["nested"].Unset.Status != StatusCovered {
		t.Errorf("expected nested UNSET to be COVERED, got %+v", info.Fields["nested"].Unset)
	}
	if info.Fields["nested.opt_a"].Unset.Status != StatusMissing {
		t.Errorf("expected nested.opt_a UNSET to be MISSING, got %+v", info.Fields["nested.opt_a"].Unset)
	}
	if info.Fields["nested.opt_b"].Unset.Status != StatusMissing {
		t.Errorf("expected nested.opt_b UNSET to be MISSING, got %+v", info.Fields["nested.opt_b"].Unset)
	}
}

func TestDetectExhaustiveMissingTests_SingleOptionalChildParentUnset(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"encryption_config": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     &schema.Resource{},
					},
				},
				"encryption_config.kms_key": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"encryption_config.state": {
					New: &schema.Schema{Type: schema.TypeString, Computed: true},
				},
				"oc_parent": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem:     &schema.Resource{},
					},
				},
				"oc_parent.opt_child": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"text_format": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     &schema.Resource{},
					},
				},
				"text_format.delimiter": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, Default: "\n"},
				},
				"timeouts_config": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem:     &schema.Resource{},
					},
				},
				"timeouts_config.connect_timeout": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"timeouts_config.read_timeout": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"fn_parent": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						ForceNew: true,
						Elem:     &schema.Resource{},
					},
				},
				"fn_parent.fn_sole_child": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
			},
		},
	}

	// Case 1: encryption_config is purely Optional with a single defaultless configurable
	// child (kms_key), so removing encryption_config in Step 4 covers kms_key UNSET.
	// By contrast:
	// - oc_parent is Optional+Computed, so omitting oc_parent retains prior state in SDKv2
	//   rather than unsetting oc_parent.opt_child (which requires `oc_parent {}`);
	// - text_format.delimiter has a schema Default, so `text_format {}` is non-empty in
	//   SDKv2 and required to test reverting delimiter to its default before removing text_format;
	// - timeouts_config has two optional children (connect_timeout and read_timeout), so
	//   removing timeouts_config without swapping to read_timeout leaves connect_timeout UNSET MISSING.
	baseSteps := []reader.Step{
		{
			reader.ResourceBlock: {
				"google_example_resource": {
					"primary": {
						"name":                            "\"test-name\"",
						"encryption_config.kms_key":       "\"k1\"",
						"oc_parent.opt_child":             "\"oc1\"",
						"text_format.delimiter":           "\",\"",
						"timeouts_config.connect_timeout": "\"10s\"",
						"timeouts_config.read_timeout":    "\"30s\"",
						"fn_parent.fn_sole_child":         "\"fn-1\"",
					},
				},
			},
		},
		{
			reader.ResourceBlock: {
				"google_example_resource": {
					"primary": {
						"name":                            "\"test-name\"",
						"encryption_config.kms_key":       "\"k2\"",
						"oc_parent.opt_child":             "\"oc2\"",
						"text_format.delimiter":           "\";\"",
						"timeouts_config.connect_timeout": "\"20s\"",
						"timeouts_config.read_timeout":    "\"60s\"",
						"fn_parent.fn_sole_child":         "\"fn-1\"",
					},
				},
			},
		},
		{
			reader.ResourceBlock: {
				"google_example_resource": {
					"primary": {
						"name":                            "\"test-name\"",
						"encryption_config.kms_key":       "\"k2\"",
						"timeouts_config.connect_timeout": "\"20s\"",
						"fn_parent.fn_sole_child":         "\"fn-1\"",
					},
				},
			},
		},
		{
			reader.ResourceBlock: {
				"google_example_resource": {
					"primary": {
						"name":                    "\"test-name\"",
						"fn_parent.fn_sole_child": "\"fn-1\"",
					},
				},
			},
		},
	}

	basicTest := &reader.Test{
		Name: "TestAccExample_basic",
		Steps: []reader.Step{
			{
				reader.ResourceBlock: {
					"google_example_resource": {
						"primary": {
							"name": "\"test-name\"",
						},
					},
				},
			},
		},
	}

	gotWithoutSwap, err := DetectExhaustiveMissingTests(schemaDiff, nil, []*reader.Test{
		{Name: "TestAccExample_soleChild", Steps: baseSteps},
		basicTest,
	}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := gotWithoutSwap["google_example_resource"]
	if info == nil {
		t.Fatalf("expected oc_parent.opt_child, text_format.delimiter, and timeouts_config.connect_timeout UNSET to remain MISSING")
	}
	wantIncomplete := []string{
		"oc_parent.opt_child",
		"text_format.delimiter",
		"timeouts_config.connect_timeout",
	}
	if diffStr := cmp.Diff(wantIncomplete, info.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}

	// Case 2: Transitioning Step 2 -> Step 3 via `oc_parent {}` and `text_format {}`
	// (reverting delimiter to its default) and swapping timeouts_config to
	// `{ read_timeout = "60s" }` before removing the blocks covers all three.
	stepsWithSwap := []reader.Step{
		baseSteps[0],
		baseSteps[1],
		{
			reader.ResourceBlock: {
				"google_example_resource": {
					"primary": {
						"name":                            "\"test-name\"",
						"encryption_config.kms_key":       "\"k2\"",
						"oc_parent":                       "{}",
						"text_format":                     "{}",
						"timeouts_config.connect_timeout": "\"20s\"",
						"fn_parent.fn_sole_child":         "\"fn-1\"",
					},
				},
			},
		},
		{
			reader.ResourceBlock: {
				"google_example_resource": {
					"primary": {
						"name":                         "\"test-name\"",
						"timeouts_config.read_timeout": "\"60s\"",
						"fn_parent.fn_sole_child":      "\"fn-1\"",
					},
				},
			},
		},
		baseSteps[3],
	}

	gotCovered, err := DetectExhaustiveMissingTests(schemaDiff, nil, []*reader.Test{
		{Name: "TestAccExample_soleChild", Steps: stepsWithSwap},
		basicTest,
	}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotCovered) != 0 {
		t.Errorf("expected all fields covered once multi-child block swaps siblings, got:\n%s", FormatExhaustiveMarkdown(gotCovered))
	}
}

func TestDetectExhaustiveMissingTests_CoreExemptions(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"project": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
				},
				"deletion_policy": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"output_only": {
					New: &schema.Schema{Type: schema.TypeString, Computed: true},
				},
				"req_field": {
					New: &schema.Schema{Type: schema.TypeString, Required: true},
				},
				"force_new_field": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
				"opt_computed_field": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
				},
				"opt_block": {
					New: &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Resource{},
					},
				},
				"opt_block.req_child": {
					New: &schema.Schema{Type: schema.TypeString, Required: true},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := got["google_example_resource"]
	if info == nil {
		t.Fatalf("expected google_example_resource in results")
	}

	for _, skipped := range []string{"project", "deletion_policy", "output_only"} {
		if _, exists := info.Fields[skipped]; exists {
			t.Errorf("expected field %q to be skipped, but it was present", skipped)
		}
	}

	assertExempt := func(fieldPath string, tr TransitionResult, wantCode string) {
		t.Helper()
		if tr.Status != StatusExempt || tr.ExemptCode != wantCode {
			t.Errorf("%s: got %+v, want EXEMPT (%s)", fieldPath, tr, wantCode)
		}
	}
	assertMissing := func(fieldPath string, tr TransitionResult) {
		t.Helper()
		if tr.Status != StatusMissing {
			t.Errorf("%s: got %+v, want MISSING", fieldPath, tr)
		}
	}

	assertExempt("req_field.Unset", info.Fields["req_field"].Unset, "required")
	assertExempt("force_new_field.Modify", info.Fields["force_new_field"].Modify, "force_new")
	assertMissing("force_new_field.Unset", info.Fields["force_new_field"].Unset)
	assertExempt("opt_computed_field.Unset", info.Fields["opt_computed_field"].Unset, "optional_computed")
	assertMissing("opt_computed_field.Modify", info.Fields["opt_computed_field"].Modify)
	assertExempt("opt_block.req_child.Unset", info.Fields["opt_block.req_child"].Unset, "required_in_parent")
}

func TestDetectExhaustiveMissingTests_CoreForceNewOmissionUnset(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_updatable_with_forcenew": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"labels": {
					New: &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
				},
				"description": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
			},
		},
	}

	// Case 1: Only a _full test exists where optional ForceNew description is always set.
	fullOnlyTests := []*reader.Test{
		{
			Name: "TestAccUpdatableWithForcenew_full",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_updatable_with_forcenew": {
							"primary": {
								"name":        "\"n1\"",
								"labels.env":  "\"dev\"",
								"description": "\"desc\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_updatable_with_forcenew": {
							"primary": {
								"name":        "\"n1\"",
								"labels.env":  "\"prod\"",
								"description": "\"desc\"",
							},
						},
					},
				},
			},
		},
	}

	gotMissing, err := DetectExhaustiveMissingTests(schemaDiff, nil, fullOnlyTests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	infoMissing := gotMissing["google_updatable_with_forcenew"]
	if infoMissing == nil {
		t.Fatalf("expected incomplete fields when optional ForceNew fields are never omitted")
	}
	for _, fp := range []string{"description", "labels"} {
		if infoMissing.Fields[fp].Unset.Status != StatusMissing {
			t.Errorf("expected %s.Unset to be MISSING, got %+v", fp, infoMissing.Fields[fp].Unset)
		}
	}

	// Case 2: Add _basic test omitting ForceNew description -> description UNSET becomes COVERED,
	// while mutable labels still requires an in-place update UNSET step.
	withBasicNoLabelsUnset := append([]*reader.Test{}, fullOnlyTests...)
	withBasicNoLabelsUnset = append(withBasicNoLabelsUnset, &reader.Test{
		Name: "TestAccUpdatableWithForcenew_basic",
		Steps: []reader.Step{
			{
				reader.ResourceBlock: {
					"google_updatable_with_forcenew": {
						"primary": {
							"name": "\"n1\"",
						},
					},
				},
			},
		},
	})

	gotPartial, err := DetectExhaustiveMissingTests(schemaDiff, nil, withBasicNoLabelsUnset, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	infoPartial := gotPartial["google_updatable_with_forcenew"]
	if infoPartial == nil {
		t.Fatalf("expected mutable labels UNSET to still be missing without an in-place unset step")
	}
	if diffStr := cmp.Diff([]string{"labels"}, infoPartial.IncompleteFields); diffStr != "" {
		t.Errorf("expected ONLY mutable labels to remain incomplete (-want +got):\n%s", diffStr)
	}
	if infoPartial.Fields["description"].Unset.CoveredBy != "TestAccUpdatableWithForcenew_basic (step 1)" {
		t.Errorf("expected description.Unset covered by basic test, got %+v", infoPartial.Fields["description"].Unset)
	}

	// Case 3: Add Step 3 unsetting mutable labels in-place -> 0 incomplete fields!
	completeTests := []*reader.Test{
		{
			Name: "TestAccUpdatableWithForcenew_full",
			Steps: append(append([]reader.Step{}, fullOnlyTests[0].Steps...), reader.Step{
				reader.ResourceBlock: {
					"google_updatable_with_forcenew": {
						"primary": {
							"name":        "\"n1\"",
							"description": "\"desc\"",
						},
					},
				},
			}),
		},
		withBasicNoLabelsUnset[1],
	}

	gotComplete, err := DetectExhaustiveMissingTests(schemaDiff, nil, completeTests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotComplete) != 0 {
		t.Errorf("expected 0 incomplete fields once in-place labels UNSET + ForceNew omission step are present, got:\n%s", FormatExhaustiveMarkdown(gotComplete))
	}
}

func TestDetectExhaustiveMissingTests_PlanCheckRequiredForInPlaceUpdates(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_my_widget": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"description": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"labels": {
					New: &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
				},
				"immutable_tag": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
			},
		},
	}

	testFile := filepath.Join(t.TempDir(), "resource_my_widget_test.go")
	src := `package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/acctest"
)

func TestAccMyWidget_unverifiedAndReplace(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: ` + "`" + `
resource "google_my_widget" "primary" {
  name        = "w1"
  description = "desc-1"
  labels = {
    env = "dev"
  }
}
` + "`" + `,
			},
			{
				// Missing ConfigPlanChecks: cannot verify this was an in-place update.
				Config: ` + "`" + `
resource "google_my_widget" "primary" {
  name        = "w1"
  description = "desc-2"
}
` + "`" + `,
			},
			{
				// Explicit ResourceActionReplace: proves recreation rather than update,
				// though ForceNew immutable_tag SET is still valid on this apply step.
				Config: ` + "`" + `
resource "google_my_widget" "primary" {
  name          = "w2"
  description   = "desc-3"
  immutable_tag = "tag-1"
}
` + "`" + `,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_my_widget.primary", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}
`
	if err := os.WriteFile(testFile, []byte(src), 0o600); err != nil {
		t.Fatalf("error writing temp test file: %v", err)
	}
	parsedTests, errs := reader.ReadTestFiles([]string{testFile})
	if len(errs) > 0 {
		t.Fatalf("unexpected ReadTestFiles errors: %v", errs)
	}

	gotMissing, err := DetectExhaustiveMissingTests(schemaDiff, nil, parsedTests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := gotMissing["google_my_widget"]
	if info == nil {
		t.Fatalf("expected google_my_widget to have incomplete MODIFY and UNSET without ResourceActionUpdate plancheck")
	}
	if diffStr := cmp.Diff([]string{"description", "labels"}, info.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}
	if info.Fields["immutable_tag"].Set.Status != StatusCovered || info.Fields["immutable_tag"].Unset.Status != StatusCovered {
		t.Errorf("expected ForceNew immutable_tag SET and omission UNSET to be covered even on unverified/replace steps, got %+v", info.Fields["immutable_tag"])
	}
	wantSkipped := []string{
		"TestAccMyWidget_unverifiedAndReplace (step 1 -> 2: missing plancheck.ResourceActionUpdate on google_my_widget.primary)",
		"TestAccMyWidget_unverifiedAndReplace (step 2 -> 3: plancheck expects ResourceActionReplace on google_my_widget.primary)",
	}
	if diffStr := cmp.Diff(wantSkipped, info.SkippedUpdateSteps); diffStr != "" {
		t.Errorf("SkippedUpdateSteps mismatch (-want +got):\n%s", diffStr)
	}
	md := FormatExhaustiveMarkdown(gotMissing)
	if !strings.Contains(md, "plancheck.ResourceActionUpdate") || !strings.Contains(md, wantSkipped[0]) {
		t.Errorf("expected markdown to surface skipped update steps note, got:\n%s", md)
	}

	// Step 1 -> 2 has ExpectEmptyPlan (rejected for labels MODIFY), while Step 2 -> 3
	// has ResourceActionUpdate (covers description MODIFY + labels UNSET) and Step 3 -> 4
	// has ExpectNoDelete (covers description UNSET).
	emptyPlanTest := &reader.Test{
		Name: "TestAccMyWidget_verifiedUpdate",
		Steps: []reader.Step{
			{
				reader.ResourceBlock: {
					"google_my_widget": {
						"primary": {
							"name":        "\"w1\"",
							"description": "\"d1\"",
							"labels.env":  "\"dev\"",
						},
					},
				},
			},
			{
				reader.ResourceBlock: {
					"google_my_widget": {
						"primary": {
							"name":        "\"w1\"",
							"description": "\"d2\"",
							"labels.env":  "\"prod\"",
						},
					},
				},
			},
			{
				reader.ResourceBlock: {
					"google_my_widget": {
						"primary": {
							"name":        "\"w1\"",
							"description": "\"d3\"",
						},
					},
				},
			},
			{
				reader.ResourceBlock: {
					"google_my_widget": {
						"primary": {
							"name": "\"w1\"",
						},
					},
				},
			},
		},
		PlanChecks: []reader.StepPlanChecks{
			{},
			{ExpectEmptyPlan: true},
			{ResourceActions: map[string]string{"google_my_widget.primary": "ResourceActionUpdate"}},
			{ExpectNoDelete: true},
		},
	}

	gotEmptyPlan, err := DetectExhaustiveMissingTests(schemaDiff, nil, append(append([]*reader.Test{}, parsedTests...), emptyPlanTest), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	infoEmptyPlan := gotEmptyPlan["google_my_widget"]
	if infoEmptyPlan == nil {
		t.Fatalf("expected labels MODIFY to remain incomplete when step 1 -> 2 has ExpectEmptyPlan")
	}
	if diffStr := cmp.Diff([]string{"labels"}, infoEmptyPlan.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}
	wantEmptyPlanSkipped := []string{
		"TestAccMyWidget_verifiedUpdate (step 1 -> 2: plancheck expects empty plan)",
	}
	if diffStr := cmp.Diff(wantEmptyPlanSkipped, infoEmptyPlan.SkippedUpdateSteps); diffStr != "" {
		t.Errorf("SkippedUpdateSteps mismatch (-want +got):\n%s", diffStr)
	}

	// Replacing ExpectEmptyPlan with ResourceActionUpdate on Step 2 covers labels MODIFY too.
	verifiedTest := *emptyPlanTest
	verifiedTest.PlanChecks = []reader.StepPlanChecks{
		{},
		{ResourceActions: map[string]string{"google_my_widget.primary": "ResourceActionUpdate"}},
		{ResourceActions: map[string]string{"google_my_widget.primary": "ResourceActionUpdate"}},
		{ExpectNoDelete: true},
	}
	gotCovered, err := DetectExhaustiveMissingTests(schemaDiff, nil, append(append([]*reader.Test{}, parsedTests...), &verifiedTest), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotCovered) != 0 {
		t.Errorf("expected all fields covered once verified ResourceActionUpdate / ExpectNoDelete steps exist, got:\n%s", FormatExhaustiveMarkdown(gotCovered))
	}
}

func TestDetectExhaustiveMissingTests_NoUpdateResource(t *testing.T) {
	dummyCRUD := func(d *schema.ResourceData, meta interface{}) error { return nil }

	t.Run("hasUpdate classification", func(t *testing.T) {
		// 1. Update == nil -> hasUpdate == false
		noUpdateRes := &schema.Resource{
			Create: dummyCRUD,
			Read:   dummyCRUD,
			Delete: dummyCRUD,
			Schema: map[string]*schema.Schema{
				"field_a": {Type: schema.TypeString, Optional: true, ForceNew: true},
			},
		}
		if hasUpdate(noUpdateRes, noUpdateRes.Schema) {
			t.Errorf("expected resource without Update callback to have hasUpdate=false")
		}

		// 2. Stub Update callback with only ForceNew fields + deletion_policy -> hasUpdate == false
		stubUpdateRes := &schema.Resource{
			Create: dummyCRUD,
			Read:   dummyCRUD,
			Update: dummyCRUD,
			Delete: dummyCRUD,
			Schema: map[string]*schema.Schema{
				"name":            {Type: schema.TypeString, Required: true, ForceNew: true},
				"deletion_policy": {Type: schema.TypeString, Optional: true, Computed: true},
			},
		}
		if hasUpdate(stubUpdateRes, stubUpdateRes.Schema) {
			t.Errorf("expected resource with only ForceNew fields and deletion_policy stub Update to have hasUpdate=false")
		}

		// 3. Updatable singleton with ForceNew scalar + mutable empty nested blocks (`basic {}`) -> hasUpdate == true
		emptyBlockSingletonSchema := map[string]*schema.Schema{
			"region": {Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     &schema.Resource{Schema: map[string]*schema.Schema{}},
						},
					},
				},
			},
			"config.basic": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{Schema: map[string]*schema.Schema{}},
			},
		}
		emptyBlockSingletonRes := &schema.Resource{
			Create: dummyCRUD,
			Read:   dummyCRUD,
			Update: dummyCRUD,
			Delete: dummyCRUD,
		}
		if !hasUpdate(emptyBlockSingletonRes, emptyBlockSingletonSchema) {
			t.Errorf("expected resource with mutable empty nested block to have hasUpdate=true")
		}
	})

	resourceMap := map[string]*schema.Resource{
		"google_immutable_resource": {
			Create: dummyCRUD,
			Read:   dummyCRUD,
			Delete: dummyCRUD,
		},
	}

	schemaDiff := diff.SchemaDiff{
		"google_immutable_resource": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"opt_computed": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true, ForceNew: true},
				},
				"opt_covered": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
				"opt_never_unset": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
			},
		},
	}

	tests := []*reader.Test{
		{
			Name: "TestAccImmutable_full",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_immutable_resource": {
							"primary": {
								"name":            "\"n1\"",
								"opt_computed":    "\"oc1\"",
								"opt_covered":     "\"c1\"",
								"opt_never_unset": "\"nu1\"",
							},
						},
					},
				},
			},
		},
		{
			Name: "TestAccImmutable_basic",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_immutable_resource": {
							"primary": {
								"name":            "\"n1\"",
								"opt_never_unset": "\"nu1\"",
							},
						},
					},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, resourceMap, tests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := got["google_immutable_resource"]
	if info == nil {
		t.Fatalf("expected google_immutable_resource to have incomplete fields")
	}

	if diffStr := cmp.Diff([]string{"opt_never_unset"}, info.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}
	if info.Fields["opt_covered"].Modify.ExemptCode != "resource_not_updatable" {
		t.Errorf("expected opt_covered.Modify exemptCode resource_not_updatable, got %+v", info.Fields["opt_covered"].Modify)
	}
	if info.Fields["opt_covered"].Unset.Status != StatusCovered {
		t.Errorf("expected opt_covered.Unset to be COVERED by basic test, got %+v", info.Fields["opt_covered"].Unset)
	}
	if info.Fields["opt_never_unset"].Unset.Status != StatusMissing {
		t.Errorf("expected opt_never_unset.Unset to be MISSING, got %+v", info.Fields["opt_never_unset"].Unset)
	}
}

func TestDetectExhaustiveMissingTests_IncludeChangedFieldsFlag(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"changed_field": {
					Old: &schema.Schema{Type: schema.TypeString, Optional: true},
					New: &schema.Schema{Type: schema.TypeString, Optional: true, Description: "updated"},
				},
			},
		},
	}

	gotDefault, err := DetectExhaustiveMissingTests(schemaDiff, nil, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotDefault) != 0 {
		t.Errorf("expected 0 missing tests when includeChangedFields=false, got %v", gotDefault)
	}

	gotChanged, err := DetectExhaustiveMissingTests(schemaDiff, nil, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotChanged["google_example_resource"]; !ok {
		t.Errorf("expected google_example_resource when includeChangedFields=true")
	}

	md := FormatExhaustiveMarkdown(gotChanged)
	if !strings.Contains(md, "| `changed_field` | Leaf | MISSING | MISSING | MISSING |") {
		t.Errorf("unexpected markdown output:\n%s", md)
	}
}

func TestDetectExhaustiveMissingTests_OwnResourceTestsPreferred(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_my_widget": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"fn_opt": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
				"helper_only_fn": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true, ForceNew: true},
				},
				"untested_opt": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
			},
		},
	}

	// Verify 3-level priority:
	// 1. Own tests ahead of non-own helper tests (even when a helper test has a shorter name).
	// 2. Shorter test name ahead of longer test name within each sub-group
	//    (_full beats alphabetically-earlier _aardvarkLongName; _h beats _alphaLongHelper).
	// 3. Alphabetical tie-breaker at equal length (_full beats _zeta, both len 20).
	tests := []*reader.Test{
		{
			Name: "TestAccAlphaLongHelper_usesWidget",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_my_widget": {
							"helper": {
								"name":           "\"w-helper-long\"",
								"fn_opt":         "\"from-helper-long\"",
								"helper_only_fn": "\"h-long\"",
							},
						},
					},
				},
			},
		},
		{
			Name: "TestAccZeta_h",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_my_widget": {
							"helper": {
								"name":           "\"w-helper-short\"",
								"helper_only_fn": "\"h-short\"",
							},
						},
					},
				},
			},
		},
		{
			Name: "TestAccMyWidget_aardvarkLongName",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_my_widget": {
							"primary": {
								"name":   "\"w-aardvark\"",
								"fn_opt": "\"v-aardvark\"",
							},
						},
					},
				},
			},
		},
		{
			Name: "TestAccMyWidget_zeta",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_my_widget": {
							"primary": {
								"name":   "\"w-zeta\"",
								"fn_opt": "\"v-zeta\"",
							},
						},
					},
				},
			},
		},
		{
			Name: "TestAccMyWidget_full",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_my_widget": {
							"primary": {
								"name":   "\"w1\"",
								"fn_opt": "\"v1\"",
							},
						},
					},
				},
			},
		},
		{
			Name: "TestAccMyWidget_basic",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_my_widget": {
							"primary": {
								"name": "\"w2\"",
							},
						},
					},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, tests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := got["google_my_widget"]
	if info == nil {
		t.Fatalf("expected google_my_widget in results (untested_opt is incomplete)")
	}
	if diffStr := cmp.Diff([]string{"untested_opt"}, info.IncompleteFields); diffStr != "" {
		t.Errorf("expected ONLY untested_opt to be incomplete (-want +got):\n%s", diffStr)
	}
	if info.Fields["name"].Set.CoveredBy != "TestAccMyWidget_full (step 1)" {
		t.Errorf("expected name.Set covered by shortest/alphabetical own test TestAccMyWidget_full (step 1), got %+v", info.Fields["name"].Set)
	}
	if info.Fields["fn_opt"].Set.CoveredBy != "TestAccMyWidget_full (step 1)" {
		t.Errorf("expected fn_opt.Set covered by TestAccMyWidget_full (step 1) ahead of longer _aardvarkLongName and equal-length _zeta, got %+v", info.Fields["fn_opt"].Set)
	}
	if info.Fields["fn_opt"].Unset.CoveredBy != "TestAccMyWidget_basic (step 1)" {
		t.Errorf("expected fn_opt.Unset covered by own test TestAccMyWidget_basic (step 1) ahead of shorter helper TestAccZeta_h, got %+v", info.Fields["fn_opt"].Unset)
	}
	if info.Fields["helper_only_fn"].Set.CoveredBy != "TestAccZeta_h (step 1)" {
		t.Errorf("expected helper_only_fn.Set covered by shorter helper test TestAccZeta_h (step 1), got %+v", info.Fields["helper_only_fn"].Set)
	}
}

func TestDetectExhaustiveMissingTests_NullLifecycleUnsetSetChangedUnset(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"name": {
					New: &schema.Schema{Type: schema.TypeString, Required: true, ForceNew: true},
				},
				"display_name": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"always_null": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
				"null_then_set_only": {
					New: &schema.Schema{Type: schema.TypeString, Optional: true},
				},
			},
		},
	}

	// Step 1: display_name = null (unset)
	// Step 2: display_name = "v1" (set), null_then_set_only = "x" (null -> set must NOT count as MODIFY)
	// Step 3: display_name = "v2" (changed)
	// Step 4: display_name = null (unset via explicit null)
	tests := []*reader.Test{
		{
			Name: "TestAccExample_nullLifecycle",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":               "\"n1\"",
								"display_name":       "null",
								"always_null":        "null",
								"null_then_set_only": "null",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":               "\"n1\"",
								"display_name":       "\"v1\"",
								"always_null":        "null",
								"null_then_set_only": "\"x\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":               "\"n1\"",
								"display_name":       "\"v2\"",
								"always_null":        "null",
								"null_then_set_only": "\"x\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"name":               "\"n1\"",
								"display_name":       "null",
								"always_null":        "null",
								"null_then_set_only": "\"x\"",
							},
						},
					},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, tests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := got["google_example_resource"]
	if info == nil {
		t.Fatalf("expected incomplete fields for always_null and null_then_set_only")
	}
	if diffStr := cmp.Diff([]string{"always_null", "null_then_set_only"}, info.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}
	dn := info.Fields["display_name"]
	if dn.Set.CoveredBy != "TestAccExample_nullLifecycle (step 2)" {
		t.Errorf("display_name.Set = %+v, want step 2", dn.Set)
	}
	if dn.Modify.CoveredBy != "TestAccExample_nullLifecycle (step 2 -> 3)" {
		t.Errorf("display_name.Modify = %+v, want step 2 -> 3", dn.Modify)
	}
	if dn.Unset.CoveredBy != "TestAccExample_nullLifecycle (step 3 -> 4)" || dn.Unset.UnsetMode != UnsetModeNull {
		t.Errorf("display_name.Unset = %+v, want step 3 -> 4 unsetMode=null", dn.Unset)
	}
	if info.Fields["always_null"].Set.Status != StatusMissing {
		t.Errorf("always_null.Set = %+v, want MISSING", info.Fields["always_null"].Set)
	}
	if info.Fields["null_then_set_only"].Modify.Status != StatusMissing {
		t.Errorf("null_then_set_only.Modify = %+v, want MISSING (null -> set is not modify)", info.Fields["null_then_set_only"].Modify)
	}
}

func TestDetectExhaustiveMissingTests_SingleValueEnum(t *testing.T) {
	schemaDiff := diff.SchemaDiff{
		"google_example_resource": {
			Fields: map[string]diff.FieldDiff{
				"required_single_enum": {
					New: &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"ONLY_MODE"}, false),
					},
				},
				"default_single_enum": {
					New: &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "ONLY_MODE",
						ValidateFunc: validation.StringInSlice([]string{"ONLY_MODE", ""}, false),
					},
				},
				"optional_single_enum": {
					New: &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"ONLY_MODE", ""}, false),
					},
				},
				"optional_computed_single_enum": {
					New: &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringInSlice([]string{"ONLY_MODE", ""}, false),
					},
				},
				"multi_enum": {
					New: &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"MODE_A", "MODE_B", ""}, false),
					},
				},
				"list_of_single_enum": {
					New: &schema.Schema{
						Type:        schema.TypeList,
						Optional:    true,
						Description: `Modes. Possible values: ["ONLY_MODE"]`,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"ONLY_MODE"}, false),
						},
					},
				},
				"forcenew_single_enum": {
					New: &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice([]string{"ONLY_MODE", ""}, false),
					},
				},
			},
		},
	}

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	info := got["google_example_resource"]
	if info == nil {
		t.Fatalf("expected google_example_resource in results")
	}

	if gotMod := info.Fields["required_single_enum"].Modify; gotMod.Status != StatusExempt || gotMod.ExemptCode != "single_value_enum" {
		t.Errorf("required_single_enum.Modify = %+v, want EXEMPT (single_value_enum)", gotMod)
	}
	if gotMod := info.Fields["default_single_enum"].Modify; gotMod.Status != StatusExempt || gotMod.ExemptCode != "single_value_enum" {
		t.Errorf("default_single_enum.Modify = %+v, want EXEMPT (single_value_enum)", gotMod)
	}
	if gotUnset := info.Fields["default_single_enum"].Unset; gotUnset.Status != StatusMissing {
		t.Errorf("default_single_enum.Unset = %+v, want MISSING", gotUnset)
	}
	if gotMod := info.Fields["optional_single_enum"].Modify; gotMod.Status != StatusMissing {
		t.Errorf("optional_single_enum.Modify = %+v, want MISSING (unset vs ONLY_MODE are distinct)", gotMod)
	}
	if gotMod := info.Fields["optional_computed_single_enum"].Modify; gotMod.Status != StatusMissing {
		t.Errorf("optional_computed_single_enum.Modify = %+v, want MISSING", gotMod)
	}
	if gotMulti := info.Fields["multi_enum"].Modify; gotMulti.Status != StatusMissing {
		t.Errorf("multi_enum.Modify = %+v, want MISSING", gotMulti)
	}
	if gotList := info.Fields["list_of_single_enum"].Modify; gotList.Status != StatusMissing {
		t.Errorf("list_of_single_enum.Modify = %+v, want MISSING", gotList)
	}
	if gotFN := info.Fields["forcenew_single_enum"].Modify; gotFN.Status != StatusExempt || gotFN.ExemptCode != "force_new" {
		t.Errorf("forcenew_single_enum.Modify = %+v, want EXEMPT (force_new)", gotFN)
	}

	// Verify that an in-place unset -> ONLY_MODE transition satisfies MODIFY on
	// optional 1-value enums (both plain Optional and Optional+Computed), while
	// NOT satisfying MODIFY on multi-value enums.
	tests := []*reader.Test{
		{
			Name: "TestAccExample_singleValueEnumLifecycle",
			Steps: []reader.Step{
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"required_single_enum": "\"ONLY_MODE\"",
								"default_single_enum":  "\"ONLY_MODE\"",
								"forcenew_single_enum": "\"ONLY_MODE\"",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"required_single_enum":          "\"ONLY_MODE\"",
								"optional_single_enum":          "\"ONLY_MODE\"",
								"optional_computed_single_enum": "\"ONLY_MODE\"",
								"multi_enum":                    "\"MODE_A\"",
								"list_of_single_enum":           "[\"ONLY_MODE\"]",
							},
						},
					},
				},
				{
					reader.ResourceBlock: {
						"google_example_resource": {
							"primary": {
								"required_single_enum":          "\"ONLY_MODE\"",
								"optional_computed_single_enum": "\"ONLY_MODE\"",
							},
						},
					},
				},
			},
		},
	}

	gotAfter, err := DetectExhaustiveMissingTests(schemaDiff, nil, tests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	afterInfo := gotAfter["google_example_resource"]
	if afterInfo == nil {
		t.Fatalf("expected multi_enum and list_of_single_enum MODIFY to remain MISSING")
	}
	wantIncomplete := []string{"list_of_single_enum", "multi_enum"}
	if diffStr := cmp.Diff(wantIncomplete, afterInfo.IncompleteFields); diffStr != "" {
		t.Errorf("IncompleteFields mismatch (-want +got):\n%s", diffStr)
	}
	if gotMod := afterInfo.Fields["optional_single_enum"].Modify; gotMod.Status != StatusCovered || gotMod.CoveredBy != "TestAccExample_singleValueEnumLifecycle (step 1 -> 2)" {
		t.Errorf("optional_single_enum.Modify = %+v, want COVERED (step 1 -> 2)", gotMod)
	}
	if gotMod := afterInfo.Fields["optional_computed_single_enum"].Modify; gotMod.Status != StatusCovered || gotMod.CoveredBy != "TestAccExample_singleValueEnumLifecycle (step 1 -> 2)" {
		t.Errorf("optional_computed_single_enum.Modify = %+v, want COVERED (step 1 -> 2)", gotMod)
	}
}

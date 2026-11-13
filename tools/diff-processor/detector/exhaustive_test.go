package detector

import (
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	got, err := DetectExhaustiveMissingTests(schemaDiff, tests, false)
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

	got, err := DetectExhaustiveMissingTests(schemaDiff, tests, false)
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

	got, err := DetectExhaustiveMissingTests(schemaDiff, nil, false)
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

	gotMissing, err := DetectExhaustiveMissingTests(schemaDiff, fullOnlyTests, false)
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

	gotPartial, err := DetectExhaustiveMissingTests(schemaDiff, withBasicNoLabelsUnset, false)
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

	gotComplete, err := DetectExhaustiveMissingTests(schemaDiff, completeTests, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotComplete) != 0 {
		t.Errorf("expected 0 incomplete fields once in-place labels UNSET + ForceNew omission step are present, got:\n%s", FormatExhaustiveMarkdown(gotComplete))
	}
}

func TestDetectExhaustiveMissingTests_NoUpdateResource(t *testing.T) {
	dummyCRUD := func(d *schema.ResourceData, meta interface{}) error { return nil }
	computedDiff := diff.ComputeSchemaDiff(nil, map[string]*schema.Resource{
		"google_all_forcenew_with_deletion_policy": {
			Create: dummyCRUD,
			Read:   dummyCRUD,
			Update: dummyCRUD,
			Delete: dummyCRUD,
			Schema: map[string]*schema.Schema{
				"name":            {Type: schema.TypeString, Required: true, ForceNew: true},
				"deletion_policy": {Type: schema.TypeString, Optional: true, Computed: true},
			},
		},
	})
	if !computedDiff["google_all_forcenew_with_deletion_policy"].NoUpdate {
		t.Fatalf("expected google_all_forcenew_with_deletion_policy to have NoUpdate=true")
	}

	schemaDiff := diff.SchemaDiff{
		"google_immutable_resource": {
			NoUpdate: true,
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

	got, err := DetectExhaustiveMissingTests(schemaDiff, tests, false)
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

	gotDefault, err := DetectExhaustiveMissingTests(schemaDiff, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotDefault) != 0 {
		t.Errorf("expected 0 missing tests when includeChangedFields=false, got %v", gotDefault)
	}

	gotChanged, err := DetectExhaustiveMissingTests(schemaDiff, nil, true)
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

	got, err := DetectExhaustiveMissingTests(schemaDiff, tests, false)
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

package api_test

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
)

func TestResourceMinVersionObj(t *testing.T) {
	t.Parallel()
	p := api.Product{
		Name: "test",
		Versions: []*product.Version{
			&product.Version{
				Name:    "beta",
				BaseUrl: "beta_url",
			},
			&product.Version{
				Name:    "ga",
				BaseUrl: "ga_url",
			},
			&product.Version{
				Name:    "alpha",
				BaseUrl: "alpha_url",
			},
		},
	}

	cases := []struct {
		description string
		obj         api.Resource
		expected    string
	}{
		{
			description: "resource minVersion is empty",
			obj: api.Resource{
				Name:            "test",
				MinVersion:      "",
				ProductMetadata: &p,
			},
			expected: "ga",
		},
		{
			description: "resource minVersion is not empty",
			obj: api.Resource{
				Name:            "test",
				MinVersion:      "beta",
				ProductMetadata: &p,
			},
			expected: "beta",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			versionObj := tc.obj.MinVersionObj()

			if got, want := versionObj.Name, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestResourceNotInVersion(t *testing.T) {
	t.Parallel()
	p := api.Product{
		Name: "test",
		Versions: []*product.Version{
			&product.Version{
				Name:    "beta",
				BaseUrl: "beta_url",
			},
			&product.Version{
				Name:    "ga",
				BaseUrl: "ga_url",
			},
			&product.Version{
				Name:    "alpha",
				BaseUrl: "alpha_url",
			},
		},
	}

	cases := []struct {
		description string
		obj         api.Resource
		input       *product.Version
		expected    bool
	}{
		{
			description: "ga is in version if MinVersion is empty",
			obj: api.Resource{
				Name:            "test",
				MinVersion:      "",
				ProductMetadata: &p,
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: false,
		},
		{
			description: "ga is not in version if MinVersion is beta",
			obj: api.Resource{
				Name:            "test",
				MinVersion:      "beta",
				ProductMetadata: &p,
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := tc.obj.NotInVersion(tc.input), tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestResourceServiceVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         api.Resource
		expected    string
	}{
		{
			description: "BaseUrl does not start with a version",
			obj: api.Resource{
				BaseUrl: "test",
			},
			expected: "",
		},
		{
			description: "BaseUrl starts with / and does not include a version",
			obj: api.Resource{
				BaseUrl: "/test",
			},
			expected: "",
		},
		{
			description: "BaseUrl starts with a version",
			obj: api.Resource{
				BaseUrl: "v3/test",
			},
			expected: "v3",
		},
		{
			description: "BaseUrl starts with a / followed by version",
			obj: api.Resource{
				BaseUrl: "/v3/test",
			},
			expected: "v3",
		},
		{
			description: "CaiBaseUrl does not start with a version",
			obj: api.Resource{
				BaseUrl:    "apis/serving.knative.dev/v1/namespaces/{{project}}/services",
				CaiBaseUrl: "projects/{{project}}/locations/{{location}}/services",
			},
			expected: "",
		},
		{
			description: "CaiBaseUrl starts with a version",
			obj: api.Resource{
				BaseUrl:    "apis/serving.knative.dev/v1/namespaces/{{project}}/services",
				CaiBaseUrl: "v1/projects/{{project}}/locations/{{location}}/services",
			},
			expected: "v1",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			if got, want := tc.obj.ServiceVersion(), tc.expected; got != want {
				t.Errorf("expected %q to be %q", got, want)
			}
		})
	}
}

func TestProviderDefaultFieldsAreSynthesizedAndDeduplicated(t *testing.T) {
	t.Parallel()

	version := &product.Version{Name: "ga", BaseUrl: "https://example.googleapis.com/v1/"}

	cases := []struct {
		description  string
		obj          api.Resource
		listScope    bool
		expectedName map[string]int
	}{
		{
			description: "identity synthesizes project without duplication",
			obj: api.Resource{
				BaseUrl: "projects/{{project}}/foos",
				Parameters: []*api.Type{
					{Name: "project", Type: "String"},
				},
			},
			expectedName: map[string]int{"project": 1},
		},
		{
			description: "list scope synthesizes missing defaults and deduplicates",
			obj: api.Resource{
				BaseUrl: "projects/{{project}}/zones/{{zone}}/foos",
				Parameters: []*api.Type{
					{Name: "project", Type: "String"},
					{Name: "zone", Type: "String", IgnoreRead: true, Exclude: true},
				},
				ProductMetadata: &api.Product{
					Versions: []*product.Version{version},
					Version:  version,
				},
			},
			listScope:    true,
			expectedName: map[string]int{"project": 1, "zone": 1},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			var got []*api.Type
			if tc.listScope {
				got = tc.obj.ListScopeProperties()
			} else {
				got = tc.obj.IdentityProperties()
			}

			counts := map[string]int{}
			for _, p := range got {
				counts[p.Name]++
			}

			for name, want := range tc.expectedName {
				if gotCount := counts[name]; gotCount != want {
					t.Fatalf("expected %s exactly %d time(s), got %d", name, want, gotCount)
				}
			}
		})
	}
}

// TestMagicianLocation verifies that the current package is being executed from within
// the RELATIVE_MAGICIAN_LOCATION ("mmv1/") directory structure. This ensures that references
// to files relative to this location will remain valid even if the repository structure
// changes or the source is downloaded without git metadata.
func TestMagicianLocation(t *testing.T) {
	// Get the path where this test file is located
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current test file path")
	}
	pwd := filepath.Dir(testFilePath)

	// Walk up directories until we either:
	// 1. Find the mmv1 directory
	// 2. Hit the root directory
	dir := pwd
	for {
		// Check if we're in the directory containing mmv1
		if _, err := os.Stat(filepath.Join(dir, "mmv1")); err == nil {
			break
		}

		// When running under bazel runtime paths are relative
		if dir == "." {
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			t.Fatal("Could not find mmv1 directory in parent directories")
		}
		dir = parentDir
	}

	// Check if package is under mmv1
	magicianPath := filepath.Join(dir, api.RELATIVE_MAGICIAN_LOCATION)
	relPath, err := filepath.Rel(magicianPath, pwd)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}
	if strings.HasPrefix(relPath, "..") {
		t.Errorf("Current package is not under %s. Path from magician dir to current dir: %s", api.RELATIVE_MAGICIAN_LOCATION, relPath)
	}
}

func TestHasPostCreateComputedFields(t *testing.T) {
	cases := []struct {
		name, description string
		resource          api.Resource
		want              bool
	}{
		{
			name: "no properties",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
			},
			want: false,
		},
		{
			name: "no computed properties",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name: "resource",
					},
				},
			},
			want: false,
		},
		{
			name: "output-only property",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name:   "field",
						Output: true,
					},
				},
			},
			want: false,
		},
		{
			name: "output-only property in id_format",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name:   "resource",
						Output: true,
					},
				},
			},
			want: true,
		},
		{
			name: "output-only property in id_format with ignore_read",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name:       "resource",
						Output:     true,
						IgnoreRead: true,
					},
				},
			},
			want: false,
		},
		{
			name: "default_from_api property",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name:           "field",
						DefaultFromApi: true,
					},
				},
			},
			want: false,
		},
		{
			name: "default_from_api property in id_format",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name:           "resource",
						DefaultFromApi: true,
					},
				},
			},
			want: true,
		},
		{
			name: "default_from_api property in id_format with ignore_read",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*api.Type{
					{
						Name:           "resource",
						DefaultFromApi: true,
						IgnoreRead:     true,
					},
				},
			},
			want: false,
		},
		{
			name: "converts prop.name to snake case",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource_id}}",
				Properties: []*api.Type{
					{
						Name:   "resourceId",
						Output: true,
					},
				},
			},
			want: true,
		},
		{
			name: "includes fields in self link that aren't in id format",
			resource: api.Resource{
				IdFormat: "projects/{{project}}/resource/{{resource_id}}",
				SelfLink: "{{name}}",
				Properties: []*api.Type{
					{
						Name:   "name",
						Output: true,
					},
				},
			},
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.resource.HasPostCreateComputedFields()
			if got != tc.want {
				t.Errorf("HasPostCreateComputedFields(%q) returned unexpected value. got %t; want %t.", tc.name, got, tc.want)
			}
		})
	}
}

func TestResourceAddExtraFields(t *testing.T) {
	t.Parallel()

	createTestResource := func(name, pn string) *api.Resource {
		r := &api.Resource{
			Name: name,
			ProductMetadata: &api.Product{
				Name: "testproduct",
			},
		}
		r.ProductMetadata.SetCompiler(pn)
		return r
	}

	createTestType := func(name, typeStr string, options ...func(*api.Type)) *api.Type {
		t := &api.Type{
			Name: name,
			Type: typeStr,
		}
		for _, option := range options {
			option(t)
		}
		if t.ResourceMetadata == nil {
			t.ResourceMetadata = &api.Resource{
				Immutable: false,
			}
		}
		return t
	}

	withWriteOnly := func(writeOnly bool) func(*api.Type) {
		return func(t *api.Type) { t.WriteOnly = writeOnly }
	}
	withRequired := func(required bool) func(*api.Type) {
		return func(t *api.Type) { t.Required = required }
	}
	withDescription := func(desc string) func(*api.Type) {
		return func(t *api.Type) { t.Description = desc }
	}
	withProperties := func(props []*api.Type) func(*api.Type) {
		return func(t *api.Type) { t.Properties = props }
	}

	t.Run("WriteOnly property adds companion fields", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		writeOnlyProp := createTestType("password", "String",
			withWriteOnly(true),
			withRequired(true),
			withDescription("A password field"),
		)

		props := []*api.Type{writeOnlyProp}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 3 {
			t.Errorf("Expected 3 properties after adding WriteOnly fields, got %d", len(result))
		}

		if writeOnlyProp.WriteOnly {
			t.Error("Original WriteOnly property should have WriteOnly set to false after processing")
		}
		if writeOnlyProp.Required {
			t.Error("Original WriteOnly property should have Required set to false after processing")
		}

		var foundWoField, foundVersionField bool
		for _, prop := range result {
			if prop.Name == "passwordWo" {
				foundWoField = true
				if !prop.WriteOnly {
					t.Error("passwordWo field should have WriteOnly=true")
				}
			}
			if prop.Name == "passwordWoVersion" {
				foundVersionField = true
				if !prop.ClientSide {
					t.Error("passwordWoVersion field should have ClientSide=true")
				}
			}
		}

		if !foundWoField {
			t.Error("Expected to find passwordWo field")
		}
		if !foundVersionField {
			t.Error("Expected to find passwordWoVersion field")
		}
	})

	t.Run("WriteOnly property doesn't add companion fields for tgc", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraformgoogleconversionnext")
		writeOnlyProp := createTestType("password", "String",
			withWriteOnly(true),
			withRequired(true),
			withDescription("A password field"),
		)

		props := []*api.Type{writeOnlyProp}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 1 {
			t.Errorf("Expected 1 property as WriteOnly fields should not be added, got %d", len(result))
		}

		if writeOnlyProp.WriteOnly {
			t.Error("Original WriteOnly property should have WriteOnly set to false after processing")
		}
	})

	t.Run("KeyValueLabels property adds terraform and effective labels", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		labelsType := &api.Type{
			Name:        "labels",
			Type:        "KeyValueLabels",
			Description: "Resource labels",
		}

		props := []*api.Type{labelsType}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 3 {
			t.Errorf("Expected 3 properties after adding labels fields, got %d", len(result))
		}

		if !labelsType.IgnoreWrite {
			t.Error("Original labels field should have IgnoreWrite=true after processing")
		}
		if !strings.Contains(labelsType.Description, "**Note**") {
			t.Error("Original labels field description should contain note after processing")
		}

		var foundTerraformLabels, foundEffectiveLabels bool
		for _, prop := range result {
			if prop.Name == "terraformLabels" {
				foundTerraformLabels = true
				if prop.Type != "KeyValueTerraformLabels" {
					t.Errorf("terraformLabels should have type KeyValueTerraformLabels, got %s", prop.Type)
				}
			}
			if prop.Name == "effectiveLabels" {
				foundEffectiveLabels = true
				if prop.Type != "KeyValueEffectiveLabels" {
					t.Errorf("effectiveLabels should have type KeyValueEffectiveLabels, got %s", prop.Type)
				}
			}
		}

		if !foundTerraformLabels {
			t.Error("Expected to find terraformLabels field")
		}
		if !foundEffectiveLabels {
			t.Error("Expected to find effectiveLabels field")
		}

		expectedDiff := "tpgresource.SetLabelsDiff"
		if !slices.Contains(resource.CustomDiff, expectedDiff) {
			t.Errorf("Expected CustomDiff to contain %s", expectedDiff)
		}
	})

	t.Run("KeyValueLabels with ExcludeAttributionLabel adds different CustomDiff", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		resource.ExcludeAttributionLabel = true

		labelsType := &api.Type{
			Name: "labels",
			Type: "KeyValueLabels",
		}

		props := []*api.Type{labelsType}
		resource.AddExtraFields(props, nil)

		expectedDiff := "tpgresource.SetLabelsDiffWithoutAttributionLabel"
		if !slices.Contains(resource.CustomDiff, expectedDiff) {
			t.Errorf("Expected CustomDiff to contain %s", expectedDiff)
		}
	})

	t.Run("KeyValueLabels with metadata parent adds metadata CustomDiff", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		parent := &api.Type{Name: "metadata"}

		labelsType := &api.Type{
			Name: "labels",
			Type: "KeyValueLabels",
		}

		props := []*api.Type{labelsType}
		resource.AddExtraFields(props, parent)

		expectedDiff := "tpgresource.SetMetadataLabelsDiff"
		if !slices.Contains(resource.CustomDiff, expectedDiff) {
			t.Errorf("Expected CustomDiff to contain %s", expectedDiff)
		}
	})

	t.Run("KeyValueAnnotations property adds effective annotations", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		annotationsType := &api.Type{
			Name:        "annotations",
			Type:        "KeyValueAnnotations",
			Description: "Resource annotations",
		}

		props := []*api.Type{annotationsType}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 2 {
			t.Errorf("Expected 2 properties after adding annotations fields, got %d", len(result))
		}

		if !annotationsType.IgnoreWrite {
			t.Error("Original annotations field should have IgnoreWrite=true after processing")
		}

		var foundEffectiveAnnotations bool
		for _, prop := range result {
			if prop.Name == "effectiveAnnotations" {
				foundEffectiveAnnotations = true
				if prop.Type != "KeyValueEffectiveLabels" {
					t.Errorf("effectiveAnnotations should have type KeyValueEffectiveLabels, got %s", prop.Type)
				}
			}
		}

		if !foundEffectiveAnnotations {
			t.Error("Expected to find effectiveAnnotations field")
		}

		expectedDiff := "tpgresource.SetAnnotationsDiff"
		if !slices.Contains(resource.CustomDiff, expectedDiff) {
			t.Errorf("Expected CustomDiff to contain %s", expectedDiff)
		}
	})

	t.Run("NestedObject with properties processes recursively", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")

		nestedWriteOnly := createTestType("nestedPassword", "String", withWriteOnly(true))
		nestedObject := createTestType("config", "NestedObject", withProperties([]*api.Type{nestedWriteOnly}))

		props := []*api.Type{nestedObject}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 1 {
			t.Errorf("Expected 1 top-level property, got %d", len(result))
		}

		if len(nestedObject.Properties) != 3 {
			t.Errorf("Expected 3 nested properties after recursive processing, got %d", len(nestedObject.Properties))
		}

		if nestedWriteOnly.WriteOnly {
			t.Error("Nested WriteOnly property should have WriteOnly=false after processing")
		}
	})

	t.Run("Empty NestedObject properties are not processed", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		emptyNestedObject := createTestType("config", "NestedObject", withProperties([]*api.Type{}))

		props := []*api.Type{emptyNestedObject}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 1 {
			t.Errorf("Expected 1 property, got %d", len(result))
		}
		if len(emptyNestedObject.Properties) != 0 {
			t.Errorf("Expected 0 nested properties, got %d", len(emptyNestedObject.Properties))
		}
	})

	t.Run("WriteOnly property already ending with Wo is skipped", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		woProperty := createTestType("passwordWo", "String", withWriteOnly(true))

		props := []*api.Type{woProperty}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 1 {
			t.Errorf("Expected 1 property for Wo-suffixed field, got %d", len(result))
		}

		if !woProperty.WriteOnly {
			t.Error("Wo-suffixed property should remain WriteOnly=true")
		}
	})

	t.Run("Regular properties are passed through unchanged", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")
		regularProp := createTestType("name", "String", withRequired(true))

		props := []*api.Type{regularProp}
		result := resource.AddExtraFields(props, nil)

		if len(result) != 1 {
			t.Errorf("Expected 1 property for regular field, got %d", len(result))
		}

		if result[0] != regularProp {
			t.Error("Regular property should be passed through unchanged")
		}
		if !regularProp.Required {
			t.Error("Regular property Required should be unchanged")
		}
	})

	t.Run("Multiple property types processed correctly", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource", "terraform")

		regularProp := createTestType("name", "String")
		writeOnlyProp := createTestType("password", "String", withWriteOnly(true))
		labelsType := &api.Type{Name: "labels", Type: "KeyValueLabels"}

		props := []*api.Type{regularProp, writeOnlyProp, labelsType}
		result := resource.AddExtraFields(props, nil)

		// Should have: name + password + passwordWo + passwordWoVersion + labels + terraformLabels + effectiveLabels = 7
		if len(result) != 7 {
			t.Errorf("Expected 7 properties total, got %d", len(result))
		}

		names := make(map[string]bool)
		for _, prop := range result {
			names[prop.Name] = true
		}

		expectedNames := []string{"name", "password", "passwordWo", "passwordWoVersion", "labels", "terraformLabels", "effectiveLabels"}
		for _, expected := range expectedNames {
			if !names[expected] {
				t.Errorf("Expected to find property named %s", expected)
			}
		}
	})
}

func TestResource_TestDependencies(t *testing.T) {
	cases := []struct {
		name     string
		resource api.Resource
		want     map[string]string
	}{
		{
			name: "empty",
			resource: api.Resource{
				ProductMetadata: &api.Product{Name: "Apigee"},
				Samples: []*resource.Sample{
					{
						Steps: []*resource.Step{
							{
								TestContextVars: map[string]string{},
								TestHCLText:     "",
							},
						},
					},
				},
				Runtime: api.Runtime{
					ResourcePrefixPkgMap: map[string]string{},
				},
			},
			want: map[string]string{},
		},
		{
			name: "no conflict",
			resource: api.Resource{
				ProductMetadata: &api.Product{Name: "Apigee"},
				Samples: []*resource.Sample{
					{
						Steps: []*resource.Step{
							{
								TestContextVars: map[string]string{
									"network": "compute.BootstrapSubnet",
								},
								TestHCLText: "",
							},
						},
					},
					{
						Steps: []*resource.Step{
							{
								TestContextVars: map[string]string{
									"network": "compute.BootstrapSubnet",
								},
								TestHCLText: "",
							},
						},
					},
				},
				Runtime: api.Runtime{
					ResourcePrefixPkgMap: map[string]string{},
				},
			},
			want: map[string]string{
				"services/compute": "",
			},
		},
		{
			name: "underscore vs empty string conflict",
			resource: api.Resource{
				ProductMetadata: &api.Product{Name: "Apigee"},
				Samples: []*resource.Sample{
					{
						Steps: []*resource.Step{
							{
								TestContextVars: map[string]string{
									"network": "compute.BootstrapSubnet",
								},
								TestHCLText: "",
							},
						},
					},
					{
						Steps: []*resource.Step{
							{
								TestContextVars: map[string]string{},
								TestHCLText:     `resource "google_compute_instance"`,
							},
						},
					},
				},
				Runtime: api.Runtime{
					ResourcePrefixPkgMap: map[string]string{
						"google_compute_": "services/compute",
					},
				},
			},
			want: map[string]string{
				"services/compute": "",
			},
		},
		{
			name: "remove current product",
			resource: api.Resource{
				ProductMetadata: &api.Product{Name: "Compute"},
				Samples: []*resource.Sample{
					{
						Steps: []*resource.Step{
							{
								TestContextVars: map[string]string{
									"network": "compute.BootstrapSubnet",
								},
								TestHCLText: "",
							},
						},
					},
				},
				Runtime: api.Runtime{
					ResourcePrefixPkgMap: map[string]string{},
				},
			},
			want: map[string]string{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.resource.TestDependencies()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestDependencies() mismatch (-want +got:\n%s", diff)
			}
		})
	}
}

// TestIdentityPropertiesFlattenObject ensures that identifiers nested under a
// property marked with flatten_object (e.g. datasetReference.datasetId ->
// dataset_id) are collapsed and included in the identity schema. Without this,
// importing by resource identity panics because the identifier is missing from
// the generated identity schema.
func TestIdentityPropertiesFlattenObject(t *testing.T) {
	t.Parallel()

	res := &api.Resource{
		Name:            "Dataset",
		BaseUrl:         "projects/{{project}}/datasets",
		ImportFormat:    []string{"projects/{{project}}/datasets/{{dataset_id}}"},
		ProductMetadata: &api.Product{Name: "BigQuery"},
	}
	res.Properties = []*api.Type{
		{
			Name:             "datasetReference",
			Type:             "NestedObject",
			FlattenObject:    true,
			ResourceMetadata: res,
			Properties: []*api.Type{
				{
					Name:     "datasetId",
					Type:     "String",
					Required: true,
				},
			},
		},
	}

	got := make([]string, 0)
	for _, p := range res.IdentityProperties() {
		got = append(got, p.Name)
	}

	if !slices.Contains(got, "datasetId") {
		t.Errorf("expected IdentityProperties to include flattened identifier \"datasetId\", got %v", got)
	}
}

func TestSamplePrimaryResourceId(t *testing.T) {
	t.Parallel()

	p := &api.Product{
		Name: "test",
		Versions: []*product.Version{
			{
				Name:    "ga",
				BaseUrl: "ga_url",
			},
			{
				Name:    "beta",
				BaseUrl: "beta_url",
			},
		},
	}

	cases := []struct {
		description string
		resource    api.Resource
		expected    string
	}{
		{
			description: "empty samples returns empty string",
			resource: api.Resource{
				Samples:           []*resource.Sample{},
				ProductMetadata:   p,
				TargetVersionName: "ga",
			},
			expected: "",
		},
		{
			description: "samples with higher min_version returns empty string",
			resource: api.Resource{
				Samples: []*resource.Sample{
					{
						PrimaryResourceId: "beta-res",
						MinVersion:        "beta",
					},
				},
				ProductMetadata:   p,
				TargetVersionName: "ga",
			},
			expected: "",
		},
		{
			description: "valid sample returns primary resource id",
			resource: api.Resource{
				Samples: []*resource.Sample{
					{
						PrimaryResourceId: "ga-res",
					},
				},
				ProductMetadata:   p,
				TargetVersionName: "ga",
			},
			expected: "ga-res",
		},
		{
			description: "only the first sample should be used",
			resource: api.Resource{
				Samples: []*resource.Sample{
					{
						PrimaryResourceId: "first-res",
						MinVersion:        "ga",
					},
					{
						PrimaryResourceId: "second-res",
						MinVersion:        "ga",
					},
				},
				ProductMetadata:   p,
				TargetVersionName: "ga",
			},
			expected: "first-res",
		},
		{
			description: "excludetest should be honored using first non-excluded sample",
			resource: api.Resource{
				Samples: []*resource.Sample{
					{
						PrimaryResourceId: "excluded-res",
						ExcludeTest:       true,
					},
					{
						PrimaryResourceId: "included-res",
						ExcludeTest:       false,
					},
				},
				ProductMetadata:   p,
				TargetVersionName: "ga",
			},
			expected: "included-res",
		},
		{
			description: "fallback to first matching excluded sample when all are excluded",
			resource: api.Resource{
				Samples: []*resource.Sample{
					{
						PrimaryResourceId: "excluded-first",
						ExcludeTest:       true,
					},
					{
						PrimaryResourceId: "excluded-second",
						ExcludeTest:       true,
					},
				},
				ProductMetadata:   p,
				TargetVersionName: "ga",
			},
			expected: "excluded-first",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			if got := tc.resource.SamplePrimaryResourceId(); got != tc.expected {
				t.Errorf("SamplePrimaryResourceId() = %q, want %q", got, tc.expected)
			}
		})
	}
}

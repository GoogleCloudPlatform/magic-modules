package api

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
)

func TestResourceMinVersionObj(t *testing.T) {
	t.Parallel()
	p := Product{
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
		obj         Resource
		expected    string
	}{
		{
			description: "resource minVersion is empty",
			obj: Resource{
				Name:            "test",
				MinVersion:      "",
				ProductMetadata: &p,
			},
			expected: "ga",
		},
		{
			description: "resource minVersion is not empty",
			obj: Resource{
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
	p := Product{
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
		obj         Resource
		input       *product.Version
		expected    bool
	}{
		{
			description: "ga is in version if MinVersion is empty",
			obj: Resource{
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
			obj: Resource{
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
		obj         Resource
		expected    string
	}{
		{
			description: "BaseUrl does not start with a version",
			obj: Resource{
				BaseUrl: "test",
			},
			expected: "",
		},
		{
			description: "BaseUrl starts with / and does not include a version",
			obj: Resource{
				BaseUrl: "/test",
			},
			expected: "",
		},
		{
			description: "BaseUrl starts with a version",
			obj: Resource{
				BaseUrl: "v3/test",
			},
			expected: "v3",
		},
		{
			description: "BaseUrl starts with a / followed by version",
			obj: Resource{
				BaseUrl: "/v3/test",
			},
			expected: "v3",
		},
		{
			description: "CaiBaseUrl does not start with a version",
			obj: Resource{
				BaseUrl:    "apis/serving.knative.dev/v1/namespaces/{{project}}/services",
				CaiBaseUrl: "projects/{{project}}/locations/{{location}}/services",
			},
			expected: "",
		},
		{
			description: "CaiBaseUrl starts with a version",
			obj: Resource{
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

func TestLeafProperties(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Resource
		expected    Type
	}{
		{
			description: "non-nested type",
			obj: Resource{
				BaseUrl: "test",
				Properties: []*Type{
					{
						Name: "basic",
						Type: "String",
					},
				},
			},
			expected: Type{
				Name: "basic",
			},
		},
		{
			description: "nested type",
			obj: Resource{
				BaseUrl: "test",
				Properties: []*Type{
					{
						Name: "root",
						Type: "NestedObject",
						Properties: []*Type{
							{
								Name: "foo",
								Type: "NestedObject",
								Properties: []*Type{
									{
										Name: "bars",
										Type: "Array",
										ItemType: &Type{
											Type: "NestedObject",
											Properties: []*Type{
												{
													Name: "fooBar",
													Type: "String",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: Type{
				Name: "fooBar",
			},
		},
		{
			description: "nested virtual",
			obj: Resource{
				BaseUrl: "test",
				VirtualFields: []*Type{
					{
						Name: "root",
						Type: "NestedObject",
						Properties: []*Type{
							{
								Name: "foo",
								Type: "String",
							},
						},
					},
				},
			},
			expected: Type{
				Name: "foo",
			},
		},
		{
			description: "nested param",
			obj: Resource{
				BaseUrl: "test",
				Parameters: []*Type{
					{
						Name: "root",
						Type: "NestedObject",
						Properties: []*Type{
							{
								Name: "foo",
								Type: "String",
							},
						},
					},
				},
			},
			expected: Type{
				Name: "foo",
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			tc.obj.SetDefault(nil)
			if got, want := tc.obj.LeafProperties(), tc.expected; got[0].Name != want.Name {
				t.Errorf("expected %q to be %q", got[0].Name, want.Name)
			}
		})
	}
}

// TestMagicianLocation verifies that the current package is being executed from within
// the RELATIVE_MAGICIAN_LOCATION ("mmv1/") directory structure. This ensures that references
// to files relative to this location will remain valid even if the repository structure
// changes or the source is downloaded without git metadata.
func TestMagicianLocation(t *testing.T) {
	// Get the current working directory of the test
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Walk up directories until we either:
	// 1. Find the mmv1 directory
	// 2. Hit the root directory
	dir := pwd
	for {
		// Check if we're in the directory containing mmv1
		if _, err := os.Stat(filepath.Join(dir, "mmv1")); err == nil {
			break
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			t.Fatal("Could not find mmv1 directory in parent directories")
		}
		dir = parentDir
	}

	// Check if package is under mmv1
	magicianPath := filepath.Join(dir, RELATIVE_MAGICIAN_LOCATION)
	relPath, err := filepath.Rel(magicianPath, pwd)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}
	if strings.HasPrefix(relPath, "..") {
		t.Errorf("Current package is not under %s. Path from magician dir to current dir: %s", RELATIVE_MAGICIAN_LOCATION, relPath)
	}
}

func TestHasPostCreateComputedFields(t *testing.T) {
	cases := []struct {
		name, description string
		resource          Resource
		want              bool
	}{
		{
			name: "no properties",
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
			},
			want: false,
		},
		{
			name: "no computed properties",
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
					{
						Name: "resource",
					},
				},
			},
			want: false,
		},
		{
			name: "output-only property",
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource_id}}",
				Properties: []*Type{
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
			resource: Resource{
				IdFormat: "projects/{{project}}/resource/{{resource_id}}",
				SelfLink: "{{name}}",
				Properties: []*Type{
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

	createTestResource := func(name string) *Resource {
		return &Resource{
			Name: name,
			ProductMetadata: &Product{
				Name: "testproduct",
			},
		}
	}

	createTestType := func(name, typeStr string, options ...func(*Type)) *Type {
		t := &Type{
			Name: name,
			Type: typeStr,
		}
		for _, option := range options {
			option(t)
		}
		return t
	}

	withWriteOnly := func(writeOnly bool) func(*Type) {
		return func(t *Type) { t.WriteOnly = writeOnly }
	}
	withRequired := func(required bool) func(*Type) {
		return func(t *Type) { t.Required = required }
	}
	withDescription := func(desc string) func(*Type) {
		return func(t *Type) { t.Description = desc }
	}
	withProperties := func(props []*Type) func(*Type) {
		return func(t *Type) { t.Properties = props }
	}

	t.Run("WriteOnly property adds companion fields", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource")
		writeOnlyProp := createTestType("password", "String",
			withWriteOnly(true),
			withRequired(true),
			withDescription("A password field"),
		)

		props := []*Type{writeOnlyProp}
		result := resource.AddExtraFields(props, nil, "")

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

	t.Run("KeyValueLabels property adds terraform and effective labels", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource")
		labelsType := &Type{
			Name:        "labels",
			Type:        "KeyValueLabels",
			Description: "Resource labels",
		}

		props := []*Type{labelsType}
		result := resource.AddExtraFields(props, nil, "")

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

		resource := createTestResource("testresource")
		resource.ExcludeAttributionLabel = true

		labelsType := &Type{
			Name: "labels",
			Type: "KeyValueLabels",
		}

		props := []*Type{labelsType}
		resource.AddExtraFields(props, nil, "")

		expectedDiff := "tpgresource.SetLabelsDiffWithoutAttributionLabel"
		if !slices.Contains(resource.CustomDiff, expectedDiff) {
			t.Errorf("Expected CustomDiff to contain %s", expectedDiff)
		}
	})

	t.Run("KeyValueLabels with metadata parent adds metadata CustomDiff", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource")
		parent := &Type{Name: "metadata"}

		labelsType := &Type{
			Name: "labels",
			Type: "KeyValueLabels",
		}

		props := []*Type{labelsType}
		resource.AddExtraFields(props, parent, "")

		expectedDiff := "tpgresource.SetMetadataLabelsDiff"
		if !slices.Contains(resource.CustomDiff, expectedDiff) {
			t.Errorf("Expected CustomDiff to contain %s", expectedDiff)
		}
	})

	t.Run("KeyValueAnnotations property adds effective annotations", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource")
		annotationsType := &Type{
			Name:        "annotations",
			Type:        "KeyValueAnnotations",
			Description: "Resource annotations",
		}

		props := []*Type{annotationsType}
		result := resource.AddExtraFields(props, nil, "")

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

		resource := createTestResource("testresource")

		nestedWriteOnly := createTestType("nestedPassword", "String", withWriteOnly(true))
		nestedObject := createTestType("config", "NestedObject", withProperties([]*Type{nestedWriteOnly}))

		props := []*Type{nestedObject}
		result := resource.AddExtraFields(props, nil, "")

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

		resource := createTestResource("testresource")
		emptyNestedObject := createTestType("config", "NestedObject", withProperties([]*Type{}))

		props := []*Type{emptyNestedObject}
		result := resource.AddExtraFields(props, nil, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 property, got %d", len(result))
		}
		if len(emptyNestedObject.Properties) != 0 {
			t.Errorf("Expected 0 nested properties, got %d", len(emptyNestedObject.Properties))
		}
	})

	t.Run("WriteOnly property already ending with Wo is skipped", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource")
		woProperty := createTestType("passwordWo", "String", withWriteOnly(true))

		props := []*Type{woProperty}
		result := resource.AddExtraFields(props, nil, "")

		if len(result) != 1 {
			t.Errorf("Expected 1 property for Wo-suffixed field, got %d", len(result))
		}

		if !woProperty.WriteOnly {
			t.Error("Wo-suffixed property should remain WriteOnly=true")
		}
	})

	t.Run("Regular properties are passed through unchanged", func(t *testing.T) {
		t.Parallel()

		resource := createTestResource("testresource")
		regularProp := createTestType("name", "String", withRequired(true))

		props := []*Type{regularProp}
		result := resource.AddExtraFields(props, nil, "")

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

		resource := createTestResource("testresource")

		regularProp := createTestType("name", "String")
		writeOnlyProp := createTestType("password", "String", withWriteOnly(true))
		labelsType := &Type{Name: "labels", Type: "KeyValueLabels"}

		props := []*Type{regularProp, writeOnlyProp, labelsType}
		result := resource.AddExtraFields(props, nil, "")

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

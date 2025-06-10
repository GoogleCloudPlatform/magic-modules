package api

import (
	"os"
	"path/filepath"
	"reflect"
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

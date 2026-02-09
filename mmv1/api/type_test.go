package api

import (
	"reflect"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
)

func TestTypeMinVersionObj(t *testing.T) {
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
		obj         Type
		expected    string
	}{
		{
			description: "type minVersion is empty and resource minVersion is empty",
			obj: Type{
				Name:       "test",
				MinVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
			},
			expected: "ga",
		},
		{
			description: "type minVersion is empty and resource minVersion is beta",
			obj: Type{
				Name:       "test",
				MinVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "beta",
					ProductMetadata: &p,
				},
			},
			expected: "beta",
		},
		{
			description: "type minVersion is not empty",
			obj: Type{
				Name:       "test",
				MinVersion: "beta",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
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

func TestTypeFieldType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		obj         Type
		expected    []string
	}{
		{
			description: "Required field",
			obj: Type{
				Required: true,
			},
			expected: []string{"Required"},
		},
		{
			description: "Optional field",
			obj: Type{
				Required: false,
				Output:   false,
			},
			expected: []string{"Optional"},
		},
		{
			description: "Output field with parent",
			obj: Type{
				Required:       false,
				Output:         true,
				ParentMetadata: &Type{},
			},
			expected: []string{"Output"},
		},
		{
			description: "Output field without parent",
			obj: Type{
				Required:       false,
				Output:         true,
				ParentMetadata: nil,
			},
			expected: []string{},
		},
		{
			description: "WriteOnlyLegacy field",
			obj: Type{
				WriteOnlyLegacy: true,
			},
			expected: []string{"Optional", "Write-Only"},
		},
		{
			description: "WriteOnly field",
			obj: Type{
				WriteOnly: true,
			},
			expected: []string{"Optional", "Write-Only"},
		},
		{
			description: "Beta field in GA resource",
			obj: Type{
				MinVersion: "beta",
				ResourceMetadata: &Resource{
					MinVersion: "ga",
				},
			},
			expected: []string{"Optional", "[Beta](../guides/provider_versions.html.markdown)"},
		},
		{
			description: "Beta field in Beta resource",
			obj: Type{
				MinVersion: "beta",
				ResourceMetadata: &Resource{
					MinVersion: "beta",
				},
			},
			expected: []string{"Optional"},
		},
		{
			description: "GA field in GA resource",
			obj: Type{
				MinVersion: "ga",
				ResourceMetadata: &Resource{
					MinVersion: "ga",
				},
			},
			expected: []string{"Optional"},
		},
		{
			description: "Deprecated field",
			obj: Type{
				DeprecationMessage: "This field is deprecated.",
			},
			expected: []string{"Optional", "Deprecated"},
		},
		{
			description: "All fields set for a required property",
			obj: Type{
				Required:           true,
				WriteOnly:          true,
				MinVersion:         "beta",
				ResourceMetadata:   &Resource{MinVersion: "ga"},
				DeprecationMessage: "This field is deprecated.",
			},
			expected: []string{"Required", "Write-Only", "[Beta](../guides/provider_versions.html.markdown)", "Deprecated"},
		},
		{
			description: "All fields set for an optional property",
			obj: Type{
				WriteOnly:          true,
				MinVersion:         "beta",
				ResourceMetadata:   &Resource{MinVersion: "ga"},
				DeprecationMessage: "This field is deprecated.",
			},
			expected: []string{"Optional", "Write-Only", "[Beta](../guides/provider_versions.html.markdown)", "Deprecated"},
		},
		{
			description: "Output and deprecated",
			obj: Type{
				Output:             true,
				ParentMetadata:     &Type{},
				DeprecationMessage: "This field is deprecated.",
			},
			expected: []string{"Output", "Deprecated"},
		},
		{
			description: "Required and deprecated",
			obj: Type{
				Required:           true,
				ParentMetadata:     &Type{},
				DeprecationMessage: "This field is deprecated.",
			},
			expected: []string{"Required", "Deprecated"},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			fieldType := tc.obj.FieldType()

			if got, want := fieldType, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestTypeExcludeIfNotInVersion(t *testing.T) {
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
		obj         Type
		input       *product.Version
		expected    bool
	}{
		{
			description: "type has Exclude true",
			obj: Type{
				Name:       "test",
				Exclude:    true,
				MinVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: true,
		},
		{
			description: "type has Exclude false and not empty ExactVersion",
			obj: Type{
				Name:         "test",
				MinVersion:   "",
				Exclude:      false,
				ExactVersion: "beta",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "beta",
					ProductMetadata: &p,
				},
			},
			input: &product.Version{
				Name: "ga",
			},
			expected: true,
		},
		{
			description: "type has Exclude false and empty ExactVersion",
			obj: Type{
				Name:         "test",
				MinVersion:   "beta",
				Exclude:      false,
				ExactVersion: "",
				ResourceMetadata: &Resource{
					Name:            "test",
					MinVersion:      "",
					ProductMetadata: &p,
				},
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

			tc.obj.ExcludeIfNotInVersion(tc.input)
			if got, want := tc.obj.Exclude, tc.expected; !reflect.DeepEqual(got, want) {
				t.Errorf("expected %v to be %v", got, want)
			}
		})
	}
}

func TestLineage(t *testing.T) {
	t.Parallel()

	root := Type{
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
	}
	root.SetDefault(&Resource{})

	cases := []struct {
		description string
		obj         Type
		expected    string
	}{
		{
			description: "root type",
			obj:         root,
			expected:    "root",
		},
		{
			description: "sub type",
			obj:         *root.Properties[0],
			expected:    "root.foo",
		},
		{
			description: "array",
			obj:         *root.Properties[0].Properties[0],
			expected:    "root.foo.bars",
		},
		{
			description: "array of objects",
			obj:         *root.Properties[0].Properties[0].ItemType.Properties[0],
			expected:    "root.foo.bars.foo_bar",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			got := strings.Join(tc.obj.Lineage(), ".")
			if got != tc.expected {
				t.Errorf("expected %q to be %q", got, tc.expected)
			}
		})
	}
}

func TestApiLineage(t *testing.T) {
	t.Parallel()

	root := Type{
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
			{
				Name:    "baz",
				ApiName: "bazbaz",
				Type:    "String",
			},
		},
	}
	root.SetDefault(&Resource{})

	fineGrainedRoot := Type{
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
			{
				Name:    "baz",
				ApiName: "bazbaz",
				Type:    "String",
			},
		},
	}
	fineGrainedRoot.SetDefault(&Resource{ApiResourceField: "whatever"})

	cases := []struct {
		description string
		obj         Type
		expected    string
	}{
		{
			description: "root type",
			obj:         root,
			expected:    "root",
		},
		{
			description: "sub type",
			obj:         *root.Properties[0],
			expected:    "root.foo",
		},
		{
			description: "array",
			obj:         *root.Properties[0].Properties[0],
			expected:    "root.foo.bars",
		},
		{
			description: "array of objects",
			obj:         *root.Properties[0].Properties[0].ItemType.Properties[0],
			expected:    "root.foo.bars.fooBar",
		},
		{
			description: "with api name",
			obj:         *root.Properties[1],
			expected:    "root.bazbaz",
		},
		{
			description: "fine-grained root type",
			obj:         fineGrainedRoot,
			expected:    "whatever.root",
		},
		{
			description: "fine-grained sub type",
			obj:         *fineGrainedRoot.Properties[0],
			expected:    "whatever.root.foo",
		},
		{
			description: "fine-grained array",
			obj:         *fineGrainedRoot.Properties[0].Properties[0],
			expected:    "whatever.root.foo.bars",
		},
		{
			description: "fine-grained array of objects",
			obj:         *fineGrainedRoot.Properties[0].Properties[0].ItemType.Properties[0],
			expected:    "whatever.root.foo.bars.fooBar",
		},
		{
			description: "fine-grained with api name",
			obj:         *fineGrainedRoot.Properties[1],
			expected:    "whatever.root.bazbaz",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			got := strings.Join(tc.obj.ApiLineage(), ".")
			if got != tc.expected {
				t.Errorf("expected %q to be %q", got, tc.expected)
			}
		})
	}
}

func TestProviderOnly(t *testing.T) {
	t.Parallel()

	nested := Type{
		Name:       "foo",
		ClientSide: true,
		Type:       "NestedObject",
		Properties: []*Type{
			{
				Name: "bar",
			},
		},
	}
	nested.SetDefault(&Resource{})

	labeled := Resource{
		BaseUrl: "test",
		Properties: []*Type{
			{
				Name: "labels",
				Type: "KeyValueLabels",
			},
		},
	}
	labeled.Properties = labeled.AddExtraFields(labeled.PropertiesWithExcluded(), nil)
	labeled.SetDefault(nil)

	cases := []struct {
		description string
		obj         Type
		expected    bool
	}{
		{
			description: "normal",
			obj: Type{
				Name: "foo",
			},
			expected: false,
		},
		{
			description: "url param",
			obj: Type{
				Name:         "foo",
				UrlParamOnly: true,
			},
			expected: true,
		},
		{
			description: "virtual",
			obj: Type{
				Name: "foo",
				// Virtual fields will have this field set during SetDefault()
				ClientSide: true,
			},
			expected: true,
		},
		{
			description: "child of virtual",
			obj:         *nested.Properties[0],
			expected:    true,
		},
		{
			description: "terraform labels",
			// Terraform labels are added first
			obj:      *labeled.Properties[1],
			expected: true,
		},
		{
			description: "effective labels",
			// Effective labels are added second
			obj:      *labeled.Properties[2],
			expected: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			got := tc.obj.ProviderOnly()
			if got != tc.expected {
				t.Errorf("expected %t to be %t", got, tc.expected)
			}
		})
	}
}

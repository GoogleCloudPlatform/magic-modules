package api

import (
	"reflect"
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

func TestMetadataLineage(t *testing.T) {
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

			got := tc.obj.MetadataLineage()
			if got != tc.expected {
				t.Errorf("expected %q to be %q", got, tc.expected)
			}
		})
	}
}

func TestMetadataApiLineage(t *testing.T) {
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
		{
			description: "with api name",
			obj:         *root.Properties[1],
			expected:    "root.bazbaz",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			got := tc.obj.MetadataApiLineage()
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
	labeled.Properties = labeled.AddLabelsRelatedFields(labeled.PropertiesWithExcluded(), nil)
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

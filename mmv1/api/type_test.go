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

	t.Run("map value_type nested property is excluded if not in version", func(t *testing.T) {
		t.Parallel()

		makeMapProp := func() (*Type, *Type) {
			res := &Resource{
				Name:            "test",
				ProductMetadata: &p,
			}
			nestedBetaProp := &Type{
				Name:             "betaProp",
				Type:             "String",
				MinVersion:       "beta",
				ResourceMetadata: res,
			}
			mapProp := &Type{
				Name:             "mapProp",
				Type:             "Map",
				ResourceMetadata: res,
				ValueType: &Type{
					Type:             "NestedObject",
					ResourceMetadata: res,
					Properties:       []*Type{nestedBetaProp},
				},
			}
			return mapProp, nestedBetaProp
		}

		gaMapProp, gaNestedBetaProp := makeMapProp()
		gaMapProp.ExcludeIfNotInVersion(&product.Version{Name: "ga"})
		if gaMapProp.Exclude {
			t.Errorf("expected mapProp.Exclude to be false in ga")
		}
		if !gaNestedBetaProp.Exclude {
			t.Errorf("expected nestedBetaProp.Exclude to be true in ga")
		}

		betaMapProp, betaNestedBetaProp := makeMapProp()
		betaMapProp.ExcludeIfNotInVersion(&product.Version{Name: "beta"})
		if betaMapProp.Exclude {
			t.Errorf("expected mapProp.Exclude to be false in beta")
		}
		if betaNestedBetaProp.Exclude {
			t.Errorf("expected nestedBetaProp.Exclude to be false in beta")
		}
	})
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

func TestTypeValidateConflicts(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
			{Name: "beta", BaseUrl: "beta_url"},
		},
	}

	t.Run("valid conflicts path passes validation", func(t *testing.T) {
		t.Parallel()

		res := &Resource{
			Name: "TestResource",
			Properties: []*Type{
				{
					Name:      "foo",
					Type:      "String",
					Conflicts: []string{"bar"},
				},
				{
					Name:      "bar",
					Type:      "String",
					Conflicts: []string{"foo"},
				},
			},
		}
		res.SetDefault(p)

		if errs := res.Properties[0].Validate(res.Name); len(errs) != 0 {
			t.Errorf("expected no validation errors, got %v", errs)
		}
	})

	t.Run("bad conflicts path that resolves to empty ConflictsWith returns error", func(t *testing.T) {
		t.Parallel()

		res := &Resource{
			Name: "TestResource",
			Properties: []*Type{
				{
					Name: "nested",
					Type: "NestedObject",
					Properties: []*Type{
						{
							Name:      "foo",
							Type:      "String",
							Conflicts: []string{"bar"},
						},
						{
							Name:      "bar",
							Type:      "String",
							Conflicts: []string{"foo"},
						},
					},
				},
			},
		}
		res.SetDefault(p)

		errs := res.Properties[0].Validate(res.Name)
		if len(errs) != 2 {
			t.Fatalf("expected 2 validation errors, got %d: %v", len(errs), errs)
		}
		for _, err := range errs {
			if !strings.Contains(err.Error(), "resolved to an empty ConflictsWith") {
				t.Errorf("expected error to mention empty ConflictsWith, got %v", err)
			}
		}
	})

	t.Run("conflicts with one valid and one invalid entry returns error", func(t *testing.T) {
		t.Parallel()

		res := &Resource{
			Name: "TestResource",
			Properties: []*Type{
				{
					Name: "nested",
					Type: "NestedObject",
					Properties: []*Type{
						{
							Name:      "foo",
							Type:      "String",
							Conflicts: []string{"nested.0.bar", "nonexistent"},
						},
						{
							Name:      "bar",
							Type:      "String",
							Conflicts: []string{"nested.0.foo"},
						},
					},
				},
			},
		}
		res.SetDefault(p)

		errs := res.Properties[0].Validate(res.Name)
		if len(errs) != 1 {
			t.Fatalf("expected 1 validation error, got %d: %v", len(errs), errs)
		}
		if !strings.Contains(errs[0].Error(), "resolved to an empty ConflictsWith") {
			t.Errorf("expected error to mention empty ConflictsWith, got %v", errs[0])
		}
		if !strings.Contains(errs[0].Error(), "nonexistent") {
			t.Errorf("expected error to mention the invalid entry 'nonexistent', got %v", errs[0])
		}
	})

	t.Run("GA field with mutual beta conflict is culled in GA and preserved in beta", func(t *testing.T) {
		t.Parallel()

		newRes := func() *Resource {
			res := &Resource{
				Name: "TestResource",
				Properties: []*Type{
					{
						Name:      "gaField",
						Type:      "String",
						Conflicts: []string{"beta_field"},
					},
					{
						Name:       "betaField",
						Type:       "String",
						MinVersion: "beta",
						Conflicts:  []string{"ga_field"},
					},
				},
			}
			res.SetDefault(p)
			return res
		}

		gaRes := newRes()
		gaRes.ExcludeIfNotInVersion(&product.Version{Name: "ga"})
		if errs := gaRes.Properties[0].Validate(gaRes.Name); len(errs) != 0 {
			t.Errorf("expected no validation error for gaField with valid mutual beta conflict in GA, got %v", errs)
		}
		if got := gaRes.Properties[0].Conflicting(); len(got) != 0 {
			t.Errorf("expected gaField.Conflicting() to be culled to empty in GA, got %v", got)
		}
		if errs := gaRes.Properties[1].Validate(gaRes.Name); len(errs) != 0 {
			t.Errorf("expected no validation error for excluded betaField in GA, got %v", errs)
		}

		betaRes := newRes()
		betaRes.ExcludeIfNotInVersion(&product.Version{Name: "beta"})
		if got := betaRes.Properties[0].Conflicting(); !reflect.DeepEqual(got, []string{"beta_field"}) {
			t.Errorf("expected gaField.Conflicting() to be [beta_field] in beta, got %v", got)
		}
		if errs := betaRes.Properties[0].Validate(betaRes.Name); len(errs) != 0 {
			t.Errorf("expected no validation error for gaField in beta, got %v", errs)
		}
		if errs := betaRes.Properties[1].Validate(betaRes.Name); len(errs) != 0 {
			t.Errorf("expected no validation error for betaField in beta, got %v", errs)
		}
	})

	t.Run("GA virtual field with mutual beta virtual field conflict is culled in GA and preserved in beta", func(t *testing.T) {
		t.Parallel()

		newRes := func() *Resource {
			res := &Resource{
				Name:        "TestResource",
				Description: "Test resource",
				CreateVerb:  "POST",
				ReadVerb:    "GET",
				UpdateVerb:  "POST",
				DeleteVerb:  "DELETE",
				Properties: []*Type{
					{
						Name:        "name",
						Type:        "String",
						Description: "Resource name",
					},
				},
				VirtualFields: []*Type{
					{
						Name:        "allowed_actions",
						Type:        "Array",
						Description: "Beta virtual field",
						MinVersion:  "beta",
						ItemType: &Type{
							Type: "String",
						},
						Conflicts: []string{"minimal_action"},
					},
					{
						Name:         "minimal_action",
						Type:         "String",
						Description:  "GA virtual field",
						DefaultValue: "NONE",
						Conflicts:    []string{"allowed_actions"},
					},
				},
			}
			res.SetDefault(p)
			return res
		}

		gaRes := newRes()
		// Call ExcludeIfNotInVersion twice (mirroring CountProductsAndResources + GenerateObjects)
		// to ensure populateMaxVersionSchema preserves the unfiltered virtual fields.
		gaRes.ExcludeIfNotInVersion(&product.Version{Name: "ga"})
		gaRes.ExcludeIfNotInVersion(&product.Version{Name: "ga"})

		if !gaRes.VirtualFields[0].Exclude {
			t.Errorf("expected beta virtual field allowed_actions to have Exclude=true in GA")
		}
		if gaRes.VirtualFields[1].Exclude {
			t.Errorf("expected GA virtual field minimal_action to have Exclude=false in GA")
		}
		if got := len(gaRes.UserVirtualFields()); got != 1 || gaRes.UserVirtualFields()[0].Name != "minimal_action" {
			t.Errorf("expected UserVirtualFields() to only contain minimal_action in GA, got %v", gaRes.UserVirtualFields())
		}
		if errs := gaRes.Validate(); len(errs) != 0 {
			t.Errorf("expected no validation errors for gaRes in GA, got %v", errs)
		}
		if got := gaRes.VirtualFields[1].Conflicting(); len(got) != 0 {
			t.Errorf("expected minimal_action.Conflicting() to be culled to empty in GA, got %v", got)
		}

		betaRes := newRes()
		betaRes.ExcludeIfNotInVersion(&product.Version{Name: "beta"})
		betaRes.ExcludeIfNotInVersion(&product.Version{Name: "beta"})

		if betaRes.VirtualFields[0].Exclude || betaRes.VirtualFields[1].Exclude {
			t.Errorf("expected both virtual fields to have Exclude=false in beta")
		}
		if got := len(betaRes.UserVirtualFields()); got != 2 {
			t.Errorf("expected 2 UserVirtualFields in beta, got %d", got)
		}
		if errs := betaRes.Validate(); len(errs) != 0 {
			t.Errorf("expected no validation errors for betaRes in beta, got %v", errs)
		}
		if got := betaRes.VirtualFields[0].Conflicting(); !reflect.DeepEqual(got, []string{"minimal_action"}) {
			t.Errorf("expected allowed_actions.Conflicting() to be [minimal_action] in beta, got %v", got)
		}
		if got := betaRes.VirtualFields[1].Conflicting(); !reflect.DeepEqual(got, []string{"allowed_actions"}) {
			t.Errorf("expected minimal_action.Conflicting() to be [allowed_actions] in beta, got %v", got)
		}
	})

	t.Run("GA field conflicting with non-mutual beta field returns error in GA", func(t *testing.T) {
		t.Parallel()

		res := &Resource{
			Name: "TestResource",
			Properties: []*Type{
				{
					Name:      "gaField",
					Type:      "String",
					Conflicts: []string{"beta_field"},
				},
				{
					Name:       "otherGaField",
					Type:       "String",
					MinVersion: "ga",
				},
				{
					Name:       "betaField",
					Type:       "String",
					MinVersion: "beta",
					Conflicts:  []string{"other_ga_field"},
				},
			},
		}
		res.SetDefault(p)
		res.ExcludeIfNotInVersion(&product.Version{Name: "ga"})

		if errs := res.Properties[0].Validate(res.Name); len(errs) == 0 {
			t.Errorf("expected validation error for gaField conflicting with non-mutual betaField in GA, got none")
		}
	})

	t.Run("field conflicting with unconditionally excluded field returns error", func(t *testing.T) {
		t.Parallel()

		res := &Resource{
			Name: "TestResource",
			Properties: []*Type{
				{
					Name:      "gaField",
					Type:      "String",
					Conflicts: []string{"excluded_field"},
				},
				{
					Name:      "excludedField",
					Type:      "String",
					Exclude:   true,
					Conflicts: []string{"ga_field"},
				},
			},
		}
		res.SetDefault(p)
		res.ExcludeIfNotInVersion(&product.Version{Name: "ga"})

		if errs := res.Properties[0].Validate(res.Name); len(errs) == 0 {
			t.Errorf("expected validation error for gaField conflicting with unconditionally excluded field, got none")
		}
	})
}

func TestTypeValidateSiblingConstraints(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
			{Name: "beta", BaseUrl: "beta_url"},
		},
	}

	kinds := []struct {
		yamlAttr   string
		schemaAttr string
		set        func(prop *Type, list []string)
		get        func(prop *Type) []string
	}{
		{
			yamlAttr:   "at_least_one_of",
			schemaAttr: "AtLeastOneOf",
			set:        func(prop *Type, list []string) { prop.AtLeastOneOf = list },
			get:        func(prop *Type) []string { return prop.AtLeastOneOfList() },
		},
		{
			yamlAttr:   "exactly_one_of",
			schemaAttr: "ExactlyOneOf",
			set:        func(prop *Type, list []string) { prop.ExactlyOneOf = list },
			get:        func(prop *Type) []string { return prop.ExactlyOneOfList() },
		},
		{
			yamlAttr:   "required_with",
			schemaAttr: "RequiredWith",
			set:        func(prop *Type, list []string) { prop.RequiredWith = list },
			get:        func(prop *Type) []string { return prop.RequiredWithList() },
		},
	}

	for _, kind := range kinds {
		kind := kind
		t.Run(kind.yamlAttr, func(t *testing.T) {
			t.Parallel()

			t.Run("valid nested path passes validation", func(t *testing.T) {
				t.Parallel()

				foo := &Type{Name: "foo", Type: "String"}
				bar := &Type{Name: "bar", Type: "String"}
				kind.set(foo, []string{"nested.0.foo", "nested.0.bar"})
				kind.set(bar, []string{"nested.0.foo", "nested.0.bar"})

				res := &Resource{
					Name: "TestResource",
					Properties: []*Type{
						{
							Name:       "nested",
							Type:       "NestedObject",
							Properties: []*Type{foo, bar},
						},
					},
				}
				res.SetDefault(p)

				if errs := res.Properties[0].Validate(res.Name); len(errs) != 0 {
					t.Errorf("expected no validation errors, got %v", errs)
				}
			})

			t.Run("unqualified nested relative path returns error", func(t *testing.T) {
				t.Parallel()

				foo := &Type{Name: "foo", Type: "String"}
				bar := &Type{Name: "bar", Type: "String"}
				kind.set(foo, []string{"foo", "bar"})
				kind.set(bar, []string{"foo", "bar"})

				res := &Resource{
					Name: "TestResource",
					Properties: []*Type{
						{
							Name:       "nested",
							Type:       "NestedObject",
							Properties: []*Type{foo, bar},
						},
					},
				}
				res.SetDefault(p)

				errs := res.Properties[0].Validate(res.Name)
				if len(errs) == 0 {
					t.Fatalf("expected validation errors for unqualified nested %s, got none", kind.yamlAttr)
				}
				wantMsg := "resolved to an empty " + kind.schemaAttr
				for _, err := range errs {
					if !strings.Contains(err.Error(), wantMsg) {
						t.Errorf("expected error containing %q, got %v", wantMsg, err)
					}
				}
			})

			t.Run("cross-version beta sibling is culled in GA and preserved in beta", func(t *testing.T) {
				t.Parallel()

				newRes := func() *Resource {
					gaOne := &Type{Name: "gaOne", Type: "String"}
					gaTwo := &Type{Name: "gaTwo", Type: "String"}
					betaThree := &Type{Name: "betaThree", Type: "String", MinVersion: "beta"}
					all := []string{"ga_one", "ga_two", "beta_three"}
					kind.set(gaOne, all)
					kind.set(gaTwo, all)
					kind.set(betaThree, all)

					res := &Resource{
						Name:       "TestResource",
						Properties: []*Type{gaOne, gaTwo, betaThree},
					}
					res.Properties = res.AddExtraFields(res.PropertiesWithExcluded(), nil)
					res.SetDefault(p)
					return res
				}

				gaRes := newRes()
				gaRes.ExcludeIfNotInVersion(&product.Version{Name: "ga"})
				for _, prop := range gaRes.Properties {
					if errs := prop.Validate(gaRes.Name); len(errs) != 0 {
						t.Errorf("expected no GA validation errors on %s, got %v", prop.Name, errs)
					}
				}
				if got := kind.get(gaRes.Properties[0]); !reflect.DeepEqual(got, []string{"ga_one", "ga_two"}) {
					t.Errorf("expected GA %s list [ga_one ga_two], got %v", kind.yamlAttr, got)
				}

				betaRes := newRes()
				betaRes.ExcludeIfNotInVersion(&product.Version{Name: "beta"})
				for _, prop := range betaRes.Properties {
					if errs := prop.Validate(betaRes.Name); len(errs) != 0 {
						t.Errorf("expected no beta validation errors on %s, got %v", prop.Name, errs)
					}
				}
				if got := kind.get(betaRes.Properties[0]); !reflect.DeepEqual(got, []string{"beta_three", "ga_one", "ga_two"}) {
					t.Errorf("expected beta %s list [beta_three ga_one ga_two], got %v", kind.yamlAttr, got)
				}
			})

			t.Run("unconditionally excluded target returns error", func(t *testing.T) {
				t.Parallel()

				gaField := &Type{Name: "gaField", Type: "String"}
				exclField := &Type{Name: "excludedField", Type: "String", Exclude: true}
				kind.set(gaField, []string{"ga_field", "excluded_field"})
				kind.set(exclField, []string{"ga_field", "excluded_field"})

				res := &Resource{
					Name:       "TestResource",
					Properties: []*Type{gaField, exclField},
				}
				res.SetDefault(p)
				res.ExcludeIfNotInVersion(&product.Version{Name: "ga"})

				if errs := res.Properties[0].Validate(res.Name); len(errs) == 0 {
					t.Errorf("expected validation error when referencing unconditionally excluded field in %s, got none", kind.yamlAttr)
				}
			})

			t.Run("one valid and one invalid entry triggers error even with valid beta sibling", func(t *testing.T) {
				t.Parallel()

				gaOne := &Type{Name: "gaOne", Type: "String"}
				gaTwo := &Type{Name: "gaTwo", Type: "String"}
				betaThree := &Type{Name: "betaThree", Type: "String", MinVersion: "beta"}
				list := []string{"ga_one", "ga_two", "beta_three", "nonexistent_field"}
				kind.set(gaOne, list)
				kind.set(gaTwo, list)
				kind.set(betaThree, list)

				res := &Resource{
					Name:       "TestResource",
					Properties: []*Type{gaOne, gaTwo, betaThree},
				}
				res.Properties = res.AddExtraFields(res.PropertiesWithExcluded(), nil)
				res.SetDefault(p)
				res.ExcludeIfNotInVersion(&product.Version{Name: "ga"})

				if errs := res.Properties[0].Validate(res.Name); len(errs) == 0 {
					t.Errorf("expected validation error for mixed valid + invalid %s entries, got none", kind.yamlAttr)
				}
			})
		})
	}
}

func TestTypeValidateDefaultFromApi(t *testing.T) {
	t.Parallel()

	cases := []struct {
		description string
		prop        *Type
		resource    string
		wantError   string
	}{
		{
			description: "default_from_api and required not allowlisted returns error",
			prop: &Type{
				Name:           "customField",
				Type:           "String",
				DefaultFromApi: true,
				Required:       true,
			},
			resource:  "MyResource",
			wantError: "property custom_field cannot set both default_from_api and required in resource MyResource",
		},
		{
			description: "default_from_api and output not allowlisted returns error",
			prop: &Type{
				Name:           "customField",
				Type:           "String",
				DefaultFromApi: true,
				Output:         true,
			},
			resource:  "MyResource",
			wantError: "property custom_field cannot set both default_from_api and output in resource MyResource",
		},
		{
			description: "default_from_api and output allowlisted succeeds",
			prop: &Type{
				Name:           "ruleId",
				Type:           "String",
				DefaultFromApi: true,
				Output:         true,
			},
			resource:  "Rule",
			wantError: "",
		},
	}

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       tc.resource,
				Properties: []*Type{tc.prop},
			}
			res.SetDefault(p)

			errs := tc.prop.Validate(tc.resource)
			if tc.wantError == "" {
				if len(errs) > 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), tc.wantError) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected error containing %q, got: %v", tc.wantError, errs)
			}
		})
	}
}

func TestTypeValidateOutputOnlyEnum(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	cases := []struct {
		description string
		prop        *Type
		wantError   string
	}{
		{
			description: "output-only Enum returns error",
			prop: &Type{
				Name:       "state",
				Type:       "Enum",
				Output:     true,
				EnumValues: []string{"ACTIVE", "INACTIVE"},
			},
			wantError: "property state cannot be an output-only Enum in resource TestResource; use String instead",
		},
		{
			description: "output-only Array of Enum returns error",
			prop: &Type{
				Name:   "states",
				Type:   "Array",
				Output: true,
				ItemType: &Type{
					Type:       "Enum",
					EnumValues: []string{"ACTIVE", "INACTIVE"},
				},
			},
			wantError: "property states cannot be an output-only Enum in resource TestResource; use String instead",
		},
		{
			description: "input Enum succeeds",
			prop: &Type{
				Name:       "mode",
				Type:       "Enum",
				EnumValues: []string{"FAST", "SLOW"},
			},
			wantError: "",
		},
		{
			description: "output String succeeds",
			prop: &Type{
				Name:   "state",
				Type:   "String",
				Output: true,
			},
			wantError: "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       "TestResource",
				Properties: []*Type{tc.prop},
			}
			res.SetDefault(p)

			errs := tc.prop.Validate(res.Name)
			if tc.wantError == "" {
				if len(errs) > 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), tc.wantError) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected error containing %q, got: %v", tc.wantError, errs)
			}
		})
	}
}

func TestTypeJsonValue(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	prop := &Type{
		Name: "payload",
		Type: "JsonValue",
	}
	res := &Resource{
		Name:       "TestResource",
		Properties: []*Type{prop},
	}
	res.SetDefault(p)

	if errs := prop.Validate(res.Name); len(errs) > 0 {
		t.Fatalf("expected no errors for JsonValue, got: %v", errs)
	}
	if got := prop.TFType(prop.Type); got != "schema.TypeString" {
		t.Errorf("TFType(JsonValue) = %q, want %q", got, "schema.TypeString")
	}
	if got := prop.GetFWType(); got != "String" {
		t.Errorf("GetFWType() = %q, want %q", got, "String")
	}
	if !prop.IsJsonField() {
		t.Errorf("IsJsonField() = false, want true")
	}

	disallowCases := []struct {
		description string
		prop        *Type
	}{
		{
			description: "custom_expand json_value.tmpl disallowed",
			prop: &Type{
				Name:         "payload",
				Type:         "String",
				CustomExpand: "templates/terraform/custom_expand/json_value.tmpl",
			},
		},
		{
			description: "custom_expand json_schema.tmpl disallowed",
			prop: &Type{
				Name:         "payload",
				Type:         "String",
				CustomExpand: "templates/terraform/custom_expand/json_schema.tmpl",
			},
		},
		{
			description: "custom_flatten json_schema.tmpl disallowed",
			prop: &Type{
				Name:          "payload",
				Type:          "String",
				Output:        true,
				CustomFlatten: "templates/terraform/custom_flatten/json_schema.tmpl",
			},
		},
	}

	for _, tc := range disallowCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			r := &Resource{
				Name:       "TestResource",
				Properties: []*Type{tc.prop},
			}
			r.SetDefault(p)
			errs := tc.prop.Validate(r.Name)
			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), "use type JsonValue instead") {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected error containing %q, got: %v", "use type JsonValue instead", errs)
			}
		})
	}
}

func TestTypeValidateRequiredUnderOutput(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	cases := []struct {
		description string
		prop        *Type
		wantError   string
	}{
		{
			description: "direct required child under output NestedObject returns error",
			prop: &Type{
				Name:   "status",
				Type:   "NestedObject",
				Output: true,
				Properties: []*Type{
					{Name: "code", Type: "String", Required: true},
				},
			},
			wantError: "property status.code cannot be required because it has an output ancestor in resource TestResource",
		},
		{
			description: "multi-level nested required descendant under output Array returns error",
			prop: &Type{
				Name:   "runtimeConfigs",
				Type:   "Array",
				Output: true,
				ItemType: &Type{
					Type: "NestedObject",
					Properties: []*Type{
						{
							Name: "gkeWorkload",
							Type: "NestedObject",
							Properties: []*Type{
								{Name: "cluster", Type: "String", Required: true},
							},
						},
					},
				},
			},
			wantError: "property runtime_configs.gke_workload.cluster cannot be required because it has an output ancestor in resource TestResource",
		},
		{
			description: "required descendant under output Map ValueType returns error",
			prop: &Type{
				Name:    "entries",
				Type:    "Map",
				KeyName: "key",
				Output:  true,
				ValueType: &Type{
					Type: "NestedObject",
					Properties: []*Type{
						{Name: "value", Type: "String", Required: true},
					},
				},
			},
			wantError: "property entries.value cannot be required because it has an output ancestor in resource TestResource",
		},
		{
			description: "output descendant under output parent succeeds",
			prop: &Type{
				Name:   "status",
				Type:   "NestedObject",
				Output: true,
				Properties: []*Type{
					{Name: "code", Type: "String", Output: true},
				},
			},
		},
		{
			description: "required child with no output ancestor succeeds",
			prop: &Type{
				Name: "config",
				Type: "NestedObject",
				Properties: []*Type{
					{Name: "name", Type: "String", Required: true},
				},
			},
		},
		{
			description: "non-required non-output descendant under output parent succeeds",
			prop: &Type{
				Name:   "status",
				Type:   "NestedObject",
				Output: true,
				Properties: []*Type{
					{Name: "message", Type: "String"},
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       "TestResource",
				Properties: []*Type{tc.prop},
			}
			res.SetDefault(p)

			errs := tc.prop.Validate(res.Name)
			if tc.wantError == "" {
				if len(errs) > 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			if len(errs) != 1 || errs[0].Error() != tc.wantError {
				t.Fatalf("expected error %q, got: %v", tc.wantError, errs)
			}
		})
	}

	t.Run("hasOutputAncestor is nil-safe and cycle-safe", func(t *testing.T) {
		t.Parallel()
		var nilType *Type
		if nilType.hasOutputAncestor() {
			t.Errorf("expected nil Type to return false")
		}
		a := &Type{Name: "a", Type: "NestedObject"}
		b := &Type{Name: "b", Type: "NestedObject", ParentMetadata: a}
		a.ParentMetadata = b
		if a.hasOutputAncestor() {
			t.Errorf("expected cycle without output ancestor to return false")
		}
		b.Output = true
		if !a.hasOutputAncestor() {
			t.Errorf("expected cycle with output ancestor to return true")
		}
	})
}

func TestTypeValidateInternalRefs(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name:     "TestProduct",
		Versions: []*product.Version{{Name: "ga", BaseUrl: "ga_url"}},
	}

	cases := []struct {
		name      string
		prop      *Type
		wantError string
	}{
		{
			name:      "corp.google.com subdomain in description",
			prop:      &Type{Name: "field", Type: "String", Description: "See https://clouddocs.devsite.corp.google.com/docs."},
			wantError: `property field ` + "`description`" + ` in resource TestResource must not contain internal Google reference "corp.google.com"`,
		},
		{
			name:      "goto.google.com in description",
			prop:      &Type{Name: "field", Type: "String", Description: "See https://goto.google.com/my-doc."},
			wantError: `internal Google reference "goto.google.com"`,
		},
		{
			name:      "http://go/ in description",
			prop:      &Type{Name: "field", Type: "String", Description: "See http://go/internal-doc."},
			wantError: `internal Google reference "go/"`,
		},
		{
			name:      "https://go/ in description",
			prop:      &Type{Name: "field", Type: "String", Description: "See https://go/internal-doc."},
			wantError: `internal Google reference "go/"`,
		},
		{
			name:      "bare go/ in description",
			prop:      &Type{Name: "field", Type: "String", Description: "Drawn from go/cloud-ai-gcp-pool."},
			wantError: `internal Google reference "go/"`,
		},
		{
			name:      "bare b/ bug ID in description",
			prop:      &Type{Name: "field", Type: "String", Description: "See b/360405077#comment8."},
			wantError: `internal Google reference "b/"`,
		},
		{
			name:      "https://b/ bug URL in description",
			prop:      &Type{Name: "field", Type: "String", Description: "See https://b/123456."},
			wantError: `internal Google reference "b/"`,
		},
		{
			name:      "//depot/ in description",
			prop:      &Type{Name: "field", Type: "String", Description: "Defined at //depot/foo/bar.proto."},
			wantError: `internal Google reference "//depot/"`,
		},
		{
			name:      "google3/ in description",
			prop:      &Type{Name: "field", Type: "String", Description: "Defined at google3/cloud/foo/bar.proto."},
			wantError: `internal Google reference "google3/"`,
		},
		{
			name:      "case-insensitive corp.google.com",
			prop:      &Type{Name: "field", Type: "String", Description: "See HTTPS://B.CORP.GOOGLE.COM/issues/123."},
			wantError: `internal Google reference "corp.google.com"`,
		},
		{
			name:      "case-insensitive HTTPS://GO/",
			prop:      &Type{Name: "field", Type: "String", Description: "See HTTPS://GO/my-link."},
			wantError: `internal Google reference "go/"`,
		},
		{
			name:      "case-insensitive GOOGLE3/",
			prop:      &Type{Name: "field", Type: "String", Description: "Path GOOGLE3/cloud/api.proto."},
			wantError: `internal Google reference "google3/"`,
		},
		{
			name:      "deprecation_message with internal URL",
			prop:      &Type{Name: "field", Type: "String", DeprecationMessage: "See https://goto.google.com/deprecation."},
			wantError: `property field ` + "`deprecation_message`" + ` in resource TestResource must not contain internal Google reference "goto.google.com"`,
		},
		{
			name:      "removed_message with internal path",
			prop:      &Type{Name: "field", Type: "String", RemovedMessage: "Removed per google3/cloud/removed.proto."},
			wantError: `property field ` + "`removed_message`" + ` in resource TestResource must not contain internal Google reference "google3/"`,
		},
		{
			name: "nested object child description",
			prop: &Type{
				Name: "parent",
				Type: "NestedObject",
				Properties: []*Type{
					{Name: "childField", Type: "String", Description: "See https://g3doc.corp.google.com/foo."},
				},
			},
			wantError: `property parent.child_field ` + "`description`" + ` in resource TestResource must not contain internal Google reference "corp.google.com"`,
		},
		{
			name: "array item_type nested description",
			prop: &Type{
				Name: "items",
				Type: "Array",
				ItemType: &Type{
					Type: "NestedObject",
					Properties: []*Type{
						{Name: "entry", Type: "String", Description: "See //depot/foo/bar.proto."},
					},
				},
			},
			wantError: `property items.entry ` + "`description`" + ` in resource TestResource must not contain internal Google reference "//depot/"`,
		},
		{
			name: "valid public URLs and prose pass",
			prop: &Type{
				Name: "field",
				Type: "String",
				Description: "See https://cloud.google.com/go/docs, https://go.dev/doc, " +
					"https://beyondcorp.googleapis.com/v1/, beyondcorp.google.com, " +
					"https://example.com/depot/foo, https://example.com/b/123456, " +
					"django/app, mygoogle3/path, and 100 b/s.",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := &Resource{Name: "TestResource", Properties: []*Type{tc.prop}}
			res.SetDefault(p)

			errs := tc.prop.Validate(res.Name)
			if tc.wantError == "" {
				if len(errs) != 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			for _, err := range errs {
				if strings.Contains(err.Error(), tc.wantError) {
					return
				}
			}
			t.Fatalf("expected error containing %q, got: %v", tc.wantError, errs)
		})
	}
}

func TestTypeValidateEnumAndSendEmptyValue(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	cases := []struct {
		description string
		prop        *Type
		resource    string
		wantError   string
	}{
		{
			description: "Enum with empty enum_values fails",
			prop: &Type{
				Name: "mode",
				Type: "Enum",
			},
			resource:  "MyResource",
			wantError: "property mode of type Enum must specify non-empty enum_values in resource MyResource",
		},
		{
			description: "Enum with manual Possible values in description fails",
			prop: &Type{
				Name:        "mode",
				Type:        "Enum",
				Description: "Operating mode. Possible values: FAST, SLOW.",
				EnumValues:  []string{"FAST", "SLOW"},
			},
			resource:  "MyResource",
			wantError: "property mode cannot manually include 'Possible values:' in description in resource MyResource",
		},
		{
			description: "Array of Enum with manual Possible values in description fails",
			prop: &Type{
				Name:        "modes",
				Type:        "Array",
				Description: "Modes. Possible values: [FAST, SLOW]",
				ItemType: &Type{
					Type:       "Enum",
					EnumValues: []string{"FAST", "SLOW"},
				},
			},
			resource:  "MyResource",
			wantError: "property modes cannot manually include 'Possible values:' in description in resource MyResource",
		},
		{
			description: "Valid Enum without manual Possible values succeeds",
			prop: &Type{
				Name:        "mode",
				Type:        "Enum",
				Description: "Operating mode.",
				EnumValues:  []string{"FAST", "SLOW"},
			},
			resource:  "MyResource",
			wantError: "",
		},
		{
			description: "send_empty_value combined with output fails",
			prop: &Type{
				Name:           "state",
				Type:           "String",
				SendEmptyValue: true,
				Output:         true,
			},
			resource:  "MyResource",
			wantError: "property state cannot set both send_empty_value and output in resource MyResource",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       tc.resource,
				Properties: []*Type{tc.prop},
			}
			res.SetDefault(p)
			errs := tc.prop.Validate(tc.resource)
			if tc.wantError == "" {
				if len(errs) > 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), tc.wantError) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected error containing %q, got: %v", tc.wantError, errs)
			}
		})
	}

	t.Run("send_empty_value allowlist is scoped to product folder", func(t *testing.T) {
		t.Parallel()
		makeRes := func(pkgPath, prodName string) (*Resource, *Type) {
			prod := &Product{
				Name:        prodName,
				PackagePath: pkgPath,
				Versions:    []*product.Version{{Name: "ga", BaseUrl: "ga_url"}},
			}
			leaf := &Type{
				Name:             "geAuthProvider",
				Type:             "NestedObject",
				AllowEmptyObject: true,
				SendEmptyValue:   true,
				Output:           true,
				Properties:       []*Type{},
			}
			parent := &Type{
				Name:       "authProviderTypeParams",
				Type:       "NestedObject",
				Properties: []*Type{leaf},
			}
			res := &Resource{
				Name:       "AuthProvider",
				Properties: []*Type{parent},
			}
			res.SetDefault(prod)
			return res, parent
		}

		resAllowed, parentAllowed := makeRes("products/agentidentity", "AgentIdentity")
		if errs := parentAllowed.Validate(resAllowed.Name); len(errs) > 0 {
			t.Fatalf("expected allowlisted agentidentity field to pass, got: %v", errs)
		}

		resOther, parentOther := makeRes("products/otherproduct", "OtherProduct")
		errs := parentOther.Validate(resOther.Name)
		want := "property auth_provider_type_params.ge_auth_provider cannot set both send_empty_value and output in resource AuthProvider"
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), want) {
			t.Fatalf("expected error containing %q for non-allowlisted product, got: %v", want, errs)
		}
	})
}

func TestTypeValidateDefaultValue(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	cases := []struct {
		description string
		prop        *Type
		wantErrors  []string
	}{
		{
			description: "Boolean with true succeeds",
			prop:        &Type{Name: "flag", Type: "Boolean", DefaultValue: true},
		},
		{
			description: "Boolean with zero value false succeeds",
			prop:        &Type{Name: "flag", Type: "Boolean", DefaultValue: false},
		},
		{
			description: "Integer with int succeeds",
			prop:        &Type{Name: "count", Type: "Integer", DefaultValue: 42},
		},
		{
			description: "Integer with int64 succeeds",
			prop:        &Type{Name: "count", Type: "Integer", DefaultValue: int64(42)},
		},
		{
			description: "Integer with zero value 0 succeeds",
			prop:        &Type{Name: "count", Type: "Integer", DefaultValue: 0},
		},
		{
			description: "Double with float64 succeeds",
			prop:        &Type{Name: "ratio", Type: "Double", DefaultValue: 3.14},
		},
		{
			description: "Double with zero float64 succeeds",
			prop:        &Type{Name: "ratio", Type: "Double", DefaultValue: 0.0},
		},
		{
			description: "Double with int succeeds",
			prop:        &Type{Name: "ratio", Type: "Double", DefaultValue: 1},
		},
		{
			description: "Double with int64 succeeds",
			prop:        &Type{Name: "ratio", Type: "Double", DefaultValue: int64(2)},
		},
		{
			description: "String with non-empty string succeeds",
			prop:        &Type{Name: "label", Type: "String", DefaultValue: "hello"},
		},
		{
			description: "String with empty string zero value succeeds",
			prop:        &Type{Name: "label", Type: "String", DefaultValue: ""},
		},
		{
			description: "ResourceRef with string succeeds",
			prop:        &Type{Name: "network", Type: "ResourceRef", DefaultValue: "default"},
		},
		{
			description: "Enum with valid member succeeds",
			prop:        &Type{Name: "mode", Type: "Enum", DefaultValue: "FOO", EnumValues: []string{"FOO", "BAR"}},
		},
		{
			description: "Enum with empty string member succeeds when empty string is in enum_values",
			prop:        &Type{Name: "mode", Type: "Enum", DefaultValue: "", EnumValues: []string{"", "FOO"}},
		},
		{
			description: "excluded Enum with string default and empty enum_values succeeds",
			prop:        &Type{Name: "mode", Type: "Enum", DefaultValue: "FOO", Exclude: true},
		},
		{
			description: "required field with nil default_value succeeds",
			prop:        &Type{Name: "req", Type: "String", Required: true},
		},
		{
			description: "output field with nil default_value succeeds",
			prop:        &Type{Name: "out", Type: "String", Output: true},
		},
		{
			description: "Boolean with string default_value returns error",
			prop:        &Type{Name: "flag", Type: "Boolean", DefaultValue: "true"},
			wantErrors:  []string{"property flag has default_value true of type string, expected bool for Boolean in resource TestResource"},
		},
		{
			description: "Boolean with zero int default_value returns error",
			prop:        &Type{Name: "flag", Type: "Boolean", DefaultValue: 0},
			wantErrors:  []string{"property flag has default_value 0 of type int, expected bool for Boolean in resource TestResource"},
		},
		{
			description: "Integer with string default_value returns error",
			prop:        &Type{Name: "count", Type: "Integer", DefaultValue: "1"},
			wantErrors:  []string{"property count has default_value 1 of type string, expected int or int64 for Integer in resource TestResource"},
		},
		{
			description: "Integer with float64 default_value returns error",
			prop:        &Type{Name: "count", Type: "Integer", DefaultValue: 1.5},
			wantErrors:  []string{"property count has default_value 1.5 of type float64, expected int or int64 for Integer in resource TestResource"},
		},
		{
			description: "Double with string default_value returns error",
			prop:        &Type{Name: "ratio", Type: "Double", DefaultValue: "1.5"},
			wantErrors:  []string{"property ratio has default_value 1.5 of type string, expected float64, int, or int64 for Double in resource TestResource"},
		},
		{
			description: "Double with bool default_value returns error",
			prop:        &Type{Name: "ratio", Type: "Double", DefaultValue: false},
			wantErrors:  []string{"property ratio has default_value false of type bool, expected float64, int, or int64 for Double in resource TestResource"},
		},
		{
			description: "String with int default_value returns error",
			prop:        &Type{Name: "label", Type: "String", DefaultValue: 123},
			wantErrors:  []string{"property label has default_value 123 of type int, expected string for String in resource TestResource"},
		},
		{
			description: "Enum with int default_value returns error even when enum_values is empty",
			prop:        &Type{Name: "mode", Type: "Enum", DefaultValue: 1},
			wantErrors:  []string{"property mode has default_value 1 of type int, expected string for Enum in resource TestResource"},
		},
		{
			description: "ResourceRef with int default_value returns error",
			prop:        &Type{Name: "network", Type: "ResourceRef", DefaultValue: 123},
			wantErrors:  []string{"property network has default_value 123 of type int, expected string for ResourceRef in resource TestResource"},
		},
		{
			description: "NestedObject with default_value returns error",
			prop: &Type{
				Name:         "config",
				Type:         "NestedObject",
				DefaultValue: map[string]interface{}{},
				Properties:   []*Type{{Name: "a", Type: "String"}},
			},
			wantErrors: []string{"property config cannot set default_value for type NestedObject in resource TestResource"},
		},
		{
			description: "Array with default_value returns error",
			prop: &Type{
				Name:         "items",
				Type:         "Array",
				DefaultValue: []string{"a"},
				ItemType:     &Type{Type: "String"},
			},
			wantErrors: []string{"property items cannot set default_value for type Array in resource TestResource"},
		},
		{
			description: "Map with default_value returns error",
			prop: &Type{
				Name:         "entries",
				Type:         "Map",
				DefaultValue: map[string]interface{}{},
				ValueType:    &Type{Type: "NestedObject"},
			},
			wantErrors: []string{"property entries cannot set default_value for type Map in resource TestResource"},
		},
		{
			description: "KeyValuePairs with default_value returns error",
			prop:        &Type{Name: "meta", Type: "KeyValuePairs", DefaultValue: "a"},
			wantErrors:  []string{"property meta cannot set default_value for type KeyValuePairs in resource TestResource"},
		},
		{
			description: "Time with default_value returns error",
			prop:        &Type{Name: "createTime", Type: "Time", DefaultValue: "2024-01-01T00:00:00Z"},
			wantErrors:  []string{"property create_time cannot set default_value for type Time in resource TestResource"},
		},
		{
			description: "JsonValue with default_value returns error",
			prop:        &Type{Name: "payload", Type: "JsonValue", DefaultValue: "{}"},
			wantErrors:  []string{"property payload cannot set default_value for type JsonValue in resource TestResource"},
		},
		{
			description: "Enum with default_value not in enum_values returns error",
			prop:        &Type{Name: "mode", Type: "Enum", DefaultValue: "BAZ", EnumValues: []string{"FOO", "BAR"}},
			wantErrors:  []string{`property mode has default_value "BAZ" which is not in enum_values [FOO BAR] in resource TestResource`},
		},
		{
			description: "Enum with empty string default_value not in enum_values returns error",
			prop:        &Type{Name: "mode", Type: "Enum", DefaultValue: "", EnumValues: []string{"FOO", "BAR"}},
			wantErrors:  []string{`property mode has default_value "" which is not in enum_values [FOO BAR] in resource TestResource`},
		},
		{
			description: "required with non-zero default_value returns error",
			prop:        &Type{Name: "field", Type: "String", Required: true, DefaultValue: "abc"},
			wantErrors:  []string{"property field cannot set both default_value and required in resource TestResource"},
		},
		{
			description: "required with false zero default_value returns error",
			prop:        &Type{Name: "flag", Type: "Boolean", Required: true, DefaultValue: false},
			wantErrors:  []string{"property flag cannot set both default_value and required in resource TestResource"},
		},
		{
			description: "required with 0 zero default_value returns error",
			prop:        &Type{Name: "count", Type: "Integer", Required: true, DefaultValue: 0},
			wantErrors:  []string{"property count cannot set both default_value and required in resource TestResource"},
		},
		{
			description: "required with empty string zero default_value returns error",
			prop:        &Type{Name: "label", Type: "String", Required: true, DefaultValue: ""},
			wantErrors:  []string{"property label cannot set both default_value and required in resource TestResource"},
		},
		{
			description: "output with non-zero default_value returns error",
			prop:        &Type{Name: "field", Type: "String", Output: true, DefaultValue: "abc"},
			wantErrors:  []string{"property field cannot set both default_value and output in resource TestResource"},
		},
		{
			description: "output with false zero default_value returns error",
			prop:        &Type{Name: "flag", Type: "Boolean", Output: true, DefaultValue: false},
			wantErrors:  []string{"property flag cannot set both default_value and output in resource TestResource"},
		},
		{
			description: "output with 0 zero default_value returns error",
			prop:        &Type{Name: "count", Type: "Integer", Output: true, DefaultValue: 0},
			wantErrors:  []string{"property count cannot set both default_value and output in resource TestResource"},
		},
		{
			description: "nested property default_value error uses full lineage path",
			prop: &Type{
				Name: "nested",
				Type: "NestedObject",
				Properties: []*Type{
					{Name: "childField", Type: "Integer", DefaultValue: "wrong"},
				},
			},
			wantErrors: []string{"property nested.child_field has default_value wrong of type string, expected int or int64 for Integer in resource TestResource"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       "TestResource",
				Properties: []*Type{tc.prop},
			}
			res.SetDefault(p)

			errs := tc.prop.Validate(res.Name)
			if len(tc.wantErrors) == 0 {
				if len(errs) > 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			for _, want := range tc.wantErrors {
				found := false
				for _, err := range errs {
					if strings.Contains(err.Error(), want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got: %v", want, errs)
				}
			}
		})
	}
}

func TestTypeValidateEnumValues(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	cases := []struct {
		description    string
		prop           *Type
		wantErrCount   int
		wantSubstrings []string
	}{
		{
			description: "distinct enum_values succeeds",
			prop:        &Type{Name: "mode", Type: "Enum", EnumValues: []string{"A", "B", "C"}},
		},
		{
			description: "excluded Enum with empty enum_values succeeds (not enforced by validateEnumValues)",
			prop:        &Type{Name: "mode", Type: "Enum", Exclude: true},
		},
		{
			description: "Array of Enum with distinct enum_values succeeds",
			prop: &Type{
				Name: "modes",
				Type: "Array",
				ItemType: &Type{
					Type:       "Enum",
					EnumValues: []string{"A", "B"},
				},
			},
		},
		{
			description:    "duplicate enum_values returns error",
			prop:           &Type{Name: "mode", Type: "Enum", EnumValues: []string{"A", "B", "A"}},
			wantErrCount:   1,
			wantSubstrings: []string{`property mode has duplicate enum_values entry "A" in resource TestResource`},
		},
		{
			description:    "triple duplicate enum_values reports each duplicate value once",
			prop:           &Type{Name: "mode", Type: "Enum", EnumValues: []string{"A", "A", "A"}},
			wantErrCount:   1,
			wantSubstrings: []string{`property mode has duplicate enum_values entry "A" in resource TestResource`},
		},
		{
			description:    "multiple distinct duplicate enum_values reports all",
			prop:           &Type{Name: "mode", Type: "Enum", EnumValues: []string{"A", "B", "A", "B"}},
			wantErrCount:   2,
			wantSubstrings: []string{`duplicate enum_values entry "A"`, `duplicate enum_values entry "B"`},
		},
		{
			description: "Array of Enum with duplicate enum_values returns error",
			prop: &Type{
				Name: "modes",
				Type: "Array",
				ItemType: &Type{
					Type:       "Enum",
					EnumValues: []string{"X", "Y", "X"},
				},
			},
			wantErrCount:   1,
			wantSubstrings: []string{`property modes has duplicate enum_values entry "X" in resource TestResource`},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       "TestResource",
				Properties: []*Type{tc.prop},
			}
			res.SetDefault(p)

			errs := tc.prop.Validate(res.Name)
			if len(errs) != tc.wantErrCount {
				t.Fatalf("expected %d error(s), got %d: %v", tc.wantErrCount, len(errs), errs)
			}
			for _, want := range tc.wantSubstrings {
				found := false
				for _, err := range errs {
					if strings.Contains(err.Error(), want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got: %v", want, errs)
				}
			}
		})
	}
}

func TestTypeValidateRequiredConstraints(t *testing.T) {
	t.Parallel()

	p := &Product{
		Name: "TestProduct",
		Versions: []*product.Version{
			{Name: "ga", BaseUrl: "ga_url"},
		},
	}

	strSlicePtr := func(vals ...string) *[]string {
		s := vals
		return &s
	}

	cases := []struct {
		description string
		props       []*Type
		wantErrors  []string
	}{
		{
			description: "required without constraints succeeds",
			props: []*Type{
				{Name: "foo", Type: "String", Required: true},
			},
		},
		{
			description: "required with empty constraint slices succeeds",
			props: []*Type{
				{
					Name:         "foo",
					Type:         "String",
					Required:     true,
					Conflicts:    []string{},
					ExactlyOneOf: []string{},
					AtLeastOneOf: []string{},
				},
			},
		},
		{
			description: "required with required_with succeeds",
			props: []*Type{
				{Name: "foo", Type: "String", Required: true, RequiredWith: []string{"bar"}},
				{Name: "bar", Type: "String"},
			},
		},
		{
			description: "optional with conflicts, exactly_one_of, and at_least_one_of succeeds",
			props: []*Type{
				{
					Name:         "foo",
					Type:         "String",
					Conflicts:    []string{"bar"},
					ExactlyOneOf: []string{"foo", "bar"},
					AtLeastOneOf: []string{"foo", "bar"},
				},
				{Name: "bar", Type: "String"},
			},
		},
		{
			description: "required with conflicts returns error",
			props: []*Type{
				{Name: "foo", Type: "String", Required: true, Conflicts: []string{"bar"}},
				{Name: "bar", Type: "String"},
			},
			wantErrors: []string{"property foo cannot set both required and conflicts in resource TestResource"},
		},
		{
			description: "required with exactly_one_of returns error",
			props: []*Type{
				{Name: "foo", Type: "String", Required: true, ExactlyOneOf: []string{"foo", "bar"}},
				{Name: "bar", Type: "String"},
			},
			wantErrors: []string{"property foo cannot set both required and exactly_one_of in resource TestResource"},
		},
		{
			description: "required with at_least_one_of returns error",
			props: []*Type{
				{Name: "foo", Type: "String", Required: true, AtLeastOneOf: []string{"foo", "bar"}},
				{Name: "bar", Type: "String"},
			},
			wantErrors: []string{"property foo cannot set both required and at_least_one_of in resource TestResource"},
		},
		{
			description: "required with constraint group pointers returns errors",
			props: []*Type{
				{
					Name:              "foo",
					Type:              "String",
					Required:          true,
					ConflictsGroup:    strSlicePtr("bar"),
					ExactlyOneOfGroup: strSlicePtr("foo", "bar"),
					AtLeastOneOfGroup: strSlicePtr("foo", "bar"),
				},
				{Name: "bar", Type: "String"},
			},
			wantErrors: []string{
				"property foo cannot set both required and conflicts in resource TestResource",
				"property foo cannot set both required and exactly_one_of in resource TestResource",
				"property foo cannot set both required and at_least_one_of in resource TestResource",
			},
		},
		{
			description: "nested required property with conflicts reports full lineage path",
			props: []*Type{
				{
					Name: "nested",
					Type: "NestedObject",
					Properties: []*Type{
						{Name: "foo", Type: "String", Required: true, Conflicts: []string{"nested.0.bar"}},
						{Name: "bar", Type: "String"},
					},
				},
			},
			wantErrors: []string{"property nested.foo cannot set both required and conflicts in resource TestResource"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			res := &Resource{
				Name:       "TestResource",
				Properties: tc.props,
			}
			res.SetDefault(p)

			errs := tc.props[0].Validate(res.Name)
			if len(tc.wantErrors) == 0 {
				if len(errs) > 0 {
					t.Fatalf("expected no errors, got: %v", errs)
				}
				return
			}
			for _, want := range tc.wantErrors {
				found := false
				for _, err := range errs {
					if strings.Contains(err.Error(), want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got: %v", want, errs)
				}
			}
		})
	}
}

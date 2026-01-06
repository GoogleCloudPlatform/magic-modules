package metadata

import (
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/google/go-cmp/cmp"
)

func TestFromProperties(t *testing.T) {
	cases := []struct {
		name             string
		resourceMetadata *api.Resource
		virtualFields    []*api.Type
		parameters       []*api.Type
		properties       []*api.Type
		wantFields       []Field
	}{
		{
			name: "json field",
			properties: []*api.Type{
				{
					Name:          "fieldName",
					CustomFlatten: "templates/terraform/custom_flatten/json_schema.tmpl",
				},
			},
			wantFields: []Field{
				{
					Json:     true,
					ApiField: "fieldName",
				},
			},
		},
		{
			name:             "fine-grained resource field",
			resourceMetadata: &api.Resource{ApiResourceField: "parentField"},
			properties: []*api.Type{
				{
					Name: "fieldName",
				},
			},
			wantFields: []Field{
				{
					ApiField: "parentField.fieldName",
					Field:    "field_name",
				},
			},
		},
		{
			name: "provider-only",
			properties: []*api.Type{
				{
					Name:         "fieldName",
					UrlParamOnly: true,
				},
			},
			wantFields: []Field{
				{
					Field:        "field_name",
					ProviderOnly: true,
				},
			},
		},
		{
			name: "nested field",
			properties: []*api.Type{
				{
					Name: "root",
					Type: "NestedObject",
					Properties: []*api.Type{
						{
							Name: "foo",
							Type: "NestedObject",
							Properties: []*api.Type{
								{
									Name: "bars",
									Type: "Array",
									ItemType: &api.Type{
										Type: "NestedObject",
										Properties: []*api.Type{
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
			wantFields: []Field{
				{
					ApiField: "root.foo.bars.fooBar",
				},
			},
		},
		{
			name: "nested virtual",
			virtualFields: []*api.Type{
				{
					Name: "root",
					Type: "NestedObject",
					Properties: []*api.Type{
						{
							Name: "foo",
							Type: "String",
						},
					},
				},
			},
			wantFields: []Field{
				{
					Field:        "root.foo",
					ProviderOnly: true,
				},
			},
		},
		{
			name: "nested param",
			parameters: []*api.Type{
				{
					Name: "root",
					Type: "NestedObject",
					Properties: []*api.Type{
						{
							Name:         "foo",
							Type:         "String",
							UrlParamOnly: true,
						},
					},
				},
			},
			wantFields: []Field{
				{
					Field:        "root.foo",
					ProviderOnly: true,
				},
			},
		},
		{
			name: "map",
			properties: []*api.Type{
				{
					Name:    "root",
					Type:    "Map",
					KeyName: "whatever",
					ValueType: &api.Type{
						Type: "NestedObject",
						Properties: []*api.Type{
							{
								Name: "foo",
								Type: "String",
							},
						},
					},
				},
			},
			wantFields: []Field{
				{
					Field:    "root.whatever",
					ApiField: "root.key",
				},
				{
					Field:    "root.foo",
					ApiField: "root.value.foo",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := tc.resourceMetadata
			if r == nil {
				r = &api.Resource{}
			}
			r.VirtualFields = tc.virtualFields
			r.Parameters = tc.parameters
			r.Properties = tc.properties
			r.SetDefault(&api.Product{})

			got := FromProperties(r.AllNestedProperties(google.Concat(r.RootProperties(), r.UserVirtualFields())))
			if diff := cmp.Diff(tc.wantFields, got); diff != "" {
				t.Errorf("FromProperties() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsDefaultLineage(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		lineage    []string
		apiLineage []string
		want       bool
	}{
		{
			name:       "empty",
			lineage:    []string{},
			apiLineage: []string{},
			want:       true,
		},
		{
			name:       "single field",
			lineage:    []string{"foo_bar"},
			apiLineage: []string{"fooBar"},
			want:       true,
		},
		{
			name:       "multiple fields",
			lineage:    []string{"foo_bar", "baz_moo"},
			apiLineage: []string{"fooBar", "bazMoo"},
			want:       true,
		},
		{
			name:       "longer lineage",
			lineage:    []string{"foo_bar", "baz_moo"},
			apiLineage: []string{"fooBar"},
			want:       false,
		},
		{
			name:       "longer apiLineage",
			lineage:    []string{"foo_bar"},
			apiLineage: []string{"fooBar", "bazMoo"},
			want:       false,
		},
		{
			name:       "parent override",
			lineage:    []string{"foo_bar", "baz_moo"},
			apiLineage: []string{"otherName", "bazMoo"},
			want:       false,
		},
		{
			name:       "child override",
			lineage:    []string{"foo_bar", "baz_moo"},
			apiLineage: []string{"fooBar", "otherName"},
			want:       false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsDefaultLineage(tc.lineage, tc.apiLineage)
			if got != tc.want {
				t.Errorf("IsDefaultLineage(%s) failed; want %t, got %t", tc.name, tc.want, got)
			}
		})
	}
}

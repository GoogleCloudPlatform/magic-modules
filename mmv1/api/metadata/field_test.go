package metadata

import (
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/google/go-cmp/cmp"
)

func TestFromProperties(t *testing.T) {
	cases := []struct {
		name             string
		resourceMetadata *api.Resource
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := tc.resourceMetadata
			if r == nil {
				r = &api.Resource{}
			}

			for _, p := range tc.properties {
				p.SetDefault(r)
			}

			got := FromProperties(tc.properties)
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

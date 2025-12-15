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

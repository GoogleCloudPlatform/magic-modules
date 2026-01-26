package metadata

import (
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/google/go-cmp/cmp"
)

func TestFromResource(t *testing.T) {
	product := &api.Product{
		Name: "Product",
		Version: &product.Version{
			BaseUrl: "https://compute.googleapis.com/beta",
		},
	}
	cases := []struct {
		name         string
		resource     api.Resource
		wantMetadata Metadata
	}{
		{
			name:     "empty resource",
			resource: api.Resource{},
			wantMetadata: Metadata{
				Resource:       "google_product_",
				GenerationType: "mmv1",
				ApiServiceName: "compute.googleapis.com",
				ApiVersion:     "beta",
			},
		},
		{
			name: "standard",
			resource: api.Resource{
				Name:           "Test",
				AutogenStatus:  "base64",
				SourceYamlFile: "Test.yaml",
				Properties: []*api.Type{
					{
						Name:    "field",
						ApiName: "field",
					},
				},
			},
			wantMetadata: Metadata{
				Resource:            "google_product_test",
				GenerationType:      "mmv1",
				SourceFile:          "Test.yaml",
				ApiServiceName:      "compute.googleapis.com",
				ApiVersion:          "beta",
				ApiResourceTypeKind: "Test",
				AutogenStatus:       true,
				Fields: []Field{
					{
						ApiField: "field",
					},
				},
			},
		},
		{
			name: "selfLink",
			resource: api.Resource{
				Name:           "Test",
				AutogenStatus:  "base64",
				SourceYamlFile: "Test.yaml",
				Properties: []*api.Type{
					{
						Name:    "field",
						ApiName: "field",
					},
				},
				HasSelfLink: true,
			},
			wantMetadata: Metadata{
				Resource:            "google_product_test",
				GenerationType:      "mmv1",
				SourceFile:          "Test.yaml",
				ApiServiceName:      "compute.googleapis.com",
				ApiVersion:          "beta",
				ApiResourceTypeKind: "Test",
				AutogenStatus:       true,
				Fields: []Field{
					{
						ApiField: "field",
					},
					{
						ApiField: "selfLink",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.resource.SetDefault(product)

			got := FromResource(tc.resource)
			if diff := cmp.Diff(tc.wantMetadata, got); diff != "" {
				t.Errorf("FromResource() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

package compute

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestExpandComputeRegionUrlMapDefaultService(t *testing.T) {
	t.Parallel()

	config := &transport_tpg.Config{}
	d := &tpgresource.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project": "test-project",
			"region":  "us-central1",
		},
	}
	baseUrl := transport_tpg.BaseUrl(Product, config)

	cases := map[string]struct {
		Input    string
		Expected string
	}{
		"empty": {
			Input:    "",
			Expected: "",
		},
		"full self link to backend bucket": {
			Input:    baseUrl + "projects/test-project/regions/us-central1/backendBuckets/my-bucket",
			Expected: baseUrl + "projects/test-project/regions/us-central1/backendBuckets/my-bucket",
		},
		// Regression test for https://github.com/hashicorp/terraform-provider-google/issues/27773:
		// backendBuckets must not be rewritten to backendServices.
		"projects/{project}/regions/{region}/backendBuckets/{name}": {
			Input:    "projects/test-project/regions/us-central1/backendBuckets/my-bucket",
			Expected: baseUrl + "projects/test-project/regions/us-central1/backendBuckets/my-bucket",
		},
		"projects/{project}/regions/{region}/backendServices/{name}": {
			Input:    "projects/test-project/regions/us-central1/backendServices/my-service",
			Expected: baseUrl + "projects/test-project/regions/us-central1/backendServices/my-service",
		},
		"regions/{region}/backendBuckets/{name}": {
			Input:    "regions/us-central1/backendBuckets/my-bucket",
			Expected: baseUrl + "projects/test-project/regions/us-central1/backendBuckets/my-bucket",
		},
		"bare name defaults to a regional backend service": {
			Input:    "my-service",
			Expected: "projects/test-project/regions/us-central1/backendServices/my-service",
		},
	}

	expandFuncs := map[string]func(interface{}, tpgresource.TerraformResourceData, *transport_tpg.Config) (interface{}, error){
		"default_service":              expandComputeRegionUrlMapDefaultService,
		"path_matcher.default_service": expandComputeRegionUrlMapPathMatcherDefaultService,
	}

	for field, expand := range expandFuncs {
		for tn, tc := range cases {
			got, err := expand(tc.Input, d, config)
			if err != nil {
				t.Errorf("%s, %s: unexpected error: %s", field, tn, err)
				continue
			}
			if got != tc.Expected {
				t.Errorf("%s, %s: expected %q, got %q", field, tn, tc.Expected, got)
			}
		}
	}
}

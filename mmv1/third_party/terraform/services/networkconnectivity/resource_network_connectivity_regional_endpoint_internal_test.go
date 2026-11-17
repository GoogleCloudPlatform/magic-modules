package networkconnectivity

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRegionalEndpointSubnetworkDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		id                 string
		address            string
		nilResourceData    bool
		old                string
		new                string
		expectDiffSuppress bool
	}{
		"existing legacy endpoint with address URI adding subnetwork": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "projects/my-proj/regions/europe-west1/addresses/my-addr",
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: true,
		},
		"existing legacy endpoint with project number address URI": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "projects/123456789/regions/europe-west1/addresses/my-addr",
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: true,
		},
		"existing legacy endpoint with domain-scoped project in address URI": {
			id:                 "projects/google.com:my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "projects/google.com:my-proj/regions/europe-west1/addresses/my-addr",
			old:                "",
			new:                "projects/google.com:my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: true,
		},
		"new resource creation is not suppressed": {
			id:                 "",
			address:            "projects/my-proj/regions/europe-west1/addresses/my-addr",
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: false,
		},
		"nil ResourceData does not panic and is not suppressed": {
			nilResourceData:    true,
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: false,
		},
		"existing endpoint with IPv4 literal address is not suppressed": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "10.128.0.5",
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: false,
		},
		"existing endpoint with empty address is not suppressed": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "",
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: false,
		},
		"existing endpoint with full compute URL address is not suppressed": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "https://www.googleapis.com/compute/v1/projects/my-proj/regions/europe-west1/addresses/my-addr",
			old:                "",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/my-subnet",
			expectDiffSuppress: false,
		},
		"existing endpoint changing non-empty subnetwork is not suppressed": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "projects/my-proj/regions/europe-west1/addresses/my-addr",
			old:                "projects/my-proj/regions/europe-west1/subnetworks/old-subnet",
			new:                "projects/my-proj/regions/europe-west1/subnetworks/new-subnet",
			expectDiffSuppress: false,
		},
		"existing endpoint clearing subnetwork is not suppressed": {
			id:                 "projects/my-proj/locations/europe-west1/regionalEndpoints/my-rep",
			address:            "projects/my-proj/regions/europe-west1/addresses/my-addr",
			old:                "projects/my-proj/regions/europe-west1/subnetworks/old-subnet",
			new:                "",
			expectDiffSuppress: false,
		},
	}

	resSchema := ResourceNetworkConnectivityRegionalEndpoint().Schema

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var d *schema.ResourceData
			if !tc.nilResourceData {
				raw := map[string]interface{}{}
				if tc.address != "" {
					raw["address"] = tc.address
				}
				d = schema.TestResourceDataRaw(t, resSchema, raw)
				if tc.id != "" {
					d.SetId(tc.id)
				}
			}

			got := regionalEndpointSubnetworkDiffSuppress("subnetwork", tc.old, tc.new, d)
			if got != tc.expectDiffSuppress {
				t.Errorf("regionalEndpointSubnetworkDiffSuppress(%q, %q) = %t, want %t", tc.old, tc.new, got, tc.expectDiffSuppress)
			}
		})
	}
}

package compute

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestDataSourceGoogleComputeNetwork_selfLinkParsing(t *testing.T) {
	tests := []struct {
		name        string
		selfLink    string
		wantProject string
		wantName    string
	}{
		{
			name:        "full self_link",
			selfLink:    "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/my-network",
			wantProject: "my-project",
			wantName:    "my-network",
		},
		{
			name:        "partial path",
			selfLink:    "projects/my-project/global/networks/my-network",
			wantProject: "my-project",
			wantName:    "my-network",
		},
		{
			name:        "beta API version",
			selfLink:    "https://www.googleapis.com/compute/beta/projects/test-proj/global/networks/test-net",
			wantProject: "test-proj",
			wantName:    "test-net",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotName := tpgresource.GetResourceNameFromSelfLink(tc.selfLink)
			if gotName != tc.wantName {
				t.Errorf("GetResourceNameFromSelfLink(%q) = %q, want %q", tc.selfLink, gotName, tc.wantName)
			}

			// Verify project extraction logic matches what's in the data source
			parts := strings.Split(tc.selfLink, "/")
			var gotProject string
			for i, part := range parts {
				if part == "projects" && i+1 < len(parts) {
					gotProject = parts[i+1]
					break
				}
			}
			if gotProject != tc.wantProject {
				t.Errorf("project from %q = %q, want %q", tc.selfLink, gotProject, tc.wantProject)
			}
		})
	}
}

func TestFlattenDataSourceComputeNetworkPeerings(t *testing.T) {
	t.Run("no peerings", func(t *testing.T) {
		if got := flattenDataSourceComputeNetworkPeerings(nil, nil, nil); len(got) != 0 {
			t.Errorf("flattenDataSourceComputeNetworkPeerings(nil) = %v, want empty", got)
		}
	})

	t.Run("maps API fields and defaults", func(t *testing.T) {
		raw := []interface{}{
			map[string]interface{}{
				"name":                           "peer-a",
				"network":                        "https://www.googleapis.com/compute/v1/projects/p/global/networks/b",
				"state":                          "ACTIVE",
				"stateDetails":                   "Connected.",
				"exportCustomRoutes":             true,
				"importCustomRoutes":             false,
				"exportSubnetRoutesWithPublicIp": true,
				"importSubnetRoutesWithPublicIp": false,
				"stackType":                      "IPV4_IPV6",
				"updateStrategy":                 "CONSENSUS",
			},
			// Older peerings may omit stackType and updateStrategy.
			map[string]interface{}{
				"name":    "peer-b",
				"network": "https://www.googleapis.com/compute/v1/projects/p/global/networks/c",
				"state":   "INACTIVE",
			},
		}

		got := flattenDataSourceComputeNetworkPeerings(raw, nil, nil)
		if len(got) != 2 {
			t.Fatalf("got %d peerings, want 2", len(got))
		}

		want0 := map[string]interface{}{
			"name":                                "peer-a",
			"network":                             "https://www.googleapis.com/compute/v1/projects/p/global/networks/b",
			"state":                               "ACTIVE",
			"state_details":                       "Connected.",
			"export_custom_routes":                true,
			"import_custom_routes":                false,
			"export_subnet_routes_with_public_ip": true,
			"import_subnet_routes_with_public_ip": false,
			"stack_type":                          "IPV4_IPV6",
			"update_strategy":                     "CONSENSUS",
		}
		for k, v := range want0 {
			if got[0][k] != v {
				t.Errorf("peerings[0][%q] = %v, want %v", k, got[0][k], v)
			}
		}

		if got[1]["stack_type"] != "IPV4_ONLY" {
			t.Errorf("peerings[1][\"stack_type\"] = %v, want IPV4_ONLY", got[1]["stack_type"])
		}
		if got[1]["update_strategy"] != "INDEPENDENT" {
			t.Errorf("peerings[1][\"update_strategy\"] = %v, want INDEPENDENT", got[1]["update_strategy"])
		}
	})
}

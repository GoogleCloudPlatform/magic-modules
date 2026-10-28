package apihub

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitApihubAddon_customizeDiff(t *testing.T) {
	t.Parallel()

	key := "config.0.gateway_plugin_addon_config.0.gateway_plugin_configs"

	cases := []struct {
		name      string
		after     map[string]interface{}
		wantError bool
	}{
		{
			name:      "all_data_addon_config without gateway_plugin_configs",
			after:     map[string]interface{}{},
			wantError: false,
		},
		{
			name: "valid apigee_x_hybrid_config only",
			after: map[string]interface{}{
				key: []interface{}{
					map[string]interface{}{
						"plugin_instance": "projects/p/locations/l/plugins/system-apigee-x-and-hybrid/instances/i1",
						"apigee_x_hybrid_config": []interface{}{
							map[string]interface{}{
								"environment_filter": []interface{}{
									map[string]interface{}{
										"all_environments": true,
									},
								},
							},
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "valid apigee_edge_config only",
			after: map[string]interface{}{
				key: []interface{}{
					map[string]interface{}{
						"plugin_instance": "projects/p/locations/l/plugins/system-edge-public-cloud/instances/i2",
						"apigee_edge_config": []interface{}{
							map[string]interface{}{
								"environment_filter": []interface{}{
									map[string]interface{}{
										"environments": []interface{}{"test"},
									},
								},
							},
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "valid apigee_opdk_config only",
			after: map[string]interface{}{
				key: []interface{}{
					map[string]interface{}{
						"plugin_instance": "projects/p/locations/l/plugins/system-edge-private-cloud/instances/i3",
						"apigee_opdk_config": []interface{}{
							map[string]interface{}{
								"environment_filter": []interface{}{
									map[string]interface{}{
										"environments": []interface{}{"prod"},
									},
								},
							},
						},
					},
				},
			},
			wantError: false,
		},
		{
			name: "invalid - no gateway type config set",
			after: map[string]interface{}{
				key: []interface{}{
					map[string]interface{}{
						"plugin_instance":        "projects/p/locations/l/plugins/system-edge-public-cloud/instances/i1",
						"apigee_x_hybrid_config": []interface{}{},
						"apigee_edge_config":     []interface{}{},
						"apigee_opdk_config":     []interface{}{},
					},
				},
			},
			wantError: true,
		},
		{
			name: "invalid - multiple gateway type configs set",
			after: map[string]interface{}{
				key: []interface{}{
					map[string]interface{}{
						"plugin_instance": "projects/p/locations/l/plugins/system-edge-public-cloud/instances/i1",
						"apigee_x_hybrid_config": []interface{}{
							map[string]interface{}{},
						},
						"apigee_edge_config": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
			wantError: true,
		},
		{
			name: "invalid - second item has zero gateway type configs set",
			after: map[string]interface{}{
				key: []interface{}{
					map[string]interface{}{
						"plugin_instance": "projects/p/locations/l/plugins/system-edge-public-cloud/instances/i1",
						"apigee_edge_config": []interface{}{
							map[string]interface{}{},
						},
					},
					map[string]interface{}{
						"plugin_instance": "projects/p/locations/l/plugins/system-edge-public-cloud/instances/i2",
					},
				},
			},
			wantError: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := &tpgresource.ResourceDiffMock{
				Before: map[string]interface{}{},
				After:  tc.after,
			}
			err := apihubAddonCustomizeDiffFunc(d)
			if tc.wantError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

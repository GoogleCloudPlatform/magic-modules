package compute

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestComputeReservationUpdateEncoder_enableEmergentMaintenance(t *testing.T) {
	t.Parallel()

	r := ResourceComputeReservation()

	cases := []struct {
		name      string
		oldVal    string
		newVal    bool
		resizeObj map[string]interface{}
		wantField bool
		wantVal   bool
	}{
		{
			name:      "enable emergent maintenance false to true",
			oldVal:    "false",
			newVal:    true,
			wantField: true,
			wantVal:   true,
		},
		{
			name:      "enable emergent maintenance true to false",
			oldVal:    "true",
			newVal:    false,
			wantField: true,
			wantVal:   false,
		},
		{
			name:   "resize request excludes enableEmergentMaintenance",
			oldVal: "false",
			newVal: true,
			resizeObj: map[string]interface{}{
				"specificReservation": map[string]interface{}{
					"count": 2,
				},
			},
			wantField: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := &terraform.InstanceState{
				ID: "projects/my-project/zones/us-central1-a/reservations/my-res",
				Attributes: map[string]string{
					"name":                        "my-res",
					"zone":                        "us-central1-a",
					"enable_emergent_maintenance": tc.oldVal,
				},
			}
			cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name":                        "my-res",
				"zone":                        "us-central1-a",
				"enable_emergent_maintenance": tc.newVal,
				"specific_reservation": []interface{}{
					map[string]interface{}{
						"count": 1,
						"instance_properties": []interface{}{
							map[string]interface{}{
								"machine_type": "n2-standard-2",
							},
						},
					},
				},
			})

			diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, cfg, nil, nil, false)
			if err != nil {
				t.Fatalf("unexpected diff error: %v", err)
			}
			attrDiff, ok := diff.Attributes["enable_emergent_maintenance"]
			if !ok {
				t.Fatalf("expected diff on enable_emergent_maintenance")
			}
			if attrDiff.RequiresNew {
				t.Fatalf("expected enable_emergent_maintenance update to be in-place, got RequiresNew=true")
			}

			d, err := schema.InternalMap(r.Schema).Data(state, diff)
			if err != nil {
				t.Fatalf("unexpected data error: %v", err)
			}

			inputObj := map[string]interface{}{
				"enableEmergentMaintenance": tc.newVal,
			}
			if tc.resizeObj != nil {
				inputObj = tc.resizeObj
			}

			got, err := resourceComputeReservationUpdateEncoder(d, &transport_tpg.Config{}, inputObj)
			if err != nil {
				t.Fatalf("resourceComputeReservationUpdateEncoder returned error: %v", err)
			}

			val, present := got["enableEmergentMaintenance"]
			if present != tc.wantField {
				t.Fatalf("enableEmergentMaintenance presence = %v, want %v (got=%#v)", present, tc.wantField, got)
			}
			if tc.wantField && val != tc.wantVal {
				t.Fatalf("enableEmergentMaintenance = %v, want %v", val, tc.wantVal)
			}
		})
	}
}

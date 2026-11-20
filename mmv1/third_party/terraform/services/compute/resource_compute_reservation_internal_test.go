package compute

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestComputeReservationUpdateEncoder_enableEmergentMaintenance(t *testing.T) {
	t.Parallel()

	r := ResourceComputeReservation()

	baseConfigMap := func(emVal interface{}, count int, withShare bool, shareProj string) map[string]interface{} {
		m := map[string]interface{}{
			"name":    "my-res",
			"zone":    "us-central1-a",
			"project": "my-project",
			"specific_reservation": []interface{}{map[string]interface{}{
				"count": count,
				"instance_properties": []interface{}{map[string]interface{}{
					"machine_type": "n2-standard-2",
				}},
			}},
		}
		if emVal != nil {
			m["enable_emergent_maintenance"] = emVal
		}
		if withShare {
			pm := []interface{}{}
			if shareProj != "" {
				pm = append(pm, map[string]interface{}{"id": shareProj, "project_id": shareProj})
			}
			m["share_settings"] = []interface{}{map[string]interface{}{
				"share_type":  "SPECIFIC_PROJECTS",
				"project_map": pm,
			}}
		}
		return m
	}

	buildDiffData := func(t *testing.T, oldEM bool, oldShare bool, newCfg map[string]interface{}) *schema.ResourceData {
		t.Helper()
		initCfg := terraform.NewResourceConfigRaw(baseConfigMap(oldEM, 1, oldShare, ""))
		initDiff, err := schema.InternalMap(r.Schema).Diff(context.Background(), nil, initCfg, nil, nil, true)
		if err != nil {
			t.Fatalf("initial diff error: %v", err)
		}
		initData, err := schema.InternalMap(r.Schema).Data(nil, initDiff)
		if err != nil {
			t.Fatalf("initial data error: %v", err)
		}
		initData.SetId("projects/my-project/zones/us-central1-a/reservations/my-res")

		diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), initData.State(), terraform.NewResourceConfigRaw(newCfg), nil, nil, false)
		if err != nil {
			t.Fatalf("update diff error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place update (!diff.RequiresNew()), got %#v", diff.Attributes)
		}
		d, err := schema.InternalMap(r.Schema).Data(initData.State(), diff)
		if err != nil {
			t.Fatalf("update data error: %v", err)
		}
		return d
	}

	for _, tc := range []struct {
		name  string
		oldEM bool
		newEM interface{}
		count int
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{"false to true", false, true, 1, map[string]interface{}{"enableEmergentMaintenance": true}, map[string]interface{}{"enableEmergentMaintenance": true}},
		{"true to false empty input", true, false, 1, map[string]interface{}{}, map[string]interface{}{"enableEmergentMaintenance": false}},
		{"true to unset sends false", true, nil, 1, map[string]interface{}{}, map[string]interface{}{"enableEmergentMaintenance": false}},
		{"resize excludes em and mask", false, true, 2, map[string]interface{}{"specificReservation": map[string]interface{}{"count": 2}}, map[string]interface{}{"specificSkuCount": 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := buildDiffData(t, tc.oldEM, false, baseConfigMap(tc.newEM, tc.count, false, ""))
			got, err := resourceComputeReservationUpdateEncoder(d, &transport_tpg.Config{}, tc.input)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got (%#v, %v), want (%#v, nil)", got, err, tc.want)
			}
		})
	}

	t.Run("combined share_settings and enable_emergent_maintenance stay isolated", func(t *testing.T) {
		type patchCall struct {
			paths      []string
			updateMask string
			body       map[string]interface{}
		}
		var patches []patchCall

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case req.Method == http.MethodPatch && req.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/reservations/my-res":
				var body map[string]interface{}
				_ = json.NewDecoder(req.Body).Decode(&body)
				patches = append(patches, patchCall{
					paths:      req.URL.Query()["paths"],
					updateMask: req.URL.Query().Get("updateMask"),
					body:       body,
				})
				_, _ = w.Write([]byte(`{"name":"op-update","zone":"projects/my-project/zones/us-central1-a","status":"DONE"}`))
			case req.Method == http.MethodGet && req.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/operations/op-update":
				_, _ = w.Write([]byte(`{"name":"op-update","zone":"projects/my-project/zones/us-central1-a","status":"DONE"}`))
			case req.Method == http.MethodGet && req.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/reservations/my-res":
				_, _ = w.Write([]byte(`{
					"name": "my-res",
					"zone": "projects/my-project/zones/us-central1-a",
					"specificReservation": {"count": "1", "instanceProperties": {"machineType": "n2-standard-2"}},
					"shareSettings": {"shareType": "SPECIFIC_PROJECTS", "projectMap": {"222222222": {"projectId": "222222222"}}}
				}`))
			case req.Method == http.MethodGet && req.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/reservations/my-res/reservationBlocks":
				_, _ = w.Write([]byte(`{}`))
			default:
				http.NotFound(w, req)
			}
		}))
		defer server.Close()

		d := buildDiffData(t, true, true, baseConfigMap(false, 1, true, "222222222"))
		tpgCfg := &transport_tpg.Config{
			Context:         context.Background(),
			Client:          server.Client(),
			Project:         "my-project",
			Zone:            "us-central1-a",
			PollInterval:    time.Millisecond,
			CustomEndpoints: map[string]string{"compute_custom_endpoint": server.URL + "/compute/v1/"},
		}
		if err := r.Update(d, tpgCfg); err != nil {
			t.Fatalf("Update failed: %v", err)
		}
		if len(patches) != 2 {
			t.Fatalf("got %d PATCH calls, want 2: %#v", len(patches), patches)
		}
		if !reflect.DeepEqual(patches[0].paths, []string{"shareSettings.projectMap.222222222"}) || patches[0].updateMask != "" {
			t.Errorf("PATCH[0] paths=%v updateMask=%q", patches[0].paths, patches[0].updateMask)
		}
		if _, leaked := patches[0].body["enableEmergentMaintenance"]; leaked {
			t.Errorf("PATCH[0] leaked enableEmergentMaintenance: %#v", patches[0].body)
		}
		if _, hasShare := patches[0].body["shareSettings"]; !hasShare {
			t.Errorf("PATCH[0] missing shareSettings: %#v", patches[0].body)
		}
		if _, leakedMask := patches[0].body["urlUpdateMask"]; leakedMask {
			t.Errorf("PATCH[0] leaked urlUpdateMask: %#v", patches[0].body)
		}
		if !reflect.DeepEqual(patches[1].paths, []string{"enableEmergentMaintenance"}) || patches[1].updateMask != "" ||
			!reflect.DeepEqual(patches[1].body, map[string]interface{}{"enableEmergentMaintenance": false}) {
			t.Errorf("PATCH[1] = %#v", patches[1])
		}
	})
}

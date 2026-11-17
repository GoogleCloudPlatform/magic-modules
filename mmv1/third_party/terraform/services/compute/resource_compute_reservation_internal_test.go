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
			"specific_reservation": []interface{}{
				map[string]interface{}{
					"count": count,
					"instance_properties": []interface{}{
						map[string]interface{}{
							"machine_type": "n2-standard-2",
						},
					},
				},
			},
		}
		if emVal != nil {
			m["enable_emergent_maintenance"] = emVal
		}
		if withShare {
			pm := []interface{}{}
			if shareProj != "" {
				pm = append(pm, map[string]interface{}{
					"id":         shareProj,
					"project_id": shareProj,
				})
			}
			m["share_settings"] = []interface{}{
				map[string]interface{}{
					"share_type":  "SPECIFIC_PROJECTS",
					"project_map": pm,
				},
			}
		}
		return m
	}

	buildState := func(t *testing.T, emVal bool, withShare bool, shareProj string) *terraform.InstanceState {
		t.Helper()
		initCfg := terraform.NewResourceConfigRaw(baseConfigMap(emVal, 1, withShare, shareProj))
		initDiff, err := schema.InternalMap(r.Schema).Diff(context.Background(), nil, initCfg, nil, nil, true)
		if err != nil {
			t.Fatalf("unexpected initial diff error: %v", err)
		}
		initData, err := schema.InternalMap(r.Schema).Data(nil, initDiff)
		if err != nil {
			t.Fatalf("unexpected initial data error: %v", err)
		}
		initData.SetId("projects/my-project/zones/us-central1-a/reservations/my-res")
		return initData.State()
	}

	t.Run("solo false to true", func(t *testing.T) {
		state := buildState(t, false, false, "")
		cfg := terraform.NewResourceConfigRaw(baseConfigMap(true, 1, false, ""))
		diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, cfg, nil, nil, false)
		if err != nil {
			t.Fatalf("unexpected diff error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place update (!diff.RequiresNew()), got diff=%#v", diff.Attributes)
		}
		d, err := schema.InternalMap(r.Schema).Data(state, diff)
		if err != nil {
			t.Fatalf("unexpected data error: %v", err)
		}

		got, err := resourceComputeReservationUpdateEncoder(d, &transport_tpg.Config{}, map[string]interface{}{
			"enableEmergentMaintenance": true,
		})
		if err != nil {
			t.Fatalf("resourceComputeReservationUpdateEncoder error: %v", err)
		}
		want := map[string]interface{}{"enableEmergentMaintenance": true}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("solo true to false with empty inputObj", func(t *testing.T) {
		state := buildState(t, true, false, "")
		cfg := terraform.NewResourceConfigRaw(baseConfigMap(false, 1, false, ""))
		diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, cfg, nil, nil, false)
		if err != nil {
			t.Fatalf("unexpected diff error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place update (!diff.RequiresNew()), got diff=%#v", diff.Attributes)
		}
		d, err := schema.InternalMap(r.Schema).Data(state, diff)
		if err != nil {
			t.Fatalf("unexpected data error: %v", err)
		}

		// Without send_empty_value, FieldSpecificUpdateMethods passes an empty map when bool is false.
		got, err := resourceComputeReservationUpdateEncoder(d, &transport_tpg.Config{}, map[string]interface{}{})
		if err != nil {
			t.Fatalf("resourceComputeReservationUpdateEncoder error: %v", err)
		}
		want := map[string]interface{}{"enableEmergentMaintenance": false}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("solo true to unset sends false without ForceNew", func(t *testing.T) {
		state := buildState(t, true, false, "")
		cfg := terraform.NewResourceConfigRaw(baseConfigMap(nil, 1, false, ""))
		diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, cfg, nil, nil, false)
		if err != nil {
			t.Fatalf("unexpected diff error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place update (!diff.RequiresNew()), got diff=%#v", diff.Attributes)
		}
		d, err := schema.InternalMap(r.Schema).Data(state, diff)
		if err != nil {
			t.Fatalf("unexpected data error: %v", err)
		}

		got, err := resourceComputeReservationUpdateEncoder(d, &transport_tpg.Config{}, map[string]interface{}{})
		if err != nil {
			t.Fatalf("resourceComputeReservationUpdateEncoder error: %v", err)
		}
		want := map[string]interface{}{"enableEmergentMaintenance": false}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("resize request excludes enableEmergentMaintenance and urlUpdateMask", func(t *testing.T) {
		state := buildState(t, false, false, "")
		cfg := terraform.NewResourceConfigRaw(baseConfigMap(true, 2, false, ""))
		diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, cfg, nil, nil, false)
		if err != nil {
			t.Fatalf("unexpected diff error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place update (!diff.RequiresNew()), got diff=%#v", diff.Attributes)
		}
		d, err := schema.InternalMap(r.Schema).Data(state, diff)
		if err != nil {
			t.Fatalf("unexpected data error: %v", err)
		}

		got, err := resourceComputeReservationUpdateEncoder(d, &transport_tpg.Config{}, map[string]interface{}{
			"specificReservation": map[string]interface{}{
				"count": 2,
			},
		})
		if err != nil {
			t.Fatalf("resourceComputeReservationUpdateEncoder error: %v", err)
		}
		want := map[string]interface{}{"specificSkuCount": 2}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
	})

	t.Run("combined share_settings and enable_emergent_maintenance stay isolated across separate PATCH calls", func(t *testing.T) {
		type patchCall struct {
			paths      []string
			updateMask string
			body       map[string]interface{}
		}
		var patches []patchCall

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case r.Method == http.MethodPatch && r.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/reservations/my-res":
				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("failed decoding PATCH body: %v", err)
				}
				patches = append(patches, patchCall{
					paths:      r.URL.Query()["paths"],
					updateMask: r.URL.Query().Get("updateMask"),
					body:       body,
				})
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":   "op-update",
					"zone":   "projects/my-project/zones/us-central1-a",
					"status": "DONE",
				})
			case r.Method == http.MethodGet && r.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/operations/op-update":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":   "op-update",
					"zone":   "projects/my-project/zones/us-central1-a",
					"status": "DONE",
				})
			case r.Method == http.MethodGet && r.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/reservations/my-res":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name": "my-res",
					"zone": "projects/my-project/zones/us-central1-a",
					"specificReservation": map[string]interface{}{
						"count": "1",
						"instanceProperties": map[string]interface{}{
							"machineType": "n2-standard-2",
						},
					},
					"shareSettings": map[string]interface{}{
						"shareType": "SPECIFIC_PROJECTS",
						"projectMap": map[string]interface{}{
							"222222222": map[string]interface{}{
								"projectId": "222222222",
							},
						},
					},
				})
			case r.Method == http.MethodGet && r.URL.Path == "/compute/v1/projects/my-project/zones/us-central1-a/reservations/my-res/reservationBlocks":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{})
			default:
				t.Errorf("unexpected HTTP request: %s %s", r.Method, r.URL.String())
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		state := buildState(t, true, true, "")
		cfg := terraform.NewResourceConfigRaw(baseConfigMap(false, 1, true, "222222222"))

		diff, err := schema.InternalMap(r.Schema).Diff(context.Background(), state, cfg, nil, nil, false)
		if err != nil {
			t.Fatalf("unexpected diff error: %v", err)
		}
		if diff.RequiresNew() {
			t.Fatalf("expected in-place update (!diff.RequiresNew()), got diff=%#v", diff.Attributes)
		}
		d, err := schema.InternalMap(r.Schema).Data(state, diff)
		if err != nil {
			t.Fatalf("unexpected data error: %v", err)
		}

		tpgCfg := &transport_tpg.Config{
			Context:      context.Background(),
			Client:       server.Client(),
			Project:      "my-project",
			Zone:         "us-central1-a",
			PollInterval: time.Millisecond,
			CustomEndpoints: map[string]string{
				"compute_custom_endpoint": server.URL + "/compute/v1/",
			},
		}

		if err := r.Update(d, tpgCfg); err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if len(patches) != 2 {
			t.Fatalf("got %d PATCH calls, want 2: %#v", len(patches), patches)
		}

		// First PATCH: share_settings ONLY (never contains enableEmergentMaintenance).
		wantSharePaths := []string{"shareSettings.projectMap.222222222"}
		if !reflect.DeepEqual(patches[0].paths, wantSharePaths) || patches[0].updateMask != "" {
			t.Errorf("PATCH[0] paths=%v updateMask=%q, want paths=%v updateMask=\"\"", patches[0].paths, patches[0].updateMask, wantSharePaths)
		}
		if _, leaked := patches[0].body["enableEmergentMaintenance"]; leaked {
			t.Errorf("PATCH[0] body leaked enableEmergentMaintenance: %#v", patches[0].body)
		}
		if _, hasShare := patches[0].body["shareSettings"]; !hasShare {
			t.Errorf("PATCH[0] body missing shareSettings: %#v", patches[0].body)
		}
		if _, leakedMask := patches[0].body["urlUpdateMask"]; leakedMask {
			t.Errorf("PATCH[0] body leaked urlUpdateMask: %#v", patches[0].body)
		}

		// Second PATCH: enableEmergentMaintenance ONLY (never contains shareSettings or urlUpdateMask).
		wantEMPaths := []string{"enableEmergentMaintenance"}
		if !reflect.DeepEqual(patches[1].paths, wantEMPaths) || patches[1].updateMask != "" {
			t.Errorf("PATCH[1] paths=%v updateMask=%q, want paths=%v updateMask=\"\"", patches[1].paths, patches[1].updateMask, wantEMPaths)
		}
		wantEMBody := map[string]interface{}{"enableEmergentMaintenance": false}
		if !reflect.DeepEqual(patches[1].body, wantEMBody) {
			t.Errorf("PATCH[1] body = %#v, want %#v", patches[1].body, wantEMBody)
		}
	})
}

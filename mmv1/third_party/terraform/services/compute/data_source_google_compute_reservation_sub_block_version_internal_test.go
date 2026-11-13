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
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestDataSourceGoogleComputeReservationSubBlockVersionRead_pollsZonalOperation(t *testing.T) {
	t.Parallel()

	const (
		project = "test-project"
		zone    = "us-central1-a"
	)

	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/test-project/zones/us-central1-a/reservations%2Fres-1%2FreservationBlocks%2Fblk-1/reservationSubBlocks/sub-1/getVersion":
			if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
				t.Errorf("failed decoding POST body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-ver-1",
				"zone":   "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
				"status": "PENDING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/test-project/zones/us-central1-a/operations/op-subblock-ver-1":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-ver-1",
				"zone":   "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
				"status": "DONE",
				"getVersionOperationMetadata": map[string]interface{}{
					"inlineSbomInfo": map[string]interface{}{
						"currentComponentVersions": map[string]interface{}{
							"nvswitch_fw": "22.04.1",
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := &transport_tpg.Config{
		Context:      context.Background(),
		Client:       server.Client(),
		Project:      project,
		PollInterval: time.Millisecond,
		CustomEndpoints: map[string]string{
			"compute_custom_endpoint": server.URL + "/compute/v1/",
		},
	}

	ds := DataSourceGoogleComputeReservationSubBlockVersion()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":              "sub-1",
		"reservation_block": "blk-1",
		"reservation":       "res-1",
		"zone":              zone,
	})

	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(receivedBody) != 0 {
		t.Errorf("expected empty POST body when sbom_selections omitted, got %#v", receivedBody)
	}
	wantID := "projects/test-project/zones/us-central1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/getVersion"
	if d.Id() != wantID {
		t.Errorf("got ID %q, want %q", d.Id(), wantID)
	}
	wantCurrent := map[string]interface{}{"nvswitch_fw": "22.04.1"}
	if !reflect.DeepEqual(d.Get("current_component_versions"), wantCurrent) {
		t.Errorf("got current_component_versions=%#v, want %#v", d.Get("current_component_versions"), wantCurrent)
	}
	if gotTarget := d.Get("target_component_versions").(map[string]interface{}); len(gotTarget) != 0 {
		t.Errorf("expected empty target_component_versions, got %#v", gotTarget)
	}
}

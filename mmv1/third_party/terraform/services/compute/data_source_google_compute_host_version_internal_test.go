package compute

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFormatComputeHostVersionAssociation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		res       string
		block     string
		subBlock  string
		want      string
		wantError bool
	}{
		{
			name: "top-level",
			want: "",
		},
		{
			name: "reservation only",
			res:  "res-1",
			want: "reservations/res-1",
		},
		{
			name:  "reservation and block",
			res:   "res-1",
			block: "blk-1",
			want:  "reservations/res-1/reservationBlocks/blk-1",
		},
		{
			name:     "reservation, block, and sub-block",
			res:      "res-1",
			block:    "blk-1",
			subBlock: "sub-1",
			want:     "reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1",
		},
		{
			name:      "block without reservation",
			block:     "blk-1",
			wantError: true,
		},
		{
			name:      "sub-block without block",
			res:       "res-1",
			subBlock:  "sub-1",
			wantError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := formatComputeHostVersionAssociation(tc.res, tc.block, tc.subBlock)
			if tc.wantError {
				if err == nil {
					t.Fatalf("expected error, got association %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDataSourceGoogleComputeHostVersionRead_pollsZonalOperation(t *testing.T) {
	t.Parallel()

	const (
		project = "test-project"
		zone    = "us-central1-a"
	)

	var opPollCount int32
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/test-project/zones/us-central1-a/reservations%2Fres-1/hosts/host-1/getVersion":
			if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
				t.Errorf("failed decoding POST body: %v", err)
			}
			// Omit "zone" intentionally to verify fallback to the requested zone for polling.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-1",
				"status": "RUNNING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/test-project/zones/us-central1-a/operations/op-host-ver-1":
			count := atomic.AddInt32(&opPollCount, 1)
			if count == 1 {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":   "op-host-ver-1",
					"zone":   "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
					"status": "RUNNING",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-1",
				"zone":   "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
				"status": "DONE",
				"getVersionOperationMetadata": map[string]interface{}{
					"inlineSbomInfo": map[string]interface{}{
						"currentComponentVersions": map[string]interface{}{
							"bios": "1.2.3",
							"bmc":  "4.5.6",
						},
						"targetComponentVersions": map[string]interface{}{
							"bios": "1.2.4",
							"bmc":  "4.5.7",
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

	ds := DataSourceGoogleComputeHostVersion()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":            "host-1",
		"reservation":     "res-1",
		"zone":            zone,
		"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT", "SBOM_SELECTION_TARGET"},
	})

	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if atomic.LoadInt32(&opPollCount) != 2 {
		t.Errorf("got %d operation polls, want 2", opPollCount)
	}
	wantSelections := []interface{}{"SBOM_SELECTION_CURRENT", "SBOM_SELECTION_TARGET"}
	if !reflect.DeepEqual(receivedBody["sbomSelections"], wantSelections) {
		t.Errorf("got POST sbomSelections=%#v, want %#v", receivedBody["sbomSelections"], wantSelections)
	}

	wantCurrent := map[string]interface{}{"bios": "1.2.3", "bmc": "4.5.6"}
	wantTarget := map[string]interface{}{"bios": "1.2.4", "bmc": "4.5.7"}
	if !reflect.DeepEqual(d.Get("current_component_versions"), wantCurrent) {
		t.Errorf("got current_component_versions=%#v, want %#v", d.Get("current_component_versions"), wantCurrent)
	}
	if !reflect.DeepEqual(d.Get("target_component_versions"), wantTarget) {
		t.Errorf("got target_component_versions=%#v, want %#v", d.Get("target_component_versions"), wantTarget)
	}
}

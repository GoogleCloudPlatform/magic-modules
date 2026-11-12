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

func TestFormatComputeHostAssociation(t *testing.T) {
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
			got, err := formatComputeHostAssociation(tc.res, tc.block, tc.subBlock)
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

func TestDataSourceGoogleComputeReservationHostsRead_paginationAndFlattening(t *testing.T) {
	t.Parallel()

	const (
		project = "test-project"
		zone    = "us-central1-a"
	)

	var seenEscapedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seenEscapedPaths = append(seenEscapedPaths, r.URL.EscapedPath())

		if r.URL.Query().Get("filter") != "state = ACTIVE" {
			t.Errorf("unexpected filter query param: %q", r.URL.Query().Get("filter"))
		}

		switch r.URL.Query().Get("pageToken") {
		case "":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"kind":          "compute#hostList",
				"nextPageToken": "page-2",
				"items": []map[string]interface{}{
					{
						"id":                "10001",
						"name":              "host-1",
						"kind":              "compute#host",
						"creationTimestamp": "2026-11-10T00:00:00Z",
						"zone":              "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a",
						"selfLink":          "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/hosts/host-1",
						"selfLinkWithId":    "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/hosts/10001",
						"aliasLinks": []string{
							"projects/123/zones/us-central1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts/host-1",
						},
						"state": "ACTIVE",
						"status": map[string]interface{}{
							"physicalTopology": map[string]interface{}{
								"cluster":  "cluster-a",
								"block":    "block-hash-1",
								"subBlock": "sub-hash-1",
								"host":     "host-hash-1",
							},
							"runningInstances": []string{
								"https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/instances/vm-1",
							},
						},
					},
				},
			})
		case "page-2":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"kind": "compute#hostList",
				"items": []map[string]interface{}{
					{
						"id":    "10002",
						"name":  "host-2",
						"kind":  "compute#host",
						"zone":  "us-central1-a",
						"state": "ACTIVE",
					},
				},
			})
		default:
			http.Error(w, "unexpected page token", http.StatusBadRequest)
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

	ds := DataSourceGoogleComputeReservationHosts()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"reservation":           "projects/test-project/zones/us-central1-a/reservations/res-1",
		"reservation_block":     "blk-1",
		"reservation_sub_block": "sub-1",
		"filter":                "state = ACTIVE",
		"zone":                  zone,
	})

	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	wantEscapedPath := "/compute/v1/projects/test-project/zones/us-central1-a/reservations%2Fres-1%2FreservationBlocks%2Fblk-1%2FreservationSubBlocks%2Fsub-1/hosts"
	if len(seenEscapedPaths) != 2 || seenEscapedPaths[0] != wantEscapedPath || seenEscapedPaths[1] != wantEscapedPath {
		t.Fatalf("unexpected escaped paths: %v, want 2 calls to %q", seenEscapedPaths, wantEscapedPath)
	}

	wantID := "projects/test-project/zones/us-central1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts"
	if d.Id() != wantID {
		t.Errorf("got ID %q, want %q", d.Id(), wantID)
	}

	hosts := d.Get("hosts").([]interface{})
	if len(hosts) != 2 {
		t.Fatalf("got %d hosts, want 2", len(hosts))
	}

	h0 := hosts[0].(map[string]interface{})
	if h0["name"] != "host-1" || h0["resource_id"] != "10001" || h0["zone"] != "us-central1-a" || h0["state"] != "ACTIVE" {
		t.Errorf("unexpected first host map: %#v", h0)
	}
	statusList := h0["status"].([]interface{})
	if len(statusList) != 1 {
		t.Fatalf("expected 1 status element, got %#v", statusList)
	}
	status0 := statusList[0].(map[string]interface{})
	ptList := status0["physical_topology"].([]interface{})
	if len(ptList) != 1 {
		t.Fatalf("expected 1 physical_topology element, got %#v", ptList)
	}
	pt0 := ptList[0].(map[string]interface{})
	if pt0["cluster"] != "cluster-a" || pt0["block"] != "block-hash-1" || pt0["sub_block"] != "sub-hash-1" || pt0["host"] != "host-hash-1" {
		t.Errorf("unexpected physical_topology: %#v", pt0)
	}
}

func TestDataSourceGoogleComputeReservationHostRead(t *testing.T) {
	t.Parallel()

	const (
		project = "test-project"
		zone    = "us-central1-a"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.EscapedPath() != "/compute/v1/projects/test-project/zones/us-central1-a/hosts/host-top" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"resource": map[string]interface{}{
				"id":                "99999",
				"name":              "host-top",
				"kind":              "compute#host",
				"creationTimestamp": "2026-11-10T01:02:03Z",
				"selfLink":          "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/hosts/host-top",
				"selfLinkWithId":    "https://www.googleapis.com/compute/v1/projects/test-project/zones/us-central1-a/hosts/99999",
				"aliasLinks":        []string{"projects/123/zones/us-central1-a/reservations/r1/hosts/host-top"},
				"state":             "ACTIVE",
				"status": map[string]interface{}{
					"physicalTopology": map[string]interface{}{
						"cluster":  "c1",
						"block":    "b1",
						"subBlock": "sb1",
						"host":     "h1",
					},
					"runningInstances": []string{"projects/test-project/zones/us-central1-a/instances/i1"},
				},
			},
		})
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

	ds := DataSourceGoogleComputeReservationHost()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name": "host-top",
		"zone": zone,
	})

	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if d.Id() != "projects/test-project/zones/us-central1-a/hosts/host-top" {
		t.Errorf("unexpected ID: %q", d.Id())
	}
	if d.Get("resource_id") != "99999" || d.Get("state") != "ACTIVE" || d.Get("kind") != "compute#host" {
		t.Errorf("unexpected scalar attributes: resource_id=%v state=%v kind=%v", d.Get("resource_id"), d.Get("state"), d.Get("kind"))
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

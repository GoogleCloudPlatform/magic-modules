package compute

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFormatComputeReservationHostsAssociation(t *testing.T) {
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
			got, err := formatComputeReservationHostsAssociation(tc.res, tc.block, tc.subBlock)
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

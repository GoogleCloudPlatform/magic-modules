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

func TestDataSourceGoogleComputeReservationHosts_schemaValidation(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeReservationHosts()

	cases := []struct {
		name      string
		raw       map[string]interface{}
		wantError bool
	}{
		{
			name: "valid top-level",
			raw:  map[string]interface{}{},
		},
		{
			name: "block without reservation rejected by RequiredWith",
			raw: map[string]interface{}{
				"reservation_block": "blk-1",
			},
			wantError: true,
		},
		{
			name: "sub_block without block rejected by RequiredWith",
			raw: map[string]interface{}{
				"reservation":           "res-1",
				"reservation_sub_block": "sub-1",
			},
			wantError: true,
		},
		{
			name: "whitespace reservation rejected",
			raw: map[string]interface{}{
				"reservation": "   ",
			},
			wantError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diags := ds.Validate(terraform.NewResourceConfigRaw(tc.raw))
			if tc.wantError && !diags.HasError() {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantError && diags.HasError() {
				t.Fatalf("unexpected validation error: %v", diags)
			}
		})
	}
}

func TestDataSourceGoogleComputeReservationHostsRead_paginationAndFlattening(t *testing.T) {
	t.Parallel()

	var seenEscapedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		seenEscapedPaths = append(seenEscapedPaths, r.URL.EscapedPath())

		switch r.URL.EscapedPath() {
		case "/compute/v1/projects/link-project/zones/us-west1-a/reservations%2Fres-1%2FreservationBlocks%2Fblk-1%2FreservationSubBlocks%2Fsub-1/hosts":
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
							"description":       "first host",
							"kind":              "compute#host",
							"creationTimestamp": "2026-11-10T00:00:00Z",
							"zone":              "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a",
							"selfLink":          "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/host-1",
							"selfLinkWithId":    "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/10001",
							"aliasLinks": []string{
								"projects/123/zones/us-west1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts/host-1",
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
									"https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/instances/vm-1",
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
							"id":          "10002",
							"name":        "host-2",
							"description": "second host",
							"kind":        "compute#host",
							"zone":        "us-west1-a",
							"state":       "ACTIVE",
						},
					},
				})
			default:
				http.Error(w, "unexpected page token", http.StatusBadRequest)
			}
		case "/compute/v1/projects/default-project/zones/us-west1-b/hosts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"kind":  "compute#hostList",
				"items": []map[string]interface{}{},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := &transport_tpg.Config{
		Context:      context.Background(),
		Client:       server.Client(),
		Project:      "default-project",
		Zone:         "us-west1-b",
		PollInterval: time.Millisecond,
		CustomEndpoints: map[string]string{
			"compute_custom_endpoint": server.URL + "/compute/v1/",
		},
	}

	ds := DataSourceGoogleComputeReservationHosts()

	// 1. Hierarchical query with reservation self-link extracting project & zone when omitted on d, plus pagination & filter in ID.
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"reservation":           "projects/link-project/zones/us-west1-a/reservations/res-1",
		"reservation_block":     "blk-1",
		"reservation_sub_block": "sub-1",
		"filter":                "state = ACTIVE",
	})

	if err := ds.Read(d, cfg); err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	wantEscapedPath := "/compute/v1/projects/link-project/zones/us-west1-a/reservations%2Fres-1%2FreservationBlocks%2Fblk-1%2FreservationSubBlocks%2Fsub-1/hosts"
	if len(seenEscapedPaths) != 2 || seenEscapedPaths[0] != wantEscapedPath || seenEscapedPaths[1] != wantEscapedPath {
		t.Fatalf("unexpected escaped paths: %v, want 2 calls to %q", seenEscapedPaths, wantEscapedPath)
	}

	wantID := "projects/link-project/zones/us-west1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts/filter=state = ACTIVE"
	if d.Id() != wantID {
		t.Errorf("got ID %q, want %q", d.Id(), wantID)
	}
	if d.Get("project") != "link-project" || d.Get("zone") != "us-west1-a" {
		t.Errorf("got project=%v zone=%v, want link-project / us-west1-a", d.Get("project"), d.Get("zone"))
	}

	hosts := d.Get("hosts").([]interface{})
	if len(hosts) != 2 {
		t.Fatalf("got %d hosts, want 2", len(hosts))
	}

	h0 := hosts[0].(map[string]interface{})
	if h0["name"] != "host-1" ||
		h0["description"] != "first host" ||
		h0["kind"] != "compute#host" ||
		h0["resource_id"] != "10001" ||
		h0["creation_timestamp"] != "2026-11-10T00:00:00Z" ||
		h0["zone"] != "us-west1-a" ||
		h0["self_link"] != "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/host-1" ||
		h0["self_link_with_id"] != "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/10001" ||
		h0["state"] != "ACTIVE" {
		t.Errorf("unexpected first host map: %#v", h0)
	}
	wantAliases := []interface{}{"projects/123/zones/us-west1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts/host-1"}
	if !reflect.DeepEqual(h0["alias_links"], wantAliases) {
		t.Errorf("got alias_links=%#v, want %#v", h0["alias_links"], wantAliases)
	}
	statusList := h0["status"].([]interface{})
	if len(statusList) != 1 {
		t.Fatalf("expected 1 status element, got %#v", statusList)
	}
	status0 := statusList[0].(map[string]interface{})
	wantRunning := []interface{}{"https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/instances/vm-1"}
	if !reflect.DeepEqual(status0["running_instances"], wantRunning) {
		t.Errorf("got running_instances=%#v, want %#v", status0["running_instances"], wantRunning)
	}
	ptList := status0["physical_topology"].([]interface{})
	if len(ptList) != 1 {
		t.Fatalf("expected 1 physical_topology element, got %#v", ptList)
	}
	pt0 := ptList[0].(map[string]interface{})
	if pt0["cluster"] != "cluster-a" || pt0["block"] != "block-hash-1" || pt0["sub_block"] != "sub-hash-1" || pt0["host"] != "host-hash-1" {
		t.Errorf("unexpected physical_topology: %#v", pt0)
	}

	// Assert page-2 host-2 fields.
	h1 := hosts[1].(map[string]interface{})
	if h1["name"] != "host-2" ||
		h1["description"] != "second host" ||
		h1["kind"] != "compute#host" ||
		h1["resource_id"] != "10002" ||
		h1["zone"] != "us-west1-a" ||
		h1["state"] != "ACTIVE" {
		t.Errorf("unexpected page-2 host-2 map: %#v", h1)
	}

	// 2. Top-level /hosts query using provider defaults without filter.
	dTop := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	if err := ds.Read(dTop, cfg); err != nil {
		t.Fatalf("Read (top-level) failed: %v", err)
	}
	if dTop.Id() != "projects/default-project/zones/us-west1-b/hosts" {
		t.Errorf("got top-level ID %q, want %q", dTop.Id(), "projects/default-project/zones/us-west1-b/hosts")
	}
}

func TestDataSourceGoogleComputeReservationHostsRead_trailingSlashErrors(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeReservationHosts()
	cfg := &transport_tpg.Config{
		Context: context.Background(),
		Project: "test-project",
		Zone:    "us-west1-a",
	}

	errorCases := []struct {
		name string
		raw  map[string]interface{}
	}{
		{
			name: "trailing slash on reservation",
			raw: map[string]interface{}{
				"reservation": "res-1/",
			},
		},
		{
			name: "collection-only reservation self-link",
			raw: map[string]interface{}{
				"reservation": "projects/test-project/zones/us-west1-a/reservations",
			},
		},
		{
			name: "trailing slash on reservation_block",
			raw: map[string]interface{}{
				"reservation":       "res-1",
				"reservation_block": "blk-1/",
			},
		},
		{
			name: "trailing slash on reservation_sub_block",
			raw: map[string]interface{}{
				"reservation":           "res-1",
				"reservation_block":     "blk-1",
				"reservation_sub_block": "sub-1/",
			},
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ds.Schema, tc.raw)
			if err := ds.Read(d, cfg); err == nil {
				t.Fatalf("expected error for %s, got nil", tc.name)
			}
		})
	}
}

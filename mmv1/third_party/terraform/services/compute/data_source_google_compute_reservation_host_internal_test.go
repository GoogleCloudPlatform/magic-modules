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

func TestFormatComputeReservationHostAssociation(t *testing.T) {
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
			got, err := formatComputeReservationHostAssociation(tc.res, tc.block, tc.subBlock)
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

func TestDataSourceGoogleComputeReservationHost_schemaValidation(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeReservationHost()

	cases := []struct {
		name      string
		raw       map[string]interface{}
		wantError bool
	}{
		{
			name: "valid minimal",
			raw: map[string]interface{}{
				"name": "host-1",
			},
		},
		{
			name: "block without reservation rejected by RequiredWith",
			raw: map[string]interface{}{
				"name":              "host-1",
				"reservation_block": "blk-1",
			},
			wantError: true,
		},
		{
			name: "sub_block without block rejected by RequiredWith",
			raw: map[string]interface{}{
				"name":                  "host-1",
				"reservation":           "res-1",
				"reservation_sub_block": "sub-1",
			},
			wantError: true,
		},
		{
			name: "whitespace name rejected",
			raw: map[string]interface{}{
				"name": "   ",
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

func TestDataSourceGoogleComputeReservationHostRead(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.EscapedPath() {
		case "/compute/v1/projects/link-project/zones/us-west1-a/hosts/host-top":
			// Real unwrapped top-level Host JSON response per compute_public_google_rest_v1.json.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                "99999",
				"name":              "host-top",
				"description":       "top-level ExR host",
				"kind":              "compute#host",
				"creationTimestamp": "2026-11-10T01:02:03Z",
				"zone":              "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a",
				"selfLink":          "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/host-top",
				"selfLinkWithId":    "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/99999",
				"aliasLinks":        []string{"projects/123/zones/us-west1-a/reservations/r1/hosts/host-top"},
				"state":             "ACTIVE",
				"status": map[string]interface{}{
					"physicalTopology": map[string]interface{}{
						"cluster":  "c1",
						"block":    "b1",
						"subBlock": "sb1",
						"host":     "h1",
					},
					"runningInstances": []string{"projects/link-project/zones/us-west1-a/instances/i1"},
				},
			})
		case "/compute/v1/projects/res-project/zones/us-west1-b/reservations%2Fres-1%2FreservationBlocks%2Fblk-1%2FreservationSubBlocks%2Fsub-1/hosts/99999":
			// Lookup by numeric resource ID returns canonical RFC1035 name "host-sub-1".
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                "99999",
				"name":              "host-sub-1",
				"description":       "sub-block host",
				"kind":              "compute#host",
				"creationTimestamp": "2026-11-10T04:05:06Z",
				"selfLink":          "https://www.googleapis.com/compute/v1/projects/res-project/zones/us-west1-b/hosts/host-sub-1",
				"selfLinkWithId":    "https://www.googleapis.com/compute/v1/projects/res-project/zones/us-west1-b/hosts/99999",
				"aliasLinks":        []string{"projects/456/zones/us-west1-b/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts/host-sub-1"},
				"state":             "ACTIVE",
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
		Zone:         "us-central1-a",
		PollInterval: time.Millisecond,
		CustomEndpoints: map[string]string{
			"compute_custom_endpoint": server.URL + "/compute/v1/",
		},
	}

	ds := DataSourceGoogleComputeReservationHost()

	// 1. Top-level host via self-link name (extracting project & zone when omitted on d) and full attribute verification.
	dTop := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name": "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/host-top",
	})
	if err := ds.Read(dTop, cfg); err != nil {
		t.Fatalf("Read (top-level) failed: %v", err)
	}
	if dTop.Id() != "projects/link-project/zones/us-west1-a/hosts/host-top" {
		t.Errorf("unexpected ID: %q", dTop.Id())
	}
	if dTop.Get("name") != "host-top" ||
		dTop.Get("project") != "link-project" ||
		dTop.Get("zone") != "us-west1-a" ||
		dTop.Get("description") != "top-level ExR host" ||
		dTop.Get("kind") != "compute#host" ||
		dTop.Get("resource_id") != "99999" ||
		dTop.Get("creation_timestamp") != "2026-11-10T01:02:03Z" ||
		dTop.Get("self_link") != "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/host-top" ||
		dTop.Get("self_link_with_id") != "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/hosts/99999" ||
		dTop.Get("state") != "ACTIVE" {
		t.Errorf("unexpected scalar attributes: name=%v project=%v zone=%v desc=%v kind=%v resource_id=%v ts=%v self_link=%v self_link_with_id=%v state=%v",
			dTop.Get("name"), dTop.Get("project"), dTop.Get("zone"), dTop.Get("description"), dTop.Get("kind"),
			dTop.Get("resource_id"), dTop.Get("creation_timestamp"), dTop.Get("self_link"), dTop.Get("self_link_with_id"), dTop.Get("state"))
	}
	wantAliases := []interface{}{"projects/123/zones/us-west1-a/reservations/r1/hosts/host-top"}
	if !reflect.DeepEqual(dTop.Get("alias_links"), wantAliases) {
		t.Errorf("got alias_links=%#v, want %#v", dTop.Get("alias_links"), wantAliases)
	}
	statusList := dTop.Get("status").([]interface{})
	if len(statusList) != 1 {
		t.Fatalf("expected 1 status item, got %#v", statusList)
	}
	status0 := statusList[0].(map[string]interface{})
	wantInstances := []interface{}{"projects/link-project/zones/us-west1-a/instances/i1"}
	if !reflect.DeepEqual(status0["running_instances"], wantInstances) {
		t.Errorf("got running_instances=%#v, want %#v", status0["running_instances"], wantInstances)
	}
	ptList := status0["physical_topology"].([]interface{})
	if len(ptList) != 1 {
		t.Fatalf("expected 1 physical_topology item, got %#v", ptList)
	}
	pt0 := ptList[0].(map[string]interface{})
	if pt0["cluster"] != "c1" || pt0["block"] != "b1" || pt0["sub_block"] != "sb1" || pt0["host"] != "h1" {
		t.Errorf("unexpected physical_topology: %#v", pt0)
	}

	// 2. Hierarchical %2F-escaped URL with reservation self-link extracting project & zone, and canonical name write-back from numeric ID.
	dScoped := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":                  "99999",
		"reservation":           "projects/res-project/zones/us-west1-b/reservations/res-1",
		"reservation_block":     "blk-1",
		"reservation_sub_block": "sub-1",
	})
	if err := ds.Read(dScoped, cfg); err != nil {
		t.Fatalf("Read (hierarchical) failed: %v", err)
	}
	wantScopedID := "projects/res-project/zones/us-west1-b/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/hosts/host-sub-1"
	if dScoped.Id() != wantScopedID {
		t.Errorf("got hierarchical ID %q, want %q", dScoped.Id(), wantScopedID)
	}
	if dScoped.Get("name") != "host-sub-1" {
		t.Errorf("got canonical name %q, want %q", dScoped.Get("name"), "host-sub-1")
	}
	if dScoped.Get("project") != "res-project" || dScoped.Get("zone") != "us-west1-b" {
		t.Errorf("got project=%v zone=%v, want res-project / us-west1-b", dScoped.Get("project"), dScoped.Get("zone"))
	}
}

func TestDataSourceGoogleComputeReservationHostRead_trailingSlashErrors(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeReservationHost()
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
			name: "trailing slash on host name",
			raw:  map[string]interface{}{"name": "host-1/"},
		},
		{
			name: "trailing slash on reservation",
			raw: map[string]interface{}{
				"name":        "host-1",
				"reservation": "res-1/",
			},
		},
		{
			name: "trailing slash on reservation_block",
			raw: map[string]interface{}{
				"name":              "host-1",
				"reservation":       "res-1",
				"reservation_block": "blk-1/",
			},
		},
		{
			name: "trailing slash on reservation_sub_block",
			raw: map[string]interface{}{
				"name":                  "host-1",
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

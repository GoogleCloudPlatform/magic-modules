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

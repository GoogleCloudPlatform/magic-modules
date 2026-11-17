package compute

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestDataSourceGoogleComputeHostVersion_schemaValidation(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeHostVersion()

	cases := []struct {
		name      string
		raw       map[string]interface{}
		wantError bool
	}{
		{
			name: "valid current selection",
			raw: map[string]interface{}{
				"name":            "host-1",
				"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "valid target selection",
			raw: map[string]interface{}{
				"name":            "host-1",
				"sbom_selections": []interface{}{"SBOM_SELECTION_TARGET"},
			},
		},
		{
			name: "missing required sbom_selections",
			raw: map[string]interface{}{
				"name": "host-1",
			},
			wantError: true,
		},
		{
			name: "more than one sbom_selections rejected",
			raw: map[string]interface{}{
				"name":            "host-1",
				"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT", "SBOM_SELECTION_TARGET"},
			},
			wantError: true,
		},
		{
			name: "empty sbom_selections element rejected",
			raw: map[string]interface{}{
				"name":            "host-1",
				"sbom_selections": []interface{}{""},
			},
			wantError: true,
		},
		{
			name: "sub_block without reservation_block rejected by RequiredWith",
			raw: map[string]interface{}{
				"name":                  "host-1",
				"reservation":           "res-1",
				"reservation_sub_block": "sub-1",
				"sbom_selections":       []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantError: true,
		},
		{
			name: "block without reservation rejected by RequiredWith",
			raw: map[string]interface{}{
				"name":              "host-1",
				"reservation_block": "blk-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantError: true,
		},
		{
			name: "whitespace name rejected",
			raw: map[string]interface{}{
				"name":            "   ",
				"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
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

func TestDataSourceGoogleComputeHostVersionRead_selectionsAndSelfLinkExtraction(t *testing.T) {
	t.Parallel()

	var currentPollCount int32
	var targetPollCount int32
	var gotCurrentBody map[string]interface{}
	var gotTargetBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/link-project/zones/us-west1-a/reservations%2Fres-1/hosts/host-1/getVersion":
			if err := json.NewDecoder(r.Body).Decode(&gotCurrentBody); err != nil {
				t.Errorf("failed decoding POST body: %v", err)
			}
			// Omit "zone" intentionally to verify fallback to the resolved zone for polling.
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-cur",
				"status": "RUNNING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/link-project/zones/us-west1-a/operations/op-host-ver-cur":
			count := atomic.AddInt32(&currentPollCount, 1)
			if count == 1 {
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":   "op-host-ver-cur",
					"zone":   "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a",
					"status": "RUNNING",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-cur",
				"zone":   "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a",
				"status": "DONE",
				"getVersionOperationMetadata": map[string]interface{}{
					"inlineSbomInfo": map[string]interface{}{
						"currentComponentVersions": map[string]interface{}{
							"bios": "1.2.3",
							"bmc":  "4.5.6",
						},
					},
				},
			})
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-west1-b/hosts/host-top/getVersion":
			if err := json.NewDecoder(r.Body).Decode(&gotTargetBody); err != nil {
				t.Errorf("failed decoding POST body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-tgt",
				"status": "RUNNING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-west1-b/operations/op-host-ver-tgt":
			atomic.AddInt32(&targetPollCount, 1)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-tgt",
				"status": "DONE",
				"getVersionOperationMetadata": map[string]interface{}{
					"inlineSbomInfo": map[string]interface{}{
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
		Project:      "default-project",
		Zone:         "us-west1-b",
		PollInterval: time.Millisecond,
		CustomEndpoints: map[string]string{
			"compute_custom_endpoint": server.URL + "/compute/v1/",
		},
	}

	ds := DataSourceGoogleComputeHostVersion()

	// 1. SBOM_SELECTION_CURRENT with reservation self-link extracting project and zone when omitted on d.
	dCurrent := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":            "host-1",
		"reservation":     "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/reservations/res-1",
		"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
	})
	if err := ds.Read(dCurrent, cfg); err != nil {
		t.Fatalf("Read (CURRENT) failed: %v", err)
	}
	if atomic.LoadInt32(&currentPollCount) != 2 {
		t.Errorf("got %d operation polls, want 2", currentPollCount)
	}
	if !reflect.DeepEqual(gotCurrentBody["sbomSelections"], []interface{}{"SBOM_SELECTION_CURRENT"}) {
		t.Errorf("got POST sbomSelections=%#v, want [SBOM_SELECTION_CURRENT]", gotCurrentBody["sbomSelections"])
	}
	wantCurrentID := "projects/link-project/zones/us-west1-a/reservations/res-1/hosts/host-1/getVersion/SBOM_SELECTION_CURRENT"
	if dCurrent.Id() != wantCurrentID {
		t.Errorf("got ID %q, want %q", dCurrent.Id(), wantCurrentID)
	}
	if dCurrent.Get("project") != "link-project" || dCurrent.Get("zone") != "us-west1-a" {
		t.Errorf("got project=%v zone=%v, want link-project / us-west1-a", dCurrent.Get("project"), dCurrent.Get("zone"))
	}
	wantCurrent := map[string]interface{}{"bios": "1.2.3", "bmc": "4.5.6"}
	if !reflect.DeepEqual(dCurrent.Get("current_component_versions"), wantCurrent) {
		t.Errorf("got current_component_versions=%#v, want %#v", dCurrent.Get("current_component_versions"), wantCurrent)
	}

	// 2. SBOM_SELECTION_TARGET on top-level host falling back to provider default project and zone.
	dTarget := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":            "host-top",
		"sbom_selections": []interface{}{"SBOM_SELECTION_TARGET"},
	})
	if err := ds.Read(dTarget, cfg); err != nil {
		t.Fatalf("Read (TARGET) failed: %v", err)
	}
	if atomic.LoadInt32(&targetPollCount) != 1 {
		t.Errorf("got %d target operation polls, want 1", targetPollCount)
	}
	if !reflect.DeepEqual(gotTargetBody["sbomSelections"], []interface{}{"SBOM_SELECTION_TARGET"}) {
		t.Errorf("got POST sbomSelections=%#v, want [SBOM_SELECTION_TARGET]", gotTargetBody["sbomSelections"])
	}
	wantTargetID := "projects/default-project/zones/us-west1-b/hosts/host-top/getVersion/SBOM_SELECTION_TARGET"
	if dTarget.Id() != wantTargetID {
		t.Errorf("got ID %q, want %q", dTarget.Id(), wantTargetID)
	}
	wantTarget := map[string]interface{}{"bios": "1.2.4", "bmc": "4.5.7"}
	if !reflect.DeepEqual(dTarget.Get("target_component_versions"), wantTarget) {
		t.Errorf("got target_component_versions=%#v, want %#v", dTarget.Get("target_component_versions"), wantTarget)
	}
}

func TestDataSourceGoogleComputeHostVersionRead_operationAndHTTPErrors(t *testing.T) {
	t.Parallel()

	var gotUserProjectHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/target-project/zones/us-central1-a/reservations%2Fres-1%2FreservationBlocks%2Fblk-1%2FreservationSubBlocks%2Fsub-1/hosts/host-err/getVersion":
			gotUserProjectHeader = r.Header.Get("X-Goog-User-Project")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-err",
				"status": "RUNNING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/target-project/zones/us-central1-a/operations/op-host-ver-err":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-host-ver-err",
				"status": "DONE",
				"error": map[string]interface{}{
					"errors": []map[string]interface{}{
						{
							"code":    "INTERNAL_ERROR",
							"message": "host SBOM retrieval failed",
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
		Context:             context.Background(),
		Client:              server.Client(),
		Project:             "target-project",
		BillingProject:      "billing-override-project",
		UserProjectOverride: true,
		Zone:                "us-central1-a",
		PollInterval:        time.Millisecond,
		CustomEndpoints: map[string]string{
			"compute_custom_endpoint": server.URL + "/compute/v1/",
		},
	}

	ds := DataSourceGoogleComputeHostVersion()

	// LRO error propagation + hierarchical sub-block URL + billing_project header.
	dOpErr := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":                  "host-err",
		"reservation":           "res-1",
		"reservation_block":     "blk-1",
		"reservation_sub_block": "sub-1",
		"sbom_selections":       []interface{}{"SBOM_SELECTION_CURRENT"},
	})
	if err := ds.Read(dOpErr, cfg); err == nil {
		t.Fatalf("expected error from failed zonal operation, got nil")
	} else if !strings.Contains(err.Error(), "host SBOM retrieval failed") {
		t.Fatalf("expected operation error message in %q", err.Error())
	}
	if dOpErr.Id() != "" {
		t.Errorf("expected empty ID on failed operation, got %q", dOpErr.Id())
	}
	if gotUserProjectHeader != "billing-override-project" {
		t.Errorf("got X-Goog-User-Project=%q, want %q", gotUserProjectHeader, "billing-override-project")
	}

	// HTTP 404 error propagation.
	dNotFound := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":            "missing-host",
		"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
	})
	if err := ds.Read(dNotFound, cfg); err == nil {
		t.Fatalf("expected error on 404 response, got nil")
	}
}

func TestDataSourceGoogleComputeHostVersionRead_trailingSlashErrors(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeHostVersion()
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
			raw: map[string]interface{}{
				"name":            "host-1/",
				"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "trailing slash on reservation",
			raw: map[string]interface{}{
				"name":            "host-1",
				"reservation":     "projects/test-project/zones/us-west1-a/reservations/res-1/",
				"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "collection-only reservation self-link",
			raw: map[string]interface{}{
				"name":            "host-1",
				"reservation":     "projects/test-project/zones/us-west1-a/reservations",
				"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "trailing slash on reservation_block",
			raw: map[string]interface{}{
				"name":              "host-1",
				"reservation":       "res-1",
				"reservation_block": "blk-1/",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "trailing slash on reservation_sub_block",
			raw: map[string]interface{}{
				"name":                  "host-1",
				"reservation":           "res-1",
				"reservation_block":     "blk-1",
				"reservation_sub_block": "sub-1/",
				"sbom_selections":       []interface{}{"SBOM_SELECTION_CURRENT"},
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

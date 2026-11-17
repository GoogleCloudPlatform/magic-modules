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

func TestDataSourceGoogleComputeReservationSubBlockVersion_schemaValidation(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeReservationSubBlockVersion()

	cases := []struct {
		name      string
		raw       map[string]interface{}
		wantError bool
	}{
		{
			name: "valid current selection",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "valid target selection",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_TARGET"},
			},
		},
		{
			name: "missing required sbom_selections",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
			},
			wantError: true,
		},
		{
			name: "more than one sbom_selections rejected",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT", "SBOM_SELECTION_TARGET"},
			},
			wantError: true,
		},
		{
			name: "empty sbom_selections element rejected",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{""},
			},
			wantError: true,
		},
		{
			name: "whitespace name rejected",
			raw: map[string]interface{}{
				"name":              "   ",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantError: true,
		},
		{
			name: "whitespace reservation_block rejected",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "   ",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantError: true,
		},
		{
			name: "whitespace reservation rejected",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "   ",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
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

func TestDataSourceGoogleComputeReservationSubBlockVersionRead_table(t *testing.T) {
	t.Parallel()

	var currentPolls int32
	var gotCurrentBody map[string]interface{}
	var gotTargetBody map[string]interface{}
	var gotBillingHeader string

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		// Case 1: SBOM_SELECTION_CURRENT + multi-poll PENDING -> RUNNING -> DONE + missing initial zone on op.
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-central1-a/reservations%2Fres-1%2FreservationBlocks%2Fblk-1/reservationSubBlocks/sub-1/getVersion":
			if err := json.NewDecoder(r.Body).Decode(&gotCurrentBody); err != nil {
				t.Errorf("failed decoding POST body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-cur",
				"status": "PENDING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-central1-a/operations/op-subblock-cur":
			n := atomic.AddInt32(&currentPolls, 1)
			switch n {
			case 1:
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":   "op-subblock-cur",
					"zone":   server.URL + "/compute/v1/projects/default-project/zones/us-central1-a",
					"status": "RUNNING",
				})
			default:
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"name":   "op-subblock-cur",
					"zone":   server.URL + "/compute/v1/projects/default-project/zones/us-central1-a",
					"status": "DONE",
					"getVersionOperationMetadata": map[string]interface{}{
						"inlineSbomInfo": map[string]interface{}{
							"currentComponentVersions": map[string]interface{}{
								"nvswitch_fw": "22.04.1",
								"bmc_fw":      "10.01.0",
								"null_entry":  nil,
							},
						},
					},
				})
			}

		// Case 2: SBOM_SELECTION_TARGET + self-link inputs extracting project & zone + billing_project override.
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/link-project/zones/us-west1-a/reservations%2Fres-link%2FreservationBlocks%2Fblk-link/reservationSubBlocks/sub-link/getVersion":
			gotBillingHeader = r.Header.Get("X-Goog-User-Project")
			if err := json.NewDecoder(r.Body).Decode(&gotTargetBody); err != nil {
				t.Errorf("failed decoding POST body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-tgt",
				"zone":   server.URL + "/compute/v1/projects/link-project/zones/us-west1-a",
				"status": "RUNNING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/link-project/zones/us-west1-a/operations/op-subblock-tgt":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-tgt",
				"zone":   server.URL + "/compute/v1/projects/link-project/zones/us-west1-a",
				"status": "DONE",
				"getVersionOperationMetadata": map[string]interface{}{
					"inlineSbomInfo": map[string]interface{}{
						"targetComponentVersions": map[string]interface{}{
							"nvswitch_fw": "22.05.0",
						},
					},
				},
			})

		// Case 3: Operation completes with error.
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-central1-a/reservations%2Fres-err%2FreservationBlocks%2Fblk-err/reservationSubBlocks/sub-err/getVersion":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-err",
				"status": "RUNNING",
			})
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-central1-a/operations/op-subblock-err":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name":   "op-subblock-err",
				"status": "DONE",
				"error": map[string]interface{}{
					"errors": []map[string]interface{}{
						{
							"code":    "INTERNAL_ERROR",
							"message": "sub-block SBOM retrieval failed",
						},
					},
				},
			})

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleComputeReservationSubBlockVersion()

	cases := []struct {
		name               string
		raw                map[string]interface{}
		billingProject     string
		userProjectOveride bool
		wantErrSubstr      string
		wantID             string
		wantProject        string
		wantZone           string
		wantCurrent        map[string]interface{}
		wantTarget         map[string]interface{}
	}{
		{
			name: "SBOM_SELECTION_CURRENT with provider default project/zone and multi-poll operation",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantID:      "projects/default-project/zones/us-central1-a/reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1/getVersion/SBOM_SELECTION_CURRENT",
			wantProject: "default-project",
			wantZone:    "us-central1-a",
			wantCurrent: map[string]interface{}{
				"nvswitch_fw": "22.04.1",
				"bmc_fw":      "10.01.0",
			},
			wantTarget: map[string]interface{}{},
		},
		{
			name: "SBOM_SELECTION_TARGET with self-link inputs and billing_project override",
			raw: map[string]interface{}{
				"name":              "projects/link-project/zones/us-west1-a/reservations/res-link/reservationBlocks/blk-link/reservationSubBlocks/sub-link",
				"reservation_block": "projects/link-project/zones/us-west1-a/reservations/res-link/reservationBlocks/blk-link",
				"reservation":       "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/reservations/res-link",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_TARGET"},
			},
			billingProject:     "billing-override-proj",
			userProjectOveride: true,
			wantID:             "projects/link-project/zones/us-west1-a/reservations/res-link/reservationBlocks/blk-link/reservationSubBlocks/sub-link/getVersion/SBOM_SELECTION_TARGET",
			wantProject:        "link-project",
			wantZone:           "us-west1-a",
			wantCurrent:        map[string]interface{}{},
			wantTarget: map[string]interface{}{
				"nvswitch_fw": "22.05.0",
			},
		},
		{
			name: "zonal operation error propagates cleanly",
			raw: map[string]interface{}{
				"name":              "sub-err",
				"reservation_block": "blk-err",
				"reservation":       "res-err",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantErrSubstr: "sub-block SBOM retrieval failed",
		},
		{
			name: "HTTP 404 propagates cleanly",
			raw: map[string]interface{}{
				"name":              "sub-404",
				"reservation_block": "blk-404",
				"reservation":       "res-404",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
			wantErrSubstr: "404",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &transport_tpg.Config{
				Context:             context.Background(),
				Client:              server.Client(),
				Project:             "default-project",
				Zone:                "us-central1-a",
				BillingProject:      tc.billingProject,
				UserProjectOverride: tc.userProjectOveride,
				PollInterval:        time.Millisecond,
				CustomEndpoints: map[string]string{
					"compute_custom_endpoint": server.URL + "/compute/v1/",
				},
			}

			d := schema.TestResourceDataRaw(t, ds.Schema, tc.raw)
			err := ds.Read(d, cfg)
			if tc.wantErrSubstr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrSubstr)
				}
				if !strings.Contains(err.Error(), tc.wantErrSubstr) {
					t.Fatalf("got error %q, want substring %q", err.Error(), tc.wantErrSubstr)
				}
				if d.Id() != "" {
					t.Errorf("expected empty ID on error, got %q", d.Id())
				}
				return
			}
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}
			if d.Id() != tc.wantID {
				t.Errorf("got ID %q, want %q", d.Id(), tc.wantID)
			}
			if d.Get("project") != tc.wantProject || d.Get("zone") != tc.wantZone {
				t.Errorf("got project=%v zone=%v, want %s / %s", d.Get("project"), d.Get("zone"), tc.wantProject, tc.wantZone)
			}
			if !reflect.DeepEqual(d.Get("current_component_versions"), tc.wantCurrent) {
				t.Errorf("got current_component_versions=%#v, want %#v", d.Get("current_component_versions"), tc.wantCurrent)
			}
			if !reflect.DeepEqual(d.Get("target_component_versions"), tc.wantTarget) {
				t.Errorf("got target_component_versions=%#v, want %#v", d.Get("target_component_versions"), tc.wantTarget)
			}
		})
	}

	if atomic.LoadInt32(&currentPolls) != 2 {
		t.Errorf("got %d polls for multi-poll operation, want 2", currentPolls)
	}
	if !reflect.DeepEqual(gotCurrentBody["sbomSelections"], []interface{}{"SBOM_SELECTION_CURRENT"}) {
		t.Errorf("got POST sbomSelections=%#v, want [SBOM_SELECTION_CURRENT]", gotCurrentBody["sbomSelections"])
	}
	if !reflect.DeepEqual(gotTargetBody["sbomSelections"], []interface{}{"SBOM_SELECTION_TARGET"}) {
		t.Errorf("got POST sbomSelections=%#v, want [SBOM_SELECTION_TARGET]", gotTargetBody["sbomSelections"])
	}
	if gotBillingHeader != "billing-override-proj" {
		t.Errorf("got X-Goog-User-Project=%q, want %q", gotBillingHeader, "billing-override-proj")
	}
}

func TestDataSourceGoogleComputeReservationSubBlockVersionRead_trailingSlashErrors(t *testing.T) {
	t.Parallel()

	ds := DataSourceGoogleComputeReservationSubBlockVersion()
	cfg := &transport_tpg.Config{
		Context: context.Background(),
		Project: "test-project",
		Zone:    "us-central1-a",
	}

	errorCases := []struct {
		name string
		raw  map[string]interface{}
	}{
		{
			name: "trailing slash on sub-block name",
			raw: map[string]interface{}{
				"name":              "sub-1/",
				"reservation_block": "blk-1",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "trailing slash on reservation_block",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1/",
				"reservation":       "res-1",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "trailing slash on reservation",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "res-1/",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
			},
		},
		{
			name: "collection-only reservation self-link",
			raw: map[string]interface{}{
				"name":              "sub-1",
				"reservation_block": "blk-1",
				"reservation":       "projects/test-project/zones/us-central1-a/reservations",
				"sbom_selections":   []interface{}{"SBOM_SELECTION_CURRENT"},
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

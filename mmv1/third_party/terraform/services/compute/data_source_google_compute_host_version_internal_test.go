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
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestFormatComputeHostVersionAssociation(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		res, block, subBlock, want string
		wantErr                    bool
	}{
		{"", "", "", "", false},
		{"res-1", "", "", "reservations/res-1", false},
		{"res-1", "blk-1", "", "reservations/res-1/reservationBlocks/blk-1", false},
		{"res-1", "blk-1", "sub-1", "reservations/res-1/reservationBlocks/blk-1/reservationSubBlocks/sub-1", false},
		{"", "blk-1", "", "", true},
		{"res-1", "", "sub-1", "", true},
	} {
		got, err := formatComputeHostVersionAssociation(tc.res, tc.block, tc.subBlock)
		if (err != nil) != tc.wantErr || got != tc.want {
			t.Errorf("formatComputeHostVersionAssociation(%q, %q, %q) = (%q, %v), want (%q, err=%v)",
				tc.res, tc.block, tc.subBlock, got, err, tc.want, tc.wantErr)
		}
	}
}

func TestDataSourceGoogleComputeHostVersionRead(t *testing.T) {
	t.Parallel()

	var currentPolls int32
	var gotCurrentBody, gotTargetBody map[string]interface{}
	var gotBillingHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/link-project/zones/us-west1-a/reservations%2Fres-1/hosts/host-1/getVersion":
			_ = json.NewDecoder(r.Body).Decode(&gotCurrentBody)
			_, _ = w.Write([]byte(`{"name":"op-cur","status":"RUNNING"}`))
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/link-project/zones/us-west1-a/operations/op-cur":
			if atomic.AddInt32(&currentPolls, 1) == 1 {
				_, _ = w.Write([]byte(`{"name":"op-cur","zone":"https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a","status":"RUNNING"}`))
				return
			}
			_, _ = w.Write([]byte(`{
				"name": "op-cur",
				"zone": "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a",
				"status": "DONE",
				"getVersionOperationMetadata": {"inlineSbomInfo": {"currentComponentVersions": {"bios": "1.2.3", "bmc": "4.5.6"}}}
			}`))
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-west1-b/hosts/host-top/getVersion":
			_ = json.NewDecoder(r.Body).Decode(&gotTargetBody)
			_, _ = w.Write([]byte(`{"name":"op-tgt","status":"RUNNING"}`))
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-west1-b/operations/op-tgt":
			_, _ = w.Write([]byte(`{
				"name": "op-tgt",
				"status": "DONE",
				"getVersionOperationMetadata": {"inlineSbomInfo": {"targetComponentVersions": {"bios": "1.2.4", "bmc": "4.5.7"}}}
			}`))
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-west1-b/reservations%2Fres-1%2FreservationBlocks%2Fblk-1%2FreservationSubBlocks%2Fsub-1/hosts/host-err/getVersion":
			gotBillingHeader = r.Header.Get("X-Goog-User-Project")
			_, _ = w.Write([]byte(`{"name":"op-err","status":"RUNNING"}`))
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/compute/v1/projects/default-project/zones/us-west1-b/operations/op-err":
			_, _ = w.Write([]byte(`{"name":"op-err","status":"DONE","error":{"errors":[{"code":"INTERNAL_ERROR","message":"host SBOM retrieval failed"}]}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	ds := DataSourceGoogleComputeHostVersion()
	newCfg := func(billing string) *transport_tpg.Config {
		return &transport_tpg.Config{
			Context:             context.Background(),
			Client:              server.Client(),
			Project:             "default-project",
			Zone:                "us-west1-b",
			BillingProject:      billing,
			UserProjectOverride: billing != "",
			PollInterval:        time.Millisecond,
			CustomEndpoints:     map[string]string{"compute_custom_endpoint": server.URL + "/compute/v1/"},
		}
	}
	cfg := newCfg("")

	checkRead := func(raw map[string]interface{}, wantID string, want map[string]interface{}) {
		t.Helper()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)
		if err := ds.Read(d, cfg); err != nil {
			t.Fatalf("Read(%v) failed: %v", raw, err)
		}
		if d.Id() != wantID {
			t.Errorf("Id() = %q, want %q", d.Id(), wantID)
		}
		for k, v := range want {
			if got := d.Get(k); !reflect.DeepEqual(got, v) {
				t.Errorf("%s = %#v, want %#v", k, got, v)
			}
		}
	}

	checkRead(
		map[string]interface{}{
			"name":            "host-1",
			"reservation":     "https://www.googleapis.com/compute/v1/projects/link-project/zones/us-west1-a/reservations/res-1",
			"sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"},
		},
		"projects/link-project/zones/us-west1-a/reservations/res-1/hosts/host-1/getVersion/SBOM_SELECTION_CURRENT",
		map[string]interface{}{
			"project":                    "link-project",
			"zone":                       "us-west1-a",
			"current_component_versions": map[string]interface{}{"bios": "1.2.3", "bmc": "4.5.6"},
			"target_component_versions":  map[string]interface{}{},
		},
	)
	if atomic.LoadInt32(&currentPolls) != 2 {
		t.Errorf("currentPolls = %d, want 2", currentPolls)
	}
	if !reflect.DeepEqual(gotCurrentBody["sbomSelections"], []interface{}{"SBOM_SELECTION_CURRENT"}) {
		t.Errorf("POST sbomSelections = %#v, want [SBOM_SELECTION_CURRENT]", gotCurrentBody["sbomSelections"])
	}

	checkRead(
		map[string]interface{}{
			"name":            "host-top",
			"sbom_selections": []interface{}{"SBOM_SELECTION_TARGET"},
		},
		"projects/default-project/zones/us-west1-b/hosts/host-top/getVersion/SBOM_SELECTION_TARGET",
		map[string]interface{}{
			"project":                    "default-project",
			"zone":                       "us-west1-b",
			"current_component_versions": map[string]interface{}{},
			"target_component_versions":  map[string]interface{}{"bios": "1.2.4", "bmc": "4.5.7"},
		},
	)
	if !reflect.DeepEqual(gotTargetBody["sbomSelections"], []interface{}{"SBOM_SELECTION_TARGET"}) {
		t.Errorf("POST sbomSelections = %#v, want [SBOM_SELECTION_TARGET]", gotTargetBody["sbomSelections"])
	}

	dOpErr := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"name":                  "host-err",
		"reservation":           "res-1",
		"reservation_block":     "blk-1",
		"reservation_sub_block": "sub-1",
		"sbom_selections":       []interface{}{"SBOM_SELECTION_CURRENT"},
	})
	if err := ds.Read(dOpErr, newCfg("billing-override-proj")); err == nil || !strings.Contains(err.Error(), "host SBOM retrieval failed") {
		t.Fatalf("err = %v, want substring %q", err, "host SBOM retrieval failed")
	}
	if dOpErr.Id() != "" || gotBillingHeader != "billing-override-proj" {
		t.Errorf("Id()=%q gotBillingHeader=%q, want empty ID and billing-override-proj", dOpErr.Id(), gotBillingHeader)
	}

	for _, raw := range []map[string]interface{}{
		{"name": "missing-host", "sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"}},
		{"name": "host-1/", "sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"}},
		{"name": "host-1", "reservation": "projects/p/zones/us-west1-a/reservations", "sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"}},
		{"name": "host-1", "reservation": "res-1", "reservation_block": "blk-1/", "sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"}},
		{"name": "host-1", "reservation": "res-1", "reservation_block": "blk-1", "reservation_sub_block": "sub-1/", "sbom_selections": []interface{}{"SBOM_SELECTION_CURRENT"}},
	} {
		if err := ds.Read(schema.TestResourceDataRaw(t, ds.Schema, raw), cfg); err == nil {
			t.Errorf("expected error for %v, got nil", raw)
		}
	}
}

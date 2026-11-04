package resourcemanagerv3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type recordedRequest struct {
	Method   string
	Path     string
	RawQuery string
	RawBody  string
}

type mockCapabilityConfigServer struct {
	mu          sync.Mutex
	requests    []recordedRequest
	boundaries  []string
	displayName string
	types       []string
	deleted     bool
	nextOp      int
	pendingOps  map[string]func()
}

func newMockCapabilityConfigServer(t *testing.T) (*mockCapabilityConfigServer, *httptest.Server, *transport_tpg.Config) {
	t.Helper()

	mock := &mockCapabilityConfigServer{
		pendingOps: make(map[string]func()),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed reading request body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mock.mu.Lock()
		defer mock.mu.Unlock()

		mock.requests = append(mock.requests, recordedRequest{
			Method:   r.Method,
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
			RawBody:  string(bodyBytes),
		})

		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && len(r.URL.Path) > len("/v3/operations/") && r.URL.Path[:len("/v3/operations/")] == "/v3/operations/":
			opName := r.URL.Path[len("/v3/"):]
			if commitFn, exists := mock.pendingOps[opName]; exists {
				commitFn()
				delete(mock.pendingOps, opName)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name": opName,
				"done": true,
			})
			return

		case r.Method == http.MethodPost && r.URL.Path == "/v3/organizations/123456789/capabilityConfigs":
			var body map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mock.nextOp++
			opName := fmt.Sprintf("operations/create-%d", mock.nextOp)
			mock.pendingOps[opName] = func() {
				mock.deleted = false
				if dn, ok := body["displayName"].(string); ok {
					mock.displayName = dn
				}
				if rawTypes, ok := body["types"].([]interface{}); ok {
					mock.types = toStringSlice(rawTypes)
				}
				if rawBoundaries, ok := body["boundaries"].([]interface{}); ok {
					mock.boundaries = toStringSlice(rawBoundaries)
				} else {
					mock.boundaries = nil
				}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name": opName,
				"done": false,
			})
			return

		case r.Method == http.MethodGet && r.URL.Path == "/v3/organizations/123456789/capabilityConfigs/tf-test-capcfg":
			if mock.deleted {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"code":    404,
						"message": "CapabilityConfig not found",
						"status":  "NOT_FOUND",
					},
				})
				return
			}
			resp := map[string]interface{}{
				"name":              "organizations/123456789/capabilityConfigs/tf-test-capcfg",
				"types":             mock.types,
				"managementProject": "projects/987654321",
				"state":             "ACTIVE",
				"etag":              "etag-1",
				"createTime":        "2026-10-20T00:00:00Z",
				"updateTime":        "2026-10-20T00:01:00Z",
			}
			if mock.displayName != "" {
				resp["displayName"] = mock.displayName
			}
			// Mirror Proto3 JSON serialization: omit empty repeated field on GET.
			if len(mock.boundaries) > 0 {
				resp["boundaries"] = mock.boundaries
			}
			_ = json.NewEncoder(w).Encode(resp)
			return

		case r.Method == http.MethodPatch && r.URL.Path == "/v3/organizations/123456789/capabilityConfigs/tf-test-capcfg":
			var body map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mock.nextOp++
			opName := fmt.Sprintf("operations/patch-%d", mock.nextOp)
			mock.pendingOps[opName] = func() {
				if r.URL.Query().Get("updateMask") == "boundaries" {
					rawBoundaries, exists := body["boundaries"]
					if !exists {
						t.Errorf("expected explicit boundaries field in PATCH body, got: %s", string(bodyBytes))
						return
					}
					mock.boundaries = toStringSlice(rawBoundaries.([]interface{}))
				}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name": opName,
				"done": false,
			})
			return

		case r.Method == http.MethodDelete && r.URL.Path == "/v3/organizations/123456789/capabilityConfigs/tf-test-capcfg":
			// Mirror CRM v3 API precondition: DeleteCapabilityConfig fails with FAILED_PRECONDITION
			// if boundaries are still associated (including if an LRO was not awaited).
			if len(mock.boundaries) > 0 {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error": map[string]interface{}{
						"code":    400,
						"message": fmt.Sprintf("CapabilityConfig organizations/123456789/capabilityConfigs/tf-test-capcfg cannot be deleted while boundaries are associated with it: %v", mock.boundaries),
						"status":  "FAILED_PRECONDITION",
					},
				})
				return
			}
			mock.nextOp++
			opName := fmt.Sprintf("operations/delete-%d", mock.nextOp)
			mock.pendingOps[opName] = func() {
				mock.deleted = true
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"name": opName,
				"done": false,
			})
			return

		default:
			t.Errorf("unexpected HTTP request: %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
			http.Error(w, "unexpected request", http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	cfg := &transport_tpg.Config{
		Client:          server.Client(),
		CustomEndpoints: map[string]string{Product.CustomEndpointField: server.URL + "/v3/"},
		PollInterval:    time.Millisecond,
		UserAgent:       "terraform-provider-google-unit-test",
	}

	return mock, server, cfg
}

func toStringSlice(raw []interface{}) []string {
	out := make([]string, len(raw))
	for i, item := range raw {
		out[i] = item.(string)
	}
	return out
}

func TestCapabilityConfig_BoundariesSchemaAndLifecycle(t *testing.T) {
	t.Parallel()

	res := ResourceResourceManagerV3CapabilityConfig()
	boundariesSchema, ok := res.Schema["boundaries"]
	if !ok {
		t.Fatal("expected boundaries field in schema")
	}
	if boundariesSchema.Type != schema.TypeList {
		t.Fatalf("expected boundaries to be TypeList (ordered), got %v", boundariesSchema.Type)
	}
	if !boundariesSchema.Optional || boundariesSchema.Required || boundariesSchema.Computed {
		t.Fatalf("expected boundaries Optional=true, Required=false, Computed=false; got Optional=%v Required=%v Computed=%v",
			boundariesSchema.Optional, boundariesSchema.Required, boundariesSchema.Computed)
	}

	mock, _, cfg := newMockCapabilityConfigServer(t)
	ctx := context.Background()

	// Step 1: Create without boundaries — POST body omits "boundaries".
	cfgStep1 := terraform.NewResourceConfigRaw(map[string]interface{}{
		"parent":               "organizations/123456789",
		"capability_config_id": "tf-test-capcfg",
		"display_name":         "Test Capability",
		"types":                []interface{}{"AGENT_MANAGEMENT"},
	})
	diff1, err := res.Diff(ctx, nil, cfgStep1, cfg)
	if err != nil {
		t.Fatalf("Step 1 Diff failed: %v", err)
	}
	state1, diags := res.Apply(ctx, nil, diff1, cfg)
	if diags.HasError() {
		t.Fatalf("Step 1 Apply failed: %v", diags)
	}
	var createBody map[string]interface{}
	if err := json.Unmarshal([]byte(mock.requests[0].RawBody), &createBody); err != nil {
		t.Fatalf("unmarshal create body: %v", err)
	}
	if _, exists := createBody["boundaries"]; exists {
		t.Fatalf("expected omitted boundaries on Create when unset, got body: %s", mock.requests[0].RawBody)
	}

	// Step 2: Update to reverse-alphabetical list — preserves exact caller order.
	wantOrdered := []string{
		"organizations/123456789/boundaries/tf-test-bnd-z",
		"organizations/123456789/boundaries/tf-test-bnd-a",
	}
	cfgStep2 := terraform.NewResourceConfigRaw(map[string]interface{}{
		"parent":               "organizations/123456789",
		"capability_config_id": "tf-test-capcfg",
		"display_name":         "Test Capability",
		"types":                []interface{}{"AGENT_MANAGEMENT"},
		"boundaries":           []interface{}{wantOrdered[0], wantOrdered[1]},
	})
	diff2, err := res.Diff(ctx, state1, cfgStep2, cfg)
	if err != nil {
		t.Fatalf("Step 2 Diff failed: %v", err)
	}
	state2, diags := res.Apply(ctx, state1, diff2, cfg)
	if diags.HasError() {
		t.Fatalf("Step 2 Apply failed: %v", diags)
	}
	gotBoundaries := []string{state2.Attributes["boundaries.0"], state2.Attributes["boundaries.1"]}
	if !reflect.DeepEqual(gotBoundaries, wantOrdered) {
		t.Fatalf("boundaries order mismatch: got %v, want %v", gotBoundaries, wantOrdered)
	}

	// Step 3: Clear boundaries = [] — sends explicit {"boundaries": []} with updateMask=boundaries,
	// and produces zero permadiff when GET omits empty boundaries.
	cfgStep3 := terraform.NewResourceConfigRaw(map[string]interface{}{
		"parent":               "organizations/123456789",
		"capability_config_id": "tf-test-capcfg",
		"display_name":         "Test Capability",
		"types":                []interface{}{"AGENT_MANAGEMENT"},
		"boundaries":           []interface{}{},
	})
	diff3, err := res.Diff(ctx, state2, cfgStep3, cfg)
	if err != nil {
		t.Fatalf("Step 3 Diff failed: %v", err)
	}
	if len(diff3.Attributes) == 0 {
		t.Fatal("expected non-empty attribute diff when clearing boundaries from 2 to 0 items")
	}
	state3, diags := res.Apply(ctx, state2, diff3, cfg)
	if diags.HasError() {
		t.Fatalf("Step 3 Apply failed: %v", diags)
	}
	if state3.Attributes["boundaries.#"] != "0" {
		t.Fatalf("expected boundaries.# = 0 after clearing, got %q", state3.Attributes["boundaries.#"])
	}
	replanDiff, err := res.Diff(ctx, state3, cfgStep3, cfg)
	if err != nil {
		t.Fatalf("Step 3 re-plan Diff failed: %v", err)
	}
	if len(replanDiff.Attributes) != 0 {
		t.Fatalf("expected zero attribute diffs after clearing boundaries = [], got: %v", replanDiff.Attributes)
	}

	// Step 4: Re-associate boundary and run Delete — verifies pre_delete detaches boundaries via PATCH
	// and awaits LRO completion before sending DELETE.
	state4, diags := res.Apply(ctx, state3, mustDiff(t, res, state3, cfgStep2, cfg), cfg)
	if diags.HasError() {
		t.Fatalf("Step 4 Apply failed: %v", diags)
	}
	mock.requests = nil
	destroyDiff := &terraform.InstanceDiff{Destroy: true}
	_, diags = res.Apply(ctx, state4, destroyDiff, cfg)
	if diags.HasError() {
		t.Fatalf("Destroy with associated boundaries failed: %v", diags)
	}

	var methods []string
	for _, req := range mock.requests {
		methods = append(methods, fmt.Sprintf("%s %s?%s", req.Method, req.Path, req.RawQuery))
	}
	wantSequence := []string{
		"PATCH /v3/organizations/123456789/capabilityConfigs/tf-test-capcfg?alt=json&updateMask=boundaries",
		"GET /v3/operations/patch-5?alt=json",
		"DELETE /v3/organizations/123456789/capabilityConfigs/tf-test-capcfg?alt=json",
		"GET /v3/operations/delete-6?alt=json",
	}
	if !reflect.DeepEqual(methods, wantSequence) {
		t.Fatalf("pre_delete + delete request sequence mismatch:\ngot:  %v\nwant: %v", methods, wantSequence)
	}
	var disassociateBody map[string]interface{}
	if err := json.Unmarshal([]byte(mock.requests[0].RawBody), &disassociateBody); err != nil {
		t.Fatalf("unmarshal pre_delete PATCH body: %v", err)
	}
	if !reflect.DeepEqual(disassociateBody, map[string]interface{}{"boundaries": []interface{}{}}) {
		t.Fatalf("pre_delete PATCH body mismatch: got %s", mock.requests[0].RawBody)
	}
}

func TestCapabilityConfig_DeleteWithoutBoundariesSkipsPatch(t *testing.T) {
	t.Parallel()

	res := ResourceResourceManagerV3CapabilityConfig()
	mock, _, cfg := newMockCapabilityConfigServer(t)

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"parent":               "organizations/123456789",
		"capability_config_id": "tf-test-capcfg",
		"types":                []interface{}{"AGENT_MANAGEMENT"},
		"deletion_policy":      "DELETE",
	})
	d.SetId("organizations/123456789/capabilityConfigs/tf-test-capcfg")

	if err := resourceResourceManagerV3CapabilityConfigDelete(d, cfg); err != nil {
		t.Fatalf("Delete without boundaries failed: %v", err)
	}

	var methods []string
	for _, req := range mock.requests {
		methods = append(methods, req.Method)
	}
	wantMethods := []string{"DELETE", "GET"}
	if !reflect.DeepEqual(methods, wantMethods) {
		t.Fatalf("expected direct DELETE + LRO GET when boundaries is empty, got: %v", methods)
	}
}

func mustDiff(t *testing.T, res *schema.Resource, prior *terraform.InstanceState, cfg *terraform.ResourceConfig, meta interface{}) *terraform.InstanceDiff {
	t.Helper()
	diff, err := res.Diff(context.Background(), prior, cfg, meta)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	return diff
}

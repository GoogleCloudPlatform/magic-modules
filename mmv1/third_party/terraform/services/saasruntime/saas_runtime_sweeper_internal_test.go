package saasruntime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var fastSaasRuntimeSweeperTiming = saasRuntimeSweeperTiming{
	pollInterval:       2 * time.Millisecond,
	deleteRetryTimeout: 100 * time.Millisecond,
	absenceTimeout:     200 * time.Millisecond,
	operationTimeout:   200 * time.Millisecond,
}

func TestSaasRuntimeSweeper_DAGOrdering(t *testing.T) {
	t.Parallel()

	filtered := sweeper.FilterSweepers("google_saas_runtime_saas", sweeper.GetSweeperInventory())
	wantNames := []string{
		"google_saas_runtime_unit_operation",
		"google_saas_runtime_unit",
		"google_saas_runtime_release",
		"google_saas_runtime_rollout_kind",
		"google_saas_runtime_unit_kind",
		"google_saas_runtime_tenant",
		"google_saas_runtime_saas",
	}
	for _, name := range wantNames {
		s, ok := filtered[name]
		if !ok {
			t.Fatalf("expected sweeper %q in filtered dependency closure, got keys %v", name, filtered)
		}
		if s.ListAndAction == nil {
			t.Errorf("sweeper %q has nil ListAndAction", name)
		}
		if s.DeleteFunction == nil {
			t.Errorf("sweeper %q has nil DeleteFunction", name)
		}
	}

	ordered, err := sweeper.ValidateAndOrderSweepersWithDependencies(filtered)
	if err != nil {
		t.Fatalf("ValidateAndOrderSweepersWithDependencies failed: %v", err)
	}

	pos := make(map[string]int, len(ordered))
	for i, s := range ordered {
		pos[s.Name] = i
	}

	assertBefore := func(first, second string) {
		t.Helper()
		if pos[first] >= pos[second] {
			t.Errorf("expected %s (idx %d) to run before %s (idx %d)", first, pos[first], second, pos[second])
		}
	}

	assertBefore("google_saas_runtime_unit_operation", "google_saas_runtime_unit")
	assertBefore("google_saas_runtime_unit_operation", "google_saas_runtime_release")
	assertBefore("google_saas_runtime_unit", "google_saas_runtime_release")
	assertBefore("google_saas_runtime_unit", "google_saas_runtime_unit_kind")
	assertBefore("google_saas_runtime_unit", "google_saas_runtime_tenant")
	assertBefore("google_saas_runtime_release", "google_saas_runtime_unit_kind")
	assertBefore("google_saas_runtime_rollout_kind", "google_saas_runtime_unit_kind")
	assertBefore("google_saas_runtime_unit_kind", "google_saas_runtime_saas")
	assertBefore("google_saas_runtime_tenant", "google_saas_runtime_saas")
}

type fakeSaasRuntimeRecord struct {
	collection string
	location   string
	obj        map[string]interface{}
	// softDeletesRemaining simulates EasySaaS async soft-delete purging:
	// while > 0, GET still returns 200 (decrementing on each GET) and parent
	// ValidateDelete still sees the child as present.
	softDeletesRemaining int
	deleted              bool
}

type fakeSaasRuntimeServer struct {
	mu               sync.Mutex
	records          map[string]*fakeSaasRuntimeRecord // key: relative path "projects/p/locations/l/col/id"
	rejectWildcard   bool
	unreachableInAll map[string]bool
	mutatedPaths     []string
}

func newFakeSaasRuntimeServer() *fakeSaasRuntimeServer {
	return &fakeSaasRuntimeServer{
		records:          make(map[string]*fakeSaasRuntimeRecord),
		unreachableInAll: make(map[string]bool),
	}
}

func (f *fakeSaasRuntimeServer) put(collection string, obj map[string]interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	name := obj["name"].(string)
	rel := saasRuntimeRelativeResourcePath(name)
	f.records[rel] = &fakeSaasRuntimeRecord{
		collection: collection,
		location:   saasRuntimeLocationFromName(rel),
		obj:        cloneMap(obj),
	}
}

func cloneMap(m map[string]interface{}) map[string]interface{} {
	raw, _ := json.Marshal(m)
	var out map[string]interface{}
	_ = json.Unmarshal(raw, &out)
	return out
}

func (f *fakeSaasRuntimeServer) activeInCollection(collection, location string) []*fakeSaasRuntimeRecord {
	var out []*fakeSaasRuntimeRecord
	for _, rec := range f.records {
		if rec.deleted || rec.collection != collection {
			continue
		}
		if location != "-" && rec.location != location {
			continue
		}
		if location == "-" && f.unreachableInAll[rec.location] {
			continue
		}
		out = append(out, rec)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].obj["name"].(string) < out[j].obj["name"].(string)
	})
	return out
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": msg,
			"status":  "INVALID_ARGUMENT",
		},
	})
}

func (f *fakeSaasRuntimeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	idx := strings.Index(r.URL.Path, "projects/")
	if idx < 0 {
		writeJSONError(w, http.StatusNotFound, "invalid path")
		return
	}
	rel := r.URL.Path[idx:]
	parts := strings.Split(rel, "/")
	// List or Create: projects/{p}/locations/{l}/{collection}
	if len(parts) == 5 && parts[0] == "projects" && parts[2] == "locations" {
		project, loc, col := parts[1], parts[3], parts[4]
		switch r.Method {
		case http.MethodGet:
			if loc == "-" && f.rejectWildcard {
				writeJSONError(w, http.StatusBadRequest, "wildcard location '-' disabled in test")
				return
			}
			recs := f.activeInCollection(col, loc)
			const pageSize = 2
			start := 0
			if tok := r.URL.Query().Get("pageToken"); tok != "" {
				start, _ = strconv.Atoi(tok)
			}
			end := start + pageSize
			nextToken := ""
			if end < len(recs) {
				nextToken = strconv.Itoa(end)
			} else {
				end = len(recs)
			}
			var page []interface{}
			if start < len(recs) {
				for _, rec := range recs[start:end] {
					page = append(page, cloneMap(rec.obj))
				}
			}
			resp := map[string]interface{}{}
			if len(page) > 0 {
				resp[col] = page
			}
			if nextToken != "" {
				resp["nextPageToken"] = nextToken
			}
			if loc == "-" && len(f.unreachableInAll) > 0 {
				var unr []string
				for uLoc := range f.unreachableInAll {
					unr = append(unr, fmt.Sprintf("projects/%s/locations/%s", project, uLoc))
				}
				sort.Strings(unr)
				resp["unreachable"] = unr
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return

		case http.MethodPost:
			if col != "unitOperations" {
				writeJSONError(w, http.StatusBadRequest, "unsupported POST collection "+col)
				return
			}
			opID := r.URL.Query().Get("unitOperationId")
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			opName := fmt.Sprintf("projects/%s/locations/%s/unitOperations/%s", project, loc, opID)
			f.mutatedPaths = append(f.mutatedPaths, "POST:"+opName)
			unitRef, _ := body["unit"].(string)
			if _, hasDeprov := body["deprovision"]; hasDeprov && unitRef != "" {
				unitRel := saasRuntimeRelativeResourcePath(unitRef)
				if unitRec, ok := f.records[unitRel]; ok && !unitRec.deleted {
					unitRec.obj["state"] = "UNIT_STATE_NOT_PROVISIONED"
					unitRec.obj["conditions"] = []interface{}{
						map[string]interface{}{
							"type":   "TYPE_PROVISIONED",
							"status": "STATUS_FALSE",
						},
					}
				}
			}
			opObj := map[string]interface{}{
				"name":  opName,
				"unit":  unitRef,
				"state": "UNIT_OPERATION_STATE_SUCCEEDED",
			}
			f.records[opName] = &fakeSaasRuntimeRecord{
				collection: "unitOperations",
				location:   loc,
				obj:        opObj,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(cloneMap(opObj))
			return
		}
	}

	// Item operations: projects/{p}/locations/{l}/{collection}/{id}
	if len(parts) == 6 && parts[0] == "projects" && parts[2] == "locations" {
		col := parts[4]
		rec, exists := f.records[rel]
		if !exists || rec.deleted {
			writeJSONError(w, http.StatusNotFound, "not found: "+rel)
			return
		}

		switch r.Method {
		case http.MethodGet:
			if rec.softDeletesRemaining > 0 {
				rec.softDeletesRemaining--
				if rec.softDeletesRemaining == 0 {
					rec.deleted = true
					writeJSONError(w, http.StatusNotFound, "not found: "+rel)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(cloneMap(rec.obj))
			return

		case http.MethodPatch:
			f.mutatedPaths = append(f.mutatedPaths, "PATCH:"+rel)
			mask := r.URL.Query().Get("updateMask")
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			if col == "unitKinds" && strings.Contains(mask, "defaultRelease") {
				val, _ := body["defaultRelease"].(string)
				if val == "" {
					delete(rec.obj, "defaultRelease")
				} else {
					rec.obj["defaultRelease"] = val
				}
			}
			if col == "unitOperations" && strings.Contains(mask, "cancel") {
				if unitEcho, _ := body["unit"].(string); unitEcho == "" {
					writeJSONError(w, http.StatusBadRequest, "unit is required on cancel")
					return
				}
				rec.obj["state"] = "UNIT_OPERATION_STATE_CANCELLED"
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(cloneMap(rec.obj))
			return

		case http.MethodDelete:
			f.mutatedPaths = append(f.mutatedPaths, "DELETE:"+rel)
			if reason := f.validateDeleteLocked(col, rel, rec); reason != "" {
				writeJSONError(w, http.StatusBadRequest, reason)
				return
			}
			rec.obj["deleteTime"] = "2026-11-05T19:30:00Z"
			rec.softDeletesRemaining = 2
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}
	}

	writeJSONError(w, http.StatusNotFound, "unhandled request "+r.Method+" "+r.URL.Path)
}

func (f *fakeSaasRuntimeServer) validateDeleteLocked(col, rel string, rec *fakeSaasRuntimeRecord) string {
	targetKey := saasRuntimeCanonicalKey(rel)
	loc := rec.location

	switch col {
	case "unitOperations":
		state, _ := rec.obj["state"].(string)
		if state != "" && !saasRuntimeTerminalUnitOpStates[state] {
			return "unit operation is executing"
		}
	case "units":
		if isSaasRuntimeUnitProvisioned(rec.obj) {
			return "unit has been provisioned, create deprovision unit operation to deprovision the unit"
		}
		for _, opRec := range f.activeInCollection("unitOperations", loc) {
			uRef, _ := opRec.obj["unit"].(string)
			opState, _ := opRec.obj["state"].(string)
			if saasRuntimeCanonicalKey(uRef) == targetKey && opState != "" && !saasRuntimeTerminalUnitOpStates[opState] {
				return fmt.Sprintf("unit has ongoing operation %q that is not completed so cannot be deleted", opRec.obj["name"])
			}
		}
		for _, otherUnit := range f.activeInCollection("units", loc) {
			if rawDeps, ok := otherUnit.obj["dependencies"].([]interface{}); ok {
				for _, rd := range rawDeps {
					if dm, ok := rd.(map[string]interface{}); ok {
						if u, _ := dm["unit"].(string); saasRuntimeCanonicalKey(u) == targetKey {
							return fmt.Sprintf("unit has dependent %q that is not in NOT_PROVISIONED state and cannot be deleted", otherUnit.obj["name"])
						}
					}
				}
			}
		}
	case "releases":
		for _, uRec := range f.activeInCollection("units", loc) {
			rRef, _ := uRec.obj["release"].(string)
			if saasRuntimeCanonicalKey(rRef) == targetKey && isSaasRuntimeUnitProvisioned(uRec.obj) {
				return fmt.Sprintf("cannot delete Release %q because it has units that are still provisioned", rel)
			}
		}
		for _, ukRec := range f.activeInCollection("unitKinds", loc) {
			dr, _ := ukRec.obj["defaultRelease"].(string)
			if dr != "" && saasRuntimeCanonicalKey(dr) == targetKey {
				return fmt.Sprintf("cannot delete Release %q because it is referenced as default_release by UnitKind %q", rel, ukRec.obj["name"])
			}
		}
	case "unitKinds":
		for _, uRec := range f.activeInCollection("units", loc) {
			ukRef, _ := uRec.obj["unitKind"].(string)
			if saasRuntimeCanonicalKey(ukRef) == targetKey && isSaasRuntimeUnitProvisioned(uRec.obj) {
				return fmt.Sprintf("ValidateDelete fails: cannot delete UnitKind %s because it has units which are still provisioned: Example: %s", rel, uRec.obj["name"])
			}
		}
		for _, otherUK := range f.activeInCollection("unitKinds", loc) {
			for _, depKey := range extractUnitKindDependencyKeys(otherUK.obj) {
				if depKey == targetKey {
					return fmt.Sprintf("ValidateDelete fails: cannot delete UnitKind %s because there is another UnitKind with a dependency on UnitKind %s: Dependent name: %s", rel, rel, otherUK.obj["name"])
				}
			}
		}
		for _, rRec := range f.activeInCollection("releases", loc) {
			ukRef, _ := rRec.obj["unitKind"].(string)
			if saasRuntimeCanonicalKey(ukRef) == targetKey {
				return fmt.Sprintf("ValidateDelete fails: cannot delete UnitKind %s because there is a Release %s with a dependency on UnitKind %s", rel, rRec.obj["name"], rel)
			}
		}
	case "tenants":
		for _, uRec := range f.activeInCollection("units", loc) {
			tRef, _ := uRec.obj["tenant"].(string)
			if saasRuntimeCanonicalKey(tRef) == targetKey {
				return fmt.Sprintf("Tenant %q has Units", rel)
			}
		}
	case "saas":
		for _, tRec := range f.activeInCollection("tenants", loc) {
			sRef, _ := tRec.obj["saas"].(string)
			if saasRuntimeCanonicalKey(sRef) == targetKey {
				return fmt.Sprintf("cannot delete SaaS %q because it has dependent Tenants", rel)
			}
		}
		for _, ukRec := range f.activeInCollection("unitKinds", loc) {
			sRef, _ := ukRec.obj["saas"].(string)
			if saasRuntimeCanonicalKey(sRef) == targetKey {
				return fmt.Sprintf("cannot delete SaaS %q because it has dependent UnitKinds", rel)
			}
		}
	}
	return ""
}

func newTestSaasRuntimeSweeperConfig(srv *httptest.Server) *transport_tpg.Config {
	return &transport_tpg.Config{
		Client:    srv.Client(),
		Project:   "test-proj",
		Region:    "us-central1",
		UserAgent: "sweeper-test",
		CustomEndpoints: map[string]string{
			Product.CustomEndpointField: srv.URL + "/v1beta1/",
		},
	}
}

func TestSaasRuntimeSweeper_EndToEndMultiLocationAndCycles(t *testing.T) {
	t.Parallel()

	fake := newFakeSaasRuntimeServer()
	// Force europe-west1 to be reported as unreachable in "locations/-" so the
	// explicit unreachable-location fallback is exercised too.
	fake.unreachableInAll["europe-west1"] = true

	// 1. Leaked us-east1 stack with defaultRelease deadlock cycle, dependent UnitKinds
	//    where "tf-test-web-unitkind" sorts AFTER "tf-test-cluster-unitkind",
	//    provisioned Unit, dependent Unit, and running UnitOperation.
	fake.put("saas", map[string]interface{}{
		"name": "projects/test-proj/locations/us-east1/saas/tf-test-example-saas1",
	})
	fake.put("tenants", map[string]interface{}{
		"name": "projects/test-proj/locations/us-east1/tenants/tf-test-example-tenant1",
		"saas": "projects/123456789/locations/us-east1/saas/tf-test-example-saas1",
	})
	fake.put("unitKinds", map[string]interface{}{
		"name":           "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-unitkind1",
		"saas":           "projects/test-proj/locations/us-east1/saas/tf-test-example-saas1",
		"defaultRelease": "projects/123456789/locations/us-east1/releases/tf-test-example-release1",
	})
	fake.put("unitKinds", map[string]interface{}{
		"name": "projects/test-proj/locations/us-east1/unitKinds/tf-test-web-unitkind1",
		"saas": "projects/test-proj/locations/us-east1/saas/tf-test-example-saas1",
		"dependencies": []interface{}{
			map[string]interface{}{
				"alias":    "cluster",
				"unitKind": "projects/123456789/locations/us-east1/unitKinds/tf-test-cluster-unitkind1",
			},
		},
	})
	fake.put("releases", map[string]interface{}{
		"name":     "projects/test-proj/locations/us-east1/releases/tf-test-example-release1",
		"unitKind": "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-unitkind1",
	})
	fake.put("units", map[string]interface{}{
		"name":     "projects/test-proj/locations/us-east1/units/tf-test-base-unit1",
		"unitKind": "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-unitkind1",
		"release":  "projects/test-proj/locations/us-east1/releases/tf-test-example-release1",
		"tenant":   "projects/test-proj/locations/us-east1/tenants/tf-test-example-tenant1",
		"state":    "UNIT_STATE_READY",
		"conditions": []interface{}{
			map[string]interface{}{
				"type":   "TYPE_PROVISIONED",
				"status": "STATUS_TRUE",
			},
		},
	})
	fake.put("units", map[string]interface{}{
		"name":     "projects/test-proj/locations/us-east1/units/tf-test-zdependent-unit1",
		"unitKind": "projects/test-proj/locations/us-east1/unitKinds/tf-test-web-unitkind1",
		"tenant":   "projects/test-proj/locations/us-east1/tenants/tf-test-example-tenant1",
		"state":    "UNIT_STATE_NOT_PROVISIONED",
		"dependencies": []interface{}{
			map[string]interface{}{
				"alias": "cluster",
				"unit":  "projects/test-proj/locations/us-east1/units/tf-test-base-unit1",
			},
		},
	})
	fake.put("unitOperations", map[string]interface{}{
		"name":  "projects/test-proj/locations/us-east1/unitOperations/tf-test-running-op1",
		"unit":  "projects/test-proj/locations/us-east1/units/tf-test-base-unit1",
		"state": "UNIT_OPERATION_STATE_RUNNING",
	})

	// 2. Leaked ApplicationTemplate-imported Saas in us-central1 + replicated global
	//    copy with system-named "u-sa-comp-1-abc123" UnitKinds and "r-sa-comp-1-1" Releases.
	for _, loc := range []string{"us-central1", "global"} {
		fake.put("saas", map[string]interface{}{
			"name": fmt.Sprintf("projects/test-proj/locations/%s/saas/tf-test-app-saas2", loc),
		})
		fake.put("unitKinds", map[string]interface{}{
			"name":           fmt.Sprintf("projects/test-proj/locations/%s/unitKinds/u-sa-comp-1-abc123", loc),
			"saas":           fmt.Sprintf("projects/test-proj/locations/%s/saas/tf-test-app-saas2", loc),
			"defaultRelease": fmt.Sprintf("projects/test-proj/locations/%s/releases/r-sa-comp-1-1", loc),
		})
		fake.put("releases", map[string]interface{}{
			"name":     fmt.Sprintf("projects/test-proj/locations/%s/releases/r-sa-comp-1-1", loc),
			"unitKind": fmt.Sprintf("projects/test-proj/locations/%s/unitKinds/u-sa-comp-1-abc123", loc),
		})
	}

	// 3. Leaked RolloutKind in unreachable location europe-west1.
	fake.put("rolloutKinds", map[string]interface{}{
		"name":     "projects/test-proj/locations/europe-west1/rolloutKinds/tf-test-rolloutkind-eu",
		"unitKind": "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-unitkind1",
	})

	// 4. Non-test & permanent bootstrap fixtures that MUST NOT be modified or deleted.
	fake.put("saas", map[string]interface{}{
		"name": "projects/test-proj/locations/us-east1/saas/prod-saas",
	})
	fake.put("unitKinds", map[string]interface{}{
		"name":           "projects/test-proj/locations/us-east1/unitKinds/prod-unitkind",
		"saas":           "projects/test-proj/locations/us-east1/saas/prod-saas",
		"defaultRelease": "projects/test-proj/locations/us-east1/releases/prod-release",
	})
	fake.put("releases", map[string]interface{}{
		"name":     "projects/test-proj/locations/us-east1/releases/prod-release",
		"unitKind": "projects/test-proj/locations/us-east1/unitKinds/prod-unitkind",
	})
	fake.put("saas", map[string]interface{}{
		"name": "projects/test-proj/locations/us-central1/saas/tf-bootstrap-saas",
	})

	srv := httptest.NewServer(fake)
	defer srv.Close()
	cfg := newTestSaasRuntimeSweeperConfig(srv)

	// Run the exact topological sequence with region="us-central1" (matching TeamCity).
	steps := []struct {
		name string
		fn   func(*transport_tpg.Config, string, saasRuntimeSweeperTiming) error
	}{
		{"unit_operation", sweepSaasRuntimeUnitOperations},
		{"rollout_kind", sweepSaasRuntimeRolloutKinds},
		{"unit", sweepSaasRuntimeUnits},
		{"release", sweepSaasRuntimeReleases},
		{"tenant", sweepSaasRuntimeTenants},
		{"unit_kind", sweepSaasRuntimeUnitKinds},
		{"saas", sweepSaasRuntimeSaas},
	}

	for _, step := range steps {
		if err := step.fn(cfg, "us-central1", fastSaasRuntimeSweeperTiming); err != nil {
			t.Fatalf("sweeper step %q returned unexpected error: %v", step.name, err)
		}
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()

	var remaining []string
	for rel, rec := range fake.records {
		if !rec.deleted {
			remaining = append(remaining, rel)
		}
	}
	sort.Strings(remaining)

	wantRemaining := []string{
		"projects/test-proj/locations/us-central1/saas/tf-bootstrap-saas",
		"projects/test-proj/locations/us-east1/releases/prod-release",
		"projects/test-proj/locations/us-east1/saas/prod-saas",
		"projects/test-proj/locations/us-east1/unitKinds/prod-unitkind",
	}
	if strings.Join(remaining, "\n") != strings.Join(wantRemaining, "\n") {
		t.Fatalf("remaining resources mismatch:\ngot:\n%s\nwant:\n%s", strings.Join(remaining, "\n"), strings.Join(wantRemaining, "\n"))
	}

	// Verify prod-unitkind's defaultRelease was never cleared and no prod/bootstrap path was mutated.
	if gotDR, _ := fake.records["projects/test-proj/locations/us-east1/unitKinds/prod-unitkind"].obj["defaultRelease"].(string); gotDR == "" {
		t.Errorf("prod-unitkind defaultRelease was unexpectedly cleared")
	}
	for _, mut := range fake.mutatedPaths {
		if strings.Contains(mut, "prod-") || strings.Contains(mut, "tf-bootstrap-") {
			t.Errorf("unexpected mutation against non-sweepable resource: %s", mut)
		}
	}
}

func TestSaasRuntimeSweeper_WildcardFallbackAndStandaloneUnitKind(t *testing.T) {
	t.Parallel()

	fake := newFakeSaasRuntimeServer()
	fake.rejectWildcard = true

	fake.put("saas", map[string]interface{}{
		"name": "projects/test-proj/locations/us-east1/saas/tf-test-fallback-saas",
	})
	fake.put("unitKinds", map[string]interface{}{
		"name":           "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-uk",
		"saas":           "projects/test-proj/locations/us-east1/saas/tf-test-fallback-saas",
		"defaultRelease": "projects/test-proj/locations/us-east1/releases/tf-test-rel",
	})
	fake.put("unitKinds", map[string]interface{}{
		"name": "projects/test-proj/locations/us-east1/unitKinds/tf-test-zapp-uk",
		"saas": "projects/test-proj/locations/us-east1/saas/tf-test-fallback-saas",
		"dependencies": []interface{}{
			map[string]interface{}{
				"alias":    "cluster",
				"unitKind": "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-uk",
			},
		},
	})
	fake.put("releases", map[string]interface{}{
		"name":     "projects/test-proj/locations/us-east1/releases/tf-test-rel",
		"unitKind": "projects/test-proj/locations/us-east1/unitKinds/tf-test-cluster-uk",
	})

	srv := httptest.NewServer(fake)
	defer srv.Close()
	cfg := newTestSaasRuntimeSweeperConfig(srv)

	// Invoke sweepSaasRuntimeUnitKinds standalone (without running sweepSaasRuntimeReleases first)
	// while wildcard "locations/-" is rejected.
	if err := sweepSaasRuntimeUnitKinds(cfg, "us-central1", fastSaasRuntimeSweeperTiming); err != nil {
		t.Fatalf("standalone sweepSaasRuntimeUnitKinds failed: %v", err)
	}
	if err := sweepSaasRuntimeSaas(cfg, "us-central1", fastSaasRuntimeSweeperTiming); err != nil {
		t.Fatalf("sweepSaasRuntimeSaas failed: %v", err)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	for rel, rec := range fake.records {
		if !rec.deleted {
			t.Errorf("expected %s to be deleted, still present", rel)
		}
	}
}

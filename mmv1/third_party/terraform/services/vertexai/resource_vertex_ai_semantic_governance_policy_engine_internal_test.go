package vertexai

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// Unit tests for the SemanticGovernancePolicyEngine helpers declared in
// mmv1/templates/terraform/constants/semantic_governance_policy_engine.go.tmpl
// (generated into resource_vertex_ai_semantic_governance_policy_engine.go).

// fakeResourceData implements sgpeResourceID and records every SetId call,
// so tests can assert whether (and how often) post_read clears the ID.
type fakeResourceData struct {
	id        string
	setIdArgs []string
}

func (f *fakeResourceData) Id() string     { return f.id }
func (f *fakeResourceData) SetId(v string) { f.setIdArgs = append(f.setIdArgs, v); f.id = v }

// post_read translates INACTIVE to SetId("") so the next plan re-creates.
func TestUnitSGPE_PostReadDispatch_INACTIVE_clears_id(t *testing.T) {
	t.Parallel()
	d := &fakeResourceData{id: "projects/p/locations/us-central1/semanticGovernancePolicyEngine"}
	res := map[string]interface{}{"state": "INACTIVE"}
	if err := sgpePostReadDispatch(d, res); err != nil {
		t.Fatalf("expected nil error on INACTIVE, got %v", err)
	}
	if d.Id() != "" {
		t.Errorf("expected SetId(\"\") to clear the ID; got %q", d.Id())
	}
	if len(d.setIdArgs) != 1 || d.setIdArgs[0] != "" {
		t.Errorf("expected exactly one SetId(\"\") call; got %v", d.setIdArgs)
	}
}

// post_read returns an error only when the state field is missing or not a
// string (a backend contract break). Unrecognized enum values are left
// tracked, not errored — see TestUnitSGPE_PostReadDispatch_inStateBranches_no_op.
func TestUnitSGPE_PostReadDispatch_malformedState_return_error(t *testing.T) {
	t.Parallel()
	cases := map[string]map[string]interface{}{
		"state field missing":         {},
		"state field is not a string": {"state": 42},
	}
	for name, res := range cases {
		t.Run(name, func(t *testing.T) {
			d := &fakeResourceData{id: "irrelevant"}
			err := sgpePostReadDispatch(d, res)
			if err == nil {
				t.Fatalf("expected non-nil error for case %q; got nil", name)
			}
			if d.Id() == "" {
				t.Errorf("error path must NOT clear the ID; SetId(\"\") was called for case %q", name)
			}
		})
	}
}

// The no-op branch of the post_read contract: for every non-INACTIVE
// state the ID is left intact, so Terraform keeps the resource and does
// not plan a re-create. (Only INACTIVE clears the ID — see the test above.)
func TestUnitSGPE_PostReadDispatch_inStateBranches_no_op(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"ACTIVE", "FAILED", "PROVISIONING", "DEPROVISIONING", "STATE_UNSPECIFIED", "FUTURE_VALUE_NOT_YET_SUPPORTED"} {
		t.Run(state, func(t *testing.T) {
			origID := "projects/p/locations/us-central1/semanticGovernancePolicyEngine"
			d := &fakeResourceData{id: origID}
			err := sgpePostReadDispatch(d, map[string]interface{}{"state": state})
			if err != nil {
				t.Fatalf("expected nil error for state %q; got %v", state, err)
			}
			if d.Id() != origID {
				t.Errorf("state %q must NOT clear the ID; ID is now %q", state, d.Id())
			}
			if len(d.setIdArgs) != 0 {
				t.Errorf("state %q must NOT call SetId; got SetId calls %v", state, d.setIdArgs)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CustomizeDiff Check 1 (sgpeFailedStateWarningCheck)
// ---------------------------------------------------------------------------
//
// Check 1 is a soft warning: prior-state FAILED + any change to a mutable
// field logs [WARN] and returns nil. It never blocks the plan. Production
// asserts log content via TF_LOG=WARN; unit tests only assert return-nil
// semantics (log capture is fragile — see the preserve-branch precedent).

// sgpeFailedStateWarningCheck returns nil (does not block) when prior
// state is FAILED and a mutable field changed. Post-scrub the sole
// mutable field on Check 1's radar is gateway_config.
func TestUnitSGPE_FailedStateWarningCheck_FAILED_with_diff_succeeds(t *testing.T) {
	t.Parallel()
	d := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"state": "FAILED",
			"gateway_config": []map[string]interface{}{
				{"name": "primary", "network": "n1"},
			},
		},
		After: map[string]interface{}{
			"state": "FAILED",
			"gateway_config": []map[string]interface{}{
				{"name": "primary", "network": "n2"},
			},
		},
	}
	if err := sgpeFailedStateWarningCheck(d); err != nil {
		t.Fatalf("expected nil error (warnings are non-blocking); got %v", err)
	}
}

// sgpeFailedStateWarningCheck no-ops on: non-FAILED prior state (any
// mutable diff), and FAILED prior state with no mutable diff.
func TestUnitSGPE_FailedStateWarningCheck_negativeCases_no_op(t *testing.T) {
	t.Parallel()
	cases := map[string]*tpgresource.ResourceDiffMock{
		"ACTIVE with gateway_config diff": {
			Before: map[string]interface{}{
				"state":          "ACTIVE",
				"gateway_config": []map[string]interface{}{{"name": "primary", "network": "n1"}},
			},
			After: map[string]interface{}{
				"state":          "ACTIVE",
				"gateway_config": []map[string]interface{}{{"name": "primary", "network": "n2"}},
			},
		},
		"PROVISIONING with gateway_config diff": {
			Before: map[string]interface{}{
				"state":          "PROVISIONING",
				"gateway_config": []map[string]interface{}{{"name": "primary", "network": "n1"}},
			},
			After: map[string]interface{}{
				"state":          "PROVISIONING",
				"gateway_config": []map[string]interface{}{{"name": "primary", "network": "n2"}},
			},
		},
		"FAILED with empty diff": {
			Before: map[string]interface{}{
				"state":          "FAILED",
				"gateway_config": []map[string]interface{}{{"name": "primary", "network": "n1"}},
			},
			After: map[string]interface{}{
				"state":          "FAILED",
				"gateway_config": []map[string]interface{}{{"name": "primary", "network": "n1"}},
			},
		},
	}
	for name, d := range cases {
		t.Run(name, func(t *testing.T) {
			if err := sgpeFailedStateWarningCheck(d); err != nil {
				t.Fatalf("expected nil error; got %v", err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CustomizeDiff Check 2 (sgpeGatewayConfigChecksFunc)
// ---------------------------------------------------------------------------
//
// Check 2 is a plan-time hard error on duplicate gateway_config.name.
// Proto3 maps silently overwrite duplicate keys server-side, so the
// duplication must be caught before the RPC is issued.
//
// This feature intentionally drops the preserve-branch's rename detector
// (HLD §2.3 Option B): with the DEPROVISIONING-tombstone backend fix
// (cl/945777282) live, rename is a normal map delta the backend teardown
// chain handles. See HLD §2.3 for the tradeoff analysis. As a
// consequence, there is NO ForceNew-on-rename test here, and rename is
// asserted as a *legitimate* pass-through case below.

// sgpeGatewayConfigChecksFunc errors on the first duplicate name it
// encounters. Message text is asserted verbatim (substring) because
// downstream tooling and user-facing docs quote it.
func TestUnitSGPE_GatewayConfigChecks_duplicateName_returns_error(t *testing.T) {
	t.Parallel()
	d := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"gateway_config": []map[string]interface{}{
				{"name": "primary", "network": "n1"},
			},
		},
		After: map[string]interface{}{
			"gateway_config": []map[string]interface{}{
				{"name": "primary", "network": "n1"},
				{"name": "primary", "network": "n2"},
			},
		},
	}
	err := sgpeGatewayConfigChecksFunc(d)
	if err == nil {
		t.Fatal("expected duplicate-name error; got nil")
	}
	if !strings.Contains(err.Error(), `gateway_config.name "primary" is duplicated`) {
		t.Errorf("error message did not match expected text;\n  got:      %q\n  expected substring: %q",
			err.Error(), `gateway_config.name "primary" is duplicated`)
	}
}

// sgpeGatewayConfigChecksFunc passes through cleanly for every
// legitimate map delta: in-place value updates, pure adds, pure
// removes, empty-to-empty, and RENAMES (drop-old-key + add-new-key).
//
// Rename is included here (not in a separate ForceNew-fires test)
// because HLD §2.3 Option B routes rename through the backend
// DEPROVISIONING-tombstone chain, not through provider-side ForceNew.
// The provider's job on rename is to pass the map delta through — the
// same as any other update.
func TestUnitSGPE_GatewayConfigChecks_legitimateUpdates_returnNil(t *testing.T) {
	t.Parallel()
	cases := map[string]*tpgresource.ResourceDiffMock{
		"in-place per-element update (same names, different fields)": {
			Before: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1"},
					{"name": "secondary", "network": "n2"},
				},
			},
			After: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1-updated"},
					{"name": "secondary", "network": "n2"},
				},
			},
		},
		"pure add (one gateway becomes two)": {
			Before: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1"},
				},
			},
			After: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1"},
					{"name": "secondary", "network": "n2"},
				},
			},
		},
		"pure remove (two gateways become one)": {
			Before: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1"},
					{"name": "secondary", "network": "n2"},
				},
			},
			After: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1"},
				},
			},
		},
		"rename (drop primary, add main; HLD §2.3 Option B pass-through)": {
			Before: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "primary", "network": "n1"},
				},
			},
			After: map[string]interface{}{
				"gateway_config": []map[string]interface{}{
					{"name": "main", "network": "n1"},
				},
			},
		},
		"empty to empty (no gateways at all)": {
			Before: map[string]interface{}{},
			After:  map[string]interface{}{},
		},
	}
	for name, d := range cases {
		t.Run(name, func(t *testing.T) {
			if err := sgpeGatewayConfigChecksFunc(d); err != nil {
				t.Fatalf("expected nil error; got %v", err)
			}
		})
	}
}

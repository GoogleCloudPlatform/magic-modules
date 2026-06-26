package securityposture

// Unit tests for securitypostureGetPolicyIdForCannedConstraint and
// securitypostureGetEnforceIsSetForOrgPolicyConstraint.
//
// These helpers rely on d.GetRawConfig() which returns the unprocessed cty value
// tree from the Terraform configuration (preserving null vs. false for optional
// bools). schema.TestResourceDataRaw does not populate the raw config, so a
// minimal mock is used instead of a real *schema.ResourceData.

import (
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tpgresource "github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// rawConfigMock is a minimal implementation of tpgresource.TerraformResourceData
// that returns a pre-built cty.Value from GetRawConfig(). All other methods are
// no-ops sufficient to satisfy the interface.
type rawConfigMock struct {
	cfg cty.Value
}

func (m *rawConfigMock) GetRawConfig() cty.Value                 { return m.cfg }
func (m *rawConfigMock) HasChange(string) bool                   { return false }
func (m *rawConfigMock) GetOkExists(string) (interface{}, bool)  { return nil, false }
func (m *rawConfigMock) GetOk(string) (interface{}, bool)        { return nil, false }
func (m *rawConfigMock) Get(string) interface{}                  { return nil }
func (m *rawConfigMock) Set(string, interface{}) error           { return nil }
func (m *rawConfigMock) SetId(string)                            {}
func (m *rawConfigMock) Id() string                              { return "" }
func (m *rawConfigMock) Identity() (*schema.IdentityData, error) { return nil, nil }
func (m *rawConfigMock) GetProviderMeta(interface{}) error       { return nil }
func (m *rawConfigMock) Timeout(string) time.Duration            { return 0 }

var _ tpgresource.TerraformResourceData = (*rawConfigMock)(nil)

// boolPtr is a test helper to get a pointer to a bool literal.
func boolPtr(b bool) *bool { return &b }

// policyRuleType is the cty object type for a single policy_rule entry.
// Only the "enforce" field is relevant for these tests.
var policyRuleType = cty.Object(map[string]cty.Type{
	"enforce": cty.Bool,
})

// opcType is the cty object type for an org_policy_constraint entry.
var opcType = cty.Object(map[string]cty.Type{
	"canned_constraint_id": cty.String,
	"policy_rules":         cty.List(policyRuleType),
})

// constraintType is the cty object type for a constraint entry.
var constraintType = cty.Object(map[string]cty.Type{
	"org_policy_constraint": cty.List(opcType),
})

// policyType is the cty object type for a policy entry.
var policyType = cty.Object(map[string]cty.Type{
	"policy_id":  cty.String,
	"constraint": cty.List(constraintType),
})

// policySetType is the cty object type for a policy_set entry.
var policySetType = cty.Object(map[string]cty.Type{
	"policies": cty.List(policyType),
})

// buildPostureCfg constructs a cty.Value that mimics the raw Terraform config
// for a posture resource, given a slice of policy_set cty.Values.
func buildPostureCfg(policySets []cty.Value) cty.Value {
	if len(policySets) == 0 {
		return cty.ObjectVal(map[string]cty.Value{
			"policy_sets": cty.ListValEmpty(policySetType),
		})
	}
	return cty.ObjectVal(map[string]cty.Value{
		"policy_sets": cty.ListVal(policySets),
	})
}

// buildPolicySet constructs a cty.Value for a single policy_set containing the
// given policies.
func buildPolicySet(policies []cty.Value) cty.Value {
	if len(policies) == 0 {
		return cty.ObjectVal(map[string]cty.Value{
			"policies": cty.ListValEmpty(policyType),
		})
	}
	return cty.ObjectVal(map[string]cty.Value{
		"policies": cty.ListVal(policies),
	})
}

// buildPolicy constructs a cty.Value for a policy with the given policy_id and
// a single org_policy_constraint identified by cannedConstraintId. enforceValues
// controls the per-rule enforce field: a nil entry means the field was omitted
// (cty.NullVal), a non-nil entry means it was explicitly set to that value.
func buildPolicy(policyId, cannedConstraintId string, enforceValues []*bool) cty.Value {
	rules := make([]cty.Value, len(enforceValues))
	for i, ev := range enforceValues {
		if ev == nil {
			rules[i] = cty.ObjectVal(map[string]cty.Value{
				"enforce": cty.NullVal(cty.Bool),
			})
		} else if *ev {
			rules[i] = cty.ObjectVal(map[string]cty.Value{
				"enforce": cty.True,
			})
		} else {
			rules[i] = cty.ObjectVal(map[string]cty.Value{
				"enforce": cty.False,
			})
		}
	}

	var policyRulesVal cty.Value
	if len(rules) == 0 {
		policyRulesVal = cty.ListValEmpty(policyRuleType)
	} else {
		policyRulesVal = cty.ListVal(rules)
	}

	opcVal := cty.ObjectVal(map[string]cty.Value{
		"canned_constraint_id": cty.StringVal(cannedConstraintId),
		"policy_rules":         policyRulesVal,
	})
	constraintVal := cty.ObjectVal(map[string]cty.Value{
		"org_policy_constraint": cty.ListVal([]cty.Value{opcVal}),
	})
	return cty.ObjectVal(map[string]cty.Value{
		"policy_id":  cty.StringVal(policyId),
		"constraint": cty.ListVal([]cty.Value{constraintVal}),
	})
}

// buildPolicyWithNullOrgConstraint builds a policy whose org_policy_constraint
// list is null (simulating a policy that uses a different constraint type).
func buildPolicyWithNullOrgConstraint(policyId string) cty.Value {
	constraintVal := cty.ObjectVal(map[string]cty.Value{
		"org_policy_constraint": cty.NullVal(cty.List(opcType)),
	})
	return cty.ObjectVal(map[string]cty.Value{
		"policy_id":  cty.StringVal(policyId),
		"constraint": cty.ListVal([]cty.Value{constraintVal}),
	})
}

// ---------------------------------------------------------------------------
// Tests for securitypostureGetPolicyIdForCannedConstraint
// ---------------------------------------------------------------------------

func TestSecuritypostureGetPolicyIdForCannedConstraint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		cfg                cty.Value
		cannedConstraintId string
		want               string
	}{
		{
			name: "matching policy found",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("my-policy", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(true)}),
				}),
			}),
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               "my-policy",
		},
		{
			name: "no matching constraint returns empty string",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("my-policy", "constraints/compute.vmExternalIpAccess", []*bool{nil}),
				}),
			}),
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               "",
		},
		{
			name: "multiple policies, second one matches",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("policy-a", "constraints/compute.vmExternalIpAccess", []*bool{boolPtr(true)}),
					buildPolicy("policy-b", "constraints/iam.allowedPolicyMemberDomains", []*bool{nil}),
				}),
			}),
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               "policy-b",
		},
		{
			name: "null org_policy_constraint is skipped",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicyWithNullOrgConstraint("policy-a"),
					buildPolicy("policy-b", "constraints/iam.allowedPolicyMemberDomains", []*bool{nil}),
				}),
			}),
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               "policy-b",
		},
		{
			name:               "null raw config returns empty string",
			cfg:                cty.NullVal(cty.EmptyObject),
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               "",
		},
		{
			name: "duplicate canned_constraint_id returns first match",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("policy-first", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(true)}),
					buildPolicy("policy-second", "constraints/iam.allowedPolicyMemberDomains", []*bool{nil}),
				}),
			}),
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               "policy-first",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := &rawConfigMock{cfg: tc.cfg}
			got := securitypostureGetPolicyIdForCannedConstraint(d, tc.cannedConstraintId)
			if got != tc.want {
				t.Errorf("securitypostureGetPolicyIdForCannedConstraint() = %q, want %q", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Tests for securitypostureGetEnforceIsSetForOrgPolicyConstraint
// ---------------------------------------------------------------------------

func TestSecuritypostureGetEnforceIsSetForOrgPolicyConstraint(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		cfg                cty.Value
		policyId           string
		cannedConstraintId string
		want               []bool
	}{
		{
			name: "enforce explicitly set to true",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("p1", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(true)}),
				}),
			}),
			policyId:           "p1",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               []bool{true},
		},
		{
			name: "enforce explicitly set to false",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("p1", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(false)}),
				}),
			}),
			policyId:           "p1",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			// false is explicitly set, so enforce IS included in the request
			want: []bool{true},
		},
		{
			name: "enforce omitted (null) — list constraint",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("p1", "constraints/gcp.resourceLocations", []*bool{nil}),
				}),
			}),
			policyId:           "p1",
			cannedConstraintId: "constraints/gcp.resourceLocations",
			// nil means the user did not set enforce, so it must not be sent
			want: []bool{false},
		},
		{
			name: "multiple rules with mixed enforce settings",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("p1", "constraints/iam.allowedPolicyMemberDomains",
						[]*bool{boolPtr(true), nil, boolPtr(false)}),
				}),
			}),
			policyId:           "p1",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               []bool{true, false, true},
		},
		{
			name: "two policies with same constraint — correct policy selected by policy_id",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					// policy-a has enforce = true
					buildPolicy("policy-a", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(true)}),
					// policy-b has enforce omitted (list constraint)
					buildPolicy("policy-b", "constraints/iam.allowedPolicyMemberDomains", []*bool{nil}),
				}),
			}),
			policyId:           "policy-b",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			// Must return policy-b's enforce state (omitted), not policy-a's (set)
			want: []bool{false},
		},
		{
			name: "two policies with same constraint — policy-a selected",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("policy-a", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(true)}),
					buildPolicy("policy-b", "constraints/iam.allowedPolicyMemberDomains", []*bool{nil}),
				}),
			}),
			policyId:           "policy-a",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               []bool{true},
		},
		{
			name: "no matching policy returns nil",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("p1", "constraints/compute.vmExternalIpAccess", []*bool{boolPtr(true)}),
				}),
			}),
			policyId:           "p1",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               nil,
		},
		{
			name:               "null raw config returns nil",
			cfg:                cty.NullVal(cty.EmptyObject),
			policyId:           "p1",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               nil,
		},
		{
			name: "empty policy_id falls back to canned_constraint_id matching",
			cfg: buildPostureCfg([]cty.Value{
				buildPolicySet([]cty.Value{
					buildPolicy("p1", "constraints/iam.allowedPolicyMemberDomains", []*bool{boolPtr(true)}),
				}),
			}),
			policyId:           "",
			cannedConstraintId: "constraints/iam.allowedPolicyMemberDomains",
			want:               []bool{true},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := &rawConfigMock{cfg: tc.cfg}
			got := securitypostureGetEnforceIsSetForOrgPolicyConstraint(d, tc.policyId, tc.cannedConstraintId)
			if len(got) != len(tc.want) {
				t.Fatalf("securitypostureGetEnforceIsSetForOrgPolicyConstraint() len = %d, want %d (got %v, want %v)",
					len(got), len(tc.want), got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("securitypostureGetEnforceIsSetForOrgPolicyConstraint()[%d] = %v, want %v", i, got[i], tc.want[i])
				}
			}
		})
	}
}

package compute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestHasTerminationTimeChanged(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		Old, New map[string]interface{}
		Expect   bool
	}{
		"empty": {
			Old:    map[string]interface{}{"termination_time": ""},
			New:    map[string]interface{}{"termination_time": ""},
			Expect: false,
		},
		"new": {
			Old:    map[string]interface{}{"termination_time": ""},
			New:    map[string]interface{}{"termination_time": "2025-01-31T15:04:05Z"},
			Expect: true,
		},
		"changed": {
			Old:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			New:    map[string]interface{}{"termination_time": "2025-01-31T15:04:05Z"},
			Expect: true,
		},
		"same": {
			Old:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			New:    map[string]interface{}{"termination_time": "2025-01-30T15:04:05Z"},
			Expect: false,
		},
	}
	for tn, tc := range cases {
		if hasTerminationTimeChanged(tc.Old, tc.New) != tc.Expect {
			t.Errorf("%s: expected %t for whether termination time matched for old = %q, new = %q", tn, tc.Expect, tc.Old, tc.New)
		}
	}
}

func TestSchedulingHasChange_nilElements(t *testing.T) {
	t.Parallel()

	res := ResourceComputeInstance()

	// Empty ResourceData (empty scheduling slices)
	dEmpty := res.Data(nil)
	if schedulingHasChangeRequiringReboot(dEmpty) {
		t.Errorf("expected false from schedulingHasChangeRequiringReboot on empty ResourceData")
	}
	if schedulingHasChangeWithoutReboot(dEmpty) {
		t.Errorf("expected false from schedulingHasChangeWithoutReboot on empty ResourceData")
	}

	// ResourceData with scheduling = [nil] in config diff
	state := &terraform.InstanceState{
		ID: "test-instance",
		Attributes: map[string]string{
			"id":                             "test-instance",
			"name":                           "test-instance",
			"machine_type":                   "e2-medium",
			"scheduling.#":                   "1",
			"scheduling.0.automatic_restart": "true",
		},
	}
	rawCfg := map[string]interface{}{
		"name":         "test-instance",
		"machine_type": "e2-medium",
		"scheduling":   []interface{}{nil},
	}
	diff, err := schema.InternalMap(res.Schema).Diff(t.Context(), state, terraform.NewResourceConfigRaw(rawCfg), nil, nil, false)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	dNil, err := schema.InternalMap(res.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("Data failed: %v", err)
	}
	if schedulingHasChangeRequiringReboot(dNil) {
		t.Errorf("expected false when new scheduling slice has nil element")
	}
	if schedulingHasChangeWithoutReboot(dNil) {
		t.Errorf("expected false when new scheduling slice has nil element")
	}
}

func TestExpandComputeInstanceHelpers_nilElements(t *testing.T) {
	t.Parallel()

	d := &tpgresource.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"instance_encryption_key":  []interface{}{nil},
			"source_encryption_key":    []interface{}{nil},
			"workload_identity_config": []interface{}{nil},
		},
	}

	if got := expandComputeInstanceEncryptionKey(d); got != nil {
		t.Errorf("expected nil from expandComputeInstanceEncryptionKey([nil]), got %v", got)
	}
	if got := expandComputeInstanceSourceEncryptionKey(d, "source_encryption_key"); got != nil {
		t.Errorf("expected nil from expandComputeInstanceSourceEncryptionKey([nil]), got %v", got)
	}
	if got := expandWorkloadIdentityConfig(d); got != nil {
		t.Errorf("expected nil from expandWorkloadIdentityConfig([nil]), got %v", got)
	}
}

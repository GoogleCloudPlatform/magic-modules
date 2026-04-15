package chromepolicy

import (
	"testing"
)

func TestValidatePolicyFieldValueType(t *testing.T) {
	cases := []struct {
		name      string
		fieldType string
		value     interface{}
		want      bool
	}{
		{"bool true", "TYPE_BOOL", true, true},
		{"bool false", "TYPE_BOOL", false, true},
		{"bool string", "TYPE_BOOL", "true", false},
		{"string valid", "TYPE_STRING", "hello", true},
		{"string int", "TYPE_STRING", 42.0, false},
		{"enum valid", "TYPE_ENUM", "VALUE_A", true},
		{"int64 valid", "TYPE_INT64", float64(42), true},
		{"int64 float", "TYPE_INT64", 42.5, false},
		{"int64 string", "TYPE_INT64", "42", false},
		{"int32 valid", "TYPE_INT32", float64(32), true},
		{"int32 float", "TYPE_INT32", 32.5, false},
		{"double valid", "TYPE_DOUBLE", 3.14, true},
		{"double int", "TYPE_DOUBLE", float64(3), true},
		{"float valid", "TYPE_FLOAT", 3.14, true},
		{"message valid", "TYPE_MESSAGE", map[string]interface{}{"key": "val"}, true},
		{"message string", "TYPE_MESSAGE", "not a map", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := validatePolicyFieldValueType(tc.fieldType, tc.value)
			if got != tc.want {
				t.Errorf("validatePolicyFieldValueType(%q, %v) = %v, want %v", tc.fieldType, tc.value, got, tc.want)
			}
		})
	}
}

func TestSchemaNameMatchesFilter(t *testing.T) {
	cases := []struct {
		name   string
		schema string
		filter string
		want   bool
	}{
		{"exact match", "chrome.users.MaxConnections", "chrome.users.MaxConnections", true},
		{"exact no match", "chrome.users.MaxConnections", "chrome.users.OtherPolicy", false},
		{"wildcard match", "chrome.users.MaxConnections", "chrome.users.*", true},
		{"wildcard no match nested", "chrome.users.apps.InstallType", "chrome.users.*", false},
		{"wildcard different prefix", "chrome.devices.MaxConnections", "chrome.users.*", false},
		{"wildcard apps match", "chrome.users.apps.InstallType", "chrome.users.apps.*", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := schemaNameMatchesFilter(tc.schema, tc.filter)
			if got != tc.want {
				t.Errorf("schemaNameMatchesFilter(%q, %q) = %v, want %v", tc.schema, tc.filter, got, tc.want)
			}
		})
	}
}

func TestPolicyIdentityKey(t *testing.T) {
	id := policyIdentity{
		SchemaName: "chrome.users.apps.InstallType",
		AdditionalTargetKeys: map[string]string{
			"app_id": "chrome:abcdef",
		},
	}
	got := id.key()
	want := "chrome.users.apps.InstallType\x00app_id=chrome:abcdef"
	if got != want {
		t.Errorf("policyIdentity.key() = %q, want %q", got, want)
	}
}

func TestPolicyIdentityKey_noAdditionalKeys(t *testing.T) {
	id := policyIdentity{SchemaName: "chrome.users.MaxConnections"}
	got := id.key()
	if got != "chrome.users.MaxConnections" {
		t.Errorf("policyIdentity.key() = %q, want %q", got, "chrome.users.MaxConnections")
	}
}

func TestPolicySetsEqual(t *testing.T) {
	polA := map[string]interface{}{
		"schema": "chrome.users.A",
		"value":  map[string]interface{}{"key": "val1"},
	}
	polB := map[string]interface{}{
		"schema": "chrome.users.B",
		"value":  map[string]interface{}{"key": "val2"},
	}
	polAChanged := map[string]interface{}{
		"schema": "chrome.users.A",
		"value":  map[string]interface{}{"key": "val1_changed"},
	}

	setAB := policyMapByKey([]interface{}{polA, polB})
	setBA := policyMapByKey([]interface{}{polB, polA})
	setAOnly := policyMapByKey([]interface{}{polA})
	setAChanged := policyMapByKey([]interface{}{polAChanged, polB})

	if !policySetsEqual(setAB, setBA) {
		t.Error("expected sets with same policies in different order to be equal")
	}
	if policySetsEqual(setAB, setAOnly) {
		t.Error("expected sets with different lengths to not be equal")
	}
	if policySetsEqual(setAB, setAChanged) {
		t.Error("expected sets with different values to not be equal")
	}
}

func TestPolicyValuesEqual(t *testing.T) {
	a := map[string]interface{}{"value": map[string]interface{}{"k1": "v1", "k2": "v2"}}
	b := map[string]interface{}{"value": map[string]interface{}{"k2": "v2", "k1": "v1"}}
	c := map[string]interface{}{"value": map[string]interface{}{"k1": "v1", "k2": "changed"}}
	d := map[string]interface{}{"value": map[string]interface{}{"k1": "v1"}}

	if !policyValuesEqual(a, b) {
		t.Error("expected same values in different order to be equal")
	}
	if policyValuesEqual(a, c) {
		t.Error("expected different values to not be equal")
	}
	if policyValuesEqual(a, d) {
		t.Error("expected different key count to not be equal")
	}
}

func TestIsNonFatalDeleteError(t *testing.T) {
	if isNonFatalDeleteError(nil) {
		t.Error("expected nil error to not be non-fatal")
	}
}

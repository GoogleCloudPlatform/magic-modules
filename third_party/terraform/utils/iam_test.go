package google

import (
	"encoding/json"
	"google.golang.org/api/cloudresourcemanager/v1"
	"testing"
)

func TestIamRemoveBinding(t *testing.T) {
	table := []struct {
		input    []*cloudresourcemanager.Binding
		override *cloudresourcemanager.Binding
		expect   []*cloudresourcemanager.Binding
	}{
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			override: &cloudresourcemanager.Binding{
				Role:    "role-1",
				Members: []string{"new-member"},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"new-member"},
				},
			},
		},
		{

			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			override: &cloudresourcemanager.Binding{
				Role:    "role-2",
				Members: []string{"member-3"},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-3"},
				},
			},
		},
	}

	for _, test := range table {
		got := overwriteBinding(test.input, test.override)
		if !compareBindings(got, test.expect) {
			t.Errorf("OverwriteIamBinding got unexpected value.\nActual: %+v\nExpected: %+v",
				debugPrintBindings(got),
				debugPrintBindings(test.expect))
		}
	}
}

func TestIamMergeBindings(t *testing.T) {
	table := []struct {
		input  []*cloudresourcemanager.Binding
		expect []*cloudresourcemanager.Binding
	}{
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-3"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2", "member-3"},
				},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-3", "member-4"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-2", "member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-5"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
				{Role: "empty-role", Members: []string{}},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2", "member-3", "member-4", "member-5"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1"},
				},
			},
		},
	}
	for _, test := range table {
		got := mergeBindings(test.input)
		if !compareBindings(got, test.expect) {
			t.Errorf("MergeBinding return unexpected value.\nActual: %+v\nExpected: %+v",
				debugPrintBindings(got),
				debugPrintBindings(test.expect))
		}
	}
}

// Util to deref and print auditConfigs
func debugPrintAuditConfigs(bs []*cloudresourcemanager.AuditConfig) string {
	v, _ := json.MarshalIndent(bs, "", "\t")
	return string(v)
}

// Util to deref and print bindings
func debugPrintBindings(bs []*cloudresourcemanager.Binding) string {
	v, _ := json.MarshalIndent(bs, "", "\t")
	return string(v)
}

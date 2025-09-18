package compute

import (
	"reflect"
	"testing"
)

type testCase struct {
	obj       map[string]interface{}
	oldMatch  interface{}
	newMatch  interface{}
	expectObj map[string]interface{}
	expectErr bool
}
var cases = map[string]testCase {
	"no_match": {
		obj: map[string]interface{}{
			"action": "allow",
		},
		oldMatch: nil,
		newMatch: nil,
		expectObj: map[string]interface{}{
			"action": "allow",
		},
	},
	"no_network_scope_or_context": {
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"other_field":       789,
			},
		},
		expectObj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"other_field":       123,
			},
		},
	},
	"network_scope_no_change": {
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       789,
			},
		},
		expectObj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       123,
			},
		},
	},
	"network_scope_change": {
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"src_network_scope": "NON_INTERNET",
				"dest_network_scope": "INTERNET",
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       789,
			},
		},
		expectObj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       123,
			},
		},
	},
	"network_scope_and_context_no_change": {
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTRA_VPC",
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"src_network_scope": "INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTRA_VPC",
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"src_network_scope": "INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTRA_VPC",
				"other_field":       789,
			},
		},
		expectObj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTRA_VPC",
				"other_field":       123,
			},
		},
	},
	"network_scope_and_context_change_scope": {
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTERNET",
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"src_network_scope": "INTRA_VPC",
				"src_network_context": "INTERNET",
				"dest_network_scope": "NON_INTERNET",
				"dest_network_context": "INTERNET",
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTERNET",
				"other_field":       789,
			},
		},
		expectObj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				// scope changed - context field removed from obj
				"src_network_scope": "NON_INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"other_field":       123,
			},
		},
	},
	"network_scope_and_context_both_change": {
		expectErr: true,
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTERNET",
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"src_network_scope": "INTERNET",
				"src_network_context": "INTRA_VPC",
				"dest_network_scope": "NON_INTERNET",
				"dest_network_context": "NON_INTERNET",
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTERNET",
				"other_field":       789,
			},
		},
	},
	"network_scope_and_context_change_context": {
		obj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTERNET",
				"other_field":       123,
			},
		},
		oldMatch: []any{
			map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTRA_VPC",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "NON_INTERNET",
				"other_field":       456,
			},
		},
		newMatch: []any{
			map[string]any{
				"src_network_scope": "NON_INTERNET",
				"src_network_context": "INTERNET",
				"dest_network_scope": "INTRA_VPC",
				"dest_network_context": "INTERNET",
				"other_field":       789,
			},
		},
		expectObj: map[string]interface{}{
			"action": "allow",
			"match": map[string]any{
				// context changed - scope field removed from obj
				"src_network_context": "INTERNET",
				"dest_network_context": "INTERNET",
				"other_field":       123,
			},
		},
	},
}

func TestAdjustGlobalNetworkFirewallPolicyRuleNetworkContextFields(t *testing.T) {

	for tn, tc := range cases {
		err := adjustGlobalNetworkFirewallPolicyRuleNetworkContextFields(tc.obj, tc.oldMatch, tc.newMatch)
		if err != nil && !tc.expectErr {
			t.Errorf("%s: Unexpected error: %v", tn, err)
		}
		if err != nil {
			return
		}
		if !reflect.DeepEqual(tc.obj, tc.expectObj) {
			t.Errorf("%s: Incorrect object - want: %v, got: %v", tn, tc.expectObj, tc.obj)
		}
	}
}

func TestAdjustRegionNetworkFirewallPolicyRuleNetworkContextFields(t *testing.T) {

	for tn, tc := range cases {
		err := adjustRegionNetworkFirewallPolicyRuleNetworkContextFields(tc.obj, tc.oldMatch, tc.newMatch)
		if err != nil && !tc.expectErr {
			t.Errorf("%s: Unexpected error: %v", tn, err)
		}
		if err != nil {
			return
		}
		if !reflect.DeepEqual(tc.obj, tc.expectObj) {
			t.Errorf("%s: Incorrect object - want: %v, got: %v", tn, tc.expectObj, tc.obj)
		}
	}
}

func TestAdjustOrgFirewallPolicyRuleNetworkContextFields(t *testing.T) {

	for tn, tc := range cases {
		err := adjustOrgFirewallPolicyRuleNetworkContextFields(tc.obj, tc.oldMatch, tc.newMatch)
		if err != nil && !tc.expectErr {
			t.Errorf("%s: Unexpected error: %v", tn, err)
		}
		if err != nil {
			return
		}
		if !reflect.DeepEqual(tc.obj, tc.expectObj) {
			t.Errorf("%s: Incorrect object - want: %v, got: %v", tn, tc.expectObj, tc.obj)
		}
	}
}

package google

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceGoogleIamPolicyRead(t *testing.T) {
	cases := map[string]struct {
		Bindings                   []map[string]interface{}
		ExpectedPolicyDataString   string
		ExpectedBindingCount       int
		ExpectedPolicyBindingCount int
	}{
		"members are sorted alphabetically within a single binding": {
			Bindings: []map[string]interface{}{
				{
					"role": "role/A",
					"members": []interface{}{
						"user:c",
						"user:a",
						"user:b",
					},
				},
			},
			ExpectedBindingCount:       1,
			ExpectedPolicyBindingCount: 1,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"members\":[\"user:a\",\"user:b\",\"user:c\"],\"role\":\"role/A\"}]}",
		},
		"bindings are sorted by role": {
			Bindings: []map[string]interface{}{
				{
					"role": "role/B",
					"members": []interface{}{
						"user:a",
					},
				},
				{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
			},
			ExpectedBindingCount:       2,
			ExpectedPolicyBindingCount: 2,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"members\":[\"user:a\"],\"role\":\"role/B\"}]}",
		},
		"bindings with the same role are sorted by presence of a condition": {
			Bindings: []map[string]interface{}{
				{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionA",
							"expression":  "expressionA",
							"title":       "titleA",
						},
					},
				},
				// Binding with no condition should be placed first
				{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
			},
			ExpectedBindingCount:       2,
			ExpectedPolicyBindingCount: 2,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"members in equivalent bindings are consolidated": {
			Bindings: []map[string]interface{}{
				{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
				{
					"role": "role/A",
					"members": []interface{}{
						"user:b",
					},
				},
			},
			ExpectedBindingCount:       2,
			ExpectedPolicyBindingCount: 1, // Equivalent bindings combined into one member list
			ExpectedPolicyDataString:   "{\"bindings\":[{\"members\":[\"user:a\",\"user:b\"],\"role\":\"role/A\"}]}",
		},
		// "bindings for the same role with equivalent conditions are consolidated": {
		// 	Bindings: []map[string]interface{}{
		// 		{
		// 			"role": "role/A",
		// 			"members": []interface{}{
		// 				"user:b",
		// 			},
		// 			"condition": map[string]interface{}{
		// 				"description": "my description string",
		// 				"expression":  "my expression string",
		// 				"title":       "my title string",
		// 			},
		// 		},
		// 		{
		// 			"role": "role/A",
		// 			"members": []interface{}{
		// 				"user:a",
		// 			},
		// 			"condition": map[string]interface{}{
		// 				"description": "my description string",
		// 				"expression":  "my expression string",
		// 				"title":       "my title string",
		// 			},
		// 		},
		// 		// Should not be consolidated into the above as there's no condition
		// 		{
		// 			"role": "role/A",
		// 			"members": []interface{}{
		// 				"user:c",
		// 			},
		// 		},
		// 	},
		// 	ExpectedBindingCount:       3,
		// 	ExpectedPolicyBindingCount: 2,
		// 	ExpectedPolicyDataString:   "{\"bindings\":[{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\",\"user:b\"],\"role\":\"role/A\"}]},{\"members\":[\"user:c\"],\"role\":\"role/A\"}]}",
		// },
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange - Create schema.ResourceData variable as test input
			bindings := make([]interface{}, len(tc.Bindings))

			for i, b := range tc.Bindings {
				// Handle binding conditions, if set
				if b["condition"] != nil {
					c := b["condition"].([]interface{})
					// Avoid adding zero valued condition
					if len(c) == 1 {
						conditions := make([]interface{}, len(c))
						for j, con := range c {
							conditions[j] = con
						}
						b["condition"] = conditions
					}
				}
				bindings[i] = b
			}

			rawData := map[string]interface{}{
				"binding":      bindings,
				"policy_data":  "",              // Not set
				"audit_config": []interface{}{}, // Not set
			}

			d := schema.TestResourceDataRaw(t, dataSourceGoogleIamPolicy().Schema, rawData)

			// Act - Update resource data using `dataSourceGoogleIamPolicyRead`
			var meta interface{}
			err := dataSourceGoogleIamPolicyRead(d, meta)
			if err != nil {
				t.Error(err)
			}

			// Assertions

			policyData := d.Get("policy_data").(string)
			var jsonObjs interface{}
			json.Unmarshal([]byte(policyData), &jsonObjs)
			objSlice, ok := jsonObjs.(map[string]interface{})
			if !ok {
				t.Errorf("cannot convert the JSON string")
			}
			policyDataBindings := objSlice["bindings"].([]interface{})
			if len(policyDataBindings) != tc.ExpectedPolicyBindingCount {
				t.Errorf("expected there to be %d bindings in the policy_data string, got: %d", tc.ExpectedPolicyBindingCount, len(policyDataBindings))
			}
			if policyData != tc.ExpectedPolicyDataString {
				t.Errorf("expected `policy_data` to be %s, got: %s", tc.ExpectedPolicyDataString, policyData)
			}

			bset := d.Get("binding").(*schema.Set)
			if bset.Len() != tc.ExpectedBindingCount {
				t.Errorf("expected there to be %d bindings in the data source internals, got: %d", tc.ExpectedBindingCount, bset.Len())
			}
		})
	}
}

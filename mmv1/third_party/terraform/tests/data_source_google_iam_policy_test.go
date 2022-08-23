package google

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceGoogleIamPolicyRead(t *testing.T) {
	cases := map[string]struct {
		Bindings                   []interface{}
		ExpectedPolicyDataString   string
		ExpectedBindingCount       int
		ExpectedPolicyBindingCount int
	}{
		"members are sorted alphabetically within a single binding": {
			Bindings: []interface{}{
				map[string]interface{}{
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
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/B",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
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
		"members in equivalent bindings (with no conditions) are consolidated": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
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
		"members in equivalent bindings (with equivalent conditions) are consolidated ": {
			Bindings: []interface{}{
				// Should not be consolidated into the other bindings as there's no condition
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:c",
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:b",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionA",
							"expression":  "expressionA",
							"title":       "titleA",
						},
					},
				},
				map[string]interface{}{
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
			},
			ExpectedBindingCount:       3,
			ExpectedPolicyBindingCount: 2,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"members\":[\"user:c\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\",\"user:b\"],\"role\":\"role/A\"}]}",
		},
		"bindings with the same role are sorted by presence of a condition": {
			Bindings: []interface{}{
				map[string]interface{}{
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
				// Binding with no condition should be placed first in role/A bindings
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
					"role": "role/B",
					"members": []interface{}{
						"user:b",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionB",
							"expression":  "expressionB",
							"title":       "titleB",
						},
					},
				},
			},
			ExpectedBindingCount:       3,
			ExpectedPolicyBindingCount: 3,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"descriptionB\",\"expression\":\"expressionB\",\"title\":\"titleB\"},\"members\":[\"user:b\"],\"role\":\"role/B\"}]}",
		},
		"bindings on the same role with different conditions are sorted by condition title": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"title": "C",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"title": "B",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"title": "A",
						},
					},
				},
			},
			ExpectedBindingCount:       3,
			ExpectedPolicyBindingCount: 3,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"condition\":{\"description\":\"\",\"expression\":\"\",\"title\":\"A\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"\",\"expression\":\"\",\"title\":\"B\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"\",\"expression\":\"\",\"title\":\"C\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"bindings on the same role with different conditions, with the same title, are next sorted by condition expression": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"title":      "same title",
							"expression": "C",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"title":      "same title",
							"expression": "B",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"title":      "same title",
							"expression": "A",
						},
					},
				},
			},
			ExpectedBindingCount:       3,
			ExpectedPolicyBindingCount: 3,
			ExpectedPolicyDataString:   "{\"bindings\":[{\"condition\":{\"description\":\"\",\"expression\":\"A\",\"title\":\"same title\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"\",\"expression\":\"B\",\"title\":\"same title\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"\",\"expression\":\"C\",\"title\":\"same title\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// ARRANGE - Create schema.ResourceData variable as test input
			rawData := map[string]interface{}{
				"binding":      tc.Bindings,
				"policy_data":  "",              // Not set
				"audit_config": []interface{}{}, // Not set
			}
			// Note: for TestResourceDataRaw to process rawData ok, test inputs' data types have to be
			// either primitive types, []interface{} or map[string]interface{}
			d := schema.TestResourceDataRaw(t, dataSourceGoogleIamPolicy().Schema, rawData)

			// ACT - Update resource data using `dataSourceGoogleIamPolicyRead`
			var meta interface{}
			err := dataSourceGoogleIamPolicyRead(d, meta)
			if err != nil {
				t.Error(err)
			}

			// ASSERT
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

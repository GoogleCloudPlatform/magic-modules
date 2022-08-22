package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceGoogleIamPolicyRead(t *testing.T) {
	cases := map[string]struct {
		Bindings                 []map[string]interface{}
		ExpectedPolicyDataString string
		ExpectedBindingCount     int
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
			ExpectedBindingCount:     1,
			ExpectedPolicyDataString: "{\"bindings\":[{\"members\":[\"user:a\",\"user:b\",\"user:c\"],\"role\":\"role/A\"}]}",
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
			ExpectedBindingCount:     2,
			ExpectedPolicyDataString: "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"members\":[\"user:a\"],\"role\":\"role/B\"}]}",
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
							"description": "description A",
							"expression":  "expression A",
							"title":       "title A",
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
			ExpectedBindingCount:     2,
			ExpectedPolicyDataString: "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"description A\",\"expression\":\"expression A\",\"title\":\"title A\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
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
							t.Logf("CONDITION IS %#v", con)
							conditions[j] = con
						}
						b["condition"] = conditions
					} else {
						t.Logf("Condition length != 1, is %#v", c)
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

			// Assert
			val, ok := d.GetOk("policy_data")
			if !ok {
				t.Errorf("expected `policy_data` to be gettable but it's not. Val is %#v", val)
			}
			if val != tc.ExpectedPolicyDataString {
				t.Errorf("expected `policy_data` to be %s, got: %s", tc.ExpectedPolicyDataString, val)
			}

			if len(tc.Bindings) != tc.ExpectedBindingCount {
				t.Errorf("expected there to be %d bindings, got: %d", tc.ExpectedBindingCount, len(tc.Bindings))
			}
		})
	}
}

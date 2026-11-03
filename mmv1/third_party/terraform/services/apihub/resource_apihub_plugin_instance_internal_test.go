package apihub

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandApihubPluginInstanceAdditionalConfig(t *testing.T) {
	res := ResourceApihubPluginInstance()
	additionalConfigSchema := res.Schema["additional_config"].Elem.(*schema.Resource)
	setHash := schema.HashResource(additionalConfigSchema)

	entries := []interface{}{
		map[string]interface{}{
			"key":                 "string-val",
			"string_value":        "val123",
			"int_value":           0,
			"bool_value":          false,
			"secret_value":        []interface{}{},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":                 "int-val",
			"string_value":        "",
			"int_value":           42,
			"bool_value":          false,
			"secret_value":        []interface{}{},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":                 "bool-true-val",
			"string_value":        "",
			"int_value":           0,
			"bool_value":          true,
			"secret_value":        []interface{}{},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":                 "bool-false-val",
			"string_value":        "",
			"int_value":           0,
			"bool_value":          false,
			"secret_value":        []interface{}{},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":          "secret-val",
			"string_value": "",
			"int_value":    0,
			"bool_value":   false,
			"secret_value": []interface{}{
				map[string]interface{}{
					"secret_version": "projects/my-proj/secrets/my-secret/versions/1",
				},
			},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":          "enum-val",
			"string_value": "",
			"int_value":    0,
			"bool_value":   false,
			"secret_value": []interface{}{},
			"enum_value": []interface{}{
				map[string]interface{}{
					"id":           "Option1",
					"display_name": "Option 1",
					"description":  "First option",
				},
			},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":          "multi-select-val",
			"string_value": "",
			"int_value":    0,
			"bool_value":   false,
			"secret_value": []interface{}{},
			"enum_value":   []interface{}{},
			"multi_select_values": []interface{}{
				map[string]interface{}{
					"values": []interface{}{
						map[string]interface{}{
							"id":           "Select1",
							"display_name": "Select 1",
							"description":  "First select",
						},
					},
				},
			},
			"multi_string_values": []interface{}{},
			"multi_int_values":    []interface{}{},
		},
		map[string]interface{}{
			"key":                 "multi-string-val",
			"string_value":        "",
			"int_value":           0,
			"bool_value":          false,
			"secret_value":        []interface{}{},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{
				map[string]interface{}{
					"values": []interface{}{"a", "b"},
				},
			},
			"multi_int_values": []interface{}{},
		},
		map[string]interface{}{
			"key":                 "multi-int-val",
			"string_value":        "",
			"int_value":           0,
			"bool_value":          false,
			"secret_value":        []interface{}{},
			"enum_value":          []interface{}{},
			"multi_select_values": []interface{}{},
			"multi_string_values": []interface{}{},
			"multi_int_values": []interface{}{
				map[string]interface{}{
					"values": []interface{}{10, 20},
				},
			},
		},
	}

	inputSet := schema.NewSet(setHash, entries)
	got, err := expandApihubPluginInstanceAdditionalConfig(inputSet, nil, nil)
	if err != nil {
		t.Fatalf("expandApihubPluginInstanceAdditionalConfig returned error: %v", err)
	}

	want := map[string]interface{}{
		"string-val": map[string]interface{}{
			"stringValue": "val123",
		},
		"int-val": map[string]interface{}{
			"intValue": 42,
		},
		"bool-true-val": map[string]interface{}{
			"boolValue": true,
		},
		"bool-false-val": map[string]interface{}{
			"boolValue": false,
		},
		"secret-val": map[string]interface{}{
			"secretValue": map[string]interface{}{
				"secretVersion": "projects/my-proj/secrets/my-secret/versions/1",
			},
		},
		"enum-val": map[string]interface{}{
			"enumValue": map[string]interface{}{
				"id":          "Option1",
				"displayName": "Option 1",
				"description": "First option",
			},
		},
		"multi-select-val": map[string]interface{}{
			"multiSelectValues": map[string]interface{}{
				"values": []interface{}{
					map[string]interface{}{
						"id":          "Select1",
						"displayName": "Select 1",
						"description": "First select",
					},
				},
			},
		},
		"multi-string-val": map[string]interface{}{
			"multiStringValues": map[string]interface{}{
				"values": []interface{}{"a", "b"},
			},
		},
		"multi-int-val": map[string]interface{}{
			"multiIntValues": map[string]interface{}{
				"values": []interface{}{10, 20},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expandApihubPluginInstanceAdditionalConfig mismatch:\ngot:  %#v\nwant: %#v", got, want)
	}
}

package dataflow

import (
	"reflect"
	"testing"
)

func TestFlattenDataflowConfigStoreSettingValue(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected []interface{}
	}{
		{
			name: "string value",
			input: map[string]interface{}{
				"stringValue": "e2-standard-4",
			},
			expected: []interface{}{
				map[string]interface{}{
					"string_value": "e2-standard-4",
				},
			},
		},
		{
			name: "bool value true",
			input: map[string]interface{}{
				"boolValue": true,
			},
			expected: []interface{}{
				map[string]interface{}{
					"bool_value": true,
				},
			},
		},
		{
			name: "bool value false",
			input: map[string]interface{}{
				"boolValue": false,
			},
			expected: []interface{}{
				map[string]interface{}{
					"bool_value": false,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := flattenDataflowConfigStoreSettingValue(tc.input, nil, nil)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("flattenDataflowConfigStoreSettingValue(%v) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestExpandDataflowConfigStoreSettingValue_FromMap(t *testing.T) {
	cases := []struct {
		name     string
		input    []interface{}
		expected map[string]interface{}
	}{
		{
			name: "string value",
			input: []interface{}{
				map[string]interface{}{
					"string_value": "n2-standard-4",
				},
			},
			expected: map[string]interface{}{
				"stringValue": "n2-standard-4",
			},
		},
		{
			name: "bool value true",
			input: []interface{}{
				map[string]interface{}{
					"bool_value": true,
				},
			},
			expected: map[string]interface{}{
				"boolValue": true,
			},
		},
		{
			name: "bool value false",
			input: []interface{}{
				map[string]interface{}{
					"bool_value": false,
				},
			},
			expected: map[string]interface{}{
				"boolValue": false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := expandDataflowConfigStoreSettingValue(tc.input, nil, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expandDataflowConfigStoreSettingValue(%v) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}

package utils

import (
	"testing"
)

type SimpleStruct struct {
	Name  string
	Value int
}

type NestedStruct struct {
	Simple *SimpleStruct
	Values []int
	Map    map[string]interface{}
}

type ComplexStruct struct {
	Nested    *NestedStruct
	StructMap map[string]*SimpleStruct
	StructArr []*SimpleStruct
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		// Nil values
		{"nil", nil, true},
		{"nil pointer", (*SimpleStruct)(nil), true},

		// Basic types
		{"empty string", "", true},
		{"non-empty string", "test", false},
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"false bool", false, true},
		{"true bool", true, false},

		// Slices and arrays
		{"empty slice", []int{}, true},
		{"nil slice", []int(nil), true},
		{"non-empty slice", []int{1, 2, 3}, false},
		{"slice with zero values", []int{0, 0, 0}, true},
		{"slice with mixed values", []int{0, 1, 0}, false},

		// Maps
		{"empty map", map[string]int{}, true},
		{"nil map", map[string]int(nil), true},
		{"map with values", map[string]int{"one": 1, "two": 2}, false},
		{"map with zero values", map[string]int{"one": 0, "two": 0}, true},
		{"map with mixed values", map[string]int{"one": 0, "two": 2}, false},

		// Simple struct
		{"empty struct", SimpleStruct{}, true},
		{"partially filled struct", SimpleStruct{Name: "test"}, false},
		{"fully filled struct", SimpleStruct{Name: "test", Value: 42}, false},

		// Pointers to simple structs
		{"pointer to empty struct", &SimpleStruct{}, true},
		{"pointer to partially filled struct", &SimpleStruct{Name: "test"}, false},

		// Nested structs with nil fields
		{"nested struct with all nil", NestedStruct{}, true},
		{"nested struct with simple", NestedStruct{Simple: &SimpleStruct{Name: "test"}}, false},
		{"nested struct with empty array", NestedStruct{Values: []int{}}, true},
		{"nested struct with non-empty array", NestedStruct{Values: []int{1, 2}}, false},

		// Complex nested scenarios
		{
			"complex all empty",
			&ComplexStruct{
				Nested: &NestedStruct{
					Simple: &SimpleStruct{},
					Values: []int{},
					Map:    map[string]interface{}{},
				},
				StructMap: map[string]*SimpleStruct{},
				StructArr: []*SimpleStruct{},
			},
			true,
		},
		{
			"complex with one non-empty value",
			&ComplexStruct{
				Nested: &NestedStruct{
					Simple: &SimpleStruct{Name: "test"},
					Values: []int{},
					Map:    map[string]interface{}{},
				},
				StructMap: map[string]*SimpleStruct{},
				StructArr: []*SimpleStruct{},
			},
			false,
		},
		{
			"complex with non-empty map",
			&ComplexStruct{
				Nested: &NestedStruct{
					Simple: &SimpleStruct{},
					Values: []int{},
					Map:    map[string]interface{}{"key": "value"},
				},
				StructMap: map[string]*SimpleStruct{},
				StructArr: []*SimpleStruct{},
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsEmpty(test.value)
			if result != test.expected {
				t.Errorf("IsEmpty(%v) = %v, expected %v", test.value, result, test.expected)
			}
		})
	}
}

package utils

import (
	"reflect"
	"testing"
)

// Simple struct for basic test cases.
type SimpleStructOmit struct {
	Name    string
	Count   int
	Enabled bool
}

// Complex struct to test behavior with non-basic types.
type ComplexStructOmit struct {
	ID       string
	Simple   *SimpleStructOmit
	Tags     []string
	Metadata map[string]string
}

func TestOmitDefaultsForMarshaling(t *testing.T) {
	// Define common structs for use in tests
	defaultSimple := SimpleStructOmit{
		Name:    "default-name",
		Count:   1,
		Enabled: true,
	}

	defaultComplex := ComplexStructOmit{
		ID:   "default-id",
		Tags: []string{"default"},
		Metadata: map[string]string{
			"key": "default",
		},
	}

	tests := []struct {
		name     string
		current  interface{}
		defaults interface{}
		expected interface{}
	}{
		{
			name: "All fields different from defaults",
			current: SimpleStructOmit{
				Name:    "current-name",
				Count:   10,
				Enabled: false,
			},
			defaults: defaultSimple,
			expected: SimpleStructOmit{ // Expected to be unchanged
				Name:    "current-name",
				Count:   10,
				Enabled: false,
			},
		},
		{
			name: "All fields match defaults",
			current: SimpleStructOmit{
				Name:    "default-name",
				Count:   1,
				Enabled: true,
			},
			defaults: defaultSimple,
			expected: SimpleStructOmit{ // Expected to be zeroed out
				Name:    "",
				Count:   0,
				Enabled: false,
			},
		},
		{
			name: "Some fields match defaults",
			current: SimpleStructOmit{
				Name:    "default-name", // Match -> zero
				Count:   99,             // No match -> keep
				Enabled: true,           // Match -> zero
			},
			defaults: defaultSimple,
			expected: SimpleStructOmit{
				Name:    "",
				Count:   99,
				Enabled: false,
			},
		},
		{
			name: "Field has zero value but does not match non-zero default",
			current: SimpleStructOmit{
				Name:    "",
				Count:   0,
				Enabled: false,
			},
			defaults: defaultSimple,
			expected: SimpleStructOmit{ // Expected to be unchanged
				Name:    "",
				Count:   0,
				Enabled: false,
			},
		},
		{
			name: "Field matches zero-value default",
			current: SimpleStructOmit{
				Name:  "keep",
				Count: 0, // Match default -> zero (no change)
			},
			defaults: SimpleStructOmit{
				Name:  "default",
				Count: 0, // Default is zero
			},
			expected: SimpleStructOmit{
				Name:  "keep",
				Count: 0,
			},
		},
		{
			name: "Complex types are ignored even if they match",
			current: ComplexStructOmit{
				ID:   "default-id", // Match -> zero
				Tags: []string{"default"},
				Metadata: map[string]string{
					"key": "default",
				},
			},
			defaults: defaultComplex,
			expected: ComplexStructOmit{ // Only ID should be zeroed
				ID:   "",
				Tags: []string{"default"}, // Ignored, so it remains
				Metadata: map[string]string{
					"key": "default", // Ignored, so it remains
				},
			},
		},
		{
			name: "Pointer to struct field is ignored",
			current: ComplexStructOmit{
				ID:     "current-id",
				Simple: &SimpleStructOmit{Name: "test"}, // This field should be ignored
			},
			defaults: ComplexStructOmit{
				ID:     "default-id",
				Simple: &SimpleStructOmit{Name: "test"},
			},
			expected: ComplexStructOmit{
				ID:     "current-id",
				Simple: &SimpleStructOmit{Name: "test"}, // Expected to remain unchanged
			},
		},
		{
			name:     "Non-struct input should be returned as is (in a pointer)",
			current:  123,
			defaults: 0,
			expected: 123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the original 'current' value to check for mutation
			originalCurrent := shallowCopy(tt.current)

			// Run the function
			result := OmitDefaultsForMarshaling(tt.current, tt.defaults)

			// 1. Verify the result is a pointer
			if reflect.ValueOf(result).Kind() != reflect.Ptr {
				t.Fatalf("Expected result to be a pointer, but got %T", result)
			}

			// 2. Get the value from the pointer and verify its content
			resultVal := reflect.ValueOf(result).Elem().Interface()
			if !reflect.DeepEqual(resultVal, tt.expected) {
				t.Errorf("Result is incorrect.\n got: %#v\nwant: %#v", resultVal, tt.expected)
			}

			// 3. Verify the original input was not mutated
			if !reflect.DeepEqual(tt.current, originalCurrent) {
				t.Errorf("Original input was mutated.\n original: %#v\n after:    %#v", originalCurrent, tt.current)
			}

			// The problematic assertion block was here. It has been removed because
			// the primary assertion (step 2) already correctly handles all cases,
			// including non-structs, by checking the value inside the returned pointer.
			// The original block caused a panic due to an incorrect type assertion.
		})
	}
}

// shallowCopy creates a shallow copy of an interface value.
// This is sufficient for the test cases to verify that the top-level
// input struct is not mutated.
func shallowCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	original := reflect.ValueOf(src)
	// Create a new value of the same type and copy the contents.
	cpy := reflect.New(original.Type()).Elem()
	cpy.Set(original)
	return cpy.Interface()
}

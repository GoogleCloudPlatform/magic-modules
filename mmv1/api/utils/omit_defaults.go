package utils

import (
	"fmt"
	"reflect"
)

// OmitDefaultsForMarshaling creates a clone of a struct and sets any fields
// that match the corresponding field in a 'defaults' struct to their zero value.
// This is used in custom MarshalYAML implementations to allow `omitempty` to work
// for fields that have non-zero default values.
//
// Note: This function only processes the top-level fields of the provided struct.
// It does not recurse into nested structs, slices, or maps. This is by design,
// as it is intended to be used within a tree of custom MarshalYAML functions.
// The standard YAML marshaler will handle traversing the object graph, and if a
// nested struct also has a MarshalYAML method that calls this function, the
// default-omission logic will be applied at that level as well.
//
// The `current` and `defaults` arguments must be structs of the same type.
// The function returns a pointer to the new, modified struct. If `current` is
// not a struct, it is returned unmodified in a pointer.
func OmitDefaultsForMarshaling(current, defaults interface{}) (interface{}, error) {
	// Get the reflect.Value of the current struct.
	currentVal := reflect.ValueOf(current)
	if currentVal.Kind() == reflect.Ptr {
		currentVal = currentVal.Elem()
	}

	// Ensure we are working with a struct.
	if currentVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("omitDefaultsForMarshaling expects the current object to be a struct")
	}

	// Create a new pointer to a struct of the same type as 'current'.
	// This will be our clone that we can modify safely.
	clonePtr := reflect.New(currentVal.Type())
	cloneElem := clonePtr.Elem()

	// Copy the data from the original struct to our clone.
	cloneElem.Set(currentVal)

	// Get the reflect.Value of the defaults struct.
	defaultsVal := reflect.ValueOf(defaults)

	// Iterate over the fields of the struct.
	for i := 0; i < cloneElem.NumField(); i++ {
		field := cloneElem.Field(i)
		defaultsField := defaultsVal.Field(i)

		// Ensure the field is exported and can be set.
		if !field.CanSet() {
			continue
		}

		var areEqual bool

		switch field.Kind() {
		// For basic, comparable types, use a direct comparison.
		case reflect.String, reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128:

			areEqual = (field.Interface() == defaultsField.Interface())

		// For all other complex types (slices, maps, structs, etc.), ignore.
		default:
			areEqual = false
		}

		// If the values are the same, set the field in the clone to its zero value.
		// This allows the `omitempty` tag to work during YAML marshaling.
		if areEqual {
			zeroValue := reflect.Zero(field.Type())
			field.Set(zeroValue)
		}
	}

	// Return the pointer to the modified clone.
	return clonePtr.Interface(), nil
}

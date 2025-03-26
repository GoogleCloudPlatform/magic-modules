package utils

import "reflect"

// IsEmpty checks if a value is meaningfully empty in a recursive way
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	// Handle pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true
		}
		return IsEmpty(val.Elem().Interface())
	}

	// Handle different types
	switch val.Kind() {
	case reflect.Struct:
		// Check if all fields are empty
		allEmpty := true
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.CanInterface() && !IsEmpty(field.Interface()) {
				allEmpty = false
				break
			}
		}
		return allEmpty

	case reflect.Map:
		if val.Len() == 0 {
			return true
		}
		// Check if all map values are empty
		allEmpty := true
		iter := val.MapRange()
		for iter.Next() {
			if !IsEmpty(iter.Value().Interface()) {
				allEmpty = false
				break
			}
		}
		return allEmpty

	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return true
		}
		// Check if all elements are empty
		allEmpty := true
		for i := 0; i < val.Len(); i++ {
			if !IsEmpty(val.Index(i).Interface()) {
				allEmpty = false
				break
			}
		}
		return allEmpty

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return val.IsNil()

	default:
		// For simple types (int, string, etc.), check if it's a zero value
		return val.IsZero()
	}
}

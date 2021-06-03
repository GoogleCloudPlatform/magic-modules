package google

import (
	"errors"
)

// ErrNoConversion can be returned if a conversion is unable to be returned
// because of the current state of the system.
// Example: The conversion requires that the resource has already been created
// and is now being updated).
var ErrNoConversion = errors.New("no conversion")

// ErrEmptyIdentityField can be returned when fetching a resource is not possible
// due to the identity field of that resource returning empty.
var ErrEmptyIdentityField = errors.New("empty identity field")

// Global MutexKV
var mutexKV = NewMutexKV()

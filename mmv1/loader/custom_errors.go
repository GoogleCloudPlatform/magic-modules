package loader

import "fmt"

// ErrProductVersionNotFound is returned when a product doesn't exist
// at the specified version or any lower version.
type ErrProductVersionNotFound struct {
	ProductName string
	Version     string
}

// Error implements the error interface.
func (e *ErrProductVersionNotFound) Error() string {
	return fmt.Sprintf("%s does not have a '%s' version", e.ProductName, e.Version)
}

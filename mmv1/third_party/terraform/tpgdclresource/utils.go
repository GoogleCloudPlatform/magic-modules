package tpgdclresource

import (
	"time"

	"github.com/google/go-cpy/cpy"
)

// Copy makes a deep copy of an interface.
func Copy(src any) any {
	copier := cpy.New(
		cpy.Shallow(time.Time{}),
		cpy.IgnoreAllUnexported(),
	)
	return copier.Copy(src)
}

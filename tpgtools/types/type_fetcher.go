package types

import "github.com/nasa9084/go-openapi"

// TypeFetcher fetches reused types, as marked by the $ref field being marked on an OpenAPI schema.
type TypeFetcher struct {
	doc *openapi.Document

	// Tracks if a property has already been generated.
	generates map[string]string
}

// NewTypeFetcher returns a TypeFetcher for a OpenAPI document.
func NewTypeFetcher(doc *openapi.Document) *TypeFetcher {
	return &TypeFetcher{
		doc:       doc,
		generates: make(map[string]string),
	}
}

// ResolveSchema resolves a #/components/schemas reference from a reused type.
func (r *TypeFetcher) ResolveSchema(ref string) (*openapi.Schema, error) {
	return openapi.ResolveSchema(r.doc, ref)
}

// PackagePathForReference returns either the packageName or the shared package name.
func (r *TypeFetcher) PackagePathForReference(ref, packageName string) string {
	if v, ok := r.generates[ref]; ok {
		return v
	} else {
		r.generates[ref] = packageName
		return packageName
	}
}

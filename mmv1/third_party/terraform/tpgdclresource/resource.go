package tpgdclresource

type Resource interface {
	Describe() ServiceTypeVersion
}

// ServiceTypeVersion is a tuple that can uniquely identify a
// DCL resource type.
type ServiceTypeVersion struct {
	// Service indicates the service to which this resource
	// belongs, e.g., "compute". It is roughly analogous to the
	// K8S "Group" identifier.
	Service string

	// Type identifies the particular type of this resource,
	// e.g., "ComputeInstance". It maps to the K8S "Kind".
	Type string

	// Version is the DCL version of the resource, e.g.,
	// "beta" or "ga".
	Version string
}

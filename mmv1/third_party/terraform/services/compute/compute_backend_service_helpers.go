package compute

// emptySecurityPolicyReference returns an empty request body for the
// setSecurityPolicy / setEdgeSecurityPolicy API calls. Callers should add a
// "securityPolicy" key only when a non-empty policy URL is provided; omitting
// the key clears the policy on the resource.
func emptySecurityPolicyReference() map[string]interface{} {
	return map[string]interface{}{}
}

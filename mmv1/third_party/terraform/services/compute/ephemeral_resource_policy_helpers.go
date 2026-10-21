package compute

// An existing diskConsistencyGroupPolicy object denotes enablement. Return
// the API-shaped object expected by the ephemeral ObjectAttribute converter.
func flattenEphemeralDiskConsistencyGroupPolicy(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	return map[string]interface{}{"enabled": true}, nil
}

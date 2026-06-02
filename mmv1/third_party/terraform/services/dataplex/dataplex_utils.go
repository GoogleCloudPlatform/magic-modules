package dataplex

// flattenZoneDiscoverySpecEnable flattens an instance of discovery spec from a JSON
// response object.
func flattenZoneDiscoverySpecEnable(c *Client, i any, _ *Zone) *bool {
	v, ok := i.(bool)
	if !ok {
		v = false
	}
	return &v
}

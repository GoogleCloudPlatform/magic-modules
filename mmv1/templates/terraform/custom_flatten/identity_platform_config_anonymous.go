func flatten<%= prefix -%><%= titlelize_property(property) -%>(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}

	original := v.(map[string]interface{})
	transformed := make(map[string]interface{})

	if original["enabled"] == nil {
		transformed["enabled"] = false
	} else {
		transformed["enabled"] = original["enabled"]
	}

	return []interface{}{transformed}
}

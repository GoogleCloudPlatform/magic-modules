func suppressAttachedClustersLoggingConfigDiff(_, old, new string, d *schema.ResourceData) bool {
	if old == new {
		return true
	}
	_, n := d.GetChange("logging_config.0.component_config.0.enable_components")
	if tpgresource.IsEmptyValue(reflect.ValueOf(n)) {
		return true
	}
	return false
}

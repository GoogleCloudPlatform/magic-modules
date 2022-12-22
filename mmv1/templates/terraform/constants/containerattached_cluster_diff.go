func suppressAttachedClustersFleetProjectDiff(_, old, new string, _ *schema.ResourceData) bool {
    if old == new {
        return true
    }
    // The custom expander prepends projects/ to the supplied id, but the new value has not gone
    // through that modification yet.
    new = "projects/" + new
    if old == new {
        return true
    }

    return false
}

func suppressAttachedClustersLoggingConfigDiff(_, old, new string, d *schema.ResourceData) bool {
    if old == new {
        return true
    }
    _, n := d.GetChange("logging_config.0.component_config.0.enable_components")
    if isEmptyValue(reflect.ValueOf(n)) {
        return true
    }
    return false
}

func flatten<%= prefix -%><%= titlelize_property(property) -%>(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
    if v == nil {
        return false
    }
    return v
}

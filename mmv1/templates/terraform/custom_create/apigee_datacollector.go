
userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
if err != nil {
    return err
}
obj := make(map[string]interface{})
nameProp, err := expandApigeeDataCollectorName(d.Get("name"), d, config)
if err != nil {
    return err
} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
    obj["name"] = nameProp
}
if v, ok := d.GetOkExists("description"); ok {
    obj["description"] = v
}
if v, ok := d.GetOkExists("type"); ok {
    obj["type"] = v
}
url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}/organizations/{{org_id}}/datacollectors")
if err != nil {
    return err
}
log.Printf("[DEBUG] Creating new DataCollector: %#v", obj)
billingProject := ""
if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
    billingProject = bp
}
res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
    Config:    config,
    Method:    "POST",
    Project:   billingProject,
    RawURL:    url,
    UserAgent: userAgent,
    Body:      obj,
    Timeout:   d.Timeout(schema.TimeoutCreate),
})
if err != nil {
    return fmt.Errorf("Error creating DataCollector: %s", err)
}
id, err := tpgresource.ReplaceVars(d, config, "{{org_id}}/datacollectors/{{name}}")
if err != nil {
    return fmt.Errorf("Error constructing id: %s", err)
}
d.SetId(id)
log.Printf("[DEBUG] Finished creating DataCollector %q: %#v", d.Id(), res)
return resourceApigeeDataCollectorRead(d, meta)

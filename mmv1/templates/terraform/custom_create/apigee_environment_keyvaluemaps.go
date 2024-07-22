userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
if err != nil {
	return err
}

obj := make(map[string]interface{})
nameProp, err := expandApigeeEnvironmentKeyvaluemapsName(d.Get("name"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("name"); !tpgresource.IsEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
	obj["name"] = nameProp
}

url, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{env_id}}/keyvaluemaps")
if err != nil {
	return err
}

log.Printf("[DEBUG] Creating new EnvironmentKeyvaluemaps: %#v", obj)
billingProject := ""

// err == nil indicates that the billing_project value was found
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
	return fmt.Errorf("Error creating EnvironmentKeyvaluemaps: %s", err)
}

// Store the ID now
id, err := tpgresource.ReplaceVars(d, config, "{{env_id}}/keyvaluemaps/{{name}}")
if err != nil {
	return fmt.Errorf("Error constructing id: %s", err)
}
d.SetId(id)

log.Printf("[DEBUG] Finished creating EnvironmentKeyvaluemaps %q: %#v", d.Id(), res)

return resourceApigeeEnvironmentKeyvaluemapsRead(d, meta)
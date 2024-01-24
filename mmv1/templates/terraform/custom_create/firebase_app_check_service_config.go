userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
if err != nil {
	return err
}

obj := make(map[string]interface{})
enforcementModeProp, err := expandFirebaseAppCheckServiceConfigEnforcementMode(d.Get("enforcement_mode"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("enforcement_mode"); !tpgresource.IsEmptyValue(reflect.ValueOf(enforcementModeProp)) && (ok || !reflect.DeepEqual(v, enforcementModeProp)) {
	obj["enforcementMode"] = enforcementModeProp
}
log.Printf("[DEBUG] Creating new ServiceConfig: %#v", obj)

project, err := tpgresource.GetProject(d, config)
if err != nil {
	return fmt.Errorf("Error fetching project for ServiceConfig: %s", err)
}
billingProject := project

url, err := tpgresource.ReplaceVars(d, config, "{{FirebaseAppCheckBasePath}}projects/{{project}}/services/{{service_id}}")
if err != nil {
	return err
}

// Custom logic: add all mutable fields to the updateMask
url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": "enforcementMode"})
if err != nil {
	return err
}

// err == nil indicates that the billing_project value was found
if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
	billingProject = bp
}

res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
	Config:    config,
	Method:    "PATCH",
	Project:   billingProject,
	RawURL:    url,
	UserAgent: userAgent,
	Body:      obj,
	Timeout:   d.Timeout(schema.TimeoutCreate),
})
if err != nil {
	return fmt.Errorf("Error creating ServiceConfig: %s", err)
}
if err := d.Set("name", flattenFirebaseAppCheckServiceConfigName(res["name"], d, config)); err != nil {
	return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
}

// Store the ID now
id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/services/{{service_id}}")
if err != nil {
	return fmt.Errorf("Error constructing id: %s", err)
}
d.SetId(id)

log.Printf("[DEBUG] Finished creating ServiceConfig %q: %#v", d.Id(), res)

return resourceFirebaseAppCheckServiceConfigRead(d, meta)
var obj map[string]interface{}
log.Printf("[DEBUG] Deleting ServiceConfig %q", d.Id())

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
	Timeout:   d.Timeout(schema.TimeoutDelete),
})
if err != nil {
	return transport_tpg.HandleNotFoundError(err, d, "ServiceConfig")
}

log.Printf("[DEBUG] Finished deleting ServiceConfig %q: %#v", d.Id(), res)
return nil
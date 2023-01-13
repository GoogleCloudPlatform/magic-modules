userAgent, err := generateUserAgentString(d, config.userAgent)
if err != nil {
	return err
}

url, err := replaceVars(d, config, "{{IdentityPlatformBasePath}}projects/{{project}}/identityPlatform:initializeAuth")
if err != nil {
	return err
}

billingProject := ""

project, err := getProject(d, config)
if err != nil {
	return fmt.Errorf("Error fetching project for Config: %s", err)
}
billingProject = project

// err == nil indicates that the billing_project value was found
if bp, err := getBillingProject(d, config); err == nil {
	billingProject = bp
}

res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, nil, d.Timeout(schema.TimeoutCreate))
if err != nil {
	return fmt.Errorf("Error creating Config: %s", err)
}
if err := d.Set("name", flattenIdentityPlatformConfigName(res["name"], d, config)); err != nil {
	return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
}

// Store the ID now
id, err := replaceVars(d, config, "projects/{{project}}/config")
if err != nil {
	return fmt.Errorf("Error constructing id: %s", err)
}
d.SetId(id)

// Update the resource after initializing auth to set fields.
if err := resourceIdentityPlatformConfigUpdate(d, meta); err != nil {
	return err
}

log.Printf("[DEBUG] Finished creating Config %q: %#v", d.Id(), res)

return resourceIdentityPlatformConfigRead(d, meta)

var NotebooksInstanceProvidedScopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
}

func NotebooksInstanceScopesDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
  old, new := d.GetChange("service_account_scopes")
	oldValue := old.([]interface{})
	newValue := new.([]interface{})
	oldValueList := []string{}
	newValueList := []string{}

	for _, item := range oldValue {
		oldValueList = append(oldValueList,item.(string))
	}

	for _, item := range newValue {
		newValueList = append(newValueList,item.(string))
	}
	newValueList= append(newValueList,NotebooksInstanceProvidedScopes...)

	sort.Strings(oldValueList)
	sort.Strings(newValueList)
	if reflect.DeepEqual(oldValueList, newValueList) {
		return true
	}
	return false
}

func NotebooksInstanceKmsDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	if strings.HasPrefix(old, new) {
		return true
	}
	return false
}

<% unless compiler == "terraformgoogleconversion-codegen" -%>
// waitForNotebooksInstanceActive waits for an Notebook instance to become "ACTIVE"
func waitForNotebooksInstanceActive(d *schema.ResourceData, config *transport_tpg.Config, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		if err := resourceNotebooksInstanceRead(d, config); err != nil {
			return resource.NonRetryableError(err)
		}

		name := d.Get("name").(string)
		state := d.Get("state").(string)
		if state == "ACTIVE" {
			log.Printf("[DEBUG] Notebook Instance %q has state %q.", name, state)
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("Notebook Instance %q has state %q. Waiting for ACTIVE state", name, state))
		}

	})
}
<% end -%>

func modifyNotebooksInstanceState(config *transport_tpg.Config, d *schema.ResourceData, project string, billingProject string, userAgent string, state string) (map[string]interface{}, error) {
	url, err := tpgresource.ReplaceVars(d, config, "{{NotebooksBasePath}}projects/{{project}}/locations/{{location}}/instances/{{name}}:"+state)
	if err != nil {
		return nil, err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config: config,
		Method: "POST",
		Project: billingProject,
		RawURL: url,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to %q google_notebooks_instance %q: %s", state, d.Id(), err)
	}
	return res, nil
}

<% unless compiler == "terraformgoogleconversion-codegen" -%>
func waitForNotebooksOperation(config *transport_tpg.Config, d *schema.ResourceData, project string, billingProject string, userAgent string, response map[string]interface{}) error {
	var opRes map[string]interface{}
	err := NotebooksOperationWaitTimeWithResponse(
		config, response, &opRes, project, "Modifying Notebook Instance state", userAgent,
		d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}
<% end -%>

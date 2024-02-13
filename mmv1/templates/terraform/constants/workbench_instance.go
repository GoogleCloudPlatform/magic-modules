var WorkbenchInstanceProvidedLabels = []string{
	"consumer-project-id",
	"consumer-project-number",
	"notebooks-product",
	"resource-name",
}

func WorkbenchInstanceLabelsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the labels
	for _, label := range WorkbenchInstanceProvidedLabels {
		if strings.Contains(k, label) && new == "" {
			return true
		}
	}

	// Let diff be determined by labels (above)
	if strings.Contains(k, "labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}


var WorkbenchInstanceProvidedMetadata = []string{
	"disable-swap-binaries",
	"enable-guest-attributes",
	"proxy-backend-id",
	"proxy-registration-url",
	"agent-health-check-interval-seconds",
	"agent-health-check-path",
	"container",
	"data-disk-uri",
	"dataproc-allow-custom-clusters",
	"dataproc-cluster-name",
	"dataproc-configs",
	"dataproc-default-subnet",
	"dataproc-locations-list",
	"dataproc-machine-types-list",
	"dataproc-notebooks-url",
	"dataproc-region",
	"dataproc-service-account",
	"disable-check-xsrf",
	"framework",
	"gcs-data-bucket",
	"generate-diagnostics-bucket",
	"generate-diagnostics-file",
	"generate-diagnostics-options",
	"image-url",
	"install-monitoring-agent",
	"install-nvidia-driver",
	"installed-extensions",
	"notebooks-api",
	"notebooks-api-version",
	"notebooks-examples-location",
	"notebooks-location",
	"nvidia-driver-gcs-path",
	"proxy-mode",
	"proxy-status",
	"proxy-url",
	"proxy-user-mail",
	"report-container-health",
	"report-notebook-metrics",
	"report-system-health",
	"report-system-status",
	"restriction",
	"serial-port-logging-enable",
	"shutdown-script",
	"title",
	"use-collaborative",
	"version",
	"enable-oslogin",
}

func WorkbenchInstanceMetadataDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the Metadata
	for _, metadata := range WorkbenchInstanceProvidedMetadata {
		if strings.Contains(k, metadata) && new == "" {
			return true
		}
	}

	// Let diff be determined by metadata
	if strings.Contains(k, "gce_setup.0.metadata.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}


var WorkbenchInstanceProvidedTags = []string{
	"deeplearning-vm",
	"notebook-instance",
}

func WorkbenchInstanceTagsDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
  old, new := d.GetChange("gce_setup.0.tags")
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
	newValueList= append(newValueList,WorkbenchInstanceProvidedTags...)

	sort.Strings(oldValueList)
	sort.Strings(newValueList)
	if reflect.DeepEqual(oldValueList, newValueList) {
		return true
	}
	return false
}

<% unless compiler == "terraformgoogleconversion-codegen" -%>
// waitForWorkbenchInstanceActive waits for an workbench instance to become "ACTIVE"
func waitForWorkbenchInstanceActive(d *schema.ResourceData, config *transport_tpg.Config, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		if err := resourceWorkbenchInstanceRead(d, config); err != nil {
			return resource.NonRetryableError(err)
		}

		name := d.Get("name").(string)
		state := d.Get("state").(string)
		if state == "ACTIVE" {
			log.Printf("[DEBUG] Workbench Instance %q has state %q.", name, state)
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("Workbench Instance %q has state %q. Waiting for ACTIVE state", name, state))
		}

	})
}
<% end -%>

func modifyWorkbenchInstanceState(config *transport_tpg.Config, d *schema.ResourceData, project string, billingProject string, userAgent string, state string) (map[string]interface{}, error) {
	url, err := tpgresource.ReplaceVars(d, config, "{{WorkbenchBasePath}}projects/{{project}}/locations/{{location}}/instances/{{name}}:"+state)
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
		return nil, fmt.Errorf("Unable to %q google_workbench_instance %q: %s", state, d.Id(), err)
	}
	return res, nil
}

func WorkbenchInstanceKmsDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	if strings.HasPrefix(old, new) {
		return true
	}
	return false
}

<% unless compiler == "terraformgoogleconversion-codegen" -%>
func waitForWorkbenchOperation(config *transport_tpg.Config, d *schema.ResourceData, project string, billingProject string, userAgent string, response map[string]interface{}) error {
	var opRes map[string]interface{}
	err := WorkbenchOperationWaitTimeWithResponse(
		config, response, &opRes, project, "Modifying Workbench Instance state", userAgent,
		d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	return nil
}
<% end -%>

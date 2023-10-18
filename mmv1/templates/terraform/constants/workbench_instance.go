var WorkbenchInstanceProvidedLabels = []string{
	"consumer-project-id",
	"consumer-project-number",
	"notebooks-product",
	"resource-name",
}

var WorkbenchInstanceProvidedMetadata = []string{
	"disable-swap-binaries",
	"enable-guest-attributes",
	"enable-oslogin",
	"install-nvidia-driver",
	"notebooks-api",
	"notebooks-api-version",
	"nvidia-driver-gcs-path",
	"proxy-backend-id",
	"proxy-mode",
	"proxy-registration-url",
	"proxy-url",
	"proxy-user-mail",
	"serial-port-logging-enable",
	"shutdown-script",
}

var WorkbenchInstanceProvidedTags = []string{
	"deeplearning-vm",
	"notebook-instance",
}

func WorkbenchInstanceTagsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the tags
	for _, tag := range WorkbenchInstanceProvidedTags {
		if strings.Contains(k, tag) {
			return true
		}
	}

	// For other keys, don't suppress diff.
	return false
}

func WorkbenchInstanceLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the labels
	for _, label := range WorkbenchInstanceProvidedLabels {
		if strings.Contains(k, label){
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

func WorkbenchInstanceMetadataDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the metadata
	for _, metadata := range WorkbenchInstanceProvidedMetadata {
		if strings.Contains(k, metadata) && new == "" {
			return true
		}
	}

	// Let diff be determined by metadata (above)
	if strings.Contains(k, "metadata.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

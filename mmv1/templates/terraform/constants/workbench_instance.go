var WorkbenchInstanceGoogleProvidedLabels = []string{
	"consumer-project-id",
	"consumer-project-number",
	"notebooks-product",
	"resource-name"
}

func WorkbenchInstanceLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the labels provided by Google
	for _, label := range WorkbenchInstanceGoogleProvidedLabels {
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

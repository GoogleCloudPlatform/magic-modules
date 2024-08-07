const notebooksRuntimeGoogleProvidedLabel = "goog-caip-managed-notebook"

func NotebooksRuntimeLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the label provided by Google
	if strings.Contains(k, notebooksRuntimeGoogleProvidedLabel) && new == "" {
		return true
	}

	// Let diff be determined by labels (above)
	if strings.Contains(k, "labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

// NotReturnedByAPIDiffSuppress
func NotReturnedByAPIDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return true
}
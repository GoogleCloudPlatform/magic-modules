const notebooksInstanceGoogleProvidedLabel = "goog-caip-notebook"

func NotebooksInstanceLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// Suppress diffs for the label provided by Google
	if strings.Contains(k, notebooksInstanceGoogleProvidedLabel) && new == "" {
		return true
	}

	// Let diff be determined by labels (above)
	if strings.Contains(k, "labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

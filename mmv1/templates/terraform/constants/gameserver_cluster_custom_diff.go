func suppressSuffixDiff(_, old, new string, _ *schema.ResourceData) bool {
	if strings.HasSuffix(old, new) {
		log.Printf("[INFO] suppressing diff as %s is the same as the full path of %s", new, old)
		return true
	}

	return false
}
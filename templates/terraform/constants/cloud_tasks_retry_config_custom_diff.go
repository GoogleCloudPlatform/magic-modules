func suppressOmittedMaxDuration(_, old, new string, _ *schema.ResourceData) bool {
	if old == "" && new == "0s" {
		log.Printf("[INFO] max retry is 0s and api omitted field, suppressing diff")
		return true
	}
	return false
}

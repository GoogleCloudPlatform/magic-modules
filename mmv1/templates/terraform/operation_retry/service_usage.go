// Retries errors on 403 3 times if the error message
// returned contains `has not been used in project`
maxRetries := 3
if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 403 {
	if w.retryCount < maxRetries && strings.Contains(gerr.Body, "has not been used in project") {
		w.retryCount += 1
		log.Printf("[DEBUG] retrying on 403 %v more times", w.retryCount-maxRetries-1)
		return true
	}
}
return false
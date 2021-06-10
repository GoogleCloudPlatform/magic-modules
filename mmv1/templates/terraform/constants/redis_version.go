// Is the new redis version less than the old one?
func isRedisVersionDecreasing(_ context.Context, old, new, _ interface{}) bool {
	if old == nil || new == nil {
		return false
	}
	re := regexp.MustCompile(`REDIS_(\d+)_(\d+)`)
	oldParsed := re.FindSubmatch([]byte(old.(string)))
	newParsed := re.FindSubmatch([]byte(new.(string)))

	if oldParsed == nil || newParsed == nil {
		// create new if you don't recognize the expression
		return true
	}

	oldVersion, err := strconv.ParseFloat(fmt.Sprintf("%s.%s", oldParsed[1], oldParsed[2]), 32)
	if err != nil {
		return false
	}
	newVersion, err := strconv.ParseFloat(fmt.Sprintf("%s.%s", newParsed[1], newParsed[2]), 32)
	if err != nil {
		return false
	}

	return newVersion < oldVersion
}
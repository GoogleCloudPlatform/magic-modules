
// Is the new redis version less than the old one?
func isRedisVersionDecreasing(_ context.Context, old, new, _ interface{}) bool {
	return isRedisVersionDecreasingFunc(old, new)
}

// separate function for unit testing
func isRedisVersionDecreasingFunc(old, new interface{}) bool {
	if old == nil || new == nil {
		return false
	}
	re := regexp.MustCompile(`REDIS_(\d+)_(\d+)`)
	oldParsed := re.FindSubmatch([]byte(old.(string)))
	newParsed := re.FindSubmatch([]byte(new.(string)))

	if oldParsed == nil || newParsed == nil {
		return false
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

// returns true if old=new or old='auto'
func secondaryIpDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
  if ((strings.ToLower(new) == "auto" && old != "") || old == new) {
    return true
  }
  return false
}


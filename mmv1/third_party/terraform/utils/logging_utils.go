package google

import (
	"fmt"
	"regexp"
)

// loggingSinkResourceTypes contains all the possible Stackdriver Logging resource types. Used to parse ids safely.
var loggingSinkResourceTypes = []string{
	"billingAccounts",
	"folders",
	"organizations",
	"projects",
}

// loggingSinkIdRegex matches valid logging sink canonical ids
var loggingSinkIdRegex = regexp.MustCompile("(.+)/(.+)/sinks/(.+)")

// parseLoggingSinkParentId parses a canonical id to get sink parent resource id
func parseLoggingSinkParentId(id string) (string, error) {
	parts := loggingSinkIdRegex.FindStringSubmatch(id)
	if parts == nil {
		return "", fmt.Errorf("unable to parse logging sink id %#v", id)
	}
	// If our resourceType is not a valid logging sink resource type, complain loudly
	validLoggingSinkResourceType := false
	for _, v := range loggingSinkResourceTypes {
		if v == parts[1] {
			validLoggingSinkResourceType = true
			break
		}
	}

	if !validLoggingSinkResourceType {
		return "", fmt.Errorf("Logging resource type %s is not valid. Valid resource types: %#v", parts[1],
			loggingSinkResourceTypes)
	}
	return parts[2], nil
}

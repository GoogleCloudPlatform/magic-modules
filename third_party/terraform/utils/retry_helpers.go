package google

import (
	"strings"

	"google.golang.org/api/googleapi"
)

func iamMemberMissing(err error) bool {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "permission") {
			return true
		}
	}
	return false
}

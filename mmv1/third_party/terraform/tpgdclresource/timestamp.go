package tpgdclresource

import (
	"time"
)

// ProtoToTime converts a string from a DCL proto time string to a time.Time.
func ProtoToTime(s string) time.Time {
	// Invalid time values will be picked up downstream.
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

// TimeToProto converts a time.Time to a proto time string.
func TimeToProto(t time.Time) string {
	return t.Format(time.RFC3339)
}

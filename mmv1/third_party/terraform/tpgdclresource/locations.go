package tpgdclresource

import "regexp"

// IsRegion returns true if this string refers to a GCP region.
func IsRegion(s *string) bool {
	if s == nil {
		return false
	}

	r := regexp.MustCompile(`^[a-z]+-[a-z]+[0-9]+$`)
	return r.MatchString(*s)
}

// IsZone returns true if this string refers to a GCP zone.
func IsZone(s *string) bool {
	if s == nil {
		return false
	}

	r := regexp.MustCompile(`^[a-z]+-[a-z]+[0-9]+-(ai[0-9]+)?[a-z]+$`)
	return r.MatchString(*s)
}

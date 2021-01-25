// TODO: export the DCL one or something
package main

import (
	"regexp"
	"strings"
)

// Map from initialism -> TitleCase variant
// We can assume camelCase is the same as TitleCase except that we downcase the
// first segment
var initialisms = map[string]string{
	"ip":     "IP",
	"ipv4":   "IPv4",
	"ipv6":   "IPv6",
	"oauth":  "OAuth",
	"oauth2": "OAuth2",
	"tpu":    "TPU",
	"vpc":    "VPC",
}

// snakeToTitleCase converts a snake_case string to TitleCase / Go struct case.
func snakeToTitleCase(s string) string {
	return strings.Join(snakeToTitleParts(s), "")
}

// snakeToTitleParts returns the parts of a snake_case string titlecased as an
// array, taking into account common initialisms.
func snakeToTitleParts(s string) []string {
	parts := []string{}
	segments := strings.Split(s, "_")
	for _, seg := range segments {
		if v, ok := initialisms[seg]; ok {
			parts = append(parts, v)
		} else {
			parts = append(parts, strings.ToUpper(seg[0:1])+seg[1:])
		}
	}

	return parts
}

// jsonToSnakeCase converts a jsonCase string to snake_case.
func jsonToSnakeCase(s string) string {
	for _, v := range initialisms {
		s = strings.ReplaceAll(s, v, v[0:1]+strings.ToLower(v[1:]))
	}
	result := regexp.MustCompile("(.)([A-Z][^A-Z]+)").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(result, "${1}_${2}"))
}

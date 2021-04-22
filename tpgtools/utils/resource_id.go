package utils

import "regexp"

const patternPart = "{{(\\w+)}}"

func IdParts(id string) (parts []string) {
	r := regexp.MustCompile(patternPart)

	// returns [["{{field}}", "field"] ...]
	idTmplAndParts := r.FindAllStringSubmatch(id, -1)
	for _, v := range idTmplAndParts {
		parts = append(parts, v[1])
	}

	return parts
}

// PatternToRegex formats a pattern string into a Python-compatible regex.
func PatternToRegex(s string) string {
	re := regexp.MustCompile(patternPart)
	return re.ReplaceAllString(s, "(?P<$1>[^/]+)")
}

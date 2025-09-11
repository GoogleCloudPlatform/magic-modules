package networksecurity

import "strings"

func flattenName(name string) string {
	tokens := strings.Split(name, "/")
	return tokens[len(tokens)-1]
}

func flattenProjectName(name string) string {
	tokens := strings.Split(name, "/")
	if len(tokens) < 2 || tokens[0] != "projects" {
		return ""
	}
	return tokens[1]
}

func flattenLocation(name string) string {
	tokens := strings.Split(name, "/")
	if len(tokens) < 6 || tokens[2] != "locations" {
		return ""
	}
	return tokens[3]
}

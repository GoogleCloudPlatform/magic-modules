package gotemplate

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// actionRegex extracts template actions: everything between {{ and }}
var actionRegex = regexp.MustCompile(`\{\{-?\s*(.*?)\s*-?\}\}`)

// identifierRegex matches a standalone identifier (function name) at the start of a pipeline segment
var identifierRegex = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)`)

// ValidFuncs is the registry of allowed template functions and keywords.
var ValidFuncs = map[string]bool{
	// Go built-in template functions
	"and": true, "call": true, "eq": true, "ge": true, "gt": true, "html": true,
	"index": true, "js": true, "le": true, "len": true, "lt": true, "ne": true,
	"not": true, "or": true, "print": true, "printf": true, "println": true,
	"slice": true, "urlquery": true,

	// Go template keywords
	"if": true, "else": true, "end": true, "range": true, "with": true,
	"block": true, "define": true, "template": true, "nil": true,

	// mmv1 registered functions (google/template_utils.go)
	"title": true, "replace": true, "replaceAll": true, "camelize": true,
	"underscore": true, "plural": true, "contains": true, "join": true,
	"lower": true, "upper": true, "hasSuffix": true, "dict": true,
	"format2regex": true, "hasPrefix": true, "sub": true, "plus": true,
	"firstSentence": true, "trimTemplate": true, "customTemplate": true,

	// mmv1 registered functions (provider/template_data.go)
	"TemplatePath": true,
}

// CheckInvalidFuncsForFile scans a file for invalid template function calls.
func CheckInvalidFuncsForFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		actions := actionRegex.FindAllStringSubmatch(line, -1)
		for _, action := range actions {
			body := strings.TrimSpace(action[1])
			if body == "" {
				continue
			}

			segments := strings.Split(body, "|")
			for _, seg := range segments {
				seg = strings.TrimSpace(seg)
				if seg == "" || strings.HasPrefix(seg, ".") || strings.HasPrefix(seg, "$") || strings.HasPrefix(seg, "else if") {
					continue
				}

				match := identifierRegex.FindStringSubmatch(seg)
				if match == nil {
					continue
				}

				funcName := match[1]
				if funcName == "true" || funcName == "false" || ValidFuncs[funcName] {
					continue
				}

				// Skip bare lowercase identifiers (likely Terraform interpolation)
				if seg == funcName && !hasUpperCase(funcName) {
					continue
				}

				// Skip if it appears in a string literal
				if appearsInStringLiteral(body, funcName) {
					continue
				}

				results = append(results, fmt.Sprintf("unknown function %q in action {{%s}} (line %d)", funcName, body, lineNum))
			}
		}
	}
	return results, scanner.Err()
}

func hasUpperCase(s string) bool {
	for _, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			return true
		}
	}
	return false
}

func appearsInStringLiteral(body string, identifier string) bool {
	inString, escaped := false, false
	var current strings.Builder
	for _, ch := range body {
		if escaped {
			if inString { current.WriteRune(ch) }
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			if inString { current.WriteRune(ch) }
			continue
		}
		if ch == '"' {
			if inString && strings.Contains(current.String(), identifier) { return true }
			current.Reset()
			inString = !inString
			continue
		}
		if inString { current.WriteRune(ch) }
	}
	return false
}

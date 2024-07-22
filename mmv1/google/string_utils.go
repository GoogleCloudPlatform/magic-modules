// Copyright 2024 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
)

// // Helper class to process and mutate strings.
// class StringUtils
// Converts string from camel case to underscore
func Underscore(source string) string {
	tmp := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`).ReplaceAllString(source, "${1}_${2}")
	tmp = regexp.MustCompile(`([a-z\d])([A-Z])`).ReplaceAllString(tmp, "${1}_${2}")
	tmp = strings.Replace(tmp, "-", "_", 1)
	tmp = strings.Replace(tmp, ".", "_", 1)
	tmp = strings.ToLower(tmp)
	return tmp
}

// Converts from PascalCase to Space Separated
// For example, converts "AccessApproval" to "Access approval"
func SpaceSeparated(source string) string {
	tmp := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`).ReplaceAllString(source, "${1} ${2}")
	tmp = regexp.MustCompile(`([a-z\d])([A-Z])`).ReplaceAllString(tmp, "${1} ${2}")
	tmp = strings.ToLower(tmp)

	// Capitalize the first letter
	if len(tmp) != 0 {
		r := []rune(tmp)
		r[0] = unicode.ToUpper(r[0])
		tmp = string(r)
	}
	return tmp
}

// // Converts a string to space-separated capitalized words
// def self.title(source)
func SpaceSeparatedTitle(source string) string {
	ss := SpaceSeparated(source)
	return strings.Title(ss)
}

// Returns all the characters up until the period (.) or returns text
// unchanged if there is no period.
//
//	def self.first_sentence(text)
func FirstSentence(text string) string {
	re := regexp.MustCompile(`[.?!]`)
	periodPos := re.FindStringIndex(text)
	if periodPos == nil {
		return text
	}

	return text[:periodPos[0]+1]
}

// Returns the plural form of a word
func Plural(source string) string {
	// policies -> policies
	// indices -> indices
	if strings.HasSuffix(source, "ies") || strings.HasSuffix(source, "es") {
		return source
	}

	// index -> indices
	if strings.HasSuffix(source, "ex") {
		re := regexp.MustCompile("ex$")
		result := re.ReplaceAllString(source, "")
		return fmt.Sprintf("%sices", result)
	}

	// mesh -> meshes
	if strings.HasSuffix(source, "esh") {
		return fmt.Sprintf("%ses", source)
	}

	// key -> keys
	// gateway -> gateways
	if strings.HasSuffix(source, "ey") || strings.HasSuffix(source, "ay") {
		return fmt.Sprintf("%ss", source)
	}

	// policy -> policies
	if strings.HasSuffix(source, "y") {
		re := regexp.MustCompile("y$")
		result := re.ReplaceAllString(source, "")
		return fmt.Sprintf("%sies", result)
	}

	return fmt.Sprintf("%ss", source)
}

func Camelize(term string, firstLetter string) string {
	if firstLetter != "upper" && firstLetter != "lower" {
		log.Fatalf("Invalid option, use either upper or lower")
	}

	res := term
	if firstLetter == "upper" {
		res = regexp.MustCompile(`^[a-z\d]*`).ReplaceAllStringFunc(res, func(match string) string {
			return strings.Title(match)
		})
	} else {
		// TODO: rewrite with the regular expression. Lookahead(?=) is not supported in Go
		// 	acronymsCamelizeRegex := regexp.MustCompile(`^(?:(?=a)b(?=\b|[A-Z_])|\w)`)
		// 	res = acronymsCamelizeRegex.ReplaceAllStringFunc(res, func(match string) string {
		// 		return strings.ToLower(match)
		// 	})
		if len(res) != 0 {
			r := []rune(res)
			r[0] = unicode.ToLower(r[0])
			res = string(r)
		}
	}
	// handle snake case
	re := regexp.MustCompile(`(?:_)([a-z\d]*)`)
	res = re.ReplaceAllStringFunc(res, func(match string) string {
		word := match[1:]
		word = strings.Title(word)
		return word
	})
	return res
}

/*
Transforms a format string with field markers to a regex string with capture groups.
For instance,

	projects/{{project}}/global/networks/{{name}}

is transformed to

	projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)

Values marked with % are URL-encoded, and will match any number of /'s.
Note: ?P indicates a Python-compatible named capture group. Named groups
aren't common in JS-based regex flavours, but are in Perl-based ones
*/
func Format2Regex(format string) string {
	re := regexp.MustCompile(`\{\{%([[:word:]]+)\}\}`)
	result := re.ReplaceAllStringFunc(format, func(match string) string {
		// TODO: the trims may not be needed with more effecient regex
		word := strings.TrimPrefix(match, "{{")
		word = strings.TrimSuffix(word, "}}")
		word = strings.ReplaceAll(word, "%", "")
		return fmt.Sprintf("(?P<%s>.+)", word)
	})
	re = regexp.MustCompile(`\{\{([[:word:]]+)\}\}`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		word := strings.TrimPrefix(match, "{{")
		word = strings.TrimSuffix(word, "}}")
		return fmt.Sprintf("(?P<%s>[^/]+)", word)
	})
	return result
}

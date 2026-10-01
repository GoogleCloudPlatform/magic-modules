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

// Regular expressions used by the string helpers below. These are compiled once
// at package init: several of these helpers (notably Underscore) are called
// millions of times during a full generation run, and compiling a regex on
// every call dominated the generator's CPU profile.
var (
	acronymBoundaryRegexp = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	wordBoundaryRegexp    = regexp.MustCompile(`([a-z\d])([A-Z])`)
	firstSentenceRegexp   = regexp.MustCompile(`[.?!]`)
	camelizeFirstRegexp   = regexp.MustCompile(`^[a-z\d]*`)
	camelizeSnakeRegexp   = regexp.MustCompile(`(?:_)([a-z\d]*)`)
	format2RegexPercent   = regexp.MustCompile(`\{\{%([[:word:]]+)\}\}`)
	format2RegexPlain     = regexp.MustCompile(`\{\{([[:word:]]+)\}\}`)
)

// // Helper class to process and mutate strings.
// class StringUtils
// Converts string from camel case to underscore
func Underscore(source string) string {
	tmp := acronymBoundaryRegexp.ReplaceAllString(source, "${1}_${2}")
	tmp = wordBoundaryRegexp.ReplaceAllString(tmp, "${1}_${2}")
	tmp = strings.ReplaceAll(tmp, "-", "_")
	tmp = strings.ReplaceAll(tmp, ".", "_")
	tmp = strings.ToLower(tmp)
	return tmp
}

// Converts from PascalCase to Space Separated
// For example, converts "AccessApproval" to "Access approval"
func SpaceSeparated(source string) string {
	tmp := acronymBoundaryRegexp.ReplaceAllString(source, "${1} ${2}")
	tmp = wordBoundaryRegexp.ReplaceAllString(tmp, "${1} ${2}")
	tmp = strings.ToLower(tmp)

	// Capitalize the first letter
	if len(tmp) != 0 {
		r := []rune(tmp)
		r[0] = unicode.ToUpper(r[0])
		tmp = string(r)
	}
	return tmp
}

// Converts a string to space-separated capitalized words
func SpaceSeparatedTitle(source string) string {
	ss := SpaceSeparated(source)
	return strings.Title(ss)
}

// Returns all the characters up until the period (.) or returns text
// unchanged if there is no period.
func FirstSentence(text string) string {
	periodPos := firstSentenceRegexp.FindStringIndex(text)
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
		return fmt.Sprintf("%sices", strings.TrimSuffix(source, "ex"))
	}

	// mesh -> meshes
	// messageBus -> messageBuses
	if strings.HasSuffix(source, "esh") || strings.HasSuffix(source, "s") {
		return fmt.Sprintf("%ses", source)
	}

	// key -> keys
	// gateway -> gateways
	if strings.HasSuffix(source, "ey") || strings.HasSuffix(source, "ay") {
		return fmt.Sprintf("%ss", source)
	}

	// policy -> policies
	if strings.HasSuffix(source, "y") {
		return fmt.Sprintf("%sies", strings.TrimSuffix(source, "y"))
	}

	return fmt.Sprintf("%ss", source)
}

func Camelize(term string, firstLetter string) string {
	if firstLetter != "upper" && firstLetter != "lower" {
		log.Fatalf("Invalid option, use either upper or lower")
	}

	res := term
	if firstLetter == "upper" {
		res = camelizeFirstRegexp.ReplaceAllStringFunc(res, func(match string) string {
			return strings.Title(match)
		})
	} else {
		if len(res) != 0 {
			r := []rune(res)
			r[0] = unicode.ToLower(r[0])
			res = string(r)
		}
	}
	// handle snake case
	res = camelizeSnakeRegexp.ReplaceAllStringFunc(res, func(match string) string {
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
	result := format2RegexPercent.ReplaceAllStringFunc(format, func(match string) string {
		// TODO rewrite: the trims may not be needed with more effecient regex
		word := strings.TrimPrefix(match, "{{")
		word = strings.TrimSuffix(word, "}}")
		word = strings.ReplaceAll(word, "%", "")
		return fmt.Sprintf("(?P<%s>.+)", word)
	})
	result = format2RegexPlain.ReplaceAllStringFunc(result, func(match string) string {
		word := strings.TrimPrefix(match, "{{")
		word = strings.TrimSuffix(word, "}}")
		return fmt.Sprintf("(?P<%s>[^/]+)", word)
	})
	return result
}

// Copyright 2021 Google LLC. All Rights Reserved.
//
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

package utils

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
func SnakeToTitleCase(s string) string {
	return strings.Join(SnakeToTitleParts(s), "")
}

// snakeToTitleParts returns the parts of a snake_case string titlecased as an
// array, taking into account common initialisms.
func SnakeToTitleParts(s string) []string {
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
func JsonToSnakeCase(s string) string {
	for _, v := range initialisms {
		s = strings.ReplaceAll(s, v, v[0:1]+strings.ToLower(v[1:]))
	}
	result := regexp.MustCompile("(.)([A-Z][^A-Z]+)").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(result, "${1}_${2}"))
}

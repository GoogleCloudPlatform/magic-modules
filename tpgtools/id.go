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

package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/nasa9084/go-openapi"
)

const PatternPart = "{{(\\w+)}}"

func idParts(id string) (parts []string) {
	r := regexp.MustCompile(PatternPart)

	// returns [["{{field}}", "field"] ...]
	idTmplAndParts := r.FindAllStringSubmatch(id, -1)
	for _, v := range idTmplAndParts {
		parts = append(parts, v[1])
	}

	return parts
}

// PatternToRegex formats a pattern string into a Python-compatible regex.
func PatternToRegex(s string) string {
	re := regexp.MustCompile(PatternPart)
	return re.ReplaceAllString(s, "(?P<$1>[^/]+)")
}

// Finds the correct resource id based on the schema and any overrides
func findResourceId(schema *openapi.Schema, overrides Overrides, location string) (string, error) {
	id, ok := schema.Extension["x-dcl-id"].(string)
	if !ok {
		return "", fmt.Errorf("Malformed or missing x-dcl-id: %v", schema.Extension["x-dcl-id"])
	}

	// Resource Override: Custom ID
	cid := CustomIDDetails{}
	cidOk, err := overrides.ResourceOverrideWithDetails(CustomID, &cid, location)
	if err != nil {
		return "", fmt.Errorf("failed to decode custom id details: %v", err)
	}

	if cidOk {
		id = cid.ID
	}

	for _, override := range overrides {
		if override.Type == CustomName {
			if strings.Contains(id, fmt.Sprintf("{{%s}}", *override.Field)) {
				id = strings.Replace(id, fmt.Sprintf("{{%s}}", *override.Field), fmt.Sprintf("{{%s}}", override.Details.(map[interface{}]interface{})["name"].(string)), 1)
			}
		}
	}
	return id, nil
}

// Finds all import formats for a given id. This can include short forms and
// partial forms with inferred project/region/etc
func defaultImportFormats(id string) (formats []string) {
	uniqueFormats := stringset.New()

	uniqueFormats.Add(id)

	parts := idParts(id)
	for i, v := range parts {
		parts[i] = fmt.Sprintf("{{%s}}", v)
	}

	// short form "{{project}}/{{region}}/{{name}}"
	uniqueFormats[strings.Join(parts, "/")] = struct{}{}

	// short form sans project
	var locationalParts []string
	for _, v := range parts {
		if v != "{{project}}" {
			locationalParts = append(locationalParts, v)
		}
	}
	if len(locationalParts) != 0 {
		uniqueFormats.Add(strings.Join(locationalParts, "/"))
	}

	// short form sans project, region, zone
	var resourceParts []string
	for _, v := range locationalParts {
		if v != "{{zone}}" && v != "{{region}}" {
			resourceParts = append(resourceParts, v)
		}
	}
	if len(resourceParts) != 0 {
		uniqueFormats.Add(strings.Join(resourceParts, "/"))
	}

	for _, f := range uniqueFormats.Elements() {
		formats = append(formats, f)
	}

	// formats must be ordered most to least specific
	sort.SliceStable(formats, formatComparator(formats))
	return formats
}

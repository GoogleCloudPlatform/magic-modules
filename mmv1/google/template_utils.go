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
	"errors"
	"strings"

	"text/template"
)

// Build a map(map[string]interface{}) from a list of paramerter
// The format of passed in parmeters are key1, value1, key2, value2 ...
func wrapMultipleParams(params ...interface{}) (map[string]interface{}, error) {
	if len(params)%2 != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	m := make(map[string]interface{}, len(params)/2)
	for i := 0; i < len(params); i += 2 {
		key, ok := params[i].(string)
		if !ok {
			return nil, errors.New("keys must be strings")
		}
		m[key] = params[i+1]
	}
	return m, nil
}

// subtract returns the difference between a and b
// and used in Go templates
func subtract(a, b int) int {
	return a - b
}

var TemplateFunctions = template.FuncMap{
	"title":        SpaceSeparatedTitle,
	"replace":      strings.Replace,
	"replaceAll":   strings.ReplaceAll,
	"camelize":     Camelize,
	"underscore":   Underscore,
	"plural":       Plural,
	"contains":     strings.Contains,
	"join":         strings.Join,
	"lower":        strings.ToLower,
	"upper":        strings.ToUpper,
	"dict":         wrapMultipleParams,
	"format2regex": Format2Regex,
	"hasPrefix":    strings.HasPrefix,
	"sub":          subtract,
}

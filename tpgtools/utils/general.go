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
	"strings"

	"github.com/kylelemons/godebug/pretty"
)

// Sort id formats based on the order they should be matched. This is
// most specific first, so {{project}}/{{region}}/{{name}} would be applied
// before {{region}}/{{name}}
func FormatComparator(formats []string) func(i, j int) bool {
	return func(i, j int) bool {
		l := formats[i]
		r := formats[j]

		lBrace := strings.Count(l, "{{")
		rBrace := strings.Count(r, "{{")

		lSlash := strings.Count(l, "/")
		rSlash := strings.Count(r, "/")

		if lBrace == rBrace {
			return lSlash > rSlash // > and not <, we want more to appear first
		}

		return lBrace > rBrace
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SprintResource(v interface{}) string {
	prettyConfig := &pretty.Config{
		Diffable: true,
	}
	return prettyConfig.Sprint(v)
}

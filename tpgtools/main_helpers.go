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
	"bytes"
	"fmt"
	"go/format"
	"strconv"
	"strings"

	"github.com/kylelemons/godebug/pretty"
)

// Sort id formats based on the order they should be matched. This is
// most specific first, so {{project}}/{{region}}/{{name}} would be applied
// before {{region}}/{{name}}
func formatComparator(formats []string) func(i, j int) bool {
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

func escapeDescription(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(in, `\`, `\\`), `"`, `\"`), "\n", `\n`)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stripExt(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n >= 0 {
		return s[:n]
	}
	return s
}

func sprintResource(v interface{}) string {
	prettyConfig := &pretty.Config{
		Diffable: true,
	}
	return prettyConfig.Sprint(v)
}

func formatSource(source *bytes.Buffer) ([]byte, error) {
	sourceByte := source.Bytes()
	// Replace import path based on version (beta/alpha)
	if terraformResourceDirectory != "google" {
		sourceByte = bytes.Replace(sourceByte, []byte("github.com/hashicorp/terraform-provider-google/google"), []byte(terraformProviderModule+"/"+terraformResourceDirectory), -1)
	}

	output, err := format.Source(sourceByte)
	if err != nil {
		return []byte(source.String()), err
	}

	return output, nil
}

func renderDefault(t Type, val string) (string, error) {
	switch t.String() {
	case SchemaTypeBool:
		if b, err := strconv.ParseBool(val); err == nil {
			return fmt.Sprintf("%v", b), nil
		} else {
			return "", fmt.Errorf("Failed to render default for boolean: %s", val)
		}
	case SchemaTypeFloat:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return fmt.Sprintf("%f", f), nil
		} else {
			return "", fmt.Errorf("Failed to render default for float: %s", val)
		}
	case SchemaTypeInt:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return fmt.Sprintf("%d", i), nil
		} else {
			return "", fmt.Errorf("Failed to render default for int: %s", val)
		}
	case SchemaTypeString:
		return fmt.Sprintf("%q", val), nil
	}
	return "", fmt.Errorf("Failed to find default format for type: %v", t)
}

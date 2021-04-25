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

package types

import (
	"fmt"
	"strconv"
)

// Sorts properties to be in a standard order
func propComparator(props []Property) func(i, j int) bool {
	return func(i, j int) bool {
		l := props[i]
		r := props[j]

		// required < non-required
		if l.Required && !r.Required {
			return true
		}

		// conversely, non-required > required
		if r.Required && !l.Required {
			return false
		}

		// same deal- settable (optional / O+C) fields > Computed fields
		if l.Settable && !r.Settable {
			return true
		}
		if r.Settable && !l.Settable {
			return false
		}

		// finally, sort by name
		return l.Name() < r.Name()
	}
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

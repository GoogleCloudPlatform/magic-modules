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

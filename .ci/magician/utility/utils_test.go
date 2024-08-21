/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package utility

import (
	"reflect"
	"testing"
)

func TestRemovesList(t *testing.T) {
	cases := map[string]struct {
		Original, Removal, Expected []string
	}{
		"Remove list": {
			Original: []string{"a", "b", "c"},
			Removal:  []string{"b"},
			Expected: []string{"a", "c"},
		},
		"Remove case sensitive elements": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"b", "c", "A"},
			Expected: []string{"a", "B"},
		},
		"Remove nonexistent elements": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"a", "A", "d"},
			Expected: []string{"b", "c", "B"},
		},
		"Remove none": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{},
			Expected: []string{"a", "b", "c", "A", "B"},
		},
		"Remove all": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"a", "b", "c", "A", "B"},
			Expected: []string{},
		},
		"Remove all and extra nonexistent elements": {
			Original: []string{"a", "b", "c", "A", "B"},
			Removal:  []string{"a", "b", "c", "A", "B", "D"},
			Expected: []string{},
		},
	}
	for tn, tc := range cases {
		result := Removes(tc.Original, tc.Removal)
		if !reflect.DeepEqual(result, tc.Expected) {
			t.Errorf("bad: %s, '%s' removes '%s' expect result: %s, but got: %s", tn, tc.Original, tc.Removal, tc.Expected, result)
		}
	}
}

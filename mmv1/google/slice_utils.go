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

// Returns a new slice containing all of the elements
// for which the test function returns true in the original slice
func Select[T any](S []T, test func(T) bool) (ret []T) {
	for _, s := range S {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// Returns a new slice containing all of the elements
// for which the test function returns false in the original slice
func Reject[T any](S []T, test func(T) bool) (ret []T) {
	for _, s := range S {
		if !test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// Concat two slices
func Concat[T any](S1 []T, S2 []T) (ret []T) {
	return append(S1, S2...)
}

// difference returns the elements in `S1` that aren't in `S2`.
func Diff(S1, S2 []string) []string {
	var ret []string
	mb := make(map[string]bool, len(S2))
	for _, x := range S2 {
		mb[x] = true
	}

	for _, x := range S1 {
		if _, found := mb[x]; !found {
			ret = append(ret, x)
		}
	}
	return ret
}

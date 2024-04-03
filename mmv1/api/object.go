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

package api

// require 'google/extensions'
// require 'google/logger'
// require 'google/yaml_validator'

// Represents an object that has a (mandatory) name
type NamedObject struct {
	// google.YamlValidator

	Name string

	// original value of :name before the provider override happens
	// same as :name if not overridden in provider
	ApiName string `yaml:"api_name"`
}

// func (n *Named) string_array(arr) {
// types = arr.map(&:class).uniq
// types.size == 1 && types[0] == String
// }

// func (n *Named) deep_merge(arr1, arr2) {
// // Scopes is an array of standard strings. In which case return the
// // version in the overrides. This allows scopes to be removed rather
// // than allowing for a merge of the two arrays
// if string_array?(arr1)
//   return arr2.nil? ? arr1 : arr2
// end

// // Merge any elements that exist in both
// result = arr1.map do |el1|
//   other = arr2.select { |el2| el1.name == el2.name }.first
//   other.nil? ? el1 : el1.merge(other)
// end

// // Add any elements of arr2 that don't exist in arr1
// result + arr2.reject do |el2|
//   arr1.any? { |el1| el2.name == el1.name }
// end
// }

// func (n *Named) merge(other) {
// result = self.class.new
// instance_variables.each do |v|
//   result.instance_variable_set(v, instance_variable_get(v))
// end

// other.instance_variables.each do |v|
//   if other.instance_variable_get(v).instance_of?(Array)
//     result.instance_variable_set(v, deep_merge(result.instance_variable_get(v),
//                                                 other.instance_variable_get(v)))
//   else
//     result.instance_variable_set(v, other.instance_variable_get(v))
//   end
// end

// result
// }

// func (n *Named) validate() {
// super
// check :name, type: String, required: true
// check :api_name, type: String, default: @name
// }

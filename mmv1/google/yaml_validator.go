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
	"log"

	"gopkg.in/yaml.v2"
)

// A helper class to validate contents coming from YAML files.
type YamlValidator struct{}

func (v *YamlValidator) Parse(content []byte, obj interface{}, yamlPath string) {
	// TODO(nelsonjr): Allow specifying which symbols to restrict it further.
	// But it requires inspecting all configuration files for symbol sources,
	// such as Enum values. Leaving it as a nice-to-have for the future.
	if err := yaml.UnmarshalStrict(content, obj); err != nil {
		log.Fatalf("Cannot unmarshal data from file %s: %v", yamlPath, err)
	}
}

// func (v *YamlValidator) allowed_classes() {
// ObjectSpace.each_object(Class).select do |klass|
//   klass < Google::YamlValidator
// end.push(Time, Symbol)
// }

// func (v *YamlValidator) validate() {
// Google::LOGGER.debug "Validating //{self.class} '//{@name}'"
// check_extraneous_properties
// }

// func (v *YamlValidator) set_variable(value, property) {
// Google::LOGGER.debug "Setting variable of //{value} to //{self}"
// instance_variable_set("@//{property}", value)
// }

// Does all validation checking for a particular variable.
// options:
// :default   - the default value for this variable if its nil
// :type      - the allowed types (single or array) that this value can be
// :item_type - the allowed types that all values in this array should be
//              (implied that type == array)
// :allowed   - the allowed values that this non-array variable should be.
// :required  - is the variable required? (defaults: false)
// func (v *YamlValidator) check(variable, **opts) {
// value = instance_variable_get("@//{variable}")

// // Set default value.
// if !opts[:default].nil? && value.nil?
//   instance_variable_set("@//{variable}", opts[:default])
//   value = instance_variable_get("@//{variable}")
// end

// // Check if value is required. Print nested path if available.
// lineage_path = respond_to?('lineage') ? lineage : ''
// raise "//{lineage_path} > Missing '//{variable}'" if value.nil? && opts[:required]
// return if value.nil?

// // Check type
// check_property_value(variable, value, opts[:type]) if opts[:type]

// // Check item_type
// if value.is_a?(Array)
//   raise "//{lineage_path} > //{variable} must have item_type on arrays" unless opts[:item_type]

//   value.each_with_index do |o, index|
//     check_property_value("//{variable}[//{index}]", o, opts[:item_type])
//   end
// end

// // Check if value is allowed
// return unless opts[:allowed]
// raise "//{value} on //{variable} should be one of //{opts[:allowed]}" \
//   unless opts[:allowed].include?(value)
// }

// func (v *YamlValidator) conflicts(list) {
// value_checked = false
// list.each do |item|
//   next if instance_variable_get("@//{item}").nil?
//   raise "//{list.join(',')} cannot be set at the same time" if value_checked

//   value_checked = true
// end
// }

// private

// func (v *YamlValidator) check_type(name, object, type) {
// if type == :boolean
//   return unless [TrueClass, FalseClass].find_index(object.class).nil?
// elsif type.is_a? ::Array
//   return if type.find_index(:boolean) && [TrueClass, FalseClass].find_index(object.class)
//   return unless type.find_index(object.class).nil?
// // check if class is or inherits from type
// elsif object.class <= type
//   return
// end
// raise "Property '//{name}' is '//{object.class}' instead of '//{type}'"
// }

// func (v *YamlValidator) log_check_type(object) {
// if object.respond_to?(:name)
//   Google::LOGGER.debug "Checking object //{object.name}"
// else
//   Google::LOGGER.debug "Checking object //{object}"
// end
// }

// func (v *YamlValidator) check_property_value(property, prop_value, type) {
// Google::LOGGER.debug "Checking '//{property}' on //{object_display_name}"
// check_type property, prop_value, type unless type.nil?
// prop_value.validate if prop_value.is_a?(Api::Object)
// }

// func (v *YamlValidator) check_extraneous_properties() {
// instance_variables.each do |variable|
//   var_name = variable.id2name[1..]
//   next if var_name.start_with?('__')

//   Google::LOGGER.debug "Validating '//{var_name}' on //{object_display_name}"
//   raise "Extraneous variable '//{var_name}' in //{object_display_name}" \
//     unless methods.include?(var_name.intern)
// end
// }

// func (v *YamlValidator) set_variables(objects, property) {
// return if objects.nil?

// objects.each do |object|
//   object.set_variable(self, property) if object.respond_to?(:set_variable)
// end
// }

// func (v *YamlValidator) ensure_property_does_not_exist(property) {
// raise "Conflict of property '//{property}' for object '//{self}'" \
//   unless instance_variable_get("@//{property}").nil?
// }

// func (v *YamlValidator) object_display_name() {
// "//{@name}<//{self.class.name}>"
// }

# Copyright 2018 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'tools/linter/builder/discovery'
require 'tools/linter/builder/docs'
require 'tools/linter/builder/override'
require 'api/compiler'

# Override Builder
# ruby override_builder.rb original_api.yaml discovery_api.yaml output.yaml
# This script takes in an original version of api.yaml, a discovery doc generated version
# and creates a series of overrides in output.yaml that override values from api.yaml to discovery_api.yaml

module Api
  class Object
    # Create a setter if the setter doesn't exist
    # Yes, this isn't pretty and I apologize
    def method_missing(method_name, *args)
      matches = /([a-z_]*)=/.match(method_name)
      super unless matches
      create_setter(matches[1])
      method(method_name.to_sym).call(*args)
    end

    def create_setter(variable)
      self.class.define_method("#{variable}=") { |val| instance_variable_set("@#{variable}", val) }
    end

    def validate
    end
  end
end

module Google
  class YamlValidator
    # Create a setter if the setter doesn't exist
    # Yes, this isn't pretty and I apologize
    def method_missing(method_name, *args)
      matches = /([a-z_]*)=/.match(method_name)
      super unless matches
      create_setter(matches[1])
      method(method_name.to_sym).call(*args)
    end

    def create_setter(variable)
      self.class.define_method("#{variable}=") { |val| instance_variable_set("@#{variable}", val) }
    end

    def validate
    end
  end
end

# Takes a series of properties, does diffs between them and returns a list.
def diff_properties(new_api_props, old_api_props, prefix='')
  old_api_props.map do |old_api_prop|
    all_props = []
    override = Provider::Overrides::PropertyOverride.new
    new_api_prop = new_api_props.select { |p| p.name == old_api_prop.name }.first
    # Compare the new prop values to the old prop values.
    old_api_prop.instance_variables.reject { |o| o.to_s.include?('properties') }
                                   .each do |var|
      if !values_equal(old_api_prop.instance_variable_get(var), new_api_prop.instance_variable_get(var))
        override.instance_variable_set(var, old_api_prop.instance_variable_get(var))
      end
    end

    if old_api_prop.is_a?(Api::Type::NestedObject)
      all_props.append(diff_properties(new_api_prop.properties, old_api_prop.properties, name(prefix, old_api_prop, '.')))
    # I'm not convinced that overriding Arrays of NestedObjects doesn't work.
    #elsif old_api_prop.is_a?(Api::Type::Array) && old_api_prop.item_type.is_a?(Api::Type::NestedObject)
    #  all_props.append(diff_properties(new_api_prop.item_type.properties, old_api_prop.item_type.properties,
    #                                   "#{prefix}#{old_api_prop.name}.item_type."))
    end

    all_props.append({ name(prefix, old_api_prop) => override}) if override.instance_variables.length > 0
    all_props
  end.compact
end

def name(prefix, api_prop, suffix='')
  if api_prop.api_name
    "#{prefix}#{api_prop.api_name}#{suffix}"
  else
    "#{prefix}#{api_prop.name}#{suffix}"
  end
end

# Check if two values are equal.
def values_equal(old, new)
  # nil == false
  return true if old.nil? && new == false
  # If it's in the old version, but not the new, we need it!
  return false if new.nil? && old
  # If it's an array, make sure that these things are the same.
  return old.sort == new.sort if old.is_a?(::Array)
  # If it's a normal thing, just check if it's equal.
  return old == new
end
raise "Must have 3 files" if ARGV.length < 3

original_api_file = ARGV[0]
new_api_file = ARGV[1]
output_file = ARGV[2]

original_api = Api::Compiler.new(original_api_file).run
new_api = Api::Compiler.new(new_api_file).run

overrides = Provider::Overrides::ResourceOverrides.new
original_api.objects.each do |obj|
  # Grab all of the first level things and do a diff
  override = Provider::Overrides::ResourceOverride.new
  new_api_obj = new_api.objects.select { |o| o.name == obj.name }.first
  next unless new_api_obj
  obj.instance_variables.reject { |o| o.to_s.include?('properties') || o.to_s.include?('parameters') }
                        .each do |var|
    if obj.instance_variable_get(var) != new_api_obj.instance_variable_get(var)
      override.instance_variable_set(var, obj.instance_variable_get(var))
    end
  end
  properties = diff_properties(new_api_obj.properties, obj.properties)
  override.properties = properties.flatten.reduce({}, :merge) if properties
  parameters = diff_properties(new_api_obj.parameters, obj.parameters)
  override.parameters = parameters.flatten.reduce({}, :merge) if parameters
  if override.instance_variables.length > 0
    overrides.instance_variable_set("@#{obj.name}", override)
  end
end

File.write(output_file, YAML::dump(overrides))


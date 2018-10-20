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

# Takes a series of properties, does diffs between them and returns a list.
def diff_properties(new_api_props, old_api_props, prefix='')
  old_api_props.map do |old_api_prop|
    all_props = []
    override = DiscoveryOverride::PropertyOverride.new
    new_api_prop = new_api_props.select { |p| p.name == old_api_prop.name }.first
    # Compare the new prop values to the old prop values.
    old_api_prop.instance_variables.reject { |o| o.to_s.include?('properties') }
                                   .each do |var|
      if !values_equal(old_api_prop.instance_variable_get(var), new_api_prop.instance_variable_get(var))
        override.instance_variable_set(var, old_api_prop.instance_variable_get(var))
      end
    end

    if old_api_prop.is_a?(Api::Type::NestedObject)
      require 'byebug'
      byebug if !new_api_prop.is_a?(Api::Type::NestedObject)
      all_props.append(diff_properties(new_api_prop.properties, old_api_prop.properties, "#{old_api_prop.name}."))
    end

    all_props.append({"#{prefix}#{old_api_prop.name}" => override}) if override.instance_variables.length > 0
    all_props
  end.compact
end

def values_equal(old, new)
  return true if old.nil? && new == false
  return old.sort == new.sort if old.is_a?(::Array)
  return old == new
end
raise "Must have 3 files" if ARGV.length < 3

original_api_file = ARGV[0]
new_api_file = ARGV[1]
output_file = ARGV[2]

original_api = Api::Compiler.new(original_api_file).run
new_api = Api::Compiler.new(new_api_file).run

overrides = Provider::ResourceOverrides.new
original_api.objects.each do |obj|
  # Grab all of the first level things and do a diff
  override = DiscoveryOverride::ResourceOverride.new
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
  if override.instance_variables.length > 0
    overrides.instance_variable_set("@#{obj.name}", override)
  end
end

File.write(output_file, YAML::dump(overrides))


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
raise "Must have 3 files" if ARGV.length < 3

original_api_file = ARGV[0]
new_api_file = ARGV[1]
output_file = ARGV[2]

original_api = Api::Compiler.new(original_api_file).run
new_api = Api::Compiler.new(new_api_file).run

overrides = Provider::ResourceOverrides.new
original_api.objects.each do |obj|
  override = DiscoveryOverride::ResourceOverride.new
  new_api_obj = new_api.objects.select { |o| o.name == obj.name }.first
  obj.instance_variables.reject { |o| o.to_s.include?('properties') || o.to_s.include?('parameters') }
                        .each do |var|
    if obj.instance_variable_get(var) != new_api_obj.instance_variable_get(var)
      override.instance_variable_set(var, obj.instance_variable_get(var))
    end
  end
  if override.instance_variables.length > 0
    overrides.instance_variable_set("@#{obj.name}", override)
  end
end

File.write(output_file, YAML::dump(overrides))

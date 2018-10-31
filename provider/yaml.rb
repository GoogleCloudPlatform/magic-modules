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

require 'provider/config'
require 'provider/core'
require 'provider/yaml/config'
require 'provider/yaml/property_override'
require 'provider/yaml/resource_override'

module Provider
  class Yaml < Provider::Core

    private

    def generate_resource(data)
      target_folder = data[:output_folder]
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      product_name = data[:product_name].underscore
      filepath = File.join(target_folder, "#{product_name}_#{name}.yaml")
      File.write(filepath, data[:object].to_yaml)
    end

    # rubocop:disable Layout/EmptyLineBetweenDefs
    def generate_resource_tests(data) end
    def generate_network_datas(data, object) end
    def generate_base_property(data) end
    def generate_simple_property(type, data) end
    def generate_typed_array(type, data) end
    def emit_nested_object(data) end
    def emit_resourceref_object(data) end
    # rubocop:enable Layout/EmptyLineBetweenDefs
  end
end

# Copyright 2017 Google Inc.
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

require 'google/python_utils'
require 'provider/azure/ansible/module_extension'

module Provider
  module Ansible
    # Responsible for building out the AnsibleModule code.
    # AnsibleModule is responsible for input validation.
    module Module
      include Google::PythonUtils
      include Provider::Azure::Ansible::ModuleExtension
      # Returns an array of all base options for a given property.
      def ansible_module(properties)
        properties.reject(&:output)
                  .map { |x| python_dict_for_property(x) }
                  .reduce({}, :merge)
      end

      def python_dict_for_property(prop)
        {
          prop.name.underscore => {
            'required' => (true if prop.required && !prop.default_value),
            'default' => prop.default_value,
            'type' => python_type(prop),
            'elements' => (python_type(prop.item_type) \
              if prop.is_a?(Api::Type::Array) && python_type(prop.item_type)),
            'aliases' => prop.aliases,
            'options' => (if prop.nested_properties?
                            prop.nested_properties.reject(&:output)
                                                  .map { |x| python_dict_for_property(x) }
                                                  .reduce({}, :merge)
                          end
                         )
          }.reject { |_, v| v.nil? }
        }
      end

      # GcpModule is acting as a dictionary and doesn't need the dict() notation on
      # the first level.
      def remove_outside_dict(contents)
        contents.sub('dict(', '').chomp(')')
      end
    end
  end
end

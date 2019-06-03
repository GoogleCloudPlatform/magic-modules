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

module Provider
  module Ansible
    # Responsible for building out the AnsibleModule code.
    # AnsibleModule is responsible for input validation.
    module Module
      include Google::PythonUtils
<<<<<<< HEAD
      # Returns the Python dictionary representing a simple property for
      # validation.
      def python_dict_for_property(prop, object, spaces = 0)
        if prop.is_a?(Api::Type::Array) && \
           prop.item_type.is_a?(Api::Type::NestedObject)
          nested_obj_dict(prop, object, prop.item_type.properties, spaces)
        elsif prop.is_a? Api::Type::NestedObject
          nested_obj_dict(prop, object, prop.properties, spaces)
        else
          name = python_variable_name(prop, object.azure_sdk_definition.create)
          options = prop_options(prop, object, spaces).join("\n")
          "#{name}=dict(\n#{indent_list(options, 4)}\n)"
        end
      end

      private

      # Creates a Python dictionary representing a nested object property
      # for validation.
      def nested_obj_dict(prop, object, properties, spaces)
        name = python_variable_name(prop, object.azure_sdk_definition.create)
        options = prop_options(prop, object, spaces).join("\n")
        [
          "#{name}=dict(\n#{indent_list(options, 4, true)}\n    options=dict(",
          indent_list(properties.map do |p|
            python_dict_for_property(p, object, spaces + 4)
          end, 8),
          "    )\n)"
        ]
      end

      # Returns an array of all base options for a given property.
      def prop_options(prop, _object, spaces)
        [
          ('required=True' if prop.required && !prop.default_value && !is_location?(prop)),
          ("default=#{python_literal(prop.default_value)}" \
           if prop.default_value),
          "type=#{quote_string(python_type(prop))}",
          (choices_enum(prop, spaces) if prop.is_a? Api::Type::Enum),
          ("elements=#{quote_string(python_type(prop.item_type))}" \
            if prop.is_a? Api::Type::Array),
          ("aliases=[#{prop.aliases.map { |x| quote_string(x) }.join(', ')}]" \
            if prop.aliases),
          ('updatable=False' if prop.input && !is_resource_group?(prop) && !is_resource_name?(prop)),
          ("disposition='/'" if prop.input && !is_resource_group?(prop) && !is_resource_name?(prop))
        ].compact
      end

      # Returns a formatted string represented the choices of an enum
      def choices_enum(prop, spaces)
        name = prop.out_name.underscore
        type = "type=#{quote_string(python_type(prop))}"
        # + 6 for =dict(
        choices_indent = spaces + name.length + type.length + 6
        format([
                 [
                   "choices=[#{prop.values.map do |x|
                                 quote_string(x.to_s.underscore)
                               end.join(', ')}]"
                 ],
                 [
                   "choices=['#{prop.values[0]}',",
                   prop.values[1..-2].map do |x|
                     "#{indent(quote_string(x.to_s), choices_indent + 11)},"
                   end,
                   # + 11 for ' choices='
                   indent("#{quote_string(prop.values[-1].to_s)}]",
                          choices_indent + 11)
                 ]
               ], 0, choices_indent)
=======
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
>>>>>>> master
      end
    end
  end
end

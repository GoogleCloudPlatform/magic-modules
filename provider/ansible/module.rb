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
      # Returns the Python dictionary representing a simple property for
      # validation.
      def python_dict_for_property(prop, object, spaces = 0)
        if prop.is_a?(Api::Type::Array) && \
           prop.item_type.is_a?(Api::Type::NestedObject)
          nested_obj_dict(prop, object, prop.item_type.properties, spaces)
        elsif prop.is_a? Api::Type::NestedObject
          nested_obj_dict(prop, object, prop.properties, spaces)
        else
          name = prop.out_name.underscore
          "#{name}=dict(#{prop_options(prop, object)})"
        end
      end

      private

      # Creates a Python dictionary representing a nested object property
      # for validation.
      def nested_obj_dict(prop, object, properties, spaces)
        name = prop.out_name.underscore
        options = prop_options(prop, object)
        [
          "#{name}=dict(#{options}, options=dict(",
          indent_list(properties.map do |p|
            python_dict_for_property(p, object, spaces + 4)
          end, 4),
          '))'
        ]
      end

      # Returns an array of all base options for a given property.
      def prop_options(prop, _object)
        [
          ('required=True' if prop.required && !prop.default_value),
          ("default=#{python_literal(prop.default_value)}" \
           if prop.default_value),
          ("type=#{quote_string(python_type(prop))}" if python_type(prop)),
          # Choices enum always starts on a new line.
          ("\n" + choices_enum(prop) if prop.is_a? Api::Type::Enum),
          ("elements=#{quote_string(python_type(prop.item_type))}" \
            if prop.is_a?(Api::Type::Array) && python_type(prop.item_type)),
          ("aliases=[#{prop.aliases.map { |x| quote_string(x) }.join(', ')}]" \
            if prop.aliases)
        ].compact.reduce do |prev, nxt|
          # Avoid trailing spaces if we are about to have a newline.
          nxt.start_with?("\n") ? prev + ',' + nxt : prev + ', ' + nxt
        end
      end

      # Returns a formatted string represented the choices of an enum
      def choices_enum(prop)
        if prop.values.size == 1
          "choices=[#{quote_string(prop.values.first.to_s)}]"
        else
          choices_indent = prop.out_name.underscore.length + 6
          indent(
            [
              "choices=[#{quote_string(prop.values.first.to_s)},"
            ] +
              prop.values[1..-2].map do |x|
                "#{indent(quote_string(x.to_s), 9)},"
              end +
            [
              "#{indent(quote_string(prop.values.last.to_s), 9)}]"
            ],
            choices_indent
          )
        end
      end
    end
  end
end

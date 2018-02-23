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

module Provider
  class Ansible
    # Responsible for building out the AnsibleModule code.
    # AnsibleModule is responsible for input validation.
    module Module
      # Returns the Python dictionary representing a simple property for
      # validation.
      def python_dict_for_property(prop)
        if prop.is_a?(Api::Type::Array) && \
           prop.item_type.is_a?(Api::Type::NestedObject)
          nested_obj_dict(prop, prop.item_type.properties)
        elsif prop.is_a? Api::Type::NestedObject
          nested_obj_dict(prop, prop.properties)
        else
          name = Google::StringUtils.underscore(prop.out_name)
          "#{name}=dict(#{prop_options(prop).join(', ')})"
        end
      end

      private

      # Creates a Python dictionary representing a nested object property
      # for validation.
      def nested_obj_dict(prop, properties)
        name = Google::StringUtils.underscore(prop.out_name)
        [
          "#{name}=dict(#{prop_options(prop).join(', ')}, options=dict(",
          indent_list(properties.map { |p| python_dict_for_property(p) }, 4),
          '))'
        ]
      end

      # Returns an array of all base options for a given property.
      def prop_options(prop)
        [
          ('required=True' if prop.required),
          "type=#{quote_string(python_type(prop))}",
          (if prop.is_a? Api::Type::Enum
             "choices=[#{prop.values.map do |x|
                           quote_string(x.to_s)
                         end.join(', ')}]"
           end),
          ("elements=#{quote_string(python_type(prop.item_type))}" \
            if prop.is_a? Api::Type::Array)
        ].compact
      end
    end
  end
end

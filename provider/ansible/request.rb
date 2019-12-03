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
  module Ansible
    # Responsible for building out the resource_to_request and
    # request_from_hash methods.
    module Request
      # Takes in a list of properties and outputs a python hash that takes
      # in a module and outputs a formatted JSON request.
      def request_properties(properties, hash_name = 'module.params', module_name = 'module')
        properties.map do |prop|
          {
            Google::PythonUtils::UnicodeString.new(prop.api_name) =>
            Google::PythonUtils::PythonCode.new(request_output(prop, hash_name, module_name))
          }
        end.reduce({}, :merge)
      end

      def response_properties(properties, hash_name = 'response', module_name = 'module')
        properties.map do |prop|
          {
            Google::PythonUtils::UnicodeString.new(prop.api_name) =>
            Google::PythonUtils::PythonCode.new(response_output(prop, hash_name, module_name))
          }
        end.reduce({}, :merge)
      end

      # This returns a list of properties that require classes being built out.
      def properties_with_classes(properties)
        properties.map do |p|
          [p] + properties_with_classes(p.nested_properties) if p.nested_properties?
        end.compact.flatten
      end

      private

      # This is outputting code and code is easier to read on one line.
      # rubocop:disable Metrics/LineLength
      def response_output(prop, hash_name, module_name)
        # If input true, treat like request, but use module names.
        return request_output(prop, "#{module_name}.params", module_name) \
          if prop.input

        if prop.is_a? Api::Type::NestedObject
          "#{prop.property_class[-1]}(#{hash_name}.get(#{unicode_string(prop.api_name)}, {}), #{module_name}).from_response()"
        elsif prop.is_a?(Api::Type::Array) && \
              prop.item_type.is_a?(Api::Type::NestedObject)
          "#{prop.property_class[-1]}(#{hash_name}.get(#{unicode_string(prop.api_name)}, []), #{module_name}).from_response()"
        else
          "#{hash_name}.get(#{unicode_string(prop.api_name)})"
        end
      end

      def request_output(prop, hash_name, module_name, allow_pattern = true)
        # If type has a pattern, use the function.
        return "#{prop.name.underscore}_pattern(#{request_output(prop, hash_name, module_name, false)}, module)" \
          if prop.pattern && allow_pattern

        return "response.get(#{quote_string(prop.name)})" \
          if prop.is_a? Api::Type::FetchedExternal

        if prop.is_a? Api::Type::NestedObject
          "#{prop.property_class[-1]}(#{hash_name}.get(#{quote_string(prop.out_name)}, {}), #{module_name}).to_request()"
        elsif prop.is_a?(Api::Type::Array) && \
              prop.item_type.is_a?(Api::Type::NestedObject)
          "#{prop.property_class[-1]}(#{hash_name}.get(#{quote_string(prop.out_name)}, []), #{module_name}).to_request()"
        elsif prop.is_a?(Api::Type::ResourceRef) && !prop.resource_ref.readonly
          "replace_resource_dict(#{hash_name}.get(#{unicode_string(prop.name.underscore)}, {}), #{quote_string(prop.imports)})"
        elsif prop.is_a?(Api::Type::ResourceRef) && \
              prop.resource_ref.readonly && prop.imports == 'selfLink'
          "#{prop.resource.underscore}_selflink(#{hash_name}.get(#{quote_string(prop.out_name)}), #{module_name}.params)"
        elsif prop.is_a?(Api::Type::Array) && \
              prop.item_type.is_a?(Api::Type::ResourceRef) && \
              !prop.item_type.resource_ref.readonly
          "replace_resource_dict(#{hash_name}.get(#{quote_string(prop.name.underscore)}, []), #{quote_string(prop.item_type.imports)})"
        else
          "#{hash_name}.get(#{quote_string(prop.out_name)})"
        end
      end
      # rubocop:enable Metrics/LineLength
    end
  end
end

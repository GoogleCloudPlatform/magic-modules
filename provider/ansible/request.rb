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
    # rubocop:disable Metrics/ModuleLength
    module Request
      # Takes in a list of properties and outputs a python hash that takes
      # in a module and outputs a formatted JSON request.
      def request_properties(properties, indent = 4)
        indent_list(
          properties.map do |prop|
            request_property(prop, 'module.params', 'module')
          end,
          indent
        )
      end

      def response_properties(properties, indent = 8)
        indent_list(
          properties.map do |prop|
            response_property(prop, 'response', 'module')
          end,
          indent
        )
      end

      def request_properties_in_classes(properties, indent = 4,
                                        hash_name = 'self.request',
                                        module_name = 'self.module')
        indent_list(
          properties.map do |prop|
            request_property(prop, hash_name, module_name)
          end,
          indent
        )
      end

      def response_properties_in_classes(properties, indent = 8,
                                         hash_name = 'self.request',
                                         module_name = 'self.module')
        indent_list(
          properties.map do |prop|
            response_property(prop, hash_name, module_name)
          end,
          indent
        )
      end

      # This returns a list of properties that require classes being built out.
      def properties_with_classes(properties)
        properties.map do |p|
          if p.is_a? Api::Type::NestedObject
            [p] + properties_with_classes(p.properties)
          elsif p.is_a?(Api::Type::Array) && \
                p.item_type.is_a?(Api::Type::NestedObject)
            [p] + properties_with_classes(p.item_type.properties)
          end
        end.compact.flatten
      end

      private

      def request_property(prop, hash_name, module_name)
        [
          "#{unicode_string(prop.field_name)}:",
          request_output(prop, hash_name, module_name).to_s
        ].join(' ')
      end

      def response_property(prop, hash_name, module_name)
        [
          "#{unicode_string(prop.field_name)}:",
          response_output(prop, hash_name, module_name).to_s
        ].join(' ')
      end

      def response_output(prop, hash_name, module_name)
        # If input true, treat like request, but use module names.
        return request_output(prop, "#{module_name}.params", module_name) \
          if prop.input
        if prop.is_a? Api::Type::NestedObject
          [
            "#{prop.property_class[-1]}(",
            "#{hash_name}.get(#{unicode_string(prop.name)}, {})",
            ", #{module_name}).from_response()"
          ].join
        elsif prop.is_a?(Api::Type::Array) && \
              prop.item_type.is_a?(Api::Type::NestedObject)
          [
            "#{prop.property_class[-1]}(",
            "#{hash_name}.get(#{unicode_string(prop.name)}, [])",
            ", #{module_name}).from_response()"
          ].join
        else
          "#{hash_name}.get(#{unicode_string(prop.name)})"
        end
      end
      # rubocop:enable Metrics/MethodLength

      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      def request_output(prop, hash_name, module_name)
        if prop.is_a? Api::Type::NestedObject
          [
            "#{prop.property_class[-1]}(",
            "#{hash_name}.get(#{quote_string(prop.out_name)}, {})",
            ", #{module_name}).to_request()"
          ].join
        elsif prop.is_a?(Api::Type::Array) && \
              prop.item_type.is_a?(Api::Type::NestedObject)
          [
            "#{prop.property_class[-1]}(",
            "#{hash_name}.get(#{quote_string(prop.out_name)}, [])",
            ", #{module_name}).to_request()"
          ].join
        elsif prop.is_a?(Api::Type::ResourceRef) && !prop.resources.first.resource_ref.virtual
          prop_name = Google::StringUtils.underscore(prop.name)
          [
            "replace_resource_dict(#{hash_name}",
            ".get(#{unicode_string(prop_name)}, {}), ",
            "#{quote_string(prop.resources.first.imports)})"
          ].join
        elsif prop.is_a?(Api::Type::ResourceRef) && \
              prop.resources.first.resource_ref.virtual && prop.resources.first.imports == 'selfLink'
          func_name = Google::StringUtils.underscore("#{prop.name}_selflink")
          [
            "#{func_name}(#{hash_name}.get(#{quote_string(prop.out_name)}),",
            "#{module_name}.params)"
          ].join(' ')
        elsif prop.is_a?(Api::Type::Array) && \
              prop.item_type.is_a?(Api::Type::ResourceRef) && \
              !prop.item_type.resources.first.resource_ref.virtual
          prop_name = Google::StringUtils.underscore(prop.name)
          [
            "replace_resource_dict(#{hash_name}",
            ".get(#{quote_string(prop_name)}, []), ",
            "#{quote_string(prop.item_type.resources.first.imports)})"
          ].join
        else
          "#{hash_name}.get(#{quote_string(prop.out_name)})"
        end
      end
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/PerceivedComplexity
    end
    # rubocop:enable Metrics/ModuleLength
  end
end

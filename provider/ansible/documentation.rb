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

require 'compile/core'
require 'provider/config'
require 'provider/core'
require 'provider/ansible/manifest'

# Rubocop doesn't like this file because the hashes are complicated.
# Humans like this file because the hashes are explicit and easy to read.
module Provider
  module Ansible
    # Responsible for building out YAML documentation blocks.
    module Documentation
      def to_yaml(obj)
        if obj.is_a?(::Hash)
          obj.reject { |_, v| v.nil? }.to_yaml(:indentation => 4, :line_width => -1).sub("---\n", '')
        else
          obj.to_yaml(:indentation => 4, :line_width => -1).sub("---\n", '')
        end
      end

      # Builds out the DOCUMENTATION for a property.
      # This will eventually be converted to YAML
      def documentation_for_property(prop, object)
        required = prop.required && !prop.default_value ? true : false
        {
          python_variable_name(prop, object.azure_sdk_definition.create) => {
            'description' => [
              format_description(prop.description),
              (resourceref_description(prop) if prop.is_a?(Api::Azure::Type::ResourceReference))
            ].flatten.compact,
            'required' => (true if required && !is_location?(prop)),
            'default' => (prop.default_value.to_s.underscore if prop.default_value),
            'type' => python_type(prop),
            'aliases' => (prop.aliases if prop.aliases),
            'version_added' => (prop.version_added.to_f if prop.version_added),
            'choices' => (prop.values.map{|v| v.to_s.underscore} if prop.is_a? Api::Type::Enum),
            'suboptions' => (
              if prop.is_a?(Api::Type::NestedObject)
                prop.properties.map { |p| documentation_for_property(p, object) }.reduce({}, :merge)
              elsif prop.is_a?(Api::Type::Array) && prop.item_type.is_a?(Api::Type::NestedObject)
                prop.item_type.properties
                              .map { |p| documentation_for_property(p, object) }
                              .reduce({}, :merge)
              end
            )
          }.reject { |_, v| v.nil? }
        }
      end

      # Builds out the RETURNS for a property.
      # This will eventually be converted to YAML
      def returns_for_property(prop, object)
        type = python_type(prop)
        type = 'str' if prop.is_a? Api::Azure::Type::ResourceReference
        # Complex types only mentioned in reference to RETURNS YAML block
        # Complex types are nested objects traditionally, but arrays of nested
        # objects will be included to avoid linting errors.
        type = 'complex' if prop.is_a?(Api::Type::NestedObject) \
                            || (prop.is_a?(Api::Type::Array) \
                            && prop.item_type.is_a?(Api::Type::NestedObject))
        sample = prop.document_sample_value || prop.sample_value
        {
          python_variable_name(prop, object.azure_sdk_definition.read) => {
            'description' => format_description(prop.description),
            'returned' => 'always',
            'type' => type,
            'sample' => (sample unless sample.nil?),
            'contains' => (
              if prop.is_a?(Api::Type::NestedObject)
                prop.properties.map { |p| returns_for_property(p, object) }.reduce({}, :merge)
              elsif prop.is_a?(Api::Type::Array) && prop.item_type.is_a?(Api::Type::NestedObject)
                prop.item_type.properties.map { |p| returns_for_property(p, object) }.reduce({}, :merge)
              end
            )
          }.reject { |_, v| v.nil? }
        }
      end

      def autogen_notice_contrib
        ['Please read more about how to change this file at',
         'https://github.com/Azure/magic-module-specs']
      end

      def resourceref_description(prop)
        [
          "It can be the #{prop.resource_type_name} name which is in the same resource group.",
          "It can be the #{prop.resource_type_name} ID. e.g., #{prop.document_sample_value || prop.sample_value}.",
          "It can be a dict which contains C(name) and C(resource_group) of the #{prop.resource_type_name}."
        ]
      end

      # MM puts descriptions in a text block. Ansible needs it in bullets
      def format_description(desc)
        desc.split(".\n").map do |paragraph|
          paragraph += '.' unless paragraph.end_with?('.')
          paragraph = format_url(paragraph)
          paragraph.tr("\n", ' ').strip.squeeze(' ')
        end
      end

      # Find URLs and surround with U()
      def format_url(paragraph)
        paragraph.gsub(%r{
          https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]
          [a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+
          [a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))
          [a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9]\.[^\s]{2,}
        }x, 'U(\\0)')
      end
    end
  end
end

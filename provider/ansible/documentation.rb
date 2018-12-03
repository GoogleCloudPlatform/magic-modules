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
          obj.reject { |_, v| v.nil? }.to_yaml.sub("---\n", '')
        else
          obj.to_yaml.sub("---\n", '')
        end
      end

      # Builds out the DOCUMENTATION for a property.
      # This will eventually be converted to YAML
      def documentation_for_property(prop)
        required = prop.required && !prop.default_value ? true : false
        {
          prop.name.underscore => {
            'description' => [
              format_description(prop.description),
              (resourceref_description(prop) \
               if prop.is_a?(Api::Type::ResourceRef) && !prop.resource_ref.readonly)
            ].flatten.compact,
            'required' => required,
            'default' => (
              if prop.default_value&.is_a?(::Hash)
                prop.default_value
              else
                prop.default_value&.to_s
              end),
            'type' => ('bool' if prop.is_a? Api::Type::Boolean),
            'aliases' => prop.aliases,
            'version_added' => (prop.version_added&.to_f),
            'choices' => (prop.values.map(&:to_s) if prop.is_a? Api::Type::Enum),
            'suboptions' => (
              if prop.is_a?(Api::Type::NestedObject)
                prop.properties.map { |p| documentation_for_property(p) }.reduce({}, :merge)
              elsif prop.is_a?(Api::Type::Array) && prop.item_type.is_a?(Api::Type::NestedObject)
                prop.item_type.properties
                              .map { |p| documentation_for_property(p) }
                              .reduce({}, :merge)
              end
            )
          }.reject { |_, v| v.nil? }
        }
      end

      # Builds out the RETURNS for a property.
      # This will eventually be converted to YAML
      def returns_for_property(prop)
        type = python_type(prop)
        # Type is a valid AnsibleModule type, but not a valid return type
        type = 'str' if type == 'path'
        # Complex types only mentioned in reference to RETURNS YAML block
        # Complex types are nested objects traditionally, but arrays of nested
        # objects will be included to avoid linting errors.
        type = 'complex' if prop.is_a?(Api::Type::NestedObject) \
                            || (prop.is_a?(Api::Type::Array) \
                            && prop.item_type.is_a?(Api::Type::NestedObject))
        {
          prop.name => {
            'description' => format_description(prop.description),
            'returned' => 'success',
            'type' => type,
            'contains' => (
              if prop.is_a?(Api::Type::NestedObject)
                prop.properties.map { |p| returns_for_property(p) }.reduce({}, :merge)
              elsif prop.is_a?(Api::Type::Array) && prop.item_type.is_a?(Api::Type::NestedObject)
                prop.item_type.properties.map { |p| returns_for_property(p) }.reduce({}, :merge)
              end
            )
          }.reject { |_, v| v.nil? }
        }
      end

      def autogen_notice_contrib
        ['Please read more about how to change this file at',
         'https://www.github.com/GoogleCloudPlatform/magic-modules']
      end

      def resourceref_description(prop)
        [
          "This field represents a link to a #{prop.resource_ref.name} resource in GCP.",
          'It can be specified in two ways.',
          "You can add `register: name-of-resource` to a #{module_name(prop.resource_ref)} task",
          "and then set this #{prop.name.underscore} field to \"{{ name-of-resource }}\"",
          "Alternatively, you can set this #{prop.name.underscore} to a dictionary",
          "with the #{prop.imports} key",
          "where the value is the #{prop.imports} of your #{prop.resource_ref.name}"
        ].join(' ')
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

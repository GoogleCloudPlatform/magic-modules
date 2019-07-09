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

# Rubocop doesn't like this file because the hashes are complicated.
# Humans like this file because the hashes are explicit and easy to read.
module Provider
  module Ansible
    # Responsible for building out YAML documentation blocks.
    module Documentation
      # Builds out the DOCUMENTATION for a property.
      # This will eventually be converted to YAML
      def documentation_for_property(prop)
        required = prop.required && !prop.default_value ? true : false
        {
          prop.name.underscore => {
            'description' => [
              format_description(prop.description),
              (resourceref_description(prop) \
               if prop.is_a?(Api::Type::ResourceRef) && !prop.resource_ref.readonly),
              (choices_description(prop) \
               if prop.is_a?(Api::Type::Enum))
            ].flatten.compact,
            'required' => required,
            'default' => (
              if prop.default_value&.is_a?(::Hash)
                prop.default_value
              else
                prop.default_value&.to_s
              end),
            'type' => python_type(prop),
            'aliases' => prop.aliases,
            'version_added' => (version_added(prop)&.to_f),
            'suboptions' => (
                if prop.nested_properties?
                  prop.nested_properties.reject(&:output).map { |p| documentation_for_property(p) }
                                        .reduce({}, :merge)
                end
              )
          }.reject { |_, v| v.nil? }
        }
      end

      # Builds out the RETURNS for a property.
      # This will eventually be converted to YAML
      def returns_for_property(prop)
        type = python_type(prop) || 'str'
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
              if prop.nested_properties?
                prop.nested_properties.map { |p| returns_for_property(p) }
                                      .reduce({}, :merge)
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
          "First, you can place a dictionary with key '#{prop.imports}'",
          "and value of your resource's #{prop.imports}",
          'Alternatively, you can add `register: name-of-resource` to a',
          "#{module_name(prop.resource_ref)} task",
          "and then set this #{prop.name.underscore} field to \"{{ name-of-resource }}\""
        ].join(' ')
      end

      def choices_description(prop)
        "Some valid choices include: #{prop.values.map { |x| "\"#{x}\"" }.join(', ')}"
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
      # If there's a period at the end of the URL, make sure the
      # period is outside of the ()
      def format_url(paragraph)
        paragraph.gsub(%r{
          https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]
          [a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+
          [a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))
          [a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9]\.[^\s]{2,}
        }x, 'U(\\0)').gsub('.)', ').')
      end
    end
  end
end

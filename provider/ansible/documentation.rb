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

require 'api/object'
require 'compile/core'
require 'provider/config'
require 'provider/core'
require 'provider/ansible/manifest'

module Provider
  module Ansible
    # Responsible for building out YAML documentation blocks.
    # rubocop:disable Metrics/ModuleLength
    module Documentation
      # Takes a long string and divides each string into multiple paragraphs,
      # where each paragraph is a properly indented multi-line bullet point.
      #
      # Example:
      #   - This is a paragraph
      #     that wraps under
      #     the bullet properly
      #   - This is the second
      #     paragraph.
      def bullet_lines(line, spaces)
        line.split(".\n").map { |paragraph| bullet_line(paragraph, spaces) }
      end

      # Takes in a string (representing a paragraph) and returns a multi-line
      # string, where each line is less than max_length characters long and all
      # subsequent lines are indented in by spaces characters
      #
      # Example:
      #   - This is a sentence
      #     that wraps under
      #     the bullet properly
      #
      #   - |
      #     This is a sentence
      #     that wraps under
      #     the bullet properly
      #     because of the :
      #     character
      # rubocop:disable Metrics/AbcSize
      def bullet_line(paragraph, spaces, _multiline = true, add_period = true)
        paragraph += '.' unless paragraph.end_with?('.') || !add_period
        paragraph = format_url(paragraph)
        paragraph = paragraph.tr("\n", ' ').strip

        # Paragraph placed inside array to get bullet point.
        yaml = [paragraph].to_yaml
        # YAML documentation header is not necessary.
        yaml = yaml.gsub("---\n", '') if yaml.include?("---\n")

        # YAML dumper isn't very smart about line lengths.
        # If any line is over 160 characters (with indents), build the YAML
        # block using wrap_field.
        # Using YAML.dump output ensures that all character escaping done
        if yaml.split("\n").any? { |line| line.length > (160 - spaces) }
          return wrap_field(
            yaml.tr("\n", ' ').gsub(/\s+/, ' '),
            spaces + 3
          ).each_with_index.map { |x, i| i.zero? ? x : indent(x, 2) }
        end
        yaml
      end
      # rubocop:enable Metrics/AbcSize

      # Builds out a full YAML block for DOCUMENTATION
      # This includes the YAML for the property as well as any nested props
      def doc_property_yaml(prop, object, spaces)
        block = minimal_doc_block(prop, object, spaces)
        # Ansible linter does not support nesting options this deep.
        if prop.is_a?(Api::Type::NestedObject)
          block = block[prop.name.underscore].merge(nested_doc(prop.properties, object, spaces))
        elsif prop.is_a?(Api::Type::Array) &&
              prop.item_type.is_a?(Api::Type::NestedObject)
          block = block[prop.name.underscore].merge(nested_doc(prop.item_type.properties, object, spaces))
        end
        block.to_yaml.sub("---\n", '')
      end

      # Builds out a full YAML block for RETURNS
      # This includes the YAML for the property as well as any nested props
      def return_property_yaml(prop, spaces)
        block = minimal_return_block(prop, spaces)
        if prop.is_a? Api::Type::NestedObject
          block = block[prop.name].merge(nested_return(prop.properties, spaces))
        elsif prop.is_a?(Api::Type::Array) &&
              prop.item_type.is_a?(Api::Type::NestedObject)
          block = block[prop.name].merge(nested_return(prop.item_type.properties, spaces))
        end
        block.to_yaml.sub("---\n", '')
      end

      private

      # Returns formatted nested documentation for a set of properties.
      def nested_return(properties, spaces)
        {
          'contains' => properties.map { |p| return_property_yaml(p, spaces) }
        }.reject { |_, v| v.nil? }
      end

      def nested_doc(properties, object, spaces)
        {
          'suboptions' => properties.map { |p| doc_property_yaml(p, object, spaces) }
        }.reject { |_, v| v.nil? }
      end

      # Builds out the minimal YAML block for DOCUMENTATION
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      # rubocop:disable Metrics/AbcSize
      def minimal_doc_block(prop, _object, spaces)
        required = prop.required && !prop.default_value ? true : false
        {
          prop.name.underscore => {
            'description' => [
              format_description(prop.description),
              (resourceref_description(prop) if prop.is_a?(Api::Type::ResourceRef) && !prop.resource_ref.readonly)
            ].flatten.compact,
            'required' => required,
            'default' => (prop.default_value.to_s if prop.default_value),
            'type' => ('bool' if prop.is_a? Api::Type::Boolean),
            'aliases' => ("[#{prop.aliases.join(', ')}]" if prop.aliases),
            'version_added' => (prop.version_added.to_f if prop.version_added),
            'choices' => (prop.values.map(&:to_s) if prop.is_a? Api::Type::Enum)
          }.reject { |_, v| v.nil? }
        }
      end
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/PerceivedComplexity

      # Builds out the minimal YAML block for RETURNS
      def minimal_return_block(prop, spaces)
        type = python_type(prop)
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
            'type' => type
          }
        }.reject { |_, v| v.nil? }
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
          paragraph.gsub("\n", ' ').strip.squeeze(' ')
          #paragraph = paragraph.tr("\n", ' ').strip

          ## YAML isn't very smart about keeping line lengths sane.
          ## We'll double check that the lengths will be reasonable
          ## and if they aren't, use wrap_field to do it ourselves.
          #yaml = [paragraph].to_yaml
          #yaml = yaml.gsub("---\n", '') if yaml.include?("---\n")

          #if yaml.split("\n").any? { |line| line.length > (149) }
          #  wrap_field(
          #    paragraph.tr("\n", ' ').gsub(/\s+/, ' '),
          #    11
          #  ).each_with_index.map { |x, i| i.zero? ? x : indent(x, 2) }
          #else
          #  paragraph
          #end
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

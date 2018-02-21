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
  class Ansible
    # Responsible for building out YAML documentation blocks.
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
      def bullet_line(paragraph, spaces, add_period=true)
        # - 2 for "- "
        indented = wrap_field(paragraph, spaces - 2)
        indented = indented.split("\n")
        indented[0] = indented[0].sub(/^../, '- ')
        # Add in a period at paragraph end unless there's already a period.
        if add_period
          indented[-1] += '.' unless indented.last.end_with?('.')
        end
        indented
      end

      # Builds out a full YAML block for DOCUMENTATION
      # This includes the YAML for the property as well as any nested props
      def doc_property_yaml(prop, spaces)
        block = minimal_doc_block(prop, spaces)
        return block unless prop.is_a? Api::Type::NestedObject
        block << indent('suboptions:', 4)
        block.concat(
          prop.properties.map do |p|
            indent(doc_property_yaml(p, spaces + 4), 8)
          end
        )
      end

      # Builds out a full YAML block for RETURNS
      # This includes the YAML for the property as well as any nested props
      def return_property_yaml(prop, spaces)
        block = minimal_return_block(prop, spaces)
        return block unless prop.is_a? Api::Type::NestedObject
        block << indent('contains:', 4)
        block.concat(
          prop.properties.map do |p|
            indent(return_property_yaml(p, spaces + 4), 8)
          end
        )
      end

      private

      # Builds out the minimal YAML block for DOCUMENTATION
      def minimal_doc_block(prop, spaces)
        [
          minimal_yaml(prop, spaces),
          indent("required: #{prop.required ? 'true' : 'false'}", 4)
        ]
      end

      # Builds out the minimal YAML block for RETURNS
      def minimal_return_block(prop, spaces)
        type = python_type(prop)
        # Complex types only mentioned in reference to RETURNS YAML block
        type = 'complex' if prop.is_a? Api::Type::NestedObject
        [
          minimal_yaml(prop, spaces),
          indent([
                   'returned: success',
                   "type: #{type}"
                 ], 4)
        ]
      end

      # Builds out the minimal YAML block necessary for a property.
      # This block will need to have additional information appened
      # at the end.
      def minimal_yaml(prop, spaces)
        [
          "#{Google::StringUtils.underscore(prop.name)}:",
          indent(
            [
              'description:',
              indent(bullet_lines(prop.description, spaces + 4), 4)
            ], 4
          )
        ]
      end
    end
  end
end

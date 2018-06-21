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
  module TestData
    # Creates an expectation body for POST / DELETE Unit Tests
    # rubocop:disable Metrics/ClassLength
    class CreateData
      def initialize(provider, data_gen)
        @provider = provider
        @data_gen = data_gen
      end

      # rubocop:disable Metrics/AbcSize
      #
      # Returns a hash representing the data expected in a POST call to GCP
      # Used for unit tests
      # Parameters:
      #   path: Array of strings representing a path in provider.yaml file with
      #         user-overriden create data.
      #   has_name: Boolean, true if title != name
      #   tests: Hash from provider.yaml that may have user-overriden tests
      #          (These tests would be at tests[path])
      #   object: The object being tested.
      #
      def create_expect_data(path, has_name, tests, object)
        cust_result = Google::HashUtils.navigate(tests, path)
        if cust_result.nil?
          name_prop = object.all_user_properties.select { |p| p.name == 'name' }
          expect = []
          expect << "'kind' => '#{object.kind}'" if object.kind?
          expect.concat(
            object.properties.reject(&:output)
                             .map do |prop|
                               expect_hash(prop, name_prop, has_name)
                             end
          )
          expect.concat(
            (object.parameters || []).select(&:input)
                                     .map do |prop|
                                       expect_hash(prop, name_prop, has_name)
                                     end
          )
          @provider.indent_list(expect, 0)
        else
          "load_network_result(#{quote_string(cust_result)})"
        end
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/MethodLength

      private

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      def expect_hash(prop, name_prop, has_name, seed = 0)
        if prop.name == 'name' && !has_name && !name_prop.empty?
          "'#{name_prop[0].field_name}' => 'title#{seed}'"
        elsif prop.class <= Api::Type::Array
          expect_array_hash(prop)
        elsif prop.class <= Api::Type::NestedObject
          ["'#{prop.field_name}' => {",
           @provider.indent_list(
             prop.properties.map { |p| expect_hash(p, [], false, seed) }, 2
           ),
           '}'].join("\n")
        elsif prop.is_a? Api::Type::ResourceRef
          # All ResourceRefs should expect the fetched value
          # Without this, the JSON call will be expecting the title of the
          # ResourceRef block, not a value within that block.
          ["'#{prop.field_name}'", '=>', value(prop.resources[0].property.class,
                                               prop.resources[0].property,
                                               seed)].join(' ')
        else
          ["'#{prop.field_name}'", '=>', value(prop.class, prop,
                                               seed)].join(' ')
        end
      end
      # rubocop:enable Metrics/PerceivedComplexity
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/AbcSize

      # rubocop:disable Metrics/AbcSize
      def expect_array_hash(prop, seed = 0)
        if prop.item_type.class <= Api::Type::NestedObject
          [
            ["'#{prop.field_name}'", '=>', '['].join(' '),
            expect_array_item_hash(prop, seed),
            ']'
          ]
        elsif prop.item_type.class <= Api::Type::ResourceRef
          [
            ["'#{prop.field_name}'", '=>', '['].join(' '),
            expect_array_item_rref(prop, seed),
            ']'
          ]
        elsif prop.is_a? Api::Type::ResourceRef
          "'#{prop.field_name}' => #{value(prop.property.class,
                                           prop.property, seed)}"
        else
          ["'#{prop.field_name}'", '=>', value(prop.class,
                                               prop, seed)].join(' ')
        end
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/MethodLength

      def expect_array_item_rref(item, seed = 0)
        size = @data_gen.object_size(item, seed, true)
        imports = Google::StringUtils.underscore(item.item_type.imports)
        resource = Google::StringUtils.underscore(item.item_type.resource)
        @provider.indent_list(
          (0..size - 1).map do |index|
            "'#{imports.tr('_', '')}(resource(#{resource},#{index}))'"
          end, 2
        )
      end

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/BlockLength
      # rubocop:disable Metrics/MethodLength
      def expect_array_item_hash(item, seed = 0)
        size = @data_gen.object_size(item, seed, true)
        @provider.indent_list(
          (1..size).map do |index|
            ['{',
             @provider.indent_list(
               item.item_type.properties.map do |prop|
                 if prop.is_a? Api::Type::NestedObject
                   ["'#{prop.field_name}' => {",
                    @provider.indent_list(
                      prop.properties.map do |p|
                        expect_hash(p, [], false, seed + index - 1)
                      end, 2
                    ),
                    '}'].join("\n")
                 elsif prop.is_a? Api::Type::ResourceRef
                   # All ResourceRefs should expect the fetched value
                   # Without the JSON call will be expecting the title of the
                   # ResourceRef block, not a value within that block.
                   [
                     "'#{prop.field_name}' =>",
                     value(prop.property.class,
                           prop.property, (seed + index - 1) % MAX_ARRAY_SIZE)
                   ].join(' ')
                 elsif prop.is_a? Api::Type::Array
                   expect_array_hash(prop, (seed + index - 1))
                 else
                   [
                     "'#{prop.field_name}'", '=>',
                     value(prop.class, prop, seed + index - 1)
                   ].join(' ')
                 end
               end, 2
             ),
             '}']
          end, 2
        )
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/BlockLength
      # rubocop:enable Metrics/MethodLength

      # Returns a value formatted according to its class.
      def value(prop_class, prop, seed)
        val = @data_gen.value(prop_class, prop, seed)
        format_value(val)
      end

      # Formats a value according to its class.
      def format_value(value)
        types.each do |k, v|
          return v.call(value) if value.is_a? k
        end
        value
      end

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/MethodLength
      def types
        {
          Integer => ->(value) { Google::IntegerUtils.underscore(value) },
          String => ->(value) { return_quoted(value) },
          Symbol => ->(value) { return_quoted(value.to_s) },
          Time => ->(value) { return_quoted(value.iso8601) },
          Array => lambda do |value|
            return "%w[#{value.join(' ')}]" if value[0].is_a? String
            value
          end,
          Float => lambda do |value|
            values = value.to_s.split('.')

            [Google::IntegerUtils.underscore(values[0].to_i), '.',
             values[1]].join
          end,
          Hash => lambda do |value|
            values = value.map do |k, v|
              "'#{k}' => #{format_value(v)}"
            end
            @provider.format([
                               [
                                 '{',
                                 @provider.indent_list(values, 2),
                                 '}'
                               ].flatten
                             ], 0, 0)
          end
        }
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/MethodLength

      # Return a string formatted with quotes
      def return_quoted(value)
        is_quoted = value[0] == '\'' && value[-1] == '\''
        return value if is_quoted
        "'#{value}'"
      end
    end
    # rubocop:enable Metrics/ClassLength
  end
end

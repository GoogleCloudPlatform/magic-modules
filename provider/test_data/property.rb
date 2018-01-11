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

require 'google/integer_utils'

module Provider
  module TestData
    # Class responsible for generating the per-property tests
    # rubocop:disable Metrics/ClassLength
    class Property
      def initialize(provider)
        @provider = provider
      end

      # This returns a formatted string representing a single Rspec test
      # The test will be of the form:
      #   it { is_expected.to have_attributes({prop.name}: #{expected value}) }
      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/ParameterLists
      # rubocop:disable Metrics/PerceivedComplexity
      def property(prop, index, comparator, value, start_indent = 0,
                   name_override = nil)
        Google::LOGGER.info \
          "Generating test #{prop.out_name}[#{index}] #{comparator} #{value}"

        if prop.class <= Api::Type::ResourceRef
          resourceref_property(prop, value, name_override)
        elsif prop.class <= Api::Type::NameValues
          namevalues_property(prop, value, name_override)
        elsif prop.class <= Api::Type::NestedObject
          nested_property(prop, value, name_override)
        elsif prop.class <= Api::Type::Array \
          && prop.item_type != 'Api::Type::String'
          array_property(prop, value, name_override)
        elsif prop.class <= Api::Type::NameValues
          namevalue_property(prop, value, name_override)
        else
          single_property(prop, value, start_indent, name_override)
        end
      end
      # rubocop:enable Metrics/PerceivedComplexity
      # rubocop:enable Metrics/ParameterLists
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/AbcSize

      private

      # rubocop:disable Metrics/MethodLength # long but easier to read together
      def single_property(prop, value, start_indent, name_override)
        name = name_override || prop.out_name
        value = format_value(value)
        @provider.format(
          [
            ["it { is_expected.to have_attributes(#{name}: #{value}) }"],
            [
              'it do',
              @provider.indent(
                "is_expected.to have_attributes(#{name}: #{value})",
                2
              ),
              'end'
            ],
            [
              'it do',
              @provider.indent(
                [
                  'is_expected',
                  @provider.indent(".to have_attributes(#{name}: #{value})", 2)
                ], 2
              ),
              'end'
            ],
            [
              'it do',
              @provider.indent(
                ['is_expected',
                 @provider.indent(
                   ['.to have_attributes(',
                    @provider.indent("#{name}: #{value}", 2),
                    ')'], 2
                 )], 2
              ),
              'end'
            ],
            [
              'it do',
              @provider.indent(
                [
                  'is_expected',
                  @provider.indent(
                    [
                      '.to have_attributes(',
                      @provider.indent(
                        [
                          "#{name}:",
                          value.to_s
                        ], 2
                      ),
                      ')'
                    ], 2
                  )
                ], 2
              ),
              'end'
            ]
          ], 0, start_indent
        )
      end
      # rubocop:enable Metrics/MethodLength

      def array_property(prop, _value, name_override)
        name = name_override || prop.name
        [
          '# TODO(nelsonjr): Implement complex array object test.',
          "# it '#{name}' do",
          '#   # Add test code here',
          '# end'
        ]
      end

      def namevalues_property(prop, _value, name_override)
        name = name_override || prop.name
        [
          '# TODO(nelsonjr): Implement complex namevalues property test.',
          "# it '#{name}' do",
          '#   # Add test code here',
          '# end'
        ]
      end

      def nested_property(prop, _value, name_override)
        name = name_override || prop.name
        [
          '# TODO(nelsonjr): Implement complex nested property object test.',
          "# it '#{name}' do",
          '#   # Add test code here',
          '# end'
        ]
      end

      def resourceref_property(prop, _value, name_override)
        name = name_override || prop.name
        [
          '# TODO(alexstephen): Implement resourceref test.',
          "# it '#{name}' do",
          '#   # Add test code here',
          '# end'
        ]
      end

      def namevalue_property(prop, _value, name_override)
        name = name_override || prop.name
        [
          '# TODO(alexstephen): Implement name values test.',
          "# it '#{name}' do",
          '#   # Add test code here',
          '# end'
        ]
      end

      # Returns a value formatted according to its class.
      def format_value(value)
        types.each do |k, v|
          return v.call(value) if value.is_a? k
        end
        value
      end

      def types
        {
          Integer => ->(value) { Google::IntegerUtils.underscore(value) },
          String => ->(value) { quote_string(value) },
          Symbol => ->(value) { quote_string(value.to_s) },
          Time => ->(value) { "::Time.parse('#{value.iso8601}')" },
          Array => lambda do |value|
            return "%w[#{value.join(' ')}]" if value[0].is_a? String
            # Arrays of non-String types are not fully supported.
            # All tests using non-String arrays will return a blank test.
            # This function may not return the expected value for Arrays
            # of non-String types.
            # TODO(alexstephen): Add support for testing arrays.
            raise 'Non-String arrays are not supported.'
          end
        }
      end

      def quote_string(value)
        @provider.quote_string(@provider.unquote_string(value))
      end
    end
    # rubocop:enable Metrics/ClassLength
  end
end

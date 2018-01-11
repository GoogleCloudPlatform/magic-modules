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

require 'provider/test_data/generator'

module Provider
  module TestData
    # Responsible for all actions involving the Unit Test Constant values
    # This includes building arrays of constants and building hashes of
    # references to those values.
    class Constants
      MAX_TEST_DATAS = 5

      def initialize(provider)
        @provider = provider
        @data_gen = Provider::TestData::Generator.new
      end

      # Generates a hash mapping a object value's constant and the various
      # semi-random values that can be used for tests based on a seed (that will
      # be used to index the array.
      #
      # Example:
      #
      # {
      #   A_PROJECT_DATA: [
      #     "'test project#0 data'",
      #     "'test project#1 data'",
      #     "'test project#2 data'",
      #     "'test project#3 data'",
      #     "'test project#4 data'"
      #   ]
      # }
      #
      # See how it is aggregated by aggregate_constants_by_name.
      def value_arrays(product)
        var_data = product.objects.map do |object|
          next if object.exclude
          self_link_variables(object).map do |v|
            [
              [test_prefix(object), v.upcase, 'DATA'].join('_'),
              [object.name, v],
              value_arrays_create_values(object, v)
            ]
          end
        end

        aggregate_constants_by_name(var_data.flatten(1).compact)
      end

      # Generates a series of hash mappings that map a object's self link
      # variables to the test constant file
      #
      # Used primarily for uri_data functions
      #
      # Example:
      #   project: GoogleTests::Constants::F_PROJECT_DATA[(id - 1) \
      #     % GoogleTests::Constants::F_PROJECT_DATA.size],
      def value_assign(object)
        self_link_variables(object).map do |v|
          name = ['GoogleTests', 'Constants',
                  "#{test_prefix(object)}_#{v.upcase}_DATA"].join('::')
          @provider.format(
            [
              ["#{v}: #{name}[(id - 1) \\",
               @provider.indent("% #{name}.size]", 2)],
              [
                "#{v}:",
                @provider.indent(["#{name}[(id - 1) \\",
                                  @provider.indent("% #{name}.size]", 2)], 2)
              ]
            ], 0, 6 + 1
          ) # 1 extra for trailing comma
        end
      end

      private

      # Takes an array and aggregate by constant name, but combine the values
      # that have the same data name
      #
      # [
      #   ['A', 'Animal', [....]
      #   ['A', 'Ant', [....]
      # ]
      #
      # Becomes:
      #
      # {
      #   'A': {
      #     source: [
      #       'Animal'
      #       'Ant'
      #     ]
      #     data: [....]
      #   }
      # }
      def aggregate_constants_by_name(all_data)
        all_data.each_with_object({}) do |data, result|
          result[data[0]] = {} unless result.key?(data[0])
          result[data[0]][:source] = [] unless result[data[0]].key?(:source)
          result[data[0]][:source] << data[1]
          result[data[0]][:data] = data[2]
        end
      end

      def value_arrays_create_values(object, variable)
        test_values = []
        (0..MAX_TEST_DATAS - 1).each do |index|
          property = @provider.variable_type(object, variable)
          test_values << "'#{@data_gen.value(property.class, property, index)}'"
        end
        test_values
      end

      def test_prefix(object)
        object.name.gsub(/[a-z]/, '')
      end

      def self_link_variables(object)
        @provider.extract_variables([@provider.self_link_raw_url(object)
                                              .join('/'),
                                     @provider.collection_url(object)]
                                             .join("\n")).uniq
      end
    end
  end
end

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
require 'json'
require 'time'
require 'zlib'

module Provider
  module TestData
    # rubocop:disable Metrics/ClassLength
    # Builds out network data expectations for unit tests
    class Expectations
      def initialize(provider, data_gen)
        @provider = provider
        @data_gen = data_gen
      end

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      #
      # Creates all of the network expectations for a given test that involves
      # creating an object.
      # This includes a failure for a GET request to fetch an object, a POST to
      # create the object and all resource refs.
      #
      # Requires:
      #   tests: config object for overriden tests
      #   object: An Api::Resource object.
      #   test: A hash with the following test information:
      #     path - the path on the config object with a possible test override
      #     has_name - boolean, true if title != name
      #     expected_data - a hash describing the POST request for creating
      #                     a new object.
      #  rrefs: A list of Api::Type::ResourceRefs
      def create_before_data(tests, object, test, rrefs)
        cust_before = @provider.get_code_multiline(tests, test[:path])
        if cust_before.nil?
          name_prop = object.all_user_properties.select { |p| p.name == 'name' }
          # Get failed logic
          get_failed = create_expectation('expect_network_get_failed',
                                          test[:has_name], object, 12, [],
                                          1)

          extra_props = []
          extra_props << "#{name_prop[0].field_name}: 'title0'" \
            unless test[:has_name] || name_prop.empty?

          rref_list = object.uri_properties.map do |ref|
            # We need to verify that only resourcerefs directly belonging to
            # this object are inserted into the expectation.
            next unless ref.is_a? Api::Type::ResourceRef
            name = Google::StringUtils.underscore(ref.resources[0].resource_ref.name)
            value = @data_gen.value(ref.resources[0].property.class,
                                    ref.resources[0].property, 0)
            { name => value }
          end

          extra_props.concat(
            rref_list.flatten.compact.reduce({}, :merge)
                                     .map { |k, v| "#{k}: '#{v}'" }
          )

          code = [get_failed, 'expect_network_create \\',
                  @provider.indent(
                    if extra_props.empty?
                      ['1,', test[:expected_data]]
                    else
                      ['1,', '{',
                       @provider.indent(test[:expected_data], 2),
                       '},', @provider.indent_list(extra_props, 0)]
                    end,
                    2
                  )]

          unless object.async.nil?
            code << create_expectation('expect_network_get_async',
                                       test[:has_name], object, 12, [], 1)
          end

          code.concat(create_resource_ref_get_success(object, rrefs, 12))

          add_style_exemptions code, object, test

          code.flatten.compact.uniq
        else
          # rubocop:disable Security/Eval
          eval("\"#{cust_before}\"", binding, __FILE__, __LINE__)
          # rubocop:enable Security/Eval
        end
      end
      # rubocop:enable Metrics/PerceivedComplexity
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/AbcSize

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/PerceivedComplexity
      # Creates all of the network expectations for a given test that involves
      # deleting an object.
      # This includes a success for a GET request to fetch an object, a DELETE
      # to delete the object and all resource refs.
      #
      # Requires:
      #   tests: config object for overriden tests
      #   object: An Api::Resource object.
      #   test: A hash with the following test information:
      #     exists - boolean, true if object exists
      #     path - the path on the config object with a possible test override
      #     has_name - boolean, true if title != name
      #  rrefs: A list of Api::Type::ResourceRefs
      def delete_before_data(tests, object, test, rrefs)
        cust_before = @provider.get_code_multiline(tests, test[:path])
        if cust_before.nil?
          get = "expect_network_get_#{test[:exists] ? 'success' : 'failed'}"
          get_line = create_expectation(get, test[:has_name], object, 12)
          code = [get_line]

          if test[:exists]
            has_rrefs = !object.uri_properties
                               .select { |p| p.is_a? Api::Type::ResourceRef }
                               .empty?
            # Delete specifies name as a parameter, not as part of data
            params = []
            params = ['nil'] if test[:has_name] && has_rrefs
            params = ['\'title0\''] unless test[:has_name]

            code << create_expectation('expect_network_delete',
                                       true, object, 12, params, 1)

            unless object.async.nil?
              code << create_expectation('expect_network_get_async',
                                         test[:has_name], object, 12, [],
                                         1)
            end
          end

          rrefs.each do |ref|
            next if ref.object == object
            name = Google::StringUtils.underscore(ref.object.name)
            # Puppet style refs include a seed
            code << create_expectation("expect_network_get_success_#{name}",
                                       true, ref.object, 12, [],
                                       (ref.seed % MAX_ARRAY_SIZE) + 1)
          end

          code.flatten.uniq
        else
          # rubocop:disable Security/Eval
          eval("\"#{cust_before}\"", binding, __FILE__, __LINE__)
          # rubocop:enable Security/Eval
        end
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/PerceivedComplexity

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/ParameterLists
      # Given an expectation (i.e. "expect_network_get_success"),
      # returns that expectation with all necessary parameters.
      #
      # Example input: "expect_network_get_success", false, _, 0, 2"
      # Example output: "expect_network_get_success 2, name: 'title1'
      #
      # Requires:
      #   func_name - name of expectation (i.e. "expect_network_get_success")
      #   has_name - boolean, true if title != name
      #   space_used - number of spaces used. Used to calculate indentation
      #   prop_list - List of properties appended to network request.
      #   id - the id used for loading network data yaml files
      #   rrefs - list of ResourceRefs
      def create_expectation(func_name, has_name, object, space_used,
                             prop_list = [], id = 1)
        prop_list << "name: 'title#{id - 1}'" unless has_name

        rref_list = object.uri_properties.map do |ref|
          # We need to verify that only resourcerefs directly belonging to this
          # object are inserted into the expectation.
          next unless ref.is_a? Api::Type::ResourceRef
          name = Google::StringUtils.underscore(ref.resources[0].resource_ref.name)
          value = @data_gen.value(ref.resources[0].property.class,
                                  ref.resources[0].property, id - 1)
          { name => value }
        end

        prop_list.concat(
          rref_list.flatten.compact.reduce({}, :merge)
                                   .map { |k, v| "#{k}: '#{v}'" }
        )

        return "#{func_name} #{id}" if prop_list.empty?

        prop_list.unshift id.to_s

        @provider.format([
                           ["#{func_name} #{prop_list.join(', ')}"],
                           [
                             "#{func_name} #{prop_list[0]},",
                             @provider.indent_list(prop_list.drop(1),
                                                   func_name.length + 1)
                           ],
                           [
                             "#{func_name} \\",
                             @provider.indent_list(prop_list, 2)
                           ]
                         ], 0, space_used)
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/ParameterLists

      # Creates a single "expect_network_get_success" request for a resourceref
      def create_resource_ref_get_success(object, refs, inside_indent)
        # Generate network expectations for collected resourcerefs
        refs.map do |ref|
          # If the object being tested references itself (or another object of
          # the same kind) skip generation of dependencies to avoid colliding
          # with same ID objects
          next if ref.object == object

          ref_name = Google::StringUtils.underscore(ref.object.name)

          # Find the machine resource to safity the object's dependency. If an
          # object required by the object being tested in turn requires 1+ other
          # objects they were collected by the 'manifester.collect_refs' call
          # earlier.
          #
          # We now need to find the necessary dependencies and bind them to the
          # object being emitted.
          [
            create_expectation("expect_network_get_success_#{ref_name}", true,
                               ref.object, inside_indent, [],
                               (ref.seed % MAX_ARRAY_SIZE) + 1)
          ]
        end.compact
      end

      private

      def add_style_exemptions(code, object, test)
        rubo_off = @provider.get_rubocop_exceptions(
          "spec/#{object.out_name}", :test, [object.name, test[:path]],
          :disabled
        )
        code.unshift rubo_off unless rubo_off.nil? || rubo_off.empty?

        rubo_on = @provider.get_rubocop_exceptions(
          "spec/#{object.out_name}", :test, [object.name, test[:path]], :enabled
        )
        code.push rubo_on unless rubo_on.nil? || rubo_on.empty?
      end

      def get_required_prop_data(object)
        all_props = object.parameters || []
        all_props +=
          object.properties.select { |p| p.name == 'name' || p.required }
        all_props
      end
    end
    # rubocop:enable Metrics/ClassLength
  end
end

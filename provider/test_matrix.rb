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
require 'google/logger'

module Provider
  # Helper to track and validate the test matrix for the code
  class TestMatrix # rubocop:disable Metrics/ClassLength
    include Compile::Core

    # Tracks TestMatrix and verify them.
    class Collector
      def initialize
        @matrixes = []
      end

      def add(matrix, file, object)
        Google::LOGGER.info \
          "Registering test matrix for #{object.name} @ #{file}"
        @matrixes << matrix
      end

      def verify_all
        @matrixes.each(&:verify)
      end
    end

    # Tracks globally all TestMatrix objects.
    class Registry < Collector
      include Singleton
    end

    attr_reader :level

    def initialize(file, object, provider, expected)
      Registry.instance.add self, file, object

      @file = file
      @object = object
      @provider = provider
      @expected = expected

      # A collector for test contexts being produced. It will be used at the end
      # of the file generation to ensure that all necessary tests were properly
      # created.
      @actual = []
      @level = 0
      @hierarchy = []
    end

    def push(ensurable, exists = :none, changes = :none, has_name = :none,
             success = :none)
      data = {
        ensurable: ensurable, exists: exists, changes: changes,
        has_name: has_name, success: success
      }

      bundled = has_name.class <= Array
      if bundled
        push(ensurable, exists, changes, has_name[0])
        data[:has_name] = has_name[0]
        data[:success] = has_name[1]
      end

      update_level data, bundled
      @hierarchy.push data
      emit_push data
    end

    def pop(ensurable, exists = :none, changes = :none, has_name = :none,
            success = :none)
      bundled = has_name.class <= Array
      if bundled
        bundle = has_name
        has_name = bundle[0]
        success = bundle[1]
      end

      data = {
        ensurable: ensurable, exists: exists, changes: changes,
        has_name: has_name, success: success
      }

      verify_pop data

      @hierarchy.pop

      pop(ensurable, exists, changes, bundle[0]) if bundled

      update_level data, bundled
      lines(indent('end', @level * 2))
    end

    # Ensures that all test contexts are defined
    def verify
      Google::LOGGER.info "Verifying test matrix for #{@object.name} @ #{@file}"
      verify_topics
      verify_match_expectations
      fail_if_not_all_popped unless @hierarchy.empty?
    end

    private

    def emit_push(data)
      validate_test_context_args(data)

      new_context = data.values.reject { |t| t == :ignore }
      raise "Context already exists: #{new_context}" \
        if @actual.include?(new_context)
      @actual << new_context

      format_response new_context, data
    end

    def present?(needle, haystack)
      needle.reduce(false) { |result, n| result || haystack.include?(n) }
    end

    def format_response(new_context, data)
      format_handlers = {
        %i[pass fail] => :format_response_success,
        %i[has_name no_name] => :format_response_title,
        %i[changes no_change] => :format_response_changes,
        %i[exists missing] => :format_response_exists,
        %i[present absent] => :format_response_ensurable
      }

      format_handlers.each do |values, handler|
        return send(handler, data, @level * 2) if present?(values, new_context)
      end

      raise "Unknown context: #{new_context}"
    end

    def format_response_ensurable(data, spaces)
      lines(indent("context 'ensure == #{data[:ensurable].id2name}' do",
                   spaces))
    end

    def format_response_exists(data, spaces)
      lines(indent("context 'resource #{data[:exists].id2name}' do", spaces))
    end

    def format_response_changes(data, spaces)
      prefix = data[:changes] == :changes ? '' : 'no '
      lines(indent([
                     [
                       "# Ensure #{data[:ensurable].id2name}: resource",
                       [data[:exists].id2name,
                        data[:changes].id2name].join(', ').tr('_', ' ')
                     ].join(' '),
                     "context '#{prefix}changes == #{prefix}action' do"
                   ], spaces))
    end

    def format_response_title(data, spaces)
      msg = (data[:has_name] == :has_name ? 'title != name' : 'title == name')
      lines(indent([["# Ensure #{data[:ensurable].id2name}: resource",
                     [
                       data[:exists].id2name, data[:changes].id2name,
                       data[:has_name].id2name
                     ].join(', ').tr('_', ' ')].join(' '),
                    "context '#{msg}' do"], spaces))
    end

    def format_response_success(data, spaces)
      msg = (data[:has_name] == :has_name ? 'title != name' : 'title == name')
      lines(indent([["# Ensure #{data[:ensurable].id2name}: resource",
                     [
                       data[:exists].id2name, data[:changes].id2name,
                       data[:has_name].id2name, data[:success].id2name
                     ].join(', ').tr('_', ' ')].join(' '),
                    "context '#{msg} (#{data[:success].id2name})' do"], spaces))
    end

    # Converts a tree style hash into a full array:
    #   A {
    #     B [
    #       C
    #       D
    #     ]
    #     E [
    #       F
    #     ]
    #   }
    # will become:
    #   [[A, B, C], [A, B, D], [A, E, F]
    def hash_explode(hash, tree = [], output = [])
      hash.each_with_object(output) do |(k, v), results|
        if v.class <= Hash
          hash_explode(v, [tree, k], results)
        else
          v.each { |e| results << [tree, k, e].flatten }
        end
      end
    end

    def lines(content)
      @provider.lines(content)
    end

    def validate_test_context_args(data)
      valid_values = {
        ensurable: %i[present absent ignore],
        exists: %i[none exists missing],
        changes: %i[none changes no_change ignore],
        has_name: %i[none has_name no_name ignore],
        success: %i[none pass fail]
      }

      valid_values.each do |arg, values|
        raise "Bad ensure argument #{data[arg]}" \
          unless values.include?(data[arg])
      end
    end

    def verify_topics
      processed = [%i[none none none none none]]

      @actual.each do |test|
        test = test.reject { |t| t == :ignore }
        test << :none while test.size < processed[0].size

        topic = topic_for_test(test)
        raise "Missing topic: #{topic.reject { |t| t == :none }}" \
          unless processed.include?(topic)

        processed << test
      end
    end

    def topic_for_test(test)
      topic = test.clone

      if test.index(:none).nil?
        topic[topic.size - 1] = :none
      else
        topic[test.index(:none) - 1] = :none
      end

      topic
    end

    def fail_if_not_all_popped
      Google::LOGGER.info 'Missing pop() calls:'
      @hierarchy.each { |p| Google::LOGGER.info "  - #{p}" }
      raise "Missing pop() calls: #{@hierarchy}"
    end

    def verify_match_expectations
      missing = []
      extra = @actual.reject { |t| t.include?(:none) || t.include?(:ignore) }
      hash_explode(@expected).each do |e|
        missing << e unless @actual.include?(e)
        extra.delete(e)
      end

      fail_match_expectations_if_missing(missing) unless missing.empty?
      fail_match_expectations_if_extra(extra) unless extra.empty?

      missing.empty? && extra.empty?
    end

    def fail_match_expectations_if_missing(missing)
      Google::LOGGER.info \
        'FATAL: The following tests are missing from the matrix:'
      missing.each { |t| Google::LOGGER.info "  - #{t}" }
      raise "The following tests are missing from the matrix: #{missing}"
    end

    def fail_match_expectations_if_extra(extra)
      Google::LOGGER.info \
        'FATAL: The following tests are not defined in the matrix:'
      extra.each { |t| Google::LOGGER.info "  - #{t}" }
      raise "The following tests are not defined in the matrix: #{extra}"
    end

    def verify_pop(data)
      raise "Cannot pop without a push for #{data}" if @hierarchy.empty?
      raise "Unexpected pop for #{data}. Expecting pop for #{@hierarchy.last}" \
        unless @hierarchy.last == data
    end

    def update_level(data, bundled)
      @level = data.values.reject { |t| t == :none || t == :ignore }.size
      @level -= 1 if bundled
    end
  end
end

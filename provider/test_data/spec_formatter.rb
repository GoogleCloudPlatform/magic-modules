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

require 'provider/test_data/formatter'
require 'provider/test_data/generator'

module Provider
  module TestData
    # Represents the test data
    class TestData
      def initialize(block)
        @block = block
      end

      # Creates a YAML representation of the test data
      # All keys will be sorted according to sort_by_manifest_key
      def to_yaml(opts = {})
        YAML.quick_emit(@block, opts) do |out|
          out.map(nil, to_yaml_style) do |map|
            @block.keys.sort_by { |k| sort_by_manifest_key(k) }
                  .each do |k|
              map.add(k, @block[k])
            end
          end
        end
      end

      def sort_by_manifest_key(key)
        if key == 'kind'
          'aaa_kind' # forces key to be the first item
        elsif key == 'name'
          'aab_name' # forces name to be the second item
        elsif key == 'id'
          'aac_id' # forces id to be the third item
        elsif key == 'selfLink'
          'zzz_selfLink' # forces selfLink to be the last item
        else
          key
        end
      end
    end
    # Creates the network data YAML for unit tests
    class SpecFormatter < Formatter
      def initialize(provider)
        super(Provider::TestData::Generator.new)
        @provider = provider
      end

      # rubocop:disable Metrics/MethodLength
      def generate(object, _, kind, seed, extra)
        props = object.all_user_properties
        name = Google::StringUtils.underscore(object.name)

        extra = extra.merge(
          project: "'test project\##{seed} data'",
          selfLink: "selflink(resource(#{name},#{seed}))",
          name: extra[:name]
        )

        extra = extra.merge(kind: kind) if object.kind?

        props_nodot = props.reject { |p| p.name.include? '.' }
        block = emit_manifest_block(props_nodot, seed, extra, {})

        # TODO(alexstephen): Delete once deprecated is a real nested object.
        props_with_dot = props.select { |p| p.name.include? '.' }
        unless props_with_dot.empty?
          block.merge!(emit_fake_nested_block(props_with_dot,
                                              seed, {}, {}))
        end
        TestData.new(block)
      end
      # rubocop:enable Metrics/MethodLength

      def emit_manifest_block(props, seed, extra, ctx)
        props.map do |p|
          emit_manifest_assign(p, seed, ctx, p.method(:field_name), true)
        end
             .to_h
             .merge(extra.map { |k, v| [k.id2name, v] }.to_h)
             .to_h
      end

      # TODO(alexstephen): Delete once deprecated is a real nested object.
      # Takes nested objects using a layer1.layer2 notation and converts to a
      # Nested Object style hash.
      def emit_fake_nested_block(props, seed, extra, ctx)
        unnested = props.map do |p|
          emit_manifest_assign(p, seed, ctx, p.method(:field_name), true)
        end
                        .to_h
                        .merge(extra.map { |k, v| [k.id2name, v] }.to_h)
                        .to_h
        nested = {}
        unnested.each do |k, v|
          parts = k.split('.')
          nested[parts[0]] = {} unless nested.key?(parts[0])
          nested[parts[0]][parts[1]] = v
        end
        nested
      end

      # Emits a name value
      # Should be a standard Hash, unlike other formatters.
      def emit_namevalues(prop, seed, _ctx)
        @datagen.value(prop.class, prop, seed)
      end

      private

      def quote_string(s)
        @provider.quote_string(s)
      end

      def formatter(type, value)
        raise "Unknown type '#{type}'" unless handlers.key?(type)
        handlers[type].call(value)
      end

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/MethodLength
      def handlers
        {
          Api::Type::Boolean => ->(v) { v },
          Api::Type::Constant => ->(v) { v },
          Api::Type::Double => ->(v) { v },
          Api::Type::Enum =>
            ->(v) { v.is_a?(Symbol) ? v.id2name : v.to_s },
          Api::Type::Integer => ->(v) { v },
          Api::Type::String => ->(v) { v },
          Api::Type::Time => ->(v) { v.iso8601 },
          Api::Type::Array => ->(v) { v },
          Api::Type::NestedObject => ->(v) { v },
          Api::Type::ResourceRef => ->(v) { v },
          Api::Type::Array::STRING_ARRAY_TYPE =>
            ->(v) { v.map { |e| e } },
          Api::Type::Array::RREF_ARRAY_TYPE =>
            ->(v) { v.call(exported_values: true) },
          Api::Type::NameValues => ->(v) { v }
        }.freeze
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/MethodLength
    end
  end
end

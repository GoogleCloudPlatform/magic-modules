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
require 'provider/test_data/formatter'

module Provider
  # Formats objects into Puppet manifests that are puppet-lint compliant.
  class PuppetTestManifestFormatter < Provider::TestData::Formatter
    def initialize(provider)
      super(Provider::TestData::Generator.new)
      @provider = provider
    end

    # Will generate a full Puppet manifest for an entire dependency graph.
    # Parameters:
    # graph - A DependencyGraph
    # base_name - The name of the object that has all of the dependencies. This
    #             ensures that base object has all of its properties outputted.
    def generate_all_objects(graph, base_name, kind, extra)
      graph.map do |object|
        if object.object.name == base_name
          generate_object(object.object, "title#{object.seed}", kind,
                          object.seed, extra)
        else
          generate_ref(object.object, object.seed)
        end
      end
    end

    # Will generate a Puppet block for a specific object
    def generate_object(object, title, kind, seed, extra)
      props = select_properties(object.all_user_properties, kind, extra)

      extra = {
        project: "'test project\##{seed} data'",
        credential: "'cred#{seed}'"
      }.merge(extra)

      # Puppet does not like when readonly resources have an ensure property
      extra.delete(:ensure) if object.readonly

      [
        "#{object.out_name} { '#{title}':",
        @provider.indent_list(
          emit_manifest_block(props, seed, extra, {}), 2, true
        ),
        '}'
      ]
    end

    # Generates a resource block for a resource ref.
    # Requires the ResourceRef and an index.
    def generate_ref(ref, index)
      ref_name = ref.name.underscore
      generate_object(ref, "resource(#{ref_name},#{index})", :resource,
                      index, ensure: 'present')
    end

    def emit_manifest_block(props, seed, extra, ctx)
      max_key = max_key_length(props, extra)

      props.map { |p| emit_manifest_assign(p, seed, ctx, p.method(:out_name)) }
           .to_h
           .merge(extra.map { |k, v| [k.id2name, v] }.to_h)
           .sort_by { |k, _v| sort_by_manifest_key(k) }
           .map do |k, v|
             [k, ' ' * [0, max_key - k.length].max, ' => ',
              [v].flatten.join("\n")].join
           end
    end

    private

    def sort_by_manifest_key(key)
      if key.start_with?('ensure')
        'aaa_ensure' # forces ensure to be the first item
      elsif key.start_with?('credential')
        'zzz_credential' # forces credential to be the last item
      elsif key.start_with?('project')
        'zzy_project' # forces project to be the next-to-last item
      else
        key
      end
    end

    def formatter(type, value)
      raise "Unknown type '#{type}'" unless handlers.key?(type)
      handlers[type].call(value)
    end

    def format_values(start, values, stop)
      [start, indent_list(values, 2), stop].join("\n")
    end

    # rubocop:disable Metrics/AbcSize
    def handlers
      {
        Api::Type::Boolean => ->(v) { v ? 'true' : 'false' },
        Api::Type::Constant => ->(v) { quote_string(v.upcase) },
        Api::Type::Double => ->(v) { v },
        Api::Type::Enum =>
          ->(v) { quote_string(v.is_a?(Symbol) ? v.id2name : v.to_s) },
        Api::Type::Integer => ->(v) { v },
        Api::Type::String => ->(v) { quote_string(v) },
        Api::Type::Time => ->(v) { quote_string(v.iso8601) },
        Api::Type::Array => ->(v) { format_values('[', v, ']') },
        Api::Type::NestedObject => ->(v) { format_values('{', v, '}') },
        Api::Type::ResourceRef => ->(v) { quote_string(v) },
        Api::Type::Array::STRING_ARRAY_TYPE =>
          ->(v) { ['[', v.map { |e| quote_string(e) }.join(', '), ']'].join },
        Api::Type::Array::RREF_ARRAY_TYPE =>
          ->(v) { ['[', v.call(exported_values: false).join(', '), ']'].join },
        Api::Type::NameValues => ->(v) { format_values('{', v, '}') }
      }.freeze
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize

    def quote_string(s)
      @provider.quote_string(s)
    end

    def indent_list(list, indent)
      @provider.indent_list(list, indent, true)
    end

    def max_key_length(props, extra)
      max_key_prop = props.max_by { |p| p.out_name.length }
      max_key = max_key_prop.nil? ? 0 : max_key_prop.out_name.length
      return max_key if extra.empty?
      [max_key, extra.max_by { |k, _v| k.length }.first.length].max
    end
  end
end

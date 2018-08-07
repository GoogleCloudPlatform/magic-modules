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
  # rubocop:disable Metrics/ClassLength
  # Formats objects into Puppet manifests that are puppet-lint compliant.
  class ChefTestCatalogFormatter < Provider::TestData::Formatter
    def initialize(provider)
      super(Provider::TestData::Generator.new)
      @provider = provider
    end

    def generate_all_objects(objects, base_name, kind, extra)
      objects.map do |object|
        if object.object.name == base_name
          generate_object(object.object, "title#{object.seed}", kind,
                          object.seed, extra)
        else
          generate_ref(object.object, object.seed)
        end
      end
    end

    def generate_object(object, title, kind, seed, extra)
      props = select_properties(object.all_user_properties, kind, extra)

      extra = {
        project: "'test project\##{seed} data'",
        credential: "'mycred'"
      }.merge(extra)

      [
        "#{object.out_name} '#{title}' do",
        indent_array(emit_manifest_block(props, seed, extra, {},
                                         first_level: true), 2),
        'end'
      ]
    end

    # Generates a resource block for a resource ref.
    # Requires the ResourceRef and an index.
    def generate_ref(ref, index)
      ref_name = ref.name.underscore
      generate_object(ref, "resource(#{ref_name},#{index})", :resource,
                      index, action: ':create')
    end

    # Generates a block of Chef recipe code
    # Valid options:
    #   first level: Says if this is the first level being generated.
    # rubocop:disable Metrics/AbcSize
    def emit_manifest_block(props, seed, extra, ctx, opts = {})
      manifest = props.map do |p|
        method = -> { @provider.property_out_name(p) }
        emit_manifest_assign(p, seed, ctx, method, opts)
      end
      manifest.to_h
              .merge(extra.map { |k, v| [k.id2name, v] }.to_h)
              .sort_by { |k, _v| sort_by_manifest_key(k) }
              .map do |k, v|
                if v.is_a?(String) && v[0] == '{' && opts[:first_level]
                  # First level nested objects need a () instead of {}
                  [k, '(', [v].flatten.join("\n"), ')'].join
                else
                  [k, ' ', [v].flatten.join("\n")].join
                end
              end
    end
    # rubocop:enable Metrics/AbcSize

    # Generates a key value pair for a property depending on its type
    # Valid options:
    #   first_level: Says if this is the first level being generated
    def emit_manifest_assign(prop, seed, ctx, prop_api_name, opts = {})
      # Chef name field must use label_name
      if prop.name == 'name' && opts[:first_level]
        [@provider.label_name(prop.__resource),
         formatter(prop.class, @datagen.value(prop.class, prop, seed))]
      else
        super(prop, seed, ctx, prop_api_name)
      end
    end

    def sort_by_manifest_key(key)
      if key.start_with?('action')
        '_aaa_action' # forces ensure to be the first item
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

    def format_values(start, values, stop, commas = false)
      if commas
        [start, indent_array(values, 2), stop].join("\n")
      else
        [start, @provider.indent_list(values, 2), stop].join("\n")
      end
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
        Api::Type::NameValues => ->(v) { format_values('{', v, '}') },
        Api::Type::ResourceRef => ->(v) { quote_string(v) },
        Api::Type::Array::STRING_ARRAY_TYPE =>
          ->(v) { ['[', v.map { |e| quote_string(e) }.join(', '), ']'].join },
        Api::Type::Array::RREF_ARRAY_TYPE =>
          ->(v) { ['[', v.call(exported_values: false).join(', '), ']'].join }
      }.freeze
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/MethodLength

    def quote_string(s)
      @provider.quote_string(s)
    end

    def indent_array(list, indent)
      @provider.indent_array(list, indent).join("\n")
    end

    def emit_nested(prop, seed, ctx)
      props = prop.properties
      props.map { |p| emit_manifest_assign(p, seed, ctx, p.method(:out_name)) }
           .to_h
           .sort_by { |k, _v| sort_by_manifest_key(k) }
           .map do |k, v|
        ["#{k}:", [v].flatten.join("\n")].join(' ')
      end
    end
  end
  # rubocop:enable Metrics/ClassLength
end

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
  # Formats objects into README-style outlines that are puppet-lint compliant.
  class ChefOutline
    def initialize(provider)
      @provider = provider
    end

    def generate(object)
      extra = {
        project: 'string',
        credential: 'reference to gauth_credential'
      }

      [
        "#{object.out_name} 'id-for-resource' do",
        @provider.indent_array(emit_manifest_block(object.all_user_properties,
                                                   extra), 2),
        'end'
      ]
    end

    private

    def emit_manifest_block(props, extra)
      max_key = max_key_length(props, extra)

      props.map { |p| formatter(p) }
           .to_h
           .merge(extra.map { |k, v| [k.id2name, v] }.to_h)
           .sort_by { |k, _v| sort_by_manifest_key(k) }
           .map do |k, v|
             [k, ' ' * [0, max_key - k.length].max, ' ',
              v].join
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

    def formatter(p)
      return p if p.is_a? String
      [p.out_name, handlers.fetch(p.class, ->(v) { v.type.downcase }).call(p)]
    end

    def format_values(start, values, stop, last_comma = true)
      [start, indent_list(values, 2, last_comma), stop].join("\n")
    end

    def handlers
      {
        Api::Type::NestedObject => ->(v) { emit_nested(v) },
        Api::Type::Enum => ->(v) { emit_enum(v) },
        Api::Type::Array => ->(v) { emit_array(v) },
        Api::Type::ResourceRef => ->(v) { emit_resourceref(v) }
      }.freeze
    end

    # rubocop:disable Metrics/AbcSize
    def emit_resourceref(p)
      if p.resources.length > 1
        list = p.resources.first(p.resources.size - 1).map do |x|
          x.resource_ref.out_name
        end.join(' ,')
        "reference to #{list} or #{p.resources.last.resource_ref.out_name}"
      else
        "reference to a #{p.resources.first.resource_ref.out_name}"
      end
    end

    # rubocop:enable Metrics/AbcSize
    #
    def emit_enum(p)
      return 'Enum' if p.values.empty?
      return p.values[0].to_s if p.values.length == 1

      values = p.values.map { |val| "'#{val}'" }

      "#{values.first(values.size - 1).join(', ')} or #{values.last}"
    end

    def emit_nested(p)
      format_values('{',
                    emit_manifest_block(p.properties, {}),
                    '}')
    end

    def emit_array(p)
      item = if p.item_type.is_a? Api::Type::NestedObject
               emit_nested(p.item_type)
             elsif p.item_type.is_a? Api::Type::ResourceRef
               emit_resourceref(p.item_type)
             else
               p.item_type.split('::').last.downcase
             end

      format_values('[', [item, '...'], ']', false)
    end

    def quote_string(s)
      @provider.quote_string(s)
    end

    def indent_list(list, indent, last_comma = true)
      @provider.indent_list(list, indent, last_comma)
    end

    def max_key_length(props, extra)
      max_key_prop = props.max_by { |p| p.out_name.length }
      max_key = max_key_prop.nil? ? 0 : max_key_prop.out_name.length
      return max_key if extra.empty?
      [max_key, extra.max_by { |k, _v| k.length }.first.length].max
    end
  end
end

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

require 'google/ruby_utils'
require 'provider/config'
require 'provider/core'
require 'provider/chef/manifest'
require 'provider/chef/property_override'
require 'provider/chef/resource_override'
require 'provider/chef/test_catalog'
require 'provider/chef/resource_override'
require 'provider/chef/property_override'
require 'provider/test_data/utils'

module Provider
  # Code generator for Chef Cookbooks that manage Google Cloud Platform
  # resources.
  class Chef < Provider::Core
    include Google::RubyUtils
    include Provider::TestData::TestUtils

    RESERVED_WORDS = %w[deprecated updated].freeze
    TEST_FOLDER = 'recipes'.freeze

    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      attr_reader :operating_systems
      # TODO(alexstephen): Convert this to a regular function generator
      # like Puppet.
      attr_reader :functions

      def provider
        Provider::Chef
      end

      def resource_override
        Provider::Chef::ResourceOverride
      end

      def property_override
        Provider::Chef::PropertyOverride
      end

      def validate
        super
        check_optional_property :manifest, Provider::Chef::Manifest
        check_property_list \
          :operating_systems, Provider::Config::OperatingSystem
      end
    end

    # A custom client side function for Chef
    class Function < Provider::Config::Function
      attr_reader :search_paths

      def validate
        super
        check_property_list :search_paths, Provider::Chef::SearchPath
      end
    end

    # A search path for client side functions in Chef.
    class SearchPath < Api::Object::Named
      attr_reader :path
      attr_reader :comment
    end

    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def prop_decl(prop)
      return 'kind_of: [TrueClass, FalseClass]' if prop.type == 'Boolean'
      return 'Float' if prop.type == 'Double'

      # Chef manifest files show Nested Objects as Hash, or Google::Property
      # They will be immediately coerced into a Google::Property in the end
      return "[Hash, ::#{prop.property_type.gsub(':Property:', ':Data:')}]" \
        if prop.is_a? Api::Type::NestedObject

      return "[String, ::#{prop.property_type.gsub(':Property:', ':Data:')}]" \
        if prop.is_a? Api::Type::ResourceRef

      return 'Array' if prop.is_a? Api::Type::Array

      if prop.type == 'Enum'
        return format([
          ["equal_to: %w[#{prop.values.join(' ')}]"],
          ((1..prop.values.length - 1).to_a.map do |i|
            ["equal_to: %w[#{prop.values.slice(0, i).join(' ')}",
             indent("#{prop.values.slice(i, prop.values.length).join(' ')}]",
                    13)]
          end)
        ].flatten(1), 0, 16)
      end

      return "[Hash, ::#{prop.property_type}]" \
        if prop.is_a? Api::Type::NameValues

      return ['::', prop.property_type].join \
        unless prop.is_a?(Api::Type::Primitive)

      prop.type
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity

    def label_name(product)
      return product.label_override unless product.label_override.nil?
      Google::StringUtils.underscore(product.name)
                         .split('_')
                         .map { |x| x[0] }
                         .join
                         .concat('_label')
    end

    # Returns a list of all resource types being tested
    # ChefSpec requires this list to include all ResourceRefs
    # rubocop:disable Metrics/AbcSize
    def step_into_list(object, indent, start_indent)
      props = test_resourcerefs_for_properties(object.all_user_properties,
                                               object)
      refs = props.map { |x| x.resources.first.resource_ref }
      steps = [object.out_name].concat(refs.map(&:out_name).reverse).uniq

      return indent("step_into: '#{steps[0]}',", indent) if steps.length == 1

      format(
        [
          ["step_into: %w[#{steps.join(' ')}],"],
          ["step_into: %w[#{steps.slice(0..-2).join(' ')}",
           indent("#{steps.last(1).join(' ')}],", 14)], # 14 = step_into: %w[
          ["step_into: %w[#{steps.slice(0..-3).join(' ')}",
           indent("#{steps.last(2).join(' ')}],", 14)], # 14 = step_into: %w[
          [
            "step_into: %w[#{steps[0]}",
            indent(steps.slice(1..-2), 14), # 14 = step_into: %w[
            indent("#{steps.last(1).join}],", 14)
          ]
        ], indent, start_indent
      )
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/MethodLength

    def generate_user_agent(product, file_name)
      emit_user_agent(
        product, nil,
        [
          'TODO(alexstephen): Check how to get the original Chef user agent.',
          'TODO(alexstephen): Check how to fetch cookbook version.'
        ],
        file_name
      )
    end

    def emit_coerce(product_ns, class_name, spaces_used = 0)
      type = "::Google::#{product_ns}::Property::#{class_name}"
      lines(format([
                     [
                       'def self.coerce',
                       indent("->(x) { #{type}.catalog_parse(x) }", 2),
                       'end'
                     ],
                     [
                       'def self.coerce',
                       indent('lambda do |x|', 2),
                       indent(indent("#{type}.catalog_parse(x)", 2), 2),
                       indent('end', 2),
                       'end'
                     ],
                     [
                       'def self.coerce',
                       indent('lambda do |x|', 2),
                       indent("type = #{type}", 4),
                       indent('type.catalog_parse(x)', 4),
                       indent('end', 2),
                       'end'
                     ]
                   ], spaces_used), 1)
    end
    # rubocop:enable Metrics/MethodLength

    def property_out_name(prop)
      if prop.name == 'name'
        label_name(prop.__resource)
      elsif RESERVED_WORDS.include?(prop.out_name)
        "_#{prop.out_name}" # avoid conflicts w/ reserved words with '_'
      else
        prop.out_name
      end
    end

    def compile_end2end_tests(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          # Tests go into hidden folder because we don't need to expose
          # to regular Chef users.
          ["recipes/tests~#{file}",
           "products/#{@api.prefix[1..-1]}/examples/chef/#{file}"]
        end
      )
    end

    private

    def generate_simple_property(type, data)
      {
        source: File.join('templates', 'chef', 'property', "#{type}.rb.erb"),
        target: File.join('libraries', 'google', data[:product_name],
                          'property', "#{type}.rb")
      }
    end

    def generate_base_property(data) end

    def emit_resourceref_object(data)
      target = data[:property].property_file
      {
        source: File.join('templates', 'chef', 'property',
                          'resourceref.rb.erb'),
        target: "libraries/#{target}.rb",
        overrides: data.clone.merge(
          class_name: data[:property].property_class.last
        )
      }
    end

    def generate_typed_array(data, prop)
      type = Module.const_get(prop.item_type).new(prop.name).type
      file = Google::StringUtils.underscore(type)
      prop_map = []
      prop_map << {
        source: File.join('templates', 'chef', 'property',
                          'array_typed.rb.erb'),
        target: File.join('libraries', 'google', data[:product_name],
                          'property', "#{file}_array.rb"),
        overrides: { type: type }
      }
      prop_map << generate_base_array(data)
      prop_map
    end

    def generate_base_array(data)
      {
        source: File.join('templates', 'chef', 'property', 'array.rb.erb'),
        target: File.join('libraries', 'google', data[:product_name],
                          'property', 'array.rb')
      }
    end

    def emit_nested_object(data)
      target = if data[:emit_array]
                 data[:property].item_type.property_file
               else
                 data[:property].property_file
               end
      {
        source: File.join('templates', 'chef', 'property',
                          'nested_object.rb.erb'),
        target: "libraries/#{target}.rb",
        overrides: emit_nested_object_overrides(data)
      }
    end

    def emit_nested_object_overrides(data)
      data.clone.merge(
        field_name: Google::StringUtils.camelize(data[:field], :upper),
        object_type: Google::StringUtils.camelize(data[:obj_name], :upper),
        product_ns: Google::StringUtils.camelize(data[:product_name], :upper),
        class_name: if data[:emit_array]
                      data[:property].item_type.property_class.last
                    else
                      data[:property].property_class.last
                    end
      )
    end

    def generate_resource(data)
      target_folder = File.join(data[:output_folder], 'resources')
      FileUtils.mkpath target_folder
      name = Google::StringUtils.underscore(data[:object].name)
      generate_resource_file data.clone.merge(
        default_template: provider_template_source(data),
        out_file: File.join(target_folder, "#{name}.rb")
      )
    end

    def provider_template_source(data)
      if data[:object].manual
        object_name = Google::StringUtils.underscore(data[:object].name)
        File.join('products', data[:product_name], 'files',
                  "provider~chef~#{object_name}.rb")
      else
        'templates/chef/resource.erb'
      end
    end

    def generate_resource_tests(data)
      target_folder = File.join(data[:output_folder], 'spec')
      FileUtils.mkpath target_folder
      name = Google::StringUtils.underscore(data[:object].name)
      generate_resource_file data.clone.merge(
        default_template: 'templates/chef/resource_spec.erb',
        out_file: File.join(target_folder, "#{name}_spec.rb")
      )
    end

    def compile_examples(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          ["recipes/examples~#{file}",
           "products/#{@api.prefix[1..-1]}/examples/chef/#{file}"]
        end
      )
    end

    def google_lib_basic(file, product_ns)
      google_lib_basic_files(file, product_ns, 'libraries', 'google')
    end

    def google_lib_network(file, product_ns)
      google_lib_network_files(file, product_ns, 'libraries', 'google')
    end

    def example_resource_name_prefix
      'chef-e2e-'
    end

    def test_file?(file)
      file.include? 'tests~'
    end

    # Builds the properties for a nested object of any depth
    # This returns an arrays of strings that represent Markdown formatted
    # properties for the nested object and all nested objects beneath it
    # Requires:
    #  prop: A property of type nested object.
    #  current_path: A string representing all layers above this current
    #                property.  This string will usually be the output names of
    #                all properties above the current property joined by
    #                '/' (ex. first_level/second_level) or an array denoted
    #                by [] (ex. array_of_nested_props[])
    # rubocop:disable Metrics/AbcSize
    def build_nested_object(prop, current_path)
      object_lines = []
      prop.properties.each do |nested_prop|
        next_path = "#{current_path}/#{nested_prop.out_name}"
        object_lines << lines(["* `#{next_path}`"].join(' '))

        object_lines << lines(wrap_field([
          ('Required.' if nested_prop.required),
          ('Output only.' if nested_prop.output),
          nested_prop.description
        ].compact.join(' '), 0), 1)

        if nested_prop.is_a? Api::Type::NestedObject
          object_lines.concat(build_nested_object(nested_prop, next_path))
        elsif nested_prop.is_a?(Api::Type::Array) &&
              nested_prop.item_type.is_a?(Api::Type::NestedObject)
          object_lines.concat(build_nested_object(nested_prop.item_type,
                                                  "#{next_path}[]"))
        end
      end
      object_lines
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/MethodLength

    # Emits all the Chef client functions available for use by end users.
    def generate_client_function(output_folder, fn)
      target_folder = File.join(output_folder, 'libraries', 'google',
                                'functions')
      {
        fn: fn,
        target_folder: target_folder,
        template: 'templates/chef/function.erb',
        output_folder: output_folder,
        out_file: File.join(target_folder, "#{fn.name}.rb")
      }
    end
  end
end

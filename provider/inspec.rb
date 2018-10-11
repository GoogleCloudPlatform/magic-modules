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
require 'provider/inspec/manifest'
require 'provider/inspec/property_override'
require 'provider/inspec/resource_override'
require 'provider/inspec/test_catalog'
require 'provider/inspec/resource_override'
require 'provider/inspec/property_override'

module Provider
  # Code generator for Example Cookbooks that manage Google Cloud Platform
  # resources.
  class Inspec < Provider::Core

    RESERVED_WORDS = %w[deprecated updated].freeze
    include Google::RubyUtils
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      def provider
        Provider::Inspec
      end

      def resource_override
        Provider::Inspec::ResourceOverride
      end

      def property_override
        Provider::Inspec::PropertyOverride
      end

    end

    private


    def label_name(product)
      product.name.underscore
             .split('_')
             .map { |x| x[0] }
             .join
             .concat('_label')
    end


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

    # This function uses the resource.erb template to create one file
    # per resource. The resource.erb template forms the basis of a single
    # GCP Resource on Example.
    def generate_resource(data)
      target_folder = File.join(data[:output_folder], 'inspec')
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/resource.erb',
        out_file: File.join(target_folder, "google_#{data[:product_name]}_#{name}.rb")
      )
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/resource_plural.erb',
        out_file: File.join(target_folder, "google_#{data[:product_name]}_#{name}s.rb")
      )
    end

    # This function would generate unit tests using a template
    def generate_resource_tests(data) end

    # This function would automatically generate the files used for verifying
    # network calls in unit tests. If you comment out the following line,
    # a bunch of YAML files will be created under the spec/ folder.
    def generate_network_datas(data, object) end

    def generate_typed_array(data, prop)
      type = Module.const_get(prop.item_type).new(prop.name).type
      file = type.underscore
      prop_map = []
      prop_map << {
        source: File.join('templates', 'inspec', 'property',
                          'array_typed.rb.erb'),
        target: File.join('libraries', 'google', 'property', "#{file}_array.rb"),
        overrides: { type: type }
      }
      prop_map << generate_base_array(data)
      prop_map
    end

    def generate_base_array(data)
      {
        source: File.join('templates', 'inspec', 'property', 'array.rb.erb'),
        target: File.join('libraries', 'google', 'property', 'array.rb')
      }
    end

    def generate_user_agent(product, file_name)
      emit_user_agent(
        product, nil,
        ['TODO(nelsonjr): Check how to fetch module version.'],
        file_name
      )
    end

    def generate_requires(properties, requires = [])
      nested_props = properties.select{ |type| nested_object?(type) }
      requires.concat(properties.reject{ |type| primitive?(type) || resource_ref?(type) || nested_object?(type) }.collect(&:requires))
      requires.concat(nested_props.map{|nested_prop| generate_requires(nested_prop.properties) } )
      requires.concat(nested_props.map{|nested_prop| nested_prop.property_file })
      requires
    end

    # InSpec doesn't need wrappers for primitives, so exclude them
    def emit_requires(requires)
      primitives = ['boolean', 'enum', 'string', 'time', 'integer', 'array', 'string_array', 'double']
      requires.flatten.sort.uniq.reject{|r| primitives.include?(r.split('/').last)}.map { |r| "require '#{r}'" }.join("\n")
    end

    def primitive? (property) 
      return property.is_a?(::Api::Type::Primitive) || (property.is_a?(Api::Type::Array) && !property.item_type.is_a?(::Api::Type::NestedObject))
    end

    def resource_ref? (property) 
      return property.is_a?(::Api::Type::ResourceRef)
    end

    def nested_object? (property) 
      return property.is_a?(::Api::Type::NestedObject)
    end

    def typed_array? (property) 
      return property.is_a?(::Api::Type::Array)
    end

    def google_lib_basic(file, product_ns)
      google_lib_basic_files(file, product_ns, 'inspec', 'google')
    end

    def google_lib_network(file, product_ns)
      google_lib_network_files(file, product_ns, 'inspec', 'google')
    end

    def property_out_name(prop)
      prop.out_name
    end

    def emit_link_plural(url_parts, base_url)
      plural_url = ['URI.join(',
       indent([quote_string(url_parts[0]) + ',',
               'expand_variables(',
               indent(format_expand_variables(base_url), 2),
               indent('data', 2),
               ')'], 2),
       ')'].join("\n")
      
      (params, fn_args) = emit_link_var_args(plural_url, false)
      code = ["def self.self_link(#{fn_args})",
              indent(plural_url, 2),
              'end']
      code.join("\n")
    end

    # We build a lot of property classes to help validate + coerce types.
    # The following functions would generate all of these properties.
    # Some of these property classes help us handle Strings, Times, etc.
    #
    # Others (nested objects) ensure that all Hashes contain proper values +
    # types for its nested properties.
    #
    # ResourceRef properties help ensure that links between different objects
    # (Addresses + Instances for example) work properly, are abstracted away,
    # and don't require the user to have a large knowledge base of how GCP
    # works.
    # rubocop:disable Layout/EmptyLineBetweenDefs
    def generate_base_property(data) end

    def generate_simple_property(type, data)
      
    end


    def emit_nested_object(data)
      target = if data[:emit_array]
                 data[:property].item_type.property_file
               else
                 data[:property].property_file
               end
      {
        source: File.join('templates', 'inspec', 'property',
                          'nested_object.rb.erb'),
        target: "libraries/#{target}.rb",
        overrides: emit_nested_object_overrides(data)
      }
    end

    def emit_nested_object_overrides(data)
      data.clone.merge(
        api_name: data[:api_name].camelize(:upper),
        object_type: data[:obj_name].camelize(:upper),
        product_ns: data[:product_name].camelize(:upper),
        class_name: if data[:emit_array]
                      data[:property].item_type.property_class.last
                    else
                      data[:property].property_class.last
                    end
      )
    end

    def emit_resourceref_object(data)
    end

    # rubocop:enable Layout/EmptyLineBetweenDefs
  end
end

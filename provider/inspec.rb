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
require 'overrides/inspec/resource_override'
require 'overrides/inspec/property_override'
require 'active_support/inflector'

module Provider
  # Code generator for Example Cookbooks that manage Google Cloud Platform
  # resources.
  class Inspec < Provider::Core
    include Google::RubyUtils
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      def provider
        Provider::Inspec
      end

      def resource_override
        Overrides::Inspec::ResourceOverride
      end

      def property_override
        Overrides::Inspec::PropertyOverride
      end
    end

    # This function uses the resource templates to create singular and plural
    # resources that can be used by InSpec
    def generate_resource(data)
      target_folder = File.join(data[:output_folder], 'libraries')
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/singular_resource.erb',
        out_file: File.join(target_folder, "google_#{data[:product].api_name}_#{name}.rb")
      )
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/plural_resource.erb',
        out_file: \
          File.join(target_folder, "google_#{data[:product].api_name}_#{name}".pluralize + '.rb')
      )
      generate_documentation(data, name, false)
      generate_documentation(data, name, true)
      generate_properties(data, data[:object].all_user_properties)
    end

    def generate_properties(data, props)
      nested_objects = props.select(&:nested_properties?)
      return if nested_objects.empty?

      # Create property files for any nested objects.
      prop_map = nested_objects.map { |nested_object| emit_nested_object(data, nested_object) }
      generate_property_files(prop_map, data)

      # Create property files for any deeper nested objects.
      nested_objects.each { |prop| generate_properties(data, prop.nested_properties) }
    end

    # Generate the files for the properties
    def generate_property_files(prop_map, data)
      prop_map.flatten.compact.each do |prop|
        compile_file_list(
          data[:output_folder],
          { prop[:target] => prop[:source] },
          prop
        )
      end
    end

    # Generates InSpec markdown documents for the resource
    def generate_documentation(data, base_name, plural)
      docs_folder = File.join(data[:output_folder], 'docs', 'resources')

      name = plural ? base_name.pluralize : base_name
      generate_resource_file data.clone.merge(
        name: name,
        plural: plural,
        doc_generation: true,
        default_template: 'templates/inspec/doc_template.md.erb',
        out_file: File.join(docs_folder, "google_#{data[:product].api_name}_#{name}.md")
      )
    end

    # Format a url that may be include newlines into a single line
    def format_url(url)
      return url.join('') if url.is_a?(Array)

      url.split("\n").join('')
    end

    # Copies InSpec tests to build folder
    def generate_resource_tests(data)
      target_folder = File.join(data[:output_folder], 'test')
      FileUtils.mkpath target_folder

      FileUtils.cp_r 'templates/inspec/tests/.', target_folder

      name = "google_#{data[:product].api_name}_#{data[:object].name.underscore}"

      generate_inspec_test(data, name, target_folder, name)

      # Build test for plural resource
      generate_inspec_test(data, name.pluralize, target_folder, name)
    end

    def generate_inspec_test(data, name, target_folder, attribute_file_name)
      generate_resource_file data.clone.merge(
        name: name,
        attribute_file_name: attribute_file_name,
        doc_generation: false,
        default_template: 'templates/inspec/integration_test_template.erb',
        out_file: File.join(
          target_folder,
          'integration/verify/controls',
          "#{name}.rb"
        )
      )
    end

    def emit_nested_object(property)
      target = if property.is_a?(Api::Type::Array)
                 property.item_type.property_file
               else
                 property.property_file
               end
      {
        source: File.join('templates', 'inspec', 'nested_object.erb'),
        target: "libraries/#{target}.rb",
        property: property
      }
    end

    def time?(property)
      property.is_a?(::Api::Type::Time)
    end

    # Figuring out if a property is a primitive ruby type is a hassle. But it is important
    # Fingerprints are strings, KeyValuePairs and Maps are hashes, and arrays of primitives are
    # arrays. Arrays of NestedObjects need to have their contents parsed and returned in an array
    # ResourceRefs are strings
    def primitive?(property)
      array_primitive = (property.is_a?(Api::Type::Array)\
        && !property.item_type.is_a?(::Api::Type::NestedObject))
      property.is_a?(::Api::Type::Primitive)\
        || array_primitive\
        || property.is_a?(::Api::Type::KeyValuePairs)\
        || property.is_a?(::Api::Type::Map)\
        || property.is_a?(::Api::Type::Fingerprint)\
        || property.is_a?(::Api::Type::ResourceRef)
    end

    # Arrays of nested objects need special requires statements
    def typed_array?(property)
      property.is_a?(::Api::Type::Array) && nested_object?(property.item_type)
    end

    def nested_object?(property)
      property.is_a?(::Api::Type::NestedObject)
    end

    # Only arrays of nested objects and nested object properties need require statements
    # for InSpec. Primitives are all handled natively
    def generate_requires(properties)
      nested_props = properties.select(&:nested_properties?)

      # Need to include requires statements for the requirements of a nested object
      nested_prop_requires = nested_props.map do |nested_prop|
        generate_requires(nested_prop.nested_properties) unless nested_prop.is_a?(Api::Type::Array)
      end.compact
      nested_object_requires = nested_props.map\
        { |nested_object| nested_object_requires(nested_object) }
      nested_object_requires + nested_prop_requires
    end

    def nested_object_requires(nested_object_type)
      File.join(
        'google',
        nested_object_type.__resource.__product.api_name,
        'property',
        [nested_object_type.__resource.name, nested_object_type.name.underscore].join('_')
      ).downcase
    end

    def resource_name(object, product_ns)
      "google_#{product_ns.downcase}_#{object.name.underscore}"
    end

    def sub_property_descriptions(property)
      if nested_object?(property)
        property.properties.map { |prop| markdown_format(prop) }.join("\n\n") + "\n\n"
      elsif typed_array?(property)
        property.item_type.properties.map { |prop| markdown_format(prop) }.join("\n\n") + "\n\n"
      end
    end

    def markdown_format(property)
      "    * `#{property.name}`: #{property.description.split("\n").join(' ')}"
    end

    def grab_attributes
      YAML.load_file('templates/inspec/tests/integration/configuration/mm-attributes.yml')
    end

    # Returns a variable name OR default value for that variable based on
    # defaults from the existing inspec-gcp tests that do not exist within MM
    # Default values are used within documentation to show realistic examples
    def external_attribute(attribute_name, doc_generation = false)
      return attribute_name unless doc_generation

      external_attribute_file = 'templates/inspec/examples/attributes/external_attributes.yml'
      "'#{YAML.load_file(external_attribute_file)[attribute_name]}'"
    end

    # Replaces Google module name within InSpec resources with GoogleInSpec
    # to alleviate module namespace conflicts due to dependencies on
    # Google SDKs
    def inspec_property_type(property)
      property.property_type.sub('Google::', 'GoogleInSpec::')
    end

    # Returns Ruby code that will parse the given property from a hash
    # This is used in several places that need to parse an arbitrary property
    # from a JSON representation
    def parse_code(property, hash_name)
      item_from_hash = "#{hash_name}['#{property.api_name}']"
      return "parse_time_string(#{item_from_hash})" if time?(property)

      if primitive?(property)
        return "name_from_self_link(#{item_from_hash})" \
          if property.name_from_self_link

        return item_from_hash.to_s
      elsif typed_array?(property)
        return "#{inspec_property_type(property)}.parse(#{item_from_hash})"
      end
      "#{inspec_property_type(property)}.new(#{item_from_hash})"
    end
  end
end

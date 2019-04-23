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
require 'overrides/inspec/resource_override'
require 'overrides/inspec/property_override'
require 'active_support/inflector'
require 'google/yaml_validator'

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

    # Subclass of FileTemplate with InSpec specific fields
    class InspecFileTemplate < Provider::FileTemplate
      # Used within doc template to pluralize names
      attr_accessor :plural
      # If this is a file that is being compiled for doc generation
      # This is accessed within the test templates to output an example name
      # for documentation rather than a variable name
      attr_accessor :doc_generation
      # Used to compile InSpec attributes that are used within integration tests
      attr_accessor :attribute_file_name
      # If this is a privileged resource, which will make integration tests unusable
      # unless the user is an admin of the GCP organization
      attr_accessor :privileged
    end

    # Subclass of FileTemplate with InSpec specific fields
    class NestedObjectFileTemplate < Provider::FileTemplate
      # Property to generate this file for
      attr_accessor :property
    end

    # This function uses the resource templates to create singular and plural
    # resources that can be used by InSpec
    def generate_resource(data)
      target_folder = File.join(data.output_folder, 'libraries')
      name = data.object.name.underscore

      data.generate(
        'templates/inspec/singular_resource.erb',
        File.join(target_folder, "google_#{data.product.api_name}_#{name}.rb"),
        self
      )

      data.generate(
        'templates/inspec/plural_resource.erb',
        File.join(target_folder, "google_#{data.product.api_name}_#{name}".pluralize + '.rb'),
        self
      )

      generate_documentation(data.clone, name, false)
      generate_documentation(data.clone, name, true)
      generate_properties(data.clone, data.object.all_user_properties)
    end

    # Generate the IAM policy for this object. This is used to query and test
    # IAM policies separately from the resource itself
    def generate_iam_policy(data)
      target_folder = File.join(data.output_folder, 'libraries')
      name = data.object.name.underscore

      iam_policy_resource_name = "google_#{data.product.api_name}_#{name}_iam_policy"
      data.generate(
        'templates/inspec/iam_policy/iam_policy.erb',
        File.join(target_folder, "#{iam_policy_resource_name}.rb"),
        self
      )

      markdown_target_folder = File.join(data.output_folder, 'docs/resources')
      data.generate(
        'templates/inspec/iam_policy/iam_policy.md.erb',
        File.join(markdown_target_folder, "#{iam_policy_resource_name}.md"),
        self
      )
    end

    def generate_properties(data, props)
      nested_objects = props.select(&:nested_properties?)
      return if nested_objects.empty?

      # Create property files for any nested objects.
      generate_property_files(nested_objects, data)

      # Create property files for any deeper nested objects.
      nested_objects.each { |prop| generate_properties(data, prop.nested_properties) }
    end

    # Generate the files for the properties
    def generate_property_files(properties, data)
      properties.flatten.compact.each do |property|
        nested_object_template = NestedObjectFileTemplate.new(
          data.output_folder,
          data.name,
          data.product,
          data.version,
          data.env
        )
        nested_object_template.property = property
        source = File.join('templates', 'inspec', 'nested_object.erb')
        target = File.join(
          nested_object_template.output_folder,
          "libraries/#{nested_object_requires(property)}.rb"
        )
        nested_object_template.generate(source, target, self)
      end
    end

    def build_object_data(object, output_folder, version)
      InspecFileTemplate.file_for_resource(output_folder, object, @config, version, build_env)
    end

    # Generates InSpec markdown documents for the resource
    def generate_documentation(data, base_name, plural)
      docs_folder = File.join(data.output_folder, 'docs', 'resources')

      name = plural ? base_name.pluralize : base_name

      data.name = name
      data.plural = plural
      data.doc_generation = true
      data.generate(
        'templates/inspec/doc_template.md.erb',
        File.join(docs_folder, "google_#{data.product.api_name}_#{name}.md"),
        self
      )
    end

    # Format a url that may be include newlines into a single line
    def format_url(url)
      return url.join('') if url.is_a?(Array)

      url.split("\n").join('')
    end

    # Copies InSpec tests to build folder
    def generate_resource_tests(data)
      target_folder = File.join(data.output_folder, 'test')
      FileUtils.mkpath target_folder

      FileUtils.cp_r 'templates/inspec/tests/.', target_folder

      name = "google_#{data.product.api_name}_#{data.object.name.underscore}"

      generate_inspec_test(data.clone, name, target_folder, name)

      # Build test for plural resource
      generate_inspec_test(data.clone, name.pluralize, target_folder, name)
    end

    def generate_inspec_test(data, name, target_folder, attribute_file_name)
      data.name = name
      data.attribute_file_name = attribute_file_name
      data.doc_generation = false
      data.privileged = data.object.privileged

      data.generate(
        'templates/inspec/integration_test_template.erb',
        File.join(
          target_folder,
          'integration/verify/controls',
          "#{name}.rb"
        ),
        self
      )
    end

    def emit_requires(requires)
      requires.flatten.sort.uniq.map { |r| "require '#{r}'" }.join("\n")
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
        qualified_property_class(nested_object_type)
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
      "    * `#{property.out_name}`: #{property.description.split("\n").join(' ')}"
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

    def qualified_property_class(property)
      name = property.name.underscore
      other = property.__resource.name
      until property.parent.nil?
        property = property.parent
        next if typed_array?(property)

        name = property.name.underscore + '_' + name
      end

      other + '_' + name
    end

    def modularized_property_class(property)
      class_name = qualified_property_class(property).camelize(:upper)
      product_name = property.__resource.__product.name.camelize(:upper)
      "GoogleInSpec::#{product_name}::Property::#{class_name}"
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
        class_name = modularized_property_class(property.item_type)
        return "#{class_name}Array.parse(#{item_from_hash}, to_s)"
      end
      "#{modularized_property_class(property)}.new(#{item_from_hash}, to_s)"
    end

    # Extracts identifiers of a resource in the form {{identifier}} from a url
    def extract_identifiers(url)
      url.scan(/({{)(\w+)(}})/).map { |arr| arr[1] }
    end
  end
end

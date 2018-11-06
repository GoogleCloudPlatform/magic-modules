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
require 'provider/inspec/resource_override'
require 'provider/inspec/property_override'
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
        Provider::Inspec::ResourceOverride
      end

      def property_override
        Provider::Inspec::PropertyOverride
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
        out_file: File.join(target_folder, "google_#{data[:product_name]}_#{name}.rb")
      )
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/plural_resource.erb',
        out_file: \
          File.join(target_folder, "google_#{data[:product_name]}_#{name}".pluralize + '.rb')
      )
      generate_documentation(data)
    end

    # Generates InSpec markdown documents for the resource
    def generate_documentation(data)
      name = data[:object].name.underscore
      docs_folder = File.join(data[:output_folder], 'docs', 'resources')
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/doc-template.md.erb',
        out_file: File.join(docs_folder, "google_#{data[:product_name]}_#{name}.md")
      )
    end

    # Format a url that may be include newlines into a single line
    def format_url(url)
      return url.join('') if url.is_a?(Array)
      url.split("\n").join('')
    end

    # Copies InSpec unit tests to build folder
    def generate_resource_tests(data)
      target_folder = File.join(data[:output_folder], 'test/unit')
      FileUtils.mkpath target_folder
      FileUtils.cp_r 'templates/inspec/tests/.', target_folder
    end

    def generate_base_property(data) end

    def generate_simple_property(type, data) end

    def generate_typed_array(data, prop) end

    def emit_resourceref_object(data) end

    def generate_network_datas(data, object) end

    def emit_nested_object(data)
      target = if data[:emit_array]
                 data[:property].item_type.property_file
               else
                 data[:property].property_file
               end
      {
        source: File.join('templates', 'inspec', 'nested_object.erb'),
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
      nested_props = properties.select { |type| nested_object?(type) }
      nested_object_arrays = properties.select\
        { |type| typed_array?(type) && nested_object?(type.item_type) }
      nested_array_requires = nested_object_arrays.collect { |type| array_requires(type) }
      # Need to include requires statements for the requirements of a nested object
      # TODO is this needed? Not sure how ruby works so well
      nested_prop_requires = nested_props.map\
        { |nested_prop| generate_requires(nested_prop.properties) }
      nested_object_requires = nested_props.map\
        { |nested_object| nested_object_requires(nested_object) }
      nested_object_requires + nested_prop_requires + nested_array_requires
    end

    def array_requires(type)
      File.join(
        'google',
        type.__resource.__product.prefix[1..-1],
        'property',
        [type.__resource.name.downcase, type.item_type.name.underscore].join('_')
      )
    end

    def nested_object_requires(nested_object_type)
      File.join(
        'google',
        nested_object_type.__resource.__product.prefix[1..-1],
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
  end
end

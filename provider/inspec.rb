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

require 'provider/config'
require 'provider/core'
require 'provider/inspec/manifest'
require 'provider/inspec/resource_override'
require 'provider/inspec/property_override'

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
      target_folder = File.join(data[:output_folder], 'inspec')
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/singular_resource.erb',
        out_file: File.join(target_folder, "google_#{data[:product_name]}_#{name}.rb")
      )
      generate_resource_file data.clone.merge(
        default_template: 'templates/inspec/plural_resource.erb',
        out_file: File.join(target_folder, "google_#{data[:product_name]}_#{name}s.rb")
      )
    end

    # Returns the url that this object can be retrieved from
    # based off of the self link
    def url(object)
      url = object.self_link_url[1]
      return url.join('') if url.is_a?(Array)
      url.split("\n").join('')
    end

    # TODO?
    def generate_resource_tests(data) end

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

    # Figuring out if a property is a primitive ruby type is a hassle. But it is important
    # Fingerprints are strings, NameValues are hashes, and arrays of primitives are arrays
    # Arrays of NestedObjects need to have their contents parsed and returned in an array
    def primitive?(property)
      array_primitive = (property.is_a?(Api::Type::Array)\
        && !property.item_type.is_a?(::Api::Type::NestedObject))
      property.is_a?(::Api::Type::Primitive)\
        || array_primitive\
        || property.is_a?(::Api::Type::NameValues)\
        || property.is_a?(::Api::Type::Fingerprint)
    end

    # ResourceRefs are strings
    def resource_ref?(property)
      property.is_a?(::Api::Type::ResourceRef)
    end

    # Arrays need special requires statements
    def typed_array?(property)
      property.is_a?(::Api::Type::Array)
    end

    def nested_object?(property)
      property.is_a?(::Api::Type::NestedObject)
    end

    def generate_requires(properties)
      nested_props = properties.select { |type| nested_object?(type) }
      nested_object_arrays = properties.select\
        { |type| typed_array?(type) && nested_object?(type.item_type) }
      nested_array_requires = nested_object_arrays.collect { |type| array_requires(type) }
      nested_prop_requires = nested_props.map\
        { |nested_prop| generate_requires(nested_prop.properties) }
      nested_object_requires = nested_props.map\
        { |nested_object| nested_object_requires(nested_object) }
      nested_object_requires + nested_prop_requires + nested_array_requires
    end

    # Primitives don't need requires statements.
    # Nested objects will have their requires statements handled separately
    def no_requires?(type)
      primitive?(type) || resource_ref?(type) || nested_object?(type)
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
  end
end

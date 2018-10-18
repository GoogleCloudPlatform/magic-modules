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

    def primitive?(property)
      is_primitive = property.is_a?(::Api::Type::Primitive)
      is_primitive_array = (property.is_a?(Api::Type::Array)\
        && !property.item_type.is_a?(::Api::Type::NestedObject))
      is_primitive || is_primitive_array
    end

    def resource_ref?(property)
      property.is_a?(::Api::Type::ResourceRef)
    end

    def typed_array?(property)
      property.is_a?(::Api::Type::Array)
    end

    def nested_object?(property)
      property.is_a?(::Api::Type::NestedObject)
    end

    def generate_requires(properties)
      requires = []
      nested_props = properties.select { |type| nested_object?(type) }
      requires.concat(properties.reject { |type| no_requires?(type) }\
        .collect { |type| easy_requires(type) })
      requires.concat(nested_props.map { |nested_prop| generate_requires(nested_prop.properties) })
      requires.concat(nested_props.map { |nested_object| nested_object_requires(nested_object) } )
      requires
    end

    def no_requires?(type)
      primitive?(type) || resource_ref?(type) || nested_object?(type) || type.is_a?(::Api::Type::NameValues)
    end

    def easy_requires(type)
      if typed_array?(type)
        return File.join(
          'google',
          'compute',
          'property',
          [type.__resource.name.downcase, type.item_type.name.underscore].join('_')
        )
      end
      return File.join(
        'google',
        'compute',
        'property',
        type.name.underscore
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

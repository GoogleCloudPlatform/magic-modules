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

require 'provider/abstract_core'
require 'provider/terraform/config'
require 'provider/terraform/import'
require 'provider/terraform/custom_code'
require 'provider/terraform/property_override'
require 'provider/terraform/resource_override'
require 'provider/terraform/sub_template'
require 'google/golang_utils'

require 'provider/azure/terraform'

module Provider
  # Code generator for Terraform Resources that manage Google Cloud Platform
  # resources.
  class Terraform < Provider::AbstractCore
    include Provider::Terraform::Import
    include Provider::Terraform::SubTemplate
    include Google::GolangUtils

    include Provider::Azure::Terraform

    # Sorts properties in the order they should appear in the TF schema:
    # Required, Optional, Computed
    def order_properties(properties)
      properties.select(&:required).sort_by(&:name) +
        properties.reject(&:required).reject(&:output).sort_by(&:name) +
        properties.select(&:output).sort_by(&:name)
    end

    def tf_type(property)
      tf_types[property.class]
    end

    # Converts between the Magic Modules type of an object and its type in the
    # TF schema
    def tf_types
      {
        Api::Type::Boolean => 'schema.TypeBool',
        Api::Type::Double => 'schema.TypeFloat',
        Api::Type::Integer => 'schema.TypeInt',
        Api::Type::String => 'schema.TypeString',
        # Anonymous string property used in array of strings.
        'Api::Type::String' => 'schema.TypeString',
        Api::Type::Time => 'schema.TypeString',
        Api::Type::Enum => 'schema.TypeString',
        Api::Type::ResourceRef => 'schema.TypeString',
        Api::Type::NestedObject => 'schema.TypeList',
        Api::Type::Array => 'schema.TypeList',
        Api::Type::KeyValuePairs => 'schema.TypeMap',
        Api::Type::Map => 'schema.TypeSet',
        Api::Type::Fingerprint => 'schema.TypeString'
      }
    end

    def updatable?(resource, properties)
      !resource.input || !properties.reject { |p| p.update_url.nil? }.empty?
    end

    def force_new?(property, resource)
      !property.output &&
        (property.input || (resource.input && property.update_url.nil? &&
                            (property.parent.nil? ||
                             force_new?(property.parent, resource))))
    end

    def build_url(url_parts, _extra = false)
      url_parts.flatten.join
    end

    # Transforms a format string with field markers to a regex string with
    # capture groups.
    #
    # For instance,
    #   projects/{{project}}/global/networks/{{name}}
    # is transformed to
    #   projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)
    def format2regex(format)
      format.gsub(/{{([[:word:]]+)}}/, '(?P<\1>[^/]+)')
    end

    # Capitalize the first letter of a property name.
    # E.g. "creationTimestamp" becomes "CreationTimestamp".
    def titlelize_property(property)
      p = property.name.clone
      p[0] = p[0].capitalize
      p
    end

    # Returns the nested properties. An empty list is returned if the property
    # is not a NestedObject or an Array of NestedObjects.
    def nested_properties(property)
      if property.is_a?(Api::Type::NestedObject)
        property.properties
      elsif property.is_a?(Api::Type::Array) &&
            property.item_type.is_a?(Api::Type::NestedObject)
        property.item_type.properties
      elsif property.is_a?(Api::Type::Map)
        property.value_type.properties
      else
        []
      end
    end

    private

    # This function uses the resource.erb template to create one file
    # per resource. The resource.erb template forms the basis of a single
    # GCP Resource on Terraform.
    def generate_resource(data)
      # dir = data[:version] == 'beta' ? 'google-beta' : 'google'
      dir = "azurerm"
      target_folder = File.join(data[:output_folder], dir)
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      product_name = data[:product_name].underscore
      filepath = File.join(target_folder, "resource_arm_#{name}.go")
      generate_resource_file data.clone.merge(
        default_template: 'templates/terraform/resource.erb',
        out_file: filepath
      )
      # TODO: error check goimports
      # %x(goimports -w #{filepath})
      generate_documentation(data)
    end

    def generate_documentation(data)
      target_folder = data[:output_folder]
      target_folder = File.join(target_folder, 'website', 'docs', 'r')
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      product_name = data[:product_name].underscore

      filepath =
        File.join(target_folder, "#{name}.html.markdown")
      generate_resource_file data.clone.merge(
        default_template: 'templates/terraform/resource.html.markdown.erb',
        out_file: filepath
      )
    end

    def generate_resource_tests(data)
      return if data[:object].example.reject(&:skip_test).empty?

      # dir = data[:version] == 'beta' ? 'google-beta' : 'google'
      dir = 'azurerm'
      target_folder = File.join(data[:output_folder], dir)
      FileUtils.mkpath target_folder
      name = data[:object].name.underscore
      product_name = data[:product_name].underscore
      filepath =
        File.join(
          target_folder,
          "resource_arm_#{name}_test.go"
        )
      generate_resource_file data.clone.merge(
        product: data[:product_name].camelize(:upper),
        resource_name: data[:object].name.camelize(:upper),
        default_template: 'templates/terraform/examples/base_configs/test_file.go.erb',
        out_file: filepath
      )

      # TODO: error check goimports
      # %x(goimports -w #{filepath})
    end
  end
end

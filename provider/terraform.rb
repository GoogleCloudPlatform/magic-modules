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
require 'provider/terraform/property_override'
require 'provider/terraform/resource_override'
require 'provider/terraform/sub_template'
require 'google/golang_utils'

module Provider
  # Code generator for Terraform Resources that manage Google Cloud Platform
  # resources.
  class Terraform < Provider::AbstractCore
    include Provider::Terraform::Import
    include Provider::Terraform::SubTemplate
    include Google::GolangUtils

    # Sorts properties in the order they should appear in the TF schema:
    # Required, Optional, Computed
    def order_properties(properties)
      properties.select(&:required).sort_by(&:name) +
        properties.reject(&:required).reject(&:output) +
        properties.select(&:output).sort_by(&:name)
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
        Api::Type::NameValues => 'schema.TypeMap'
      }
    end

    def updatable?(resource, properties)
      !resource.input || !properties.reject { |p| p.update_url.nil? }.empty?
    end

    def force_new?(property, resource)
      !property.output &&
        (property.input || (resource.input && property.update_url.nil?))
    end

    # Puts together the links to use to make API calls for a given resource type
    def self_link_url(resource)
      (product_url, resource_url) = self_link_raw_url(resource)
      [product_url, resource_url].flatten.join
    end

    def collection_url(resource)
      base_url = resource.base_url.split("\n").map(&:strip).compact
      [resource.__product.base_url, base_url].flatten.join
    end

    def update_url(resource, url_part)
      return self_link_url(resource) if url_part.nil?
      [resource.__product.base_url, url_part].flatten.join
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

    # Returns the resource properties without those ignored.
    def effective_properties(properties)
      properties.reject(&:exclude)
    end

    # Returns the nested properties without those ignored. An empty list is
    # returned if the property is not a NestedObject or an Array of
    # NestedObjects.
    def effective_nested_properties(property)
      if property.is_a?(Api::Type::NestedObject)
        effective_properties(property.properties)
      elsif property.is_a?(Api::Type::Array) &&
            property.item_type.is_a?(Api::Type::NestedObject)
        effective_properties(property.item_type.properties)
      else
        []
      end
    end

    # Filter the properties to keep only the ones requiring custom update
    # method and group them by update url & verb.
    def properties_by_custom_update(properties)
      update_props = properties.reject do |p|
        p.update_url.nil? || p.update_verb.nil?
      end
      update_props.group_by { |p| p.update_url + p.update_verb.to_s }.values
    end

    private

    # This function uses the resource.erb template to create one file
    # per resource. The resource.erb template forms the basis of a single
    # GCP Resource on Terraform.
    def generate_resource(data)
      target_folder = File.join(data[:output_folder], 'google')
      FileUtils.mkpath target_folder
      name = Google::StringUtils.underscore(data[:object].name)
      product_name = Google::StringUtils.underscore(data[:product_name])
      filepath = File.join(target_folder, "resource_#{product_name}_#{name}.go")
      generate_resource_file data.clone.merge(
        default_template: 'templates/terraform/resource.erb',
        out_file: filepath
      )
      # TODO: error check goimports
      %x(goimports -w #{filepath})

      generate_documentation(data)
    end

    def generate_documentation(data)
      target_folder = data[:output_folder]
      target_folder = File.join(target_folder, 'website', 'docs', 'r')
      FileUtils.mkpath target_folder
      name = Google::StringUtils.underscore(data[:object].name)
      product_name = Google::StringUtils.underscore(data[:product_name])
      filepath =
        File.join(target_folder, "#{product_name}_#{name}.html.markdown")
      generate_resource_file data.clone.merge(
        default_template: 'templates/terraform/resource.html.markdown.erb',
        out_file: filepath
      )
    end
  end
end

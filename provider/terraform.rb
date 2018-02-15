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

module Provider
  # A boilerplate provider where all methods are optional.
  class TerraformBase < Provider::Core
    private

    # rubocop:disable Layout/EmptyLineBetweenDefs
    def generate_resource(data) end
    def generate_resource_tests(data) end
    def generate_network_datas(data, object) end
    def generate_base_property(data) end
    def generate_simple_property(type, data) end
    def emit_nested_object(data) end
    def emit_resourceref_object(data) end
    def generate_typed_array(data, prop) end
    # rubocop:enable Layout/EmptyLineBetweenDefs
  end

  # Code generator for Terraform Cookbooks that manage Google Cloud Platform
  # resources.
  class Terraform < TerraformBase
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      def provider
        Provider::Terraform
      end
    end

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
        Api::Type::Array => 'schema.TypeList'
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

    # Returns a list of acceptable import id formats for a given resource.
    #
    # For instance, if the resource base url is:
    #   projects/{{project}}/global/networks
    #
    # It returns 3 formats:
    # a) self_link: projects/{{project}}/global/networks/{{name}}
    # b) short id: {{project}}/{{name}}
    # c) short id w/o defaults: {{name}}
    #
    # Fields with default values are `project`, `region` and `zone`.
    def import_id_formats(resource)
      underscored_base_url = resource.base_url
                                     .gsub(/{{[[:word:]]+}}/) do |field_name|
        Google::StringUtils.underscore(field_name)
      end

      # TODO: Add support for custom import id
      # We assume that all resources have a name field
      self_link_id_format = underscored_base_url + '/{{name}}'

      # short id: {{project}}/{{zone}}/{{name}}
      field_markers = self_link_id_format.scan(/{{[[:word:]]+}}/)
      short_id_format = field_markers.join('/')

      # short id without fields with provider-level default: {{name}}
      field_markers.delete('{{project}}')
      field_markers.delete('{{region}}')
      field_markers.delete('{{zone}}')
      short_id_default_format = field_markers.join('/')

      [self_link_id_format, short_id_format, short_id_default_format]
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

    def build_schema_property(config, property, object)
      compile_template'templates/terraform/schema_property.erb',
                      property: property,
                      config: config,
                      object: object
    end

    # Transforms a Cloud API representation of a property into a Terraform
    # schema representation.
    def build_flatten_method(config, prefix, property)
      compile_template 'templates/terraform/flatten_property_method.erb',
                       prefix: prefix,
                       property: property,
                       config: config
    end

    # Transforms a Terraform schema representation of a property into a
    # representation used by the Cloud API.
    def build_expand_method(config, prefix, property)
      compile_template 'templates/terraform/expand_property_method.erb',
                       prefix: prefix,
                       property: property,
                       config: config
    end

    def build_property_documentation(config, property)
      compile_template 'templates/terraform/property_documentation.erb',
                       property: property,
                       config: config
    end

    def build_nested_property_documentation(config, property)
      compile_template 'templates/terraform/nested_property_documentation.erb',
                       property: property,
                       config: config
    end

    # Capitalize the first letter of a property name.
    # E.g. "creationTimestamp" becomes "CreationTimestamp".
    def titlelize_property(property)
      p = property.name.clone
      p[0] = p[0].capitalize
      p
    end

    # Returns the resource properties without those ignored.
    def effective_properties(config, properties)
      ignored = get_code_multiline(config, 'ignore') || []

      properties.keep_if { |p| !ignored.include?(construct_ignore_string(p)) }
    end

    # Returns the nested properties without those ignored. An empty list is
    # returned if the property is not a NestedObject or an Array of
    # NestedObjects.
    def effective_nested_properties(config, property)
      if property.is_a?(Api::Type::NestedObject)
        effective_properties(config, property.properties)
      elsif property.is_a?(Api::Type::Array) &&
            property.item_type.is_a?(Api::Type::NestedObject)
        effective_properties(config, property.item_type.properties)
      else
        []
      end
    end

    private

    def compile_template(template_file, data)
      ctx = binding
      data.each { |name, value| ctx.local_variable_set(name, value) }
      compile_file(ctx, template_file).join("\n")
    end

    # Constructs the prefix to be used when looking for ignored properties.
    #
    # The 'ignore' list supports three formats:
    # - 'foo': Ignores top-level property 'foo'
    # - 'foo.bar': Ignores field 'bar' nested under 'foo'
    # - 'foo.*.bar': Ignores field 'bar' of all nested objects in list 'foo'
    def construct_ignore_string(property)
      return property.name if property.parent.nil?

      if property.parent.is_a?(Api::Type::Array)
        construct_ignore_string(property.parent) + '.*'
      else
        construct_ignore_string(property.parent) + '.' + property.name
      end
    end

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

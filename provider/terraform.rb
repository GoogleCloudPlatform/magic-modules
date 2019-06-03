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
require 'provider/terraform/docs'
require 'provider/terraform/examples'
require 'overrides/terraform/resource_override'
require 'overrides/terraform/property_override'
require 'provider/terraform/sub_template'
require 'google/golang_utils'

module Provider
  # Code generator for Terraform Resources that manage Google Cloud Platform
  # resources.
  class Terraform < Provider::AbstractCore
    include Provider::Terraform::Import
    include Provider::Terraform::SubTemplate
    include Google::GolangUtils

    # FileTemplate with Terraform specific fields
    class TerraformFileTemplate < Provider::FileTemplate
      # The async object used for making operations.
      # We assume that all resources share the same async properties.
      attr_accessor :async

      # When generating OiCS examples, we attach the example we're
      # generating to the data object.
      attr_accessor :example

      attr_accessor :resource_name
    end

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

    # "Namespace" - prefix with product and resource - a property with
    # information from the "object" variable
    def namespace_property_from_object(property, object)
      name = property.name.camelize
      until property.parent.nil?
        property = property.parent
        name = property.name.camelize + name
      end

      "#{property.__resource.__product.api_name.camelize(:lower)}#{object.name}#{name}"
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

    # Transforms a format string with field markers to a regex string with
    # capture groups.
    #
    # For instance,
    #   projects/{{project}}/global/networks/{{name}}
    # is transformed to
    #   projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)
    #
    # Values marked with % are URL-encoded, and will match any number of /'s.
    #
    # Note: ?P indicates a Python-compatible named capture group. Named groups
    # aren't isn't common in JS-based regex flavours, but are in Perl-based ones
    def format2regex(format)
      format
        .gsub(/{{%([[:word:]]+)}}/, '(?P<\1>.+)')
        .gsub(/{{([[:word:]]+)}}/, '(?P<\1>[^/]+)')
    end

    # Capitalize the first letter of a property name.
    # E.g. "creationTimestamp" becomes "CreationTimestamp".
    def titlelize_property(property)
      p = property.name.clone
      p[0] = p[0].capitalize
      p
    end

    private

    # This function uses the resource.erb template to create one file
    # per resource. The resource.erb template forms the basis of a single
    # GCP Resource on Terraform.
    def generate_resource(data)
      dir = data.version == 'beta' ? 'google-beta' : 'google'
      target_folder = File.join(data.output_folder, dir)

      name = data.object.name.underscore
      product_name = data.product.name.underscore
      filepath = File.join(target_folder, "resource_#{product_name}_#{name}.go")

      data.generate('templates/terraform/resource.erb', filepath, self)
      generate_documentation(data)
    end

    def generate_documentation(data)
      target_folder = data.output_folder
      target_folder = File.join(target_folder, 'website', 'docs', 'r')
      FileUtils.mkpath target_folder
      name = data.object.name.underscore
      product_name = data.product.name.underscore

      filepath =
        File.join(target_folder, "#{product_name}_#{name}.html.markdown")
      data.generate('templates/terraform/resource.html.markdown.erb', filepath, self)
    end

    def generate_resource_tests(data)
      return if data.object.examples
                    .reject(&:skip_test)
                    .reject do |e|
                  @api.version_obj_or_default(data.version) \
                < @api.version_obj_or_default(e.min_version)
                end
                    .empty?

      dir = data.version == 'beta' ? 'google-beta' : 'google'
      target_folder = File.join(data.output_folder, dir)

      name = data.object.name.underscore
      product_name = data.product.name.underscore
      filepath =
        File.join(
          target_folder,
          "resource_#{product_name}_#{name}_generated_test.go"
        )

      data.product = data.product.name
      data.resource_name = data.object.name.camelize(:upper)
      data.generate('templates/terraform/examples/base_configs/test_file.go.erb',
                    filepath, self)
    end

    def generate_operation(output_folder, _types, version_name)
      return if @api.objects.select(&:autogen_async).empty?

      product_name = @api.name.underscore
      data = build_object_data(@api.objects.first, output_folder, version_name)
      dir = data.version == 'beta' ? 'google-beta' : 'google'
      target_folder = File.join(data.output_folder, dir)

      data.object = @api.objects.select(&:autogen_async).first
      data.async = data.object.async
      data.generate('templates/terraform/operation.go.erb',
                    File.join(target_folder,
                              "#{product_name}_operation.go"),
                    self)
    end

    def build_object_data(object, output_folder, version)
      TerraformFileTemplate.file_for_resource(output_folder, object, version, @config, build_env)
    end
  end
end

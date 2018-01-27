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
  # Code generator for Terraform Cookbooks that manage Google Cloud Platform
  # resources.
  class Terraform < Provider::Core
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

    # Puts together the links to use to make API calls for a given resource type
    def self_link_url(resource)
      (product_url, resource_url) = self_link_raw_url(resource)
      [product_url, resource_url].flatten.join
    end

    def collection_url(resource)
      base_url = resource.base_url.split("\n").map(&:strip).compact
      [resource.__product.base_url, base_url].flatten.join
    end

    def build_schema_property(config, property)
      compile_template'templates/terraform/schema_property.erb',
                      property: property,
                      config: config
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

    # This function would generate unit tests using a template
    def generate_resource_tests(data) end

    # This function would automatically generate the files used for verifying
    # network calls in unit tests. If you comment out the following line,
    # a bunch of YAML files will be created under the spec/ folder.
    def generate_network_datas(data, object) end

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
    def generate_simple_property(type, data) end
    def emit_nested_object(data) end
    def emit_resourceref_object(data) end
    def generate_typed_array(data, prop) end
    # rubocop:enable Layout/EmptyLineBetweenDefs
  end
end

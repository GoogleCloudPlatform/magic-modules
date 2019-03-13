require 'provider/core'
require 'api/object'

module Provider
  module Azure
    module Terraform
      class ExampleReference < Api::Object
        attr_reader :product
        attr_reader :example

        def validate
          super
          check_property :product, String
          check_property :example, String
        end
      end

      class Example < Api::Object
        attr_reader :resource
        attr_reader :name_in_documentation
        attr_reader :name_in_test
        attr_reader :prerequisites
        attr_reader :properties

        def validate
          super
          check_property :resource, String
          check_optional_property :name_in_documentation, String
          check_optional_property :name_in_test, String
          check_optional_property :prerequisites, Array
          check_optional_property_list :prerequisites, ExampleReference
          check_property :properties, Hash
        end

        module SubTemplate
          def build_test_hcl_from_example()
          end

          def build_documentation_hcl_from_example(product_name, example_name, name_hints, with_dependencies = false)
            documentation_raw = compile_template 'templates/azure/terraform/example/documentation_hcl.erb',
                                                 example: get_example_from_file(product_name, example_name),
                                                 with_dependencies: with_dependencies,
                                                 name_hints: name_hints,
                                                 resource_id_hint: "example"
            context = {
              resource_id_hint: "example",
              is_acctest: false
            }.merge(name_hints.transform_keys{|k| "#{k.underscore}_hint"})
            compile_string context, documentation_raw
          end

          def build_documentation_import_resource_id(object, example_reference)
            compile_template 'templates/azure/terraform/example/import_resource_id.erb',
                             name_hints: example_reference.resource_name_hints,
                             object: object
          end

          def build_hcl_properties(properties_hash, indentation = 2)
            compile_template 'templates/azure/terraform/example/hcl_properties.erb',
                             properties: properties_hash,
                             indentation: indentation
          end

          private

          def get_example_from_file(product_name, example_name)
            example_yaml = "products/#{product_name}/examples/terraform/#{example_name}.yaml"
            example = Google::YamlValidator.parse(File.read(example_yaml))
            raise "#{example_yaml}(#{example.class}) is not Provider::Azure::Terraform::Example" unless example.is_a?(Provider::Azure::Terraform::Example)
            example.validate
            example
          end
        end
      end
    end
  end
end

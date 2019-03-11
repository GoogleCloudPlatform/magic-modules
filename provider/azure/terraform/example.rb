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

          def build_documentation_hcl_from_example(product_name, example_name, with_dependencies = true, id_hint = "example")
            documentation_raw = compile_template 'templates/azure/terraform/example/documentation_hcl.erb',
                                                 example: get_example_from_file(product_name, example_name),
                                                 with_dependencies: with_dependencies,
                                                 resource_id_hint: id_hint
            context = {
              resource_id_hint: id_hint,
              is_acctest: false
            }
            compile_string context, documentation_raw
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

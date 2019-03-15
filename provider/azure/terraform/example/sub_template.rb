require 'api/object'

module Provider
  module Azure
    module Terraform
      class Example < Api::Object
        module SubTemplate
          def build_test_hcl_from_example()
          end

          def build_documentation_hcl_from_example(product_name, example_name, name_hints, with_dependencies = false)
            documentation_raw = compile_template 'templates/azure/terraform/example/documentation_hcl.erb',
                                                 example: get_example_from_file(product_name, example_name),
                                                 with_dependencies: with_dependencies,
                                                 name_hints: name_hints,
                                                 resource_id_hint: "example"
            context = ExampleContextBinding.new("example", false, name_hints)
            compile_string context.get_binding, documentation_raw
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

        private

        class ExampleContextBinding
          attr_reader :my_binding
          attr_reader :name_hints

          def initialize(resource_id_hint, is_acctest, name_hints)
            @my_binding = binding
            @my_binding.local_variable_set(:resource_id_hint, resource_id_hint)
            @my_binding.local_variable_set(:is_acctest, is_acctest)
            @name_hints = name_hints.transform_keys(&:underscore)
          end

          def get_binding()
            @my_binding
          end

          def get_resource_name(name_hint)
            name_hints[name_hint] || "acctest-%d"
          end
        end
      end
    end
  end
end

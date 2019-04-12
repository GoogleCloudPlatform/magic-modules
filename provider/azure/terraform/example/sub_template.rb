require 'api/object'

module Provider
  module Azure
    module Terraform
      class Example < Api::Object
        module SubTemplate
          def build_test_hcl_from_example(product_name, example_name, random_vars, with_dependencies = false)
            build_hcl_from_example(product_name, example_name, "test", {}, random_vars, true)
          end

          def build_documentation_hcl_from_example(product_name, example_name, name_hints, with_dependencies = false)
            build_hcl_from_example(product_name, example_name, "example", name_hints, [], true)
          end

          def build_hcl_from_example(product_name, example_name, id_hint, name_hints, random_vars, with_dependencies = false)
            hcl_raw = compile_template 'templates/azure/terraform/example/example_hcl.erb',
                                       example: get_example_from_file(product_name, example_name),
                                       random_variables: random_vars,
                                       resource_id_hint: id_hint,
                                       name_hints: name_hints,
                                       with_dependencies: with_dependencies
            context = ExampleContextBinding.new(id_hint, name_hints, random_vars)
            compile_string context.get_binding, hcl_raw
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
          attr_reader :random_variables

          def initialize(resource_id_hint, name_hints, random_vars)
            @my_binding = binding
            @my_binding.local_variable_set(:resource_id_hint, resource_id_hint)
            @name_hints = name_hints
            @random_variables = random_vars
          end

          def get_binding()
            @my_binding
          end

          def get_resource_name(name_hint, postfix)
            return name_hints[name_hint] if name_hints.has_key?(name_hint)
            @random_variables <<= RandomizedVariable.new(:AccDefaultInt)
            "acctest#{postfix}-#{@random_variables.last.format_string}"
          end

          def get_location()
            return name_hints["location"] if name_hints.has_key?("location")
            @random_variables <<= RandomizedVariable.new(:AccLocation)
            @random_variables.last.format_string
          end
        end

        class RandomizedVariable
          attr_reader :variable_name
          attr_reader :parameter_name
          attr_reader :go_type
          attr_reader :create_expression
          attr_reader :format_string

          def initialize(type)
            case type
            when :AccDefaultInt
              @variable_name = "ri"
              @parameter_name = "rInt"
              @go_type = "int"
              @create_expression = "tf.AccRandTimeInt()"
              @format_string = "%d"
            when :AccLocation
              @variable_name = @parameter_name = "location"
              @go_type = "string"
              @create_expression = "testLocation()"
              @format_string = "%s"
            end
          end
        end

      end
    end
  end
end

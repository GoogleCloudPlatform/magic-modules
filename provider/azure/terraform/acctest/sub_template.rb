module Provider
  module Azure
    module Terraform
      module AccTest
        module SubTemplate

          def build_acctest_parameters_from_schema(sdk_op_def, properties, object, check_existence = false, indentation = 8)
            compile_template 'templates/azure/terraform/acctest/parameters_from_schema.erb',
                             indentation: indentation,
                             sdk_op_def: sdk_op_def,
                             check_existence: check_existence,
                             properties: properties,
                             object: object
          end

          def build_acctest_dependencies_hcl(acctest_def, object)
            compile_template 'templates/azure/terraform/acctest/dependencies_hcl.erb',
                             acctest_def: acctest_def,
                             object: object
          end

          def build_acctest_resource_test_hcl(acctest_def, object, resource_id = "test")
            compile_template 'templates/azure/terraform/acctest/resource_test_hcl.erb',
                             acctest_def: acctest_def,
                             resource_id: resource_id,
                             object: object
          end

        end
      end
    end
  end
end

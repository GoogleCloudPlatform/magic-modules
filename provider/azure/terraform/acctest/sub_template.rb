module Provider
  module Azure
    module Terraform
      module AccTest
        module SubTemplate

          def build_acctest_parameters_from_schema(sdk_op_def, properties, object, check_existence = false, indentation = 4)
            compile_template 'templates/azure/terraform/acctest/parameters_from_schema.erb',
                             indentation: indentation,
                             sdk_op_def: sdk_op_def,
                             check_existence: check_existence,
                             properties: properties,
                             object: object
          end

        end
      end
    end
  end
end

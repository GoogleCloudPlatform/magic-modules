module Provider
  module Azure
    module Terraform
      module AccTest
        module SubTemplate

          def build_acctest_parameters_from_schema(sdk_op_def, properties, indentation = 8)
            compile_template 'templates/azure/terraform/acctest/parameters_from_schema.erb',
                             indentation: indentation,
                             sdk_op_def: sdk_op_def,
                             properties: properties
          end

        end
      end
    end
  end
end

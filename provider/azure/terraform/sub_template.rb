module Provider
  module Azure
    module Terraform
      module SubTemplate
        def build_azure_id_parser(sdk_op_def, object, indentation = 4)
          compile_template 'templates/azure/terraform/sdk/azure_id_parser.erb',
                           indentation: indentation,
                           sdk_op_def: sdk_op_def,
                           object: object
        end

        def build_errorf_with_resource_name(format_string, include_error, sdk_op_def, properties, object)
          compile_template 'templates/azure/terraform/sdk/errorf_with_resource_name.erb',
                           format_string: format_string,
                           include_error: include_error,
                           sdk_op_def: sdk_op_def,
                           properties: properties,
                           object: object
        end
      end
    end
  end
end

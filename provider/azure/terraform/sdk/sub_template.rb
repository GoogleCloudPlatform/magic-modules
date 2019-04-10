module Provider
  module Azure
    module Terraform
      module SDK
        module SubTemplate
          def build_schema_property_get(input, output, property, object, indentation = 0)
            compile_template schema_property_get_template(property),
                             indentation: indentation,
                             input_var: input,
                             output_var: output,
                             prop_name: property.name.underscore,
                             property: property,
                             object: object
          end
  
          def build_schema_property_set(input, output, property, indentation = 0)
            compile_template schema_property_set_template(property),
                             indentation: indentation,
                             input_var: input,
                             output_var: output,
                             prop_name: property.name.underscore,
                             property: property
          end
  
          def build_sdk_field_assignment(property, sdk_type, resource_name, object)
            compile_template property_to_sdk_field_assignment_template(property, sdk_type),
                             property: property,
                             sdk_type: sdk_type,
                             resouce_name: resource_name,
                             object: object
          end

          def build_property_to_sdk_object(api_path, resource_name, sdk_type_defs, object, indentation = 4)
            compile_template property_to_sdk_object_template(sdk_type_defs, api_path),
                             indentation: indentation,
                             resource_name: resource_name,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             object: object
          end
  
          def build_sdk_object_to_property(input, api_path, sdk_type_defs, object, indentation = 4)
            compile_template sdk_object_to_property_template(sdk_type_defs, api_path),
                             indentation: indentation,
                             input_statement: input,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             object: object
          end
  
          def build_sdk_func_invocation(sdk_op_def)
            compile_template 'templates/azure/terraform/sdk/function_invocation.erb',
                             sdk_op_def: sdk_op_def
          end
        end
      end
    end
  end
end

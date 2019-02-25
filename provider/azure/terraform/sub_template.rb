module Provider
  module Azure
    module Terraform
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

        def build_property_to_sdk_object(output, sdk_path, sdk_type_def, sdk_package, resource_name, properties, sdk_type_defs, object, indentation = 4)
          compile_template property_to_sdk_object_template(sdk_type_def),
                           indentation: indentation,
                           output_statement: output,
                           sdk_package_name: sdk_package,
                           resource_name: resource_name,
                           sdk_obj_path: sdk_path,
                           sdk_type_def: sdk_type_def,
                           sdk_type_defs: sdk_type_defs,
                           properties: properties,
                           object: object
        end

        def build_sdk_object_to_property(input, api_path, property, sdk_type_defs, object, indentation = 4)
          compile_template sdk_object_to_property_template(sdk_type_defs, api_path),
                           indentation: indentation,
                           input_statement: input,
                           api_path: api_path,
                           property: property,
                           sdk_type_defs: sdk_type_defs,
                           object: object
        end

        def build_sdk_func_invocation(sdk_op_def)
          compile_template 'templates/azure/terraform/sdk/function_invocation.erb',
                           sdk_op_def: sdk_op_def
        end

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

module Provider
  module Azure
    module Terraform
      module SDK
        module SubTemplate
          def build_schema_property_get(input, output, property, sdk_marshal, indentation = 0)
            compile_template schema_property_get_template(property),
                             indentation: indentation,
                             input_var: input,
                             output_var: output,
                             sdk_marshal: sdk_marshal,
                             property: property,
                             prop_name: property.name.underscore
          end

          def build_schema_property_set(input, output, property, sdk_marshal, indentation = 0)
            compile_template schema_property_set_template(property),
                             indentation: indentation,
                             input_var: input,
                             output_var: output,
                             property: property,
                             sdk_marshal: sdk_marshal,
                             prop_name: property.name.underscore
          end

          def build_sdk_field_assignment(property, sdk_marshal, in_structure = true)
            compile_template property_to_sdk_field_assignment_template(property, sdk_marshal.sdktype.type_definition),
                             property: property,
                             sdk_marshal: sdk_marshal,
                             in_structure: in_structure
          end

          def build_property_to_sdk_object(sdk_marshal, indentation = 4, include_empty = false)
            compile_template 'templates/azure/terraform/sdktypes/property_to_sdkobject.erb',
                             indentation: indentation,
                             sdk_marshal: sdk_marshal,
                             include_empty: include_empty
          end

          def build_property_to_sdk_object_empty_sensitive(sdk_marshal, indentation = 4)
            compile_template 'templates/azure/terraform/sdktypes/property_to_sdkobject_empty_sensitive.erb',
                             indentation: indentation,
                             sdk_marshal: sdk_marshal
          end

          def build_schema_assignment(input, output, property, sdk_marshal)
            compile_template property_to_schema_assignment_template(property, sdk_marshal.sdktype.operation, sdk_marshal.sdktype.typedef_reference),
                             input: input,
                             output: output,
                             sdk_marshal: sdk_marshal
          end

          def build_sdk_object_to_property(input, output, sdk_marshal, indentation = 4)
            compile_template 'templates/azure/terraform/sdktypes/sdkobject_to_property.erb',
                             indentation: indentation,
                             input: input,
                             output: output,
                             sdk_marshal: sdk_marshal
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

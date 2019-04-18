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

          def build_schema_property_set(input, output, api_path, sdk_type_defs, resource_name, flatten_queue, property, indentation = 0)
            compile_template schema_property_set_template(property),
                             indentation: indentation,
                             input_var: input,
                             output_var: output,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             prop_name: property.name.underscore,
                             property: property,
                             resource_name: resource_name,
                             flatten_queue: flatten_queue
          end

          def build_sdk_field_assignment(property, api_path, sdk_type_defs, resource_name, expand_queue, properties, object)
            compile_template property_to_sdk_field_assignment_template(property, sdk_type_defs[api_path]),
                             property: property,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             resource_name: resource_name,
                             expand_queue: expand_queue,
                             properties: properties,
                             object: object
          end

          def build_property_to_sdk_object(api_path, resource_name, sdk_type_defs, expand_queue, properties, object, indentation = 4)
            compile_template 'templates/azure/terraform/sdktypes/property_to_sdkobject.erb',
                             indentation: indentation,
                             resource_name: resource_name,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             expand_queue: expand_queue,
                             properties: properties,
                             object: object
          end

          def build_sdkfield_block_assignments(resource_name, sdk_type_defs, expand_queue, properties, object, indentation = 4)
            compile_template 'templates/azure/terraform/sdktypes/sdkfield_block_assignments.erb',
                             indentation: indentation,
                             resource_name: resource_name,
                             sdk_type_defs: sdk_type_defs,
                             expand_queue: expand_queue,
                             properties: properties,
                             object: object
          end

          def build_schema_assignment(input, output, property, api_path, sdk_type_defs, resource_name, flatten_queue, properties, object)
            compile_template property_to_schema_assignment_template(property, sdk_type_defs[api_path]),
                             input_statement: input,
                             output: output,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             resource_name: resource_name,
                             flatten_queue: flatten_queue,
                             properties: properties,
                             object: object
          end

          def build_sdk_object_to_property(input, output, api_path, sdk_type_defs, resource_name, flatten_queue, properties, object, indentation = 4)
            compile_template 'templates/azure/terraform/sdktypes/sdkobject_to_property.erb',
                             indentation: indentation,
                             input_statement: input,
                             output: output,
                             api_path: api_path,
                             sdk_type_defs: sdk_type_defs,
                             resource_name: resource_name,
                             flatten_queue: flatten_queue,
                             properties: properties,
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

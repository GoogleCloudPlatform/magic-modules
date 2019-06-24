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

        def build_azure_schema_property(property, object, indentation = 0, data_source_input = [])
          compile_template schema_property_template(property, !data_source_input.empty?),
                           indentation: indentation,
                           prop_name: property.name.underscore,
                           property: property,
                           data_source_input: data_source_input,
                           object: object
        end

        # Transforms a Cloud API representation of a property into a Terraform
        # schema representation.
        def build_flatten_method(ef_desc)
          compile_template 'templates/azure/terraform/flatten_property_method.erb',
                           descriptor: ef_desc
        end

        # Transforms a Terraform schema representation of a property into a
        # representation used by the Cloud API.
        def build_expand_method(ef_desc)
          compile_template 'templates/azure/terraform/expand_property_method.erb',
                           descriptor: ef_desc
        end

        def build_azure_property_documentation(property, data_source_input = [])
          compile_template 'templates/azure/terraform/property_documentation.erb',
                           property: property,
                           data_source_input: data_source_input
        end
  
        def build_azure_nested_property_documentation(property, data_source_input = [])
          compile_template 'templates/azure/terraform/nested_property_documentation.erb',
                           property: property,
                           data_source_input: data_source_input
        end
      end

    end
  end
end

require 'provider/azure/terraform/sdk/expand_flatten_descriptor'

module Provider
  module Azure
    module Terraform
      module SDK
        module Helpers
          def get_properties_matching_sdk_reference(properties, sdk_reference, object)
            properties.select{|p| p.azure_sdk_references.include?(sdk_reference)}.sort_by{|p| [p.order, p.name]}
          end

          def get_applicable_reference(references, typedefs)
            references.each do |ref|
              return ref if typedefs.has_key?(ref)
            end
            nil
          end

          def get_sdk_typedef_by_references(references, typedefs)
            ref = get_applicable_reference(references, typedefs)
            return nil if ref.nil?
            typedefs[ref]
          end

          def expand_or_flatten_enqueue(expand_queue, property, api_path, sdk_type_defs)
            sdk_type = sdk_type_defs[api_path]
            existed = expand_queue.any?{|exp| exp.property == property && exp.sdk_type.go_type_name == sdk_type.go_type_name}
            expand_queue << ExpandFlattenDescriptor.new(property, api_path, sdk_type_defs) unless existed
          end

          def property_to_sdk_field_assignment_template(property, sdk_type)
            return property.custom_sdkfield_assign unless get_property_value(property, "custom_sdkfield_assign", nil).nil?
            return 'templates/terraform/schemas/hide_from_schema.erb' if get_property_value(property, "hide_from_schema", false)
            case sdk_type
            when Api::Azure::SDKTypeDefinition::BooleanObject, Api::Azure::SDKTypeDefinition::StringObject
              'templates/azure/terraform/sdktypes/expand_func_field_assign.erb'
            when Api::Azure::SDKTypeDefinition::EnumObject
              'templates/azure/terraform/sdktypes/enum_field_assign.erb'
            when Api::Azure::SDKTypeDefinition::ComplexObject
              return 'templates/azure/terraform/sdktypes/nested_object_field_assign.erb' if property.nil?
              'templates/azure/terraform/sdktypes/expand_func_field_assign.erb'
            else
              'templates/azure/terraform/sdktypes/unsupport.erb'
            end
          end

          def property_to_schema_assignment_template(property, sdk_operation, api_path)
            sdk_type = sdk_operation.response[api_path] || sdk_operation.request[api_path]
            case sdk_type
            when Api::Azure::SDKTypeDefinition::BooleanObject, Api::Azure::SDKTypeDefinition::StringObject
              'templates/azure/terraform/sdktypes/primitive_schema_assign.erb'
            when Api::Azure::SDKTypeDefinition::EnumObject
              'templates/azure/terraform/sdktypes/enum_schema_assign.erb'
            when Api::Azure::SDKTypeDefinition::ComplexObject
              return 'templates/azure/terraform/sdktypes/nested_object_schema_assign.erb' if property.nil?
              'templates/azure/terraform/sdktypes/primitive_schema_assign.erb'
            else
              'templates/azure/terraform/sdktypes/unsupport.erb'
            end
          end
        end
      end
    end
  end
end
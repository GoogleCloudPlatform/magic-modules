module Provider
  module Azure
    module Terraform
      module SDK
        module Helpers
          def get_properties_matching_sdk_reference(sdk_reference, object)
            object.all_user_properties
              .select{|p| p.azure_sdk_references.include?(sdk_reference)}
              .sort_by{|p| [p.order, p.name]}
          end

          def property_to_sdk_field_assignment_template(property, sdk_type)
            return 'templates/terraform/schemas/hide_from_schema.erb' if get_property_value(property, "hide_from_schema", false)
            case sdk_type
            when Api::Azure::SDKTypeDefinition::BooleanObject, Api::Azure::SDKTypeDefinition::StringObject
              'templates/azure/terraform/sdktypes/expand_func_field_assign.erb'
            when Api::Azure::SDKTypeDefinition::EnumObject
              'templates/azure/terraform/sdktypes/enum_fiield_assign.erb'
            else
              'templates/azure/terraform/sdktypes/unsupport.erb'
            end
          end

          def property_to_sdk_object_template(sdk_type_defs, api_path)
            case sdk_type_defs[api_path]
            when Api::Azure::SDKTypeDefinition::ComplexObject
              'templates/azure/terraform/sdktypes/property_to_sdkobject.erb'
            else
              'templates/azure/terraform/sdktypes/property_to_sdkfield_assign.erb'
            end
          end
  
          def sdk_object_to_property_template(sdk_type_defs, api_path)
            return 'templates/azure/terraform/sdktypes/sdkobject_to_property.erb' if api_path == ""
            case sdk_type_defs[api_path]
            when Api::Azure::SDKTypeDefinition::BooleanObject, Api::Azure::SDKTypeDefinition::StringObject
              'templates/azure/terraform/sdktypes/sdkprimitive_to_property.erb'
            when Api::Azure::SDKTypeDefinition::EnumObject
              'templates/azure/terraform/sdktypes/sdkenum_to_property.erb'
            when Api::Azure::SDKTypeDefinition::ComplexObject
              'templates/azure/terraform/sdktypes/sdkobject_to_property.erb'
            else
              'templates/azure/terraform/sdktypes/unsupport.erb'
            end
          end
        end
      end
    end
  end
end

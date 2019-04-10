module Provider
  module Azure
    module Terraform
      module Schema

        def go_type(property)
          case property
          when Api::Type::Boolean
            'bool'
          when Api::Type::Enum, Api::Type::String
            'string'
          when Api::Type::KeyValuePairs
            'map[string]interface{}'
          when Api::Type::NestedObject
            '[]interface{}'
          else
            'interface{}'
          end
        end

        def expand_func(property)
          expand_funcs[property.class]
        end

        def expand_funcs
          {
            Api::Type::Boolean => 'utils.Bool',
            Api::Type::String => 'utils.String',
            Api::Azure::Type::Location => "utils.String",
            Api::Azure::Type::Tags => 'expandTags',
          }
        end

        def schema_property_template(property)
          return property.custom_schema_definition unless get_property_value(property, "custom_schema_definition", nil).nil?
          case property
          when Api::Azure::Type::ResourceGroupName
            'templates/azure/terraform/schemas/resource_group_name.erb'
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location.erb'
          when Api::Azure::Type::Tags
            'templates/azure/terraform/schemas/tags.erb'
          when Api::Type::Boolean, Api::Type::Enum, Api::Type::String, Api::Type::KeyValuePairs, Api::Type::NestedObject
            'templates/terraform/schemas/primitive.erb'
          else
            'templates/terraform/schemas/unsupport.erb'
          end
        end

        def schema_property_get_template(property)
          return property.custom_schema_get unless get_property_value(property, "custom_schema_get", nil).nil?
          return 'templates/terraform/schemas/hide_from_schema.erb' if get_property_value(property, "hide_from_schema", false)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location_get.erb'
          when Api::Type::Boolean, Api::Type::Enum, Api::Type::String, Api::Type::KeyValuePairs, Api::Type::NestedObject
            'templates/terraform/schemas/basic_get.erb'
          else
            'templates/terraform/schemas/unsupport.erb'
          end
        end

        def schema_property_set_template(property)
          return property.custom_schema_set unless get_property_value(property, "custom_schema_set", nil).nil?
          return 'templates/terraform/schemas/hide_from_schema.erb' if get_property_value(property, "hide_from_schema", false)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location_set.erb'
          when Api::Azure::Type::Tags
            'templates/azure/terraform/schemas/tags_set.erb'
          when Api::Type::Boolean, Api::Type::Enum, Api::Type::String, Api::Type::KeyValuePairs
            'templates/terraform/schemas/basic_set.erb'
          else
            'templates/terraform/schemas/unsupport.erb'
          end
        end

      end
    end
  end
end

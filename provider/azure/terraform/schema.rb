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
          when Api::Type::Integer
            'int'
          when Api::Type::Double
            'float64'
          when Api::Type::KeyValuePairs
            'map[string]interface{}'
          when Api::Type::Array, Api::Type::NestedObject
            '[]interface{}'
          else
            'interface{}'
          end
        end

        def go_empty_value(property)
          case property
          when Api::Type::Enum, Api::Type::String
            '""'
          else
            'nil'
          end
        end

        def expand_func(property)
          expand_funcs[property.class]
        end

        def expand_funcs
          {
            Api::Type::Boolean => 'utils.Bool',
            Api::Type::String => 'utils.String',
            Api::Type::Integer => 'utils.Int',
            Api::Type::Double => 'utils.Float',
            Api::Azure::Type::Location => "utils.String",
            Api::Azure::Type::Tags => 'expandTags',
            Api::Azure::Type::ResourceReference => "utils.String"
          }
        end

        def schema_property_template(property, is_data_source)
          return property.custom_schema_definition unless get_property_value(property, "custom_schema_definition", nil).nil?
          case property
          when Api::Azure::Type::ResourceGroupName
            !is_data_source ? 'templates/azure/terraform/schemas/resource_group_name.erb' : 'templates/azure/terraform/schemas/datasource_resource_group_name.erb'
          when Api::Azure::Type::Location
            !is_data_source ? 'templates/azure/terraform/schemas/location.erb' : 'templates/azure/terraform/schemas/datasource_location.erb'
          when Api::Azure::Type::Tags
            !is_data_source ? 'templates/azure/terraform/schemas/tags.erb' : 'templates/azure/terraform/schemas/datasource_tags.erb'
          when Api::Type::Boolean, Api::Type::Enum, Api::Type::String, Api::Type::Integer, Api::Type::Double,
               Api::Type::Array, Api::Type::KeyValuePairs, Api::Type::NestedObject
            'templates/azure/terraform/schemas/primitive.erb'
          else
            'templates/azure/terraform/schemas/unsupport.erb'
          end
        end

        def schema_property_get_template(property)
          return property.custom_schema_get unless get_property_value(property, "custom_schema_get", nil).nil?
          return 'templates/azure/terraform/schemas/hide_from_schema.erb' if get_property_value(property, "hide_from_schema", false)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location_get.erb'
          when Api::Type::Boolean, Api::Type::Enum, Api::Type::String, Api::Type::Integer, Api::Type::Double,
               Api::Type::Array, Api::Type::KeyValuePairs, Api::Type::NestedObject
            'templates/azure/terraform/schemas/basic_get.erb'
          else
            'templates/azure/terraform/schemas/unsupport.erb'
          end
        end

        def schema_property_set_template(property)
          return property.custom_schema_set unless get_property_value(property, "custom_schema_set", nil).nil?
          return 'templates/azure/terraform/schemas/hide_from_schema.erb' if get_property_value(property, "hide_from_schema", false)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location_set.erb'
          when Api::Azure::Type::Tags
            'templates/azure/terraform/schemas/tags_set.erb'
          when Api::Type::Boolean, Api::Type::Enum, Api::Type::String, Api::Type::Integer, Api::Type::Double, Api::Type::KeyValuePairs
            'templates/azure/terraform/schemas/basic_set.erb'
          when Api::Type::Array, Api::Type::NestedObject
            return 'templates/azure/terraform/schemas/string_array_set.erb' if property.is_a?(Api::Type::Array) && property.item_type_class == Api::Type::String
            'templates/azure/terraform/schemas/flatten_set.erb'
          else
            'templates/azure/terraform/schemas/unsupport.erb'
          end
        end

      end
    end
  end
end

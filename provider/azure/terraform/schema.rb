module Provider
  module Azure
    module Terraform
      module Schema

        def go_type(property)
          case property
          when Api::Type::String, Api::Azure::Type::Location
            'string'
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
          }
        end

        def schema_property_template(property)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location.erb'
          else
            'templates/terraform/schemas/unsupport.erb'
          end
        end

        def schema_property_get_template(property)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location_get.erb'
          else
            'templates/terraform/schemas/unsupport.erb'
          end
        end

        def schema_property_set_template(property)
          case property
          when Api::Azure::Type::Location
            'templates/azure/terraform/schemas/location_set.erb'
          else
            'templates/terraform/schemas/unsupport.erb'
          end
        end

        def will_property_set_dereference?(output_var, property)
          return true if output_var != 'd'
          return true if property.is_a? Api::Azure::Type::Location
          false
        end

      end
    end
  end
end

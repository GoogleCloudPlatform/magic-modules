module Provider
  module Azure
    module Terraform
      module Helpers
        def get_property_value(obj, prop_name, default_value)
          return default_value unless obj.instance_variable_defined?("@#{prop_name}")
          obj.instance_variable_get("@#{prop_name}")
        end

        def order_azure_properties(properties)
          special_props = properties.select{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName' || p.name == 'resourceGroup'}
          other_props = properties.reject{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName' || p.name == 'resourceGroup'}
          special_props.sort_by(&:order) + order_properties(other_props)
        end
      end
    end
  end
end

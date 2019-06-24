module Provider
  module Azure
    module Terraform
      module Helpers
        def get_property_value(obj, prop_name, default_value)
          return default_value unless obj.instance_variable_defined?("@#{prop_name}")
          obj.instance_variable_get("@#{prop_name}")
        end

        def order_azure_properties(properties, data_source_input = [])
          special_props = properties.select{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName' || p.name == 'resourceGroup' || data_source_input.include?(p)}
          other_props = properties.reject{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName' || p.name == 'resourceGroup' || data_source_input.include?(p)}
          sorted_special = special_props.sort_by{|p| p.name == 'location' ? 2 : p.order }
          sorted_other = data_source_input.empty? ? order_properties(other_props) : other_props.sort_by(&:name)
          sorted_special + sorted_other
        end
      end
    end
  end
end

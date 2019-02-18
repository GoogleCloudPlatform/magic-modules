require 'provider/azure/terraform/schema'
require 'provider/azure/terraform/sub_template'

module Provider
  module Azure
    module Terraform
      include Provider::Azure::Terraform::Schema
      include Provider::Azure::Terraform::SubTemplate

      def order_azure_properties(properties)
        special_props = properties.select{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName'}
        other_props = properties.reject{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName'}
        special_props + order_properties(other_props)
      end
    end
  end
end
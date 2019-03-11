require 'provider/azure/terraform/schema'
require 'provider/azure/terraform/sub_template'
require 'provider/azure/terraform/example'
require 'provider/azure/terraform/acctest/sub_template'

require 'provider/azure/terraform/resource_override'
require 'provider/azure/terraform/property_override'

module Provider
  module Azure
    module Terraform
      include Provider::Azure::Terraform::Schema
      include Provider::Azure::Terraform::SubTemplate
      include Provider::Azure::Terraform::Example::SubTemplate
      include Provider::Azure::Terraform::AccTest::SubTemplate

      def azure_resource_go_package(product)
        product.azure_namespace.split('.').last.downcase
      end

      def order_azure_properties(properties)
        special_props = properties.select{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName'}
        other_props = properties.reject{|p| p.name == 'name' || p.name == 'location' || p.name == 'resourceGroupName'}
        special_props + order_properties(other_props)
      end
    end
  end
end
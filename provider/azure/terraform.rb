require 'provider/azure/terraform/helpers'
require 'provider/azure/terraform/schema'
require 'provider/azure/terraform/sub_template'
require 'provider/azure/terraform/sdk/expand_flatten_descriptor'
require 'provider/azure/terraform/sdk/sub_template'
require 'provider/azure/terraform/sdk/helpers'
require 'provider/azure/terraform/example/example'
require 'provider/azure/terraform/example/sub_template'
require 'provider/azure/terraform/example/helpers'
require 'provider/azure/terraform/acctest/sub_template'

require 'provider/azure/terraform/resource_override'
require 'provider/azure/terraform/property_override'

module Provider
  module Azure
    module Terraform
      include Provider::Azure::Terraform::Helpers
      include Provider::Azure::Terraform::Schema
      include Provider::Azure::Terraform::SubTemplate
      include Provider::Azure::Terraform::SDK::SubTemplate
      include Provider::Azure::Terraform::SDK::Helpers
      include Provider::Azure::Terraform::Example::SubTemplate
      include Provider::Azure::Terraform::Example::Helpers
      include Provider::Azure::Terraform::AccTest::SubTemplate
    end
  end
end
require 'provider/config'

module Provider
  module Azure
    module Terraform

      class Config < Provider::Config
        def provider
          Provider::Terraform
        end

        def resource_override
          Provider::Azure::Terraform::ResourceOverride
        end

        def property_override
          Provider::Azure::Terraform::PropertyOverride
        end
      end

    end
  end
end

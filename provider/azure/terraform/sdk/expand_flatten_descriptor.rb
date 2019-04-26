require 'provider/azure/terraform/sdk/sdk_type_definition_descriptor'

module Provider
  module Azure
    module Terraform
      module SDK
        class ExpandFlattenDescriptor
          attr_reader :property
          attr_reader :sdkmarshal

          def initialize(property, sdkmarshal)
            @property = property
            @sdkmarshal = sdkmarshal
          end
        end
      end
    end
  end
end

module Provider
  module Azure
    module Terraform
      module SDK
        class ExpandFlattenDescriptor
          attr_reader :property
          attr_reader :api_path
          attr_reader :sdk_type
          attr_reader :sdk_type_defs

          def initialize(property, api_path, sdk_type_defs)
            @property = property
            @api_path = api_path
            @sdk_type_defs = sdk_type_defs
            @sdk_type = sdk_type_defs[api_path]
          end
        end
      end
    end
  end
end
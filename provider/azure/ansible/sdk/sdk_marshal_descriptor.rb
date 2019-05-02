module Provider
  module Azure
    module Ansible
      module SDK
        class MarshalDescriptor
          attr_reader :properties
          attr_reader :operation
          attr_reader :parent_reference
          attr_reader :input
          attr_reader :output

          def initialize(properties, sdk_operation, input, output, parent_sdk_reference = '/')
            @properties = properties
            @operation = sdk_operation
            @input = input
            @output = output
            @parent_reference = parent_sdk_reference
          end
        end
      end
    end
  end
end

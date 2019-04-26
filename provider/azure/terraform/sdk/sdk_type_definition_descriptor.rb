module Provider
  module Azure
    module Terraform
      module SDK
        class TypeDefinitionDescriptor
          attr_reader :operation
          attr_reader :typedef_reference

          def initialize(operation, isRequest, reference = nil)
            @operation = operation
            @isRequest = isRequest
            @typedef_reference = reference || (@isRequest ? '/' : '')
          end

          def clone(typedef_reference = nil)
            TypeDefinitionDescriptor.new @operation, @isRequest, (typedef_reference || @typedef_reference)
          end

          def type_definitions
            return @operation.request if @isRequest
            @operation.response.has_key?(@typedef_reference) ? @operation.response : @operation.request
          end

          def type_definition
            type_definitions[@typedef_reference]
          end

          def go_type_name
            return type_definition.go_type_name unless type_definition.nil?
            nil
          end

          def go_field_name
            return type_definition.go_field_name unless type_definition.nil?
            nil
          end
        end
      end
    end
  end
end

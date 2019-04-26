require 'provider/azure/terraform/sdk/sdk_type_definition_descriptor'

module Provider
  module Azure
    module Terraform
      module SDK
        class MarshalDescriptor
          attr_reader :package
          attr_reader :resource
          attr_reader :queue
          attr_reader :sdktype
          attr_reader :properties

          def initialize(package, resource, queue, sdktype, properties)
            @package = package
            @resource = resource
            @queue = queue
            @sdktype = sdktype
            @properties = properties
          end

          def clone(typedef_reference = nil, properties = nil)
            sdktype = @sdktype.clone(typedef_reference)
            MarshalDescriptor.new @package, @resource, @queue, sdktype, (properties || @properties)
          end

          def enqueue(property)
            existed = @queue.any?{|e| e.property == property && e.sdktype.go_type_name == @sdktype.go_type_name}
            @queue << ExpandFlattenDescriptor.new(property, self) unless existed
          end
        end
      end
    end
  end
end

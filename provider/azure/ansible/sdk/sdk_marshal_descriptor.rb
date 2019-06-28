module Provider
  module Azure
    module Ansible
      module SDK

        class MarshalDescriptor
          attr_reader :properties
          attr_reader :operation
          attr_reader :parent_reference
          attr_reader :marshalled_references
          attr_reader :input
          attr_reader :output

          def initialize(properties, sdk_operation, input, output, parent_sdk_reference = '/', marshalled = nil)
            @properties = properties
            @operation = sdk_operation
            @input = input
            @output = output
            @parent_reference = parent_sdk_reference
            @marshalled_references = marshalled || { '/' => @output }
          end

          def add_marshalled_reference(reference, expression)
            @marshalled_references[reference] = expression
          end

          def create_child_descriptor(property, sdk_reference, sdk_type)
            input_expression = "#{@input}['#{sdk_type.python_field_name}']"
            MarshalDescriptor.new property.nested_properties, @operation, input_expression, @output, sdk_reference, @marshalled_references
          end
        end

      end
    end
  end
end

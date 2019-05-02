module Provider
    module Azure
      module Ansible
        module SDK
          class PropertyNormalizeDescriptor
            attr_reader :property
            attr_reader :input
            attr_reader :output
            attr_reader :remove_from_input
            attr_reader :python_variable
            attr_reader :parent_reference
            attr_reader :relative_references
            attr_reader :sdk_operation
  
            def initialize(property, input, output, python_var, parent_ref, relative_refs, sdk_op, remove_input = false)
              @property = property
              @input = input
              @output = output
              @python_variable = python_var
              @parent_reference = parent_ref
              @relative_references = relative_refs
              @sdk_operation = sdk_op
              @remove_from_input = remove_input
            end

            def create_child()
              PropertyNormalizeDescriptor.new @property, @input, '|', @python_variable, @parent_reference + @relative_references[0] + '/', @relative_references[1..-1], @sdk_operation, true
            end
          end
        end
      end
    end
  end
  
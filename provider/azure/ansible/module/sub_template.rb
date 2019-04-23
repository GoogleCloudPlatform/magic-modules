module Provider
  module Azure
    module Ansible
      module Module
        module SubTemplate
          def build_class_instance_variable_init(sdk_operation, object, indentation = 8)
            result = compile 'templates/azure/ansible/module/class_instance_variable_init.erb', 1
            indent result, indentation
          end

          def build_response_properties_update(properties, sdk_response_def, indentation = 16)
            result = compile 'templates/azure/ansible/module/response_properties_update.erb', 1
            indent_list result, indentation
          end
        end
      end
    end
  end
end
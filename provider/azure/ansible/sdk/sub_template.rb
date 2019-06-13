module Provider
  module Azure
    module Ansible
      module SDK
        module SubTemplate
          def build_sdk_method_invocation(sdk_client, sdk_op_def, indentation = 12)
            result = compile 'templates/azure/ansible/sdk/method_invocation.erb', 1
            indent result, indentation
          end

          def build_property_normalization(norm_desc, in_structure, indentation = 4)
            template = get_custom_template_path(norm_desc.property.custom_normalize)
            template ||= 'templates/azure/ansible/sdktypes/property_normalization.erb'
            result = compile template, 1
            indent result, indentation
          end

          def build_property_to_sdk_object(sdk_marshal, indentation = 0)
            result = compile 'templates/azure/ansible/sdktypes/property_to_sdkobject.erb', 1
            indent result, indentation
          end

          def build_property_inline_response_format(property, sdk_operation, indentation = 12)
            template = get_custom_template_path(property.inline_custom_response_format)
            template ||= 'templates/azure/ansible/sdktypes/property_inline_response_format.erb'
            result = compile template, 1
            indent result, indentation
          end
        end
      end
    end
  end
end

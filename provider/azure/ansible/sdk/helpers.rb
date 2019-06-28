module Provider
  module Azure
    module Ansible
      module SDK

        module Helpers
          def get_parent_reference(reference)
            parent = reference[0..reference.rindex('/')]
            return parent if parent == '/'
            parent.chomp '/'
          end

          def get_properties_matching_sdk_reference(sdk_reference, object)
            object.all_user_properties
              .select{|p| p.azure_sdk_references.include?(sdk_reference)}
              .sort_by{|p| [p.order, p.name]}
          end

          def get_applicable_reference(references, typedefs)
            references.each do |ref|
              return ref if typedefs.has_key?(ref)
            end
            nil
          end

          def get_sdk_typedef_by_references(references, typedefs)
            ref = get_applicable_reference(references, typedefs)
            return nil if ref.nil?
            typedefs[ref]
          end

          def self_require_type_marshal?(property, sdk_marshal)
            return true if is_location? property

            sdk_ref = get_applicable_reference(property.azure_sdk_references, sdk_marshal.operation.request)
            return false if !sdk_ref.start_with?('/')

            return true if property.is_a? Api::Type::Enum
            return true if property.is_a? Api::Azure::Type::ResourceReference

            var_name = azure_python_variable_name(property, sdk_marshal.operation)
            sdk_type = get_sdk_typedef_by_references(property.azure_sdk_references, sdk_marshal.operation.request)
            return true if var_name != sdk_type.python_field_name
            false
          end

          def descendants_require_type_marshal?(property, sdk_marshal)
            return property.nested_properties.any?{|p| require_type_marshal?(p, sdk_marshal)} if property.nested_properties?
            false
          end

          def require_type_marshal?(property, sdk_marshal)
            self_require_type_marshal?(property, sdk_marshal) || descendants_require_type_marshal?(property, sdk_marshal)
          end

          def property_normalization_template(property)
            return get_custom_template_path(property.custom_normalize) if property.custom_normalize
            return 'templates/azure/ansible/sdktypes/location_property_normalize.erb' if is_location? property
            return 'templates/azure/ansible/sdktypes/property_normalize.erb' if property.is_a? Api::Type::Enum
            return 'templates/azure/ansible/sdktypes/property_normalize.erb' if property.is_a? Api::Azure::Type::ResourceReference
            'templates/azure/ansible/sdktypes/unsupported.erb'
          end
        end

      end
    end
  end
end

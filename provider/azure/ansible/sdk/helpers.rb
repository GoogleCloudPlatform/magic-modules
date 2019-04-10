module Provider
  module Azure
    module Ansible
      module SDK
        module Helpers
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
        end
      end
    end
  end
end

require 'overrides/ansible/property_override'

module Provider
  module Azure
    module Ansible

      class PropertyOverride < Overrides::Ansible::PropertyOverride
        # Collection of fields allowed in the PropertyOverride section for
        # Ansible. All fields should be `attr_reader :<property>`
        def self.attributes
          super.concat(%i[
            resource_type_name
            document_sample_value
            custom_normalize
            inline_custom_response_format
          ])
        end

        attr_reader(*attributes)

        def validate
          super
          check :resource_type_name, type: ::String
          check :document_sample_value, type: ::String
          check :custom_normalize, type: ::String
          check :inline_custom_response_format, type: ::String
        end
      end

    end
  end
end

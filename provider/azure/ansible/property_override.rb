require 'provider/ansible/property_override'

module Provider
  module Azure
    module Ansible
      module OverrideFields
        attr_reader :resource_type_name
        attr_reader :custom_normalize
        attr_reader :inline_custom_response_format
        include Provider::Ansible::OverrideFields
      end

      class PropertyOverride < Provider::Ansible::PropertyOverride
        include Provider::Azure::Ansible::OverrideFields

        def validate
          super
          check_optional_property :resource_type_name, String
          check_optional_property :custom_normalize, String
          check_optional_property :inline_custom_response_format, String
        end

        private

        def overriden
          Provider::Azure::Ansible::OverrideFields
        end

      end
    end
  end
end

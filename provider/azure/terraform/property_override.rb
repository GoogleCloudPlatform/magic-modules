require 'provider/terraform/property_override'

module Provider
  module Azure
    module Terraform
      module OverrideFields
        attr_reader :name_in_logs
        attr_reader :hide_from_schema
        attr_reader :sdkfield_assign_type
        attr_reader :custom_schema_definition
        attr_reader :custom_schema_get
        attr_reader :custom_schema_set
        attr_reader :custom_sdkfield_assign
        include Provider::Terraform::OverrideFields
      end

      class PropertyOverride < Provider::Terraform::PropertyOverride
        include Provider::Azure::Terraform::OverrideFields

        def validate
          super
          @hide_from_schema ||= false
          check_optional_property :name_in_logs, String
          check_optional_property :hide_from_schema, :boolean
          check_optional_property :custom_schema_definition, String
          check_optional_property :custom_schema_get, String
          check_optional_property :custom_schema_set, String
          check_optional_property :custom_sdkfield_assign, String
          check_property_oneof_default :sdkfield_assign_type, ['inline', 'block'], 'inline'
        end

        private

        def overriden
          Provider::Azure::Terraform::OverrideFields
        end

      end
    end
  end
end

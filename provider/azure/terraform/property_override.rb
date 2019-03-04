require 'provider/terraform/property_override'

module Provider
  module Azure
    module Terraform
      module OverrideFields
        attr_reader :custom_schema_definition
        attr_reader :custom_schema_get
        attr_reader :custom_schema_set
        include Provider::Terraform::OverrideFields
      end

      class PropertyOverride < Provider::Terraform::PropertyOverride
        include Provider::Azure::Terraform::OverrideFields

        def validate
          super
          check_optional_property :custom_schema_definition, String
          check_optional_property :custom_schema_get, String
          check_optional_property :custom_schema_set, String
        end

        private

        def overriden
          Provider::Azure::Terraform::OverrideFields
        end

      end
    end
  end
end

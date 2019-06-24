require 'overrides/terraform/property_override'

module Provider
  module Azure
    module Terraform

      class PropertyOverride < Overrides::Terraform::PropertyOverride
        def self.attributes
          super.concat(%i[
            name_in_logs
            hide_from_schema
            custom_schema_definition
            custom_schema_get
            custom_schema_set
            custom_sdkfield_assign
          ])
        end

        attr_reader(*attributes)

        def validate
          super
          check :name_in_logs, type: ::String
          check :hide_from_schema, type: :boolean, default: false
          check :custom_schema_definition, type: ::String
          check :custom_schema_get, type: ::String
          check :custom_schema_set, type: ::String
          check :custom_sdkfield_assign, type: ::String
        end
      end

    end
  end
end

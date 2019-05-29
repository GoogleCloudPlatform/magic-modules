require 'provider/ansible/resource_override'

module Provider
  module Azure
    module Ansible
      module OverrideProperties
        attr_reader :inttests
        attr_reader :examples
        include Provider::Ansible::OverrideProperties
      end

      class ResourceOverride < Provider::Ansible::ResourceOverride
        include Provider::Azure::Ansible::OverrideProperties

        def validate
          super
          default_value_property :examples, []
          check_optional_property :examples, Array
          check_optional_property_list :examples, Provider::Azure::ExampleReference
          default_value_property :inttests, []
          check_optional_property :inttests, Array
          check_optional_property_list :inttests, IntegrationTestDefinition
        end

        class IntegrationTestDefinition < ExampleReference
          def validate
            super
          end
        end

        private

        def overriden
          Provider::Azure::Ansible::OverrideProperties
        end
      end
    end
  end
end

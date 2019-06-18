require 'provider/ansible/resource_override'
require 'provider/azure/example/example'

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
          attr_reader :delete_example
          attr_reader :info_by_name_example
          attr_reader :info_by_resource_group_example

          def validate
            super
            check_property :delete_example, String
            check_optional_property :info_by_name_example, String
            check_optional_property :info_by_resource_group_example, String
          end
        end

        class DocumentExampleReference < ExampleReference
          attr_reader :resource_name_hints
  
          def validate
            super
            check_optional_property :resource_name_hints, Hash
            check_optional_property_hash :resource_name_hints, String, String
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

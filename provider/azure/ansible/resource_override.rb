require 'overrides/ansible/resource_override'
require 'provider/azure/example/example'

module Provider
  module Azure
    module Ansible

      class IntegrationTestDefinition < ExampleReference
        attr_reader :delete_example
        attr_reader :info_by_name_example
        attr_reader :info_by_resource_group_example

        def validate
          super
          check :delete_example, type: ::String, required: true
          check :info_by_name_example, type: ::String
          check :info_by_resource_group_example, type: ::String
        end
      end

      class DocumentExampleReference < ExampleReference
        attr_reader :resource_name_hints

        def validate
          super
          check_ext :resource_name_hints, type: ::Hash, key_type: ::String, item_type: ::String
        end
      end

      class ResourceOverride < Overrides::Ansible::ResourceOverride
        def self.attributes
          super.concat(%i[
            azure_sdk_definition
            inttests
            examples
          ])
        end

        attr_reader(*attributes)

        def validate
          super
          check :examples, type: ::Array, default: [], item_type: DocumentExampleReference
          check :inttests, type: ::Array, default: [], item_type: IntegrationTestDefinition
        end

        def apply(_resource)
          filter_azure_sdk_language _resource, "python"
          merge_azure_sdk_definition _resource, @azure_sdk_definition
          @azure_sdk_definition = nil
          super
        end
      end

    end
  end
end

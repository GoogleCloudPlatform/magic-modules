require 'provider/terraform/resource_override'
require 'api/azure/sdk_definition_override'

module Provider
  module Azure
    module Terraform
      module OverrideProperties
        attr_reader :name_in_logs
        attr_reader :azure_sdk_definition
        attr_reader :acctests
        attr_reader :document_examples
        attr_reader :datasource_example_outputs
        include Provider::Terraform::OverrideProperties
      end

      class ResourceOverride < Provider::Terraform::ResourceOverride
        include Provider::Azure::Terraform::OverrideProperties

        def validate
          super
          @acctests ||= Array.new
          check_optional_property :name_in_logs, String
          check_optional_property :azure_sdk_definition, Api::Azure::SDKDefinitionOverride
          check_optional_property :acctests, Array
          check_optional_property_list :acctests, AccTestDefinition
          check_optional_property :document_examples, Array
          check_optional_property_list :document_examples, DocumentExampleReference
          check_optional_property :datasource_example_outputs, Hash
        end

        class DocumentExampleReference < Api::Object
          attr_reader :title
          attr_reader :example_name
          attr_reader :resource_name_hints

          def validate
            super
            check_property :title, String
            check_property :example_name, String
            check_optional_property :resource_name_hints, Hash
            check_optional_property_hash :resource_name_hints, String, String
          end
        end

        class DataSourceExampleReference < Api::Object
          attr_reader :title
          attr_reader :example_name

          def validate
            super
            check_property :title, String
            check_property :example_name, String
          end
        end

        class AccTestDefinition < Api::Object
          attr_reader :name
          attr_reader :steps

          def validate
            super
            @initialized = false

            check_property :name, String
            check_property :steps, Array
            check_property_list :steps, String
          end
        end

        private

        def overriden
          Provider::Azure::Terraform::OverrideProperties
        end
      end
    end
  end
end

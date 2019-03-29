require 'provider/terraform/resource_override'

module Provider
  module Azure
    module Terraform
      module OverrideProperties
        attr_reader :name_in_logs
        attr_reader :acctests
        attr_reader :document_examples
        include Provider::Terraform::OverrideProperties
      end

      class ResourceOverride < Provider::Terraform::ResourceOverride
        include Provider::Azure::Terraform::OverrideProperties

        def validate
          super
          @acctests ||= Array.new
          check_optional_property :name_in_logs, String
          check_optional_property :acctests, Array
          check_optional_property_list :acctests, AccTestDefinition
          check_optional_property :document_examples, Array
          check_optional_property_list :document_examples, DocumentExampleReference
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

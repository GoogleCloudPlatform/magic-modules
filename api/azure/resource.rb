require 'google/yaml_validator'
require 'api/azure/sdk_definition'

module Api
  module Azure
    module Resource

      # The Azure-extended properties which supplement Api::Resource::Properties
      module Properties
        attr_reader :azure_sdk_definition
      end

      # Azure-extended validate function of Api::Resource::validate
      def azure_validate
        check_property :azure_sdk_definition, Api::Azure::SDKDefinition
      end

    end
  end
end

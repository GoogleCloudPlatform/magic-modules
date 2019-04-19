require 'api/object'
require 'api/azure/sdk_operation_definition_override'

module Api
  module Azure
    class SDKDefinitionOverride < Api::Object
      attr_reader :create
      attr_reader :read
      attr_reader :update
      attr_reader :delete

      def validate
        super
        check_optional_property :create, Api::Azure::SDKOperationDefinitionOverride
        check_optional_property :read, Api::Azure::SDKOperationDefinitionOverride
        check_optional_property :update, Api::Azure::SDKOperationDefinitionOverride
        check_optional_property :delete, Api::Azure::SDKOperationDefinitionOverride
      end
    end
  end
end
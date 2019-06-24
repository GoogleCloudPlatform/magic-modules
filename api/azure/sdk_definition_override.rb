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
        check :create, type: SDKOperationDefinitionOverride
        check :read, type: SDKOperationDefinitionOverride
        check :update, type: SDKOperationDefinitionOverride
        check :delete, type: SDKOperationDefinitionOverride
      end
    end
  end
end

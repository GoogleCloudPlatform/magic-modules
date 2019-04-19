require 'api/object'
require 'api/azure/sdk_type_definition_override'

module Api
  module Azure
    class SDKOperationDefinitionOverride < Api::Object
      attr_reader :request
      attr_reader :response

      def validate
        super
        check_optional_property :request, Hash
        check_optional_property_hash :request, String, Api::Azure::SDKTypeDefinitionOverride
        check_optional_property :response, Hash
        check_optional_property_hash :response, String, Api::Azure::SDKTypeDefinitionOverride
      end
    end
  end
end
require 'api/object'
require 'api/azure/sdk_type_definition'

module Api
  module Azure
    class SDKOperationDefinition < Api::Object
      attr_reader :go_func_name
      attr_reader :async
      attr_reader :request
      attr_reader :response

      def validate
        super
        check_property :go_func_name, String
        check_optional_property :async, :boolean
        check_property :request, Hash
        check_property_hash :request, String, Api::Azure::SDKTypeDefinition
        check_optional_property :response, Hash
        check_optional_property_hash :response, String, Api::Azure::SDKTypeDefinition
      end
    end
  end
end

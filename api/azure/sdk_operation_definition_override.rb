require 'api/object'
require 'api/azure/sdk_type_definition_override'

module Api
  module Azure
    class SDKOperationDefinitionOverride < Api::Object
      attr_reader :request
      attr_reader :response

      def validate
        super
        check_ext :request, type: ::Hash, key_type: ::String, item_type: SDKTypeDefinitionOverride
        check_ext :response, type: ::Hash, key_type: ::String, item_type: SDKTypeDefinitionOverride
      end
    end
  end
end

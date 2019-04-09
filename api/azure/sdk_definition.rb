require 'api/object'
require 'api/azure/sdk_operation_definition'

module Api
  module Azure
    class SDKDefinition < Api::Object
      attr_reader :provider_name
      attr_reader :go_client_namespace
      attr_reader :go_client
      attr_reader :python_client_namespace
      attr_reader :python_client
      attr_reader :create
      attr_reader :read
      attr_reader :update
      attr_reader :delete

      def validate
        super
        check_optional_property :provider_name, String
        check_optional_property :go_client_namespace, String
        check_optional_property :go_client, String
        check_optional_property :python_client_namespace, String
        check_optional_property :python_client, String
        check_property :create, Api::Azure::SDKOperationDefinition
        check_property :read, Api::Azure::SDKOperationDefinition
        check_optional_property :update, Api::Azure::SDKOperationDefinition
        check_property :delete, Api::Azure::SDKOperationDefinition
      end

      def merge_overrides(overrides, language)
        @create.merge_overrides(nil, language) unless @create.nil?
        @read.merge_overrides(nil, language) unless @read.nil?
        @update.merge_overrides(nil, language) unless @update.nil?
        @delete.merge_overrides(nil, language) unless @delete.nil?
      end
    end
  end
end

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
      attr_reader :list_by_resource_group
      attr_reader :list_by_subscription

      def validate
        super
        check_property :provider_name, String
        check_property :go_client_namespace, String
        check_property :go_client, String
        check_property :python_client_namespace, String
        check_property :python_client, String
        check_property :create, Api::Azure::SDKOperationDefinition
        check_property :read, Api::Azure::SDKOperationDefinition
        check_optional_property :update, Api::Azure::SDKOperationDefinition
        check_property :delete, Api::Azure::SDKOperationDefinition
        check_optional_property :list_by_resource_group, Api::Azure::SDKOperationDefinition
        check_optional_property :list_by_subscription, Api::Azure::SDKOperationDefinition
      end

      def filter_language!(language)
        @create.filter_language!(language) unless @create.nil?
        @read.filter_language!(language) unless @read.nil?
        @update.filter_language!(language) unless @update.nil?
        @delete.filter_language!(language) unless @delete.nil?
      end

      def merge_overrides!(overrides)
        @create.merge_overrides!(overrides.create) if !@create.nil? && !overrides.create.nil?
        @read.merge_overrides!(nil) if !@read.nil? && !overrides.read.nil?
        @update.merge_overrides!(nil) if !@update.nil? && !overrides.update.nil?
        @delete.merge_overrides!(nil) if !@delete.nil? && !overrides.delete.nil?
      end
    end
  end
end

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
        check :provider_name, type: ::String, required: true
        check :go_client_namespace, type: ::String, required: true
        check :go_client, type: ::String, required: true
        check :python_client_namespace, type: ::String, required: true
        check :python_client, type: ::String, required: true
        check :create, type: Api::Azure::SDKOperationDefinition, required: true
        check :read, type: Api::Azure::SDKOperationDefinition, required: true
        check :update, type: Api::Azure::SDKOperationDefinition
        check :delete, type: Api::Azure::SDKOperationDefinition, required: true
        check :list_by_parent, type: Api::Azure::SDKOperationDefinition
        check :list_by_resource_group, type: Api::Azure::SDKOperationDefinition
        check :list_by_subscription, type: Api::Azure::SDKOperationDefinition
      end

      def filter_language!(language)
        @create.filter_language!(language) unless @create.nil?
        @read.filter_language!(language) unless @read.nil?
        @update.filter_language!(language) unless @update.nil?
        @delete.filter_language!(language) unless @delete.nil?
      end

      def merge_overrides!(overrides)
        @create.merge_overrides!(overrides.create) if !@create.nil? && !overrides.create.nil?
        @read.merge_overrides!(overrides.read) if !@read.nil? && !overrides.read.nil?
        @update.merge_overrides!(overrides.update) if !@update.nil? && !overrides.update.nil?
        @delete.merge_overrides!(overrides.delete) if !@delete.nil? && !overrides.delete.nil?
      end
    end
  end
end

require 'api/object'
require 'api/azure/sdk_type_definition'

module Api
  module Azure
    class SDKOperationDefinition < Api::Object
      attr_reader :go_func_name
      attr_reader :python_func_name
      attr_reader :async
      attr_reader :request
      attr_reader :response

      def validate
        super
        check_property :go_func_name, String
        check_property :python_func_name, String
        check_optional_property :async, :boolean
        check_property :request, Hash
        check_property_hash :request, String, Api::Azure::SDKTypeDefinition
        check_optional_property :response, Hash
        check_optional_property_hash :response, String, Api::Azure::SDKTypeDefinition
      end

      def merge_overrides(overrides, language)
        filter_applicable(@request, language) unless @request.nil?
        filter_applicable(@response, language) unless @response.nil?
      end

      private

      def filter_applicable(fields, language)
        fields.reject!{|name, value| !value.applicable_to.nil? && !value.applicable_to.include?(language)}
      end
    end
  end
end

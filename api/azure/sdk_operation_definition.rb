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
        @request ||= Hash.new
        @response ||= Hash.new

        check :go_func_name, type: ::String, required: true
        check :python_func_name, type: ::String, required: true
        check :async, type: :boolean
        check_ext :request, type: ::Hash, key_type: ::String, item_type: SDKTypeDefinition, required: true
        check_ext :response, type: ::Hash, key_type: ::String, item_type: SDKTypeDefinition, required: true
      end

      def filter_language!(language)
        filter_applicable! @request, language
        filter_applicable!(@response, language) unless @response.nil?
      end

      def merge_overrides!(overrides)
        merge_hash_table!(@request, overrides.request) unless overrides.request.nil?
        merge_hash_table!(@response, overrides.response) unless overrides.response.nil?
      end

      private

      def filter_applicable!(fields, language)
        fields.reject!{|name, value| !value.applicable_to.nil? && !value.applicable_to.include?(language)}
      end

      def merge_hash_table!(fields, overrides)
        overrides.each do |name, value|
          if value.remove
            fields.delete(name)
          elsif !fields.has_key?(name)
            fields[name] = value
          else
            fields[name].merge_overrides! value
          end
        end
      end
    end
  end
end

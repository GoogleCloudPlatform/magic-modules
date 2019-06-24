module Overrides
  module Azure

    module ResourceOverrideExtension
      def filter_azure_sdk_language(_resource, language)
        _resource.azure_sdk_definition.filter_language!(language) unless _resource.azure_sdk_definition.nil?
      end

      def merge_azure_sdk_definition(_resource, overrides)
        _resource.azure_sdk_definition.merge_overrides!(overrides) unless _resource.azure_sdk_definition.nil? || overrides.nil?
      end
    end

  end
end

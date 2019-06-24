require 'api/object'

module Api
  module Azure
    class SDKTypeDefinitionOverride < SDKTypeDefinition
      attr_reader :remove

      def validate
        super
        check :remove, type: :boolean, default: false
      end
    end
  end
end

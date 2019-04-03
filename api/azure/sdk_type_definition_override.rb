require 'api/object'

module Api
  module Azure
    class SDKTypeDefinitionOverride < Api::Object
      attr_reader :id_portion
      attr_reader :go_variable_name
      attr_reader :go_field_name
      attr_reader :python_variable_name

      def validate
        super
        check_optional_property :id_portion
        check_optional_property :go_variable_name
        check_optional_property :go_field_name
        check_optional_property :python_variable_name
      end
    end
  end
end
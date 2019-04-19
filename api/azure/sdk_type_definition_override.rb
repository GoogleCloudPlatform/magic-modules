require 'api/object'

module Api
  module Azure
    class SDKTypeDefinitionOverride < Api::Object
      attr_reader :remove
      attr_reader :id_portion
      attr_reader :empty_value_sensitive
      attr_reader :go_variable_name
      attr_reader :go_field_name
      attr_reader :python_parameter_name
      attr_reader :python_variable_name

      def validate
        super
        @remove ||= false

        check_optional_property :remove, :boolean
        check_optional_property :id_portion, String
        check_optional_property :empty_value_sensitive, :boolean
        check_optional_property :go_variable_name, String
        check_optional_property :go_field_name, String
        check_optional_property :python_parameter_name, String
        check_optional_property :python_variable_name, String
      end
    end
  end
end
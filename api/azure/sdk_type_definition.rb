require 'api/object'

module Api
  module Azure
    class SDKTypeDefinition < Api::Object
      attr_reader :id_portion
      attr_reader :applicable_to
      attr_reader :empty_value_sensitive
      attr_reader :go_variable_name
      attr_reader :go_field_name
      attr_reader :python_parameter_name
      attr_reader :python_variable_name

      def validate
        super
        @empty_value_sensitive ||= false

        check_optional_property :id_portion, String
        check_optional_property :applicable_to, Array
        check_property :empty_value_sensitive, :boolean
        check_optional_property_list_oneof :applicable_to, ['go', 'python'], String
        check_optional_property :go_variable_name, String
        check_optional_property :go_field_name, String
        check_optional_property :python_parameter_name, String
        check_optional_property :python_variable_name, String
      end

      def merge_overrides!(overrides)
        @id_portion = overrides.id_portion unless overrides.id_portion.nil?
        @empty_value_sensitive = overrides.empty_value_sensitive unless overrides.empty_value_sensitive.nil?
        @go_variable_name = overrides.go_variable_name unless overrides.go_variable_name.nil?
        @go_field_name = overrides.go_field_name unless overrides.go_field_name.nil?
        @python_parameter_name = overrides.python_parameter_name unless overrides.python_parameter_name.nil?
        @python_variable_name = overrides.python_variable_name unless overrides.python_variable_name.nil?
      end

      class BooleanObject < SDKTypeDefinition
      end

      class StringObject < SDKTypeDefinition
      end

      class EnumObject < SDKTypeDefinition
        attr_reader :go_enum_type_name

        def validate
          super
          check_property :go_enum_type_name, String
        end
      end

      class ComplexObject < SDKTypeDefinition
        attr_reader :go_type_name

        def validate
          super
          check_optional_property :go_type_name, String
        end
      end

    end
  end
end

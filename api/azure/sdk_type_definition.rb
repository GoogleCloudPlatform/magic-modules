require 'api/object'

module Api
  module Azure
    class SDKTypeDefinition < Api::Object
      attr_reader :id_portion
      attr_reader :go_parameter_expression
      attr_reader :go_field_name

      def validate
        super
        check_optional_property :id_portion
        check_optional_property :go_parameter_expression
        check_optional_property :go_field_name
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
          check_property :go_type_name, String
        end
      end

    end
  end
end

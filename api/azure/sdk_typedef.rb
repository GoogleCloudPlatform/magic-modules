require 'api/object'

module Api
  module Azure
    class SDKTypeDefinitions < Api::Object
      attr_reader :create
      attr_reader :read
      attr_reader :update

      def validate
        super
        check_property :create, Hash
        check_property :read, Hash
        check_property :update, Hash
        check_property_hash :create, String, Api::Azure::SDKTypeDefinition
        check_property_hash :read, String, Api::Azure::SDKTypeDefinition
        check_property_hash :update, String, Api::Azure::SDKTypeDefinition
      end
    end

    class SDKTypeDefinition < Api::Object
      attr_reader :go_field_name

      def validate
        super
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
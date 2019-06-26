require 'api/object'

module Api
  module Azure
    class SDKTypeDefinition < Api::Object
      attr_reader :id_portion
      attr_reader :applicable_to
      attr_reader :empty_value_sensitive
      attr_reader :go_variable_name
      attr_reader :go_field_name
      attr_reader :go_type_name
      attr_reader :python_parameter_name
      attr_reader :python_variable_name
      attr_reader :python_field_name

      def validate
        super
        check :id_portion, type: ::String
        check_ext :applicable_to, type: ::Array, item_type: ::String, item_allowed: ['go', 'python'], default: ['go', 'python']
        check :empty_value_sensitive, type: :boolean, default: false
        check :go_variable_name, type: ::String
        check :go_field_name, type: ::String
        check :go_type_name, type: ::String
        check :python_parameter_name, type: ::String
        check :python_variable_name, type: ::String
        check :python_field_name, type: ::String
      end

      def merge_overrides!(overrides)
        @id_portion = overrides.id_portion unless overrides.id_portion.nil?
        @empty_value_sensitive = overrides.empty_value_sensitive unless overrides.empty_value_sensitive.nil?
        @go_variable_name = overrides.go_variable_name unless overrides.go_variable_name.nil?
        @go_field_name = overrides.go_field_name unless overrides.go_field_name.nil?
        @python_parameter_name = overrides.python_parameter_name unless overrides.python_parameter_name.nil?
        @python_variable_name = overrides.python_variable_name unless overrides.python_variable_name.nil?
        @python_field_name = overrides.python_field_name unless overrides.python_field_name.nil?
      end

      class BooleanObject < SDKTypeDefinition
      end

      class IntegerObject < SDKTypeDefinition
      end

      class Integer32Object < SDKTypeDefinition
      end

      class Integer64Object < SDKTypeDefinition
      end

      class FloatObject < SDKTypeDefinition
      end

      class StringObject < SDKTypeDefinition
      end

      class EnumObject < SDKTypeDefinition
        attr_reader :go_enum_type_name
        attr_reader :go_enum_const_prefix

        def validate
          super
          check :go_enum_type_name, type: ::String
          check :go_enum_const_prefix, type: ::String, default: ''
        end
      end

      class ISO8601DurationObject < StringObject
      end

      class ISO8601DateTimeObject < SDKTypeDefinition
      end

      class ComplexObject < SDKTypeDefinition
      end

      class StringArrayObject < SDKTypeDefinition
      end

      class ComplexArrayObject < ComplexObject
      end

      class StringMapObject < SDKTypeDefinition
      end

    end
  end
end

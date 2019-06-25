require 'google/python_utils'
require 'azure/python_utils'

module Provider
  module Azure
    module Ansible
      module ModuleExtension

        include Google::PythonUtils
        include ::Azure::PythonUtils

        def azure_python_dict_for_property(prop, object, spaces = 0)
          if prop.is_a?(Api::Type::Array) && \
            prop.item_type.is_a?(Api::Type::NestedObject)
            azure_nested_obj_dict(prop, object, prop.item_type.properties, spaces)
          elsif prop.is_a? Api::Type::NestedObject
            azure_nested_obj_dict(prop, object, prop.properties, spaces)
          else
            name = azure_python_variable_name(prop, object.azure_sdk_definition.create)
            options = azure_prop_options(prop, object, spaces).join("\n")
            "#{name}=dict(\n#{indent_list(options, 4)}\n)"
          end
        end

        private

        # Creates a Python dictionary representing a nested object property
        # for validation.
        def azure_nested_obj_dict(prop, object, properties, spaces)
          name = python_variable_name(prop, object.azure_sdk_definition.create)
          options = azure_prop_options(prop, object, spaces).join("\n")
          [
            "#{name}=dict(\n#{indent_list(options, 4, true)}\n    options=dict(",
            indent_list(properties.map do |p|
              python_dict_for_property(p, object, spaces + 4)
            end, 8),
            "    )\n)"
          ]
        end

        # Returns an array of all base options for a given property.
        def azure_prop_options(prop, _object, spaces)
          [
            ('required=True' if prop.required && !prop.default_value && !is_location?(prop)),
            ("default=#{azure_python_literal(prop.default_value)}" \
            if prop.default_value),
            "type=#{quote_string(azure_python_type(prop))}",
            (azure_choices_enum(prop, spaces) if prop.is_a? Api::Type::Enum),
            ("elements=#{quote_string(azure_python_type(prop.item_type))}" \
              if prop.is_a? Api::Type::Array),
            ("aliases=[#{prop.aliases.map { |x| quote_string(x) }.join(', ')}]" \
              if prop.aliases),
            ('updatable=False' if prop.input && !is_resource_group?(prop) && !is_resource_name?(prop)),
            ("disposition='/'" if prop.input && !is_resource_group?(prop) && !is_resource_name?(prop))
          ].compact
        end

        # Returns a formatted string represented the choices of an enum
        def azure_choices_enum(prop, spaces)
          name = prop.out_name.underscore
          type = "type=#{quote_string(azure_python_type(prop))}"
          # + 6 for =dict(
          choices_indent = spaces + name.length + type.length + 6
          "choices=[#{prop.values.map do |x|
                        quote_string(x.to_s.underscore)
                      end.join(', ')}]"
        end

      end
    end
  end
end

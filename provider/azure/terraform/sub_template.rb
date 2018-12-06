module Provider
  module Azure
    module Terraform
      module SubTemplate

        def build_schema_property_get(input, output, property, object, indentation = 0)
          compile_template schema_property_get_template(property),
                           indentation: indentation,
                           input_var: input,
                           output_var: output,
                           prop_name: property.name.underscore,
                           property: property,
                           object: object
        end

        def build_schema_property_set(input, output, property, object, indentation = 0)
          compile_template schema_property_set_template(property),
                           indentation: indentation,
                           input_var: input,
                           output_var: output,
                           prop_name: property.name.underscore,
                           property: property,
                           object: object
        end

      end
    end
  end
end

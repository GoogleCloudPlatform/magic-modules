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

        def build_property_to_sdk_object(output, sdk_path, sdk_type_def, sdk_package, resource_name, properties, object, indentation = 4)
          compile_template property_to_sdk_object_template(sdk_type_def),
                           indentation: indentation,
                           output_statement: output,
                           sdk_package_name: sdk_package,
                           resource_name: resource_name,
                           sdk_obj_path: sdk_path,
                           sdk_type_def: sdk_type_def,
                           properties: properties,
                           object: object
        end

      end
    end
  end
end

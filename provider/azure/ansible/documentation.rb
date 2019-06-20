module Provider
  module Azure
    module Ansible
      module Documentation

        def azure_documentation_for_property(prop, dict, object)
          orig_name = prop.name.underscore
          new_name = python_variable_name(prop, object.azure_sdk_definition.create)
          dict[new_name] = dict.delete(orig_name)
          dict[new_name]['required'] = false if is_location?(prop)
          dict[new_name]['type'] = python_type(prop)
          dict[new_name]['choices'] = prop.values.map{|v| v.to_s.underscore} if prop.is_a? Api::Type::Enum
          dict[new_name]['description'] << azure_resource_ref_description(prop) if prop.is_a?(Api::Azure::Type::ResourceReference)
          dict
        end

        def azure_returns_for_property(prop, dict, object)
          orig_name = prop.name
          new_name = python_variable_name(prop, object.azure_sdk_definition.create)
          dict[new_name] = dict.delete(orig_name)

          dict[new_name]['type'] = 'str' if prop.is_a? Api::Azure::Type::ResourceReference
          dict[new_name]['returned'] = 'always'
          sample = prop.document_sample_value || prop.sample_value
          dict[new_name]['sample'] = sample unless sample.nil?

          dict
        end

        def azure_autogen_notic_contrib(lines)
          lines[1] = 'https://github.com/Azure/magic-module-specs'
          lines
        end

        private

        def azure_resource_ref_description(prop)
          [
            "It can be the #{prop.resource_type_name} name which is in the same resource group.",
            "It can be the #{prop.resource_type_name} ID. e.g., #{prop.document_sample_value || prop.sample_value}.",
            "It can be a dict which contains C(name) and C(resource_group) of the #{prop.resource_type_name}."
          ]
        end

      end
    end
  end
end

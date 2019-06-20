require 'azure/python_utils'

module Provider
  module Azure
    module Ansible
      module Module

        include Azure::PythonUtils

        def azure_python_dict_for_property(prop, dict, object)
          orig_name = prop.name.underscore
          new_name = python_variable_name(prop, object.azure_sdk_definition.create)
          dict[new_name] = dict.delete(orig_name)
          dict[new_name]['required'] = false if is_location?(prop)
          dict[new_name]['updatable'] = false if prop.input && !is_resource_group?(prop) && !is_resource_name?(prop)
          dict[new_name]['disposition'] = '/' if prop.input && !is_resource_group?(prop) && !is_resource_name?(prop)
          dict
        end

      end
    end
  end
end

# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

require 'api/object'

module Provider
  # Override to an Api::Resource in api.yaml
  class ResourceOverride < Api::Object
    include Api::Resource::Properties

    # Hash of properties where the key is the property path and the value is
    # a Provider::PropertyOverride.
    #
    # The property path can take of these formats:
    # - 'foo': Top-level property 'foo'
    # - 'foo.bar': Property 'bar' nested under property 'foo'
    attr_reader :properties

    # Apply this override to the given instance of Api::Resource
    def apply(api_resource)
      ensure_resource_properties
      update_overriden_properties(api_resource)
      update_overriden_azure_sdk_definition(api_resource)
      update_name_property_sort_order(api_resource)

      # TODO(nelsonjr): Enable revalidate the object to make sure we did not
      # break the object during the override process
      # | api_resource.validate # check if we did not break the object
    end

    def validate
      super

      @name_default_order = 750
      @azure_sdk_language = 'csharp'

      @properties ||= {}

      check_property :properties, Hash
    end

    private

    # Updates the resource property to a new value
    def update(resource, name, value)
      resource.instance_variable_set("@#{name}".to_sym, value)
    end

    # Attaches the overridden properties to Api::Resource and ensure they are
    # present on the class.
    def ensure_resource_properties
      Api::Resource.send(:include, overriden) # override ...
      require_module overriden
      our_override_modules.each { |mod| require_module mod } # ... and verify
    end

    # Copies all overridable properties from ResourceOverride into
    # Api::Resource.
    def update_overriden_properties(api_resource)
      our_override_modules.each do |mod|
        mod.instance_methods.each do |method|
          # If we have a variable for it, copy it.
          prop_name = "@#{method.id2name}".to_sym
          var_value = instance_variable_get(prop_name)
          api_resource.instance_variable_set(prop_name, var_value) \
            unless var_value.nil? || prop_name == :@azure_sdk_definition
        end
      end
    end

    # Returns all modules that contain overridable properties.
    def our_override_modules
      self.class.included_modules.select do |mod|
        mod == Api::Resource::Properties \
          || mod.name.split(':').last == 'OverrideProperties'
      end
    end

    # Ensures that Api::Resource includes a module.
    def require_module(clazz)
      raise "Api::Resource did not include required #{clazz} module" \
        unless Api::Resource.included_modules.include?(clazz)
      raise "#{self.class} did not include required #{clazz} module" \
        unless self.class.included_modules.include?(clazz)
    end

    def update_overriden_azure_sdk_definition(api_resource)
      unless api_resource.azure_sdk_definition.nil?
        api_resource.azure_sdk_definition.filter_language! @azure_sdk_language
        override = instance_variable_get('@azure_sdk_definition')
        api_resource.azure_sdk_definition.merge_overrides!(override) unless override.nil?
      end
    end

    def update_name_property_sort_order(api_resource)
      name_index = api_resource.properties.find_index{|p| p.name == 'name'}
      unless name_index.nil?
        name_prop = api_resource.properties[name_index]
        name_prop.instance_variable_set('@order', @name_default_order)
      end
    end

    # Returns the module that provides overriden properties for this provider.
    def overriden
      raise "overriden property should be implemented in #{self.class}"
    end

    def override_boolean(object, object_key, override_val)
      return if override_val.nil?

      object.instance_variable_set("@#{object_key}", override_val)
    end
  end
end

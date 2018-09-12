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
require 'api/type'

module Provider
  # Override a resource property (Api::Type) in api.yaml
  # TODO(rosbo): Shared common logic with ResourceOverride via a base class.
  class PropertyOverride < Api::Object

    attr_reader :new_type

    include Api::Type::Fields
    # To allow overrides for type-specific fields, include those type's
    # fields with an 'include' directive here.
    include Api::Type::NameValues::Fields
    include Api::Type::ResourceRef::Fields

    # Apply this override to property inheriting from Api::Type
    def apply(api_property)
      ensure_property_fields
      update_overriden_fields api_property

      # TODO(nelsonjr): Enable revalidate the object to make sure we did not
      # break the object during the override process
      # | api_resource.validate # check if we did not break the object
    end

    private

    # Updates a property field to a new value
    def update(property, field, value)
      property.instance_variable_set("@#{field}".to_sym, value)
    end

    # Attaches the overridden fields to the property and ensure they are
    # present on the class.
    def ensure_property_fields
      Api::Type.send(:include, overriden) # override ...
      require_module overriden
      our_override_modules.each { |mod| require_module mod } # ... and verify
    end

    # Copies all overridable properties from ResourceOverride into
    # Api::Resource.
    def update_overriden_fields(api_resource)
      our_override_modules.each do |mod|
        mod.instance_methods.each do |method|
          # If we have a variable for it, copy it.
          prop_name = "@#{method.id2name}".to_sym
          var_value = instance_variable_get(prop_name)
          api_resource.instance_variable_set(prop_name, var_value) \
            unless var_value.nil?
        end
      end
    end

    # Returns all modules that contain overridable properties.
    def our_override_modules
      self.class.included_modules.select do |mod|
        [Api::Type::Fields,
         Api::Type::NameValues::Fields,
         Api::Type::ResourceRef::Fields].include?(mod) \
          || mod.name.split(':').last == 'OverrideFields'
      end
    end

    # Ensures that Api::Type includes a module.
    def require_module(clazz)
      raise "#{self.class} did not include required #{clazz} module" \
        unless self.class.included_modules.include?(clazz)
    end

    # Returns the module that provides overriden properties for this provider.
    def overriden
      raise "overriden property should be implemented in #{self.class}"
    end
  end
end

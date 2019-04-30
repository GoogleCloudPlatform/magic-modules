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
require 'provider/property_override'
require 'provider/resource_override'

module Provider
  # A hash of Provider::ResourceOverride objects where the key is the api name
  # for that object.
  #
  # Example usage in a provider.yaml file where you want to extend a resource
  # description:
  #
  # overrides: !ruby/object:Provider::ResourceOverrides
  #   SomeResource: !ruby/object:Provider::MyProvider::ResourceOverride
  #     description: '{{description}} A tool-specific description complement'
  #     parameters:
  #       someParameter: !ruby/object:Provider::MyProvider::PropertyOverride
  #         description: 'foobar' # replaces description
  #     properties:
  #       someProperty: !ruby/object:Provider::MyProvider::PropertyOverride
  #         description: 'foobar' # replaces description
  #       anotherProperty.someNestedProperty:
  #         !ruby/object:Provider::MyProvider::PropertyOverride
  #         description: 'baz'
  #   ...
  class ResourceOverrides < Api::Object
    attr_accessor :__is_data_source

    def consume_config(api, config)
      @__api = api
      @__config = config
    end

    def validate
      return unless @__objects.nil? # allows idempotency of calling validate
      return if @__api.nil?
      @__is_data_source ||= false
      populate_nonoverridden_objects
      convert_findings_to_hash
      override_objects
      super
    end

    def [](index)
      @__objects[index]
    end

    def each
      return enum_for(:each) unless block_given?
      @__objects.each { |o| yield o }
      self
    end

    def select
      return enum_for(:select) unless block_given?
      @__objects.select { |o| yield o }
      self
    end

    def fetch(key, *args)
      # *args only holds default value. Needs to mimic ::Hash
      if args.empty?
        # KeyErorr will be thrown if key not found
        @__objects&.fetch(key)
      else
        # args[0] will be returned if key not found
        @__objects&.fetch(key, args[0])
      end
    end

    def key?(key)
      @__objects&.key?(key)
    end

    private

    # Converts every variable into @__objects
    def convert_findings_to_hash
      @__objects = {}
      instance_variables.each do |var|
        next if var.id2name.start_with?('@__')
        @__objects[var.id2name[1..-1]] = instance_variable_get(var)
        remove_instance_variable(var)
      end
    end

    # Applies the tool-specific overrides to the api objects
    def override_objects
      @__objects.each do |name, override|
        api_object = @__api.objects.find { |o| o.name == name }
        raise "The resource to override must exist #{name}" if api_object.nil?
        check_property_value 'overrides', override, Provider::ResourceOverride
        override.__is_data_source = @__is_data_source
        override.apply api_object
        populate_nonoverridden_properties api_object, override
        override_properties api_object, override
      end
    end

    def override_properties(api_object, override)
      # We apply property overrides in reverse order of level of nesting.
      # This helps us avoid a problem where we change the name of a property
      # before we have applied the overrides to its child properties (and
      # therefore can no longer find the child property, since the
      # parent name has changed)
      sorted_props = override.properties.sort_by { |path, p| p.override_order || -path.count('.') }
      sorted_props.each do |property_path, property_override|
        check_property_value "properties['#{property_path}']",
                             property_override, Provider::PropertyOverride
        api_property = find_property api_object, property_path.split('.')
        if api_property.nil?
          raise "The property to override '#{property_path}' must exist in " \
                "resource #{api_object.name}"
        end
        property_override.apply api_property
      end
    end

    def find_property(api_entity, property_path)
      property_name = property_path[0]
      properties = get_properties api_entity
      return nil if properties.nil?

      api_property = properties.find { |p| p.name == property_name }
      return nil if api_property.nil?

      property_path.shift
      if property_path.empty?
        api_property
      else
        find_property api_property, property_path
      end
    end

    def get_properties(api_entity)
      if api_entity.is_a?(Api::Resource)
        api_entity.all_properties
      elsif api_entity.is_a?(Api::Type::NestedObject)
        api_entity.all_properties
      elsif api_entity.is_a?(Api::Type::Array) &&
            api_entity.item_type.is_a?(Api::Type::NestedObject)
        api_entity.item_type.all_properties
      elsif api_entity.is_a?(Api::Type::Map)
        api_entity.value_type.all_properties
      end
    end

    def populate_nonoverridden_objects
      (@__api.objects || []).each do |object|
        var_name = "@#{object.name}".to_sym
        instance_variable_set(var_name, @__config.resource_override.new) \
          unless instance_variables.include?(var_name)
      end
    end

    def populate_nonoverridden_properties(api_entity, override)
      api_entity.all_user_properties.each do |prop|
        override.properties[prop.name] = @__config.property_override.new \
          unless override.properties.include?(prop.name)
        populate_nonoverriden_nested_properties prop.name, prop, override
      end
    end

    def populate_nonoverriden_nested_properties(prefix, property, override)
      nested_properties = get_properties(property)
      return if nested_properties.nil?

      nested_properties.each do |nested_prop|
        key = "#{prefix}.#{nested_prop.name}"
        override.properties[key] = @__config.property_override.new \
          unless override.properties.include?(key)
        populate_nonoverriden_nested_properties key, nested_prop, override
      end
    end
  end
end

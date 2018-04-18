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
  #     properties:
  #       someProperty: !ruby/object:Provider::MyProvider::PropertyOverride
  #         description: 'foobar' # replaces description
  #       anotherProperty.someNestedProperty:
  #         !ruby/object:Provider::MyProvider::PropertyOverride
  #         description: 'baz'
  #   ...
  class ResourceOverrides < Api::Object
    def consume_api(api)
      @__api = api
    end

    def validate
      return unless @__objects.nil? # allows idempotency of calling validate
      convert_findings_to_hash
      override_objects unless @__api.nil?
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
        override.apply api_object
        override_properties api_object, override
      end
    end

    def override_properties(api_object, override)
      override.properties.each do |property_path, property_override|
        check_property_value "properties['#{property_path}']",
                             property_override, Provider::PropertyOverride
        api_property = find_property api_object, property_path.split('.')
        if api_property.nil?
          raise "The property to override must exists #{property_path} " \
              "in resource #{api_object.name}"
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
      if api_entity.is_a?(Api::Resource) ||
         api_entity.is_a?(Api::Type::NestedObject)
        api_entity.properties
      elsif api_entity.is_a?(Api::Type::Array) &&
            api_entity.item_type.is_a?(Api::Type::NestedObject)
        api_entity.item_type.properties
      end
    end
  end
end

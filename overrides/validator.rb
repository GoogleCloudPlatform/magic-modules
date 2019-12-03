# Copyright 2018 Google Inc.
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

require 'api/product'
require 'overrides/resources'

module Overrides
  # Validates that the overrides will be valid as intended.
  # Throws errors if something invalid will be overridden.
  # Validation is split from Application to check for unused overrides,
  # provide better errors, split up code, and allow for validation in different places
  # than application.
  class Validator
    def initialize(api, overrides)
      @api = api
      @overrides = overrides
    end

    def run
      verify_resources(@api.objects) if @api.objects
    end

    private

    # Verify all resources in overrides exist in api
    def verify_resources(objects)
      @overrides.instance_variables.reject { |i| i == :@product }.each do |var|
        obj_array = objects.select { |o| o.name == var[1..-1] }
        raise "#{var[1..-1]} not found" if obj_array.empty?

        verify_resource(obj_array.first, @overrides[var])
      end
    end

    # Verify top-level fields exist on resource
    def verify_resource(res, overrides)
      overrides.instance_variables.reject { |i| %i[@properties @parameters].include?(i) }
               .each do |field_name|
        # Check override object.
        field_symbol = field_name[1..-1].to_sym
        next if check_if_exists(res, field_symbol, overrides.class.attributes)

        raise "#{field_name} does not exist on #{res.name}"
      end
      # Use instance_variable_get to get excluded properties
      verify_properties((res.instance_variable_get('@properties') || []) +
                        (res.instance_variable_get('@parameters') || []), overrides['properties'],
                        res.name)
    end

    # Verify a list of properties (parameters or properties on an API::Resource)
    def verify_properties(properties, overrides, res_name = '')
      overrides ||= {}
      overrides.each do |k, v|
        path = property_path(k)
        verify_property(find_property(properties, path, res_name), v)
      end
    end

    # Returns a property (or throws an error if it does not exist)
    def find_property(properties, path, res_name = '')
      # Add context for debugging property name typos
      available = []
      prefix = ''

      prop = nil
      path.each do |part|
        available = properties.map { |p| prefix + p.name } unless properties.empty?
        prefix = prefix + part + '.'

        prop = properties.select { |o| o.name == part }.first
        properties = prop&.nested_properties || []
      end
      raise "#{path.join('.')} does not exist on #{res_name}.\n Possible values #{available}" \
        unless prop

      prop
    end

    def verify_property(property, overrides)
      overrides.instance_variables
               .reject { |i| %i[@properties @item_type @type].include?(i) }
               .each do |field_name|
        # Check override object.
        field_symbol = field_name[1..-1].to_sym
        next if check_if_exists(property, field_symbol, overrides.class.attributes,
                                overrides['@type'])

        raise "#{field_name} does not exist on #{property.name}"
      end
    end

    # Check if this field exists on this object.
    # The best way (sadly) to do this is to see if a getter exists.
    def check_if_exists(obj, field, attributes = [], override_type = nil)
      # Not all types share the same values.
      # If we're changing types, the new type's getters matter, not the old type.
      return true if attributes.include?(field)

      if override_type
        Module.const_get(override_type).new.respond_to? field
      else
        obj.respond_to? field
      end
    end

    # This keeps the [] brackets in place.
    def property_path(prop_name)
      prop_name.split('.')
    end
  end
end

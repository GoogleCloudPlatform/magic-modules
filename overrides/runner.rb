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
require 'overrides/validator'

module Overrides
  # This runner takes an Api::Product and applies a set of Overrides::ResourceOverrides
  # It does this by building a brand new Api::Product object from scratch, using
  # the values from either the original Api::Product or the override values.
  # Example usage in a provider.yaml file where you want to extend a resource
  # description:
  #
  # overrides: !ruby/object:Overrides::ResourceOverrides
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
  #       anotherProperty[].someNestedPropertyInAnArray:
  #         !ruby/object:Provider::MyProvider::PropertyOverride
  #         description: 'baz'
  #   ...
  class Runner
    # Internal information for applying overrides.
    class OverrideInfo
      def initialize(all_overrides, property_class, resource_class)
        @all_overrides = all_overrides
        @property_class = property_class
        @resource_class = resource_class
      end

      def product_overrides
        return @all_overrides['product'] || {}
      end

      def resource_overrides(name)
        unless @all_overrides[name].nil? || @all_overrides[name].empty?
          @all_overrides[name]
        else
          @resource_class.new
        end
      end

      def property_overrides(res_name, property_name)
        res_override = resource_overrides(res_name)
        return @property_class.new unless res_override['properties']
        return res_override['properties'][property_name] || @property_class.new
      end
    end

    class SinglePropertyOverride
      def initialize(property_override, property_override_class)
        @property_override = property_override
        @property_override_class = property_override_class
      end

      def property_overrides(_res_name, _prop_name)
        return @property_override if @property_override
        @property_override_class.new
      end
    end

    class << self
      # Takes in the old Api::Product object, a set of overrides, and resource/property override classes
      # and returns a new Api::Product object with all overrides applied.
      def build(api, overrides, res_override_class = Overrides::ResourceOverride,
                prop_override_class = Overrides::PropertyOverride)
        override_info = OverrideInfo.new(overrides, prop_override_class, res_override_class)
        validator = Overrides::Validator.new(api, overrides)
        validator.run
        build_product(api, override_info)
      end

      # Takes in a single property, single override, and a property override class and returns
      # a brand new property with overrides applied.
      # This is used exclusively for Ansible filters, which are regular properties that live
      # outside of a Api::Product
      def build_single_property(api_property, property_override, prop_override_class)
        overrides = SinglePropertyOverride.new(property_override, prop_override_class)
        build_property(api_property, '', overrides, '')
      end

      private

      # Given a old Api::Product, and Overrides::ResourceOverrides,
      # returns a new Api::Product with overrides applied
      def build_product(old_prod, overrides)
        prod = Api::Product.new
        old_prod.instance_variables
                .reject { |o| o == :@objects }.each do |var_name|
          if overrides.product_overrides[var_name]
            prod.instance_variable_set(var_name, overrides.product_overrides[var_name])
          else
            prod.instance_variable_set(var_name, old_prod.instance_variable_get(var_name))
          end
        end
        prod.instance_variable_set('@objects',
                                   old_prod.objects
                                           .map do |o|
                                     build_resource(o, overrides)
                                   end)
        prod
      end

      # Given a Api::Resource and Provider::Override::ResourceOverride,
      # return a new Api::Resource with overrides applied.
      def build_resource(old_resource, overrides)
        # Set up the ResourceOverride object.
        res_override = overrides.resource_overrides(old_resource.name)
        res_override.validate
        res_override.apply old_resource

        # Create new Api::Resource + apply all provider-specific values.
        res = Api::Resource.new
        set_additional_values(res, res_override)

        # Loop through all values on the Resource (besides properties + parameters)
        # and replace with overriden values, if they exist.
        variables = (old_resource.instance_variables + res_override.instance_variables).uniq
        variables.reject { |o| %i[@properties @parameters].include?(o) }
                 .each do |var_name|
          if !res_override[var_name].nil?
            res.instance_variable_set(var_name, res_override[var_name])
          else
            res.instance_variable_set(var_name, old_resource.instance_variable_get(var_name))
          end
        end

        # Loop through properties + parameters and build those too.
        # Using instance_variable_get('properties') to make sure we get `exclude: true` properties
        ['@properties', '@parameters'].each do |val|
          new_props = ((old_resource.instance_variable_get(val) || [])).map do |p|
            build_property(p, old_resource.name, overrides)
          end
          res.instance_variable_set(val, new_props)
        end
        res
      end

      # Given a Api::Type property and a hash of properties, create a new Api::Type property
      # This will handle NestedObjects, Arrays of NestedObjects of arbitrary length
      def build_property(old_property, resource_name, overrides, prefix = '')
        # Build a new property, minus any nested properties.
        new_prop = build_primitive_property(old_property, overrides, resource_name,
                                            "#{prefix}#{old_property.name}")

        # Build all nested properties in a recursive manner.
        if old_property.nested_properties?
          new_props = old_property.nested_properties.map do |p|
            build_property(p, resource_name, overrides, "#{prefix}#{old_property.name}.")
          end

          if old_property.is_a?(Api::Type::NestedObject)
            new_prop.instance_variable_set('@properties', new_props)
          elsif old_property.is_a?(Api::Type::Map) && \
                old_property.value_type.is_a?(Api::Type::NestedObject)
            new_prop.instance_variable_set('@value_type', Api::Type::NestedObject.new)
            new_prop.value_type.instance_variable_set('@properties', new_props)
          elsif old_property.is_a?(Api::Type::Array) && \
                old_property.item_type.is_a?(Api::Type::NestedObject)
            new_prop.instance_variable_set('@item_type', Api::Type::NestedObject.new)
            new_prop.item_type.instance_variable_set('@properties', new_props)
          end
        end
        new_prop
      end

      # Given a primitive Api::Type (string, integers, times, etc) and override,
      # return a new Api::Type with overrides applied.
      # This will be called by build_property, which handles nesting.
      def build_primitive_property(old_property, overrides, resource_name, property_name)
        # Get the property override setup.
        prop_override = overrides.property_overrides(resource_name, property_name)

        prop_override.validate
        prop_override.apply old_property

        # If different type, create the new property from that type.
        prop = if prop_override['type']
                 Module.const_get(prop_override['type']).new
               else
                 old_property.class.new
               end

        set_additional_values(prop, prop_override)
        variables = (old_property.instance_variables + prop_override.instance_variables).uniq

        # Set api_name with old property so that a potential new name doesn't override it.
        prop.instance_variable_set('@api_name', old_property.api_name || old_property.name)

        # Loop through all values and set them according to overrides.
        variables.reject { |o| o == :@properties }
                 .each do |var_name|
          if !prop_override[var_name].nil?
            prop.instance_variable_set(var_name, prop_override[var_name])
          else
            prop.instance_variable_set(var_name, old_property.instance_variable_get(var_name))
          end
        end
        prop
      end

      # Overrides have additional values inside the override that do not regularly belong
      # on the Api::* object. These values need to be set + they need getters so they
      # can be accessed propertly in the templates.
      def set_additional_values(object, override)
        override.class.attributes.each do |o|
          object.instance_variable_set("@#{o}", override[o])
          object.define_singleton_method(o.to_sym) { instance_variable_get("@#{o}") }
        end
      end
    end
  end
end

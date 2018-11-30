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
require 'provider/overrides/resources'
require 'provider/overrides/validator'

module Provider
  module Overrides
    # This runner takes an Api::Product and applies a set of Provider::Overrides::ResourceOverrides
    # It does this by building a brand new Api::Product object from scratch, using
    # the values from either the original Api::Product or the override values.
    # Example usage in a provider.yaml file where you want to extend a resource
    # description:
    #
    # overrides: !ruby/object:Provider::Overrides::ResourceOverrides
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
      def initialize(api, overrides,
                     res_override_class = Provider::Overrides::ResourceOverride,
                     prop_override_class = Provider::Overrides::PropertyOverride)
        @api = api
        @overrides = overrides
        @res_override_class = res_override_class
        @prop_override_class = prop_override_class
      end

      def build
        validator = Provider::Overrides::Validator.new(@api, @overrides)
        validator.run
        build_product(@api, @overrides)
      end

      private

      # Given a old Api::Product, and Provider::Overrides::ResourceOverrides,
      # returns a new Api::Product with overrides applied
      def build_product(old_prod, all_overrides)
        prod = Api::Product.new
        old_prod.instance_variables
                .reject { |o| o == :@objects }.each do |var_name|
          if (all_overrides['product'] || {})[var_name]
            prod.instance_variable_set(var_name, all_overrides['product'][var_name])
          else
            prod.instance_variable_set(var_name, old_prod.instance_variable_get(var_name))
          end
        end
        prod.instance_variable_set('@objects',
                                   old_prod.objects
                                           .map { |o| build_resource(o, all_overrides[o.name]) })
        prod
      end

      # Given a Api::Resource and Provider::Override::ResourceOverride,
      # return a new Api::Resource with overrides applied.
      def build_resource(old_resource, res_override)
        res_override = @res_override_class.new if res_override.nil? || res_override.empty?
        res_override.validate
        res_override.apply old_resource

        res = Api::Resource.new
        set_additional_values(res, res_override)

        variables = (old_resource.instance_variables + res_override.instance_variables).uniq
        variables.reject { |o| o == :@properties || o == :@parameters }
                 .each do |var_name|
          if !res_override[var_name].nil?
            res.instance_variable_set(var_name, res_override[var_name])
          else
            res.instance_variable_set(var_name, old_resource.instance_variable_get(var_name))
          end
        end

        # Using instance_variable_get('properties') to make sure we get `exclude: true` properties
        ['@properties', '@parameters'].each do |val|
          new_props = (old_resource.instance_variable_get(val) || []).map do |p|
            build_property(p, res_override[val])
          end
          res.instance_variable_set(val, new_props)
        end
        res
      end
      # rubocop:enable Metrics/AbcSize

      # Given a Api::Type property and a hash of properties, create a new Api::Type property
      # This will handle NestedObjects, Arrays of NestedObjects of arbitrary length
      def build_property(old_property, property_overrides, prefix = '')
        property_overrides = {} if property_overrides.nil?
        new_prop = build_primitive_property(old_property,
                                            property_overrides["#{prefix}#{old_property.name}"])
        if old_property.is_a?(Api::Type::NestedObject)
          new_props = old_property.properties.map do |p|
            build_property(p, property_overrides, "#{prefix}#{old_property.name}.")
          end
          new_prop.instance_variable_set('@properties', new_props)
        elsif old_property.is_a?(Api::Type::Map) && \
              old_property.value_type.is_a?(Api::Type::NestedObject)
          new_prop.instance_variable_set('@value_type', Api::Type::NestedObject.new)
          new_props = old_property.value_type.properties.map do |p|
            build_property(p, property_overrides, "#{prefix}#{old_property.name}.")
          end
          new_prop.value_type.instance_variable_set('@properties', new_props)
        elsif old_property.is_a?(Api::Type::Array) && \
              old_property.item_type.is_a?(Api::Type::NestedObject)
          new_prop.instance_variable_set('@item_type', Api::Type::NestedObject.new)
          new_props = old_property.item_type.properties.map do |p|
            build_property(p, property_overrides, "#{prefix}#{old_property.name}[].")
          end
          new_prop.item_type.instance_variable_set('@properties', new_props)
        end
        new_prop
      end

      # Given a primitive Api::Type (string, integers, times, etc) and override,
      # return a new Api::Type with overrides applied.
      # This will be called by build_property, which handles nesting.
      def build_primitive_property(old_property, prop_override)
        prop_override = @prop_override_class.new if prop_override.nil? || prop_override.empty?
        prop_override.validate
        prop_override.apply old_property

        prop = if prop_override['type']
                 Module.const_get(prop_override['type']).new
               else
                 old_property.class.new
               end

        set_additional_values(prop, prop_override)
        variables = (old_property.instance_variables + prop_override.instance_variables).uniq

        # Set api_name with old property so that a potential new name doesn't override it.
        prop.instance_variable_set('@api_name', old_property.name)

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

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

require 'api/product'

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
  class OverrideRunner < Api::Object
    def initialize(api, overrides)
      @api = api
      @overrides = overrides
    end

    def build
      build_product(@api, @overrides)
    end

    private

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
                                 old_prod.objects.map { |o| build_resource(o, all_overrides[o.name]) })
      prod
    end

    def build_resource(old_resource, res_override)
      res_override = {} if res_override.nil?
      res = Api::Resource.new
      old_resource.instance_variables.reject { |o| o == :properties || o == :parameters }
                  .each do |var_name|
        if res_override[var_name]
          res.instance_variable_set(var_name, res_override[var_name])
        else
          res.instance_variable_set(var_name, old_resource.instance_variable_get(var_name))
        end
      end
      res.instance_variable_set('@properties', old_resource.properties.map { |p| build_property(p, res_override.dig('properties', p.name)) })
      res.instance_variable_set('@parameters', old_resource.parameters.map { |p| build_property(p, res_override.dig('properties', p.name)) })
      res
    end

    def build_property(old_property, prop_override)
      prop_override = {} if prop_override.nil?
      prop = if prop_override['type']
               Module.const_get(prop_override['type']).new
             else
               old_property.class.new
             end

      old_property.instance_variables.reject { |o| o == :properties }
                                     .each do |var_name|
        if prop_override[var_name]
          prop.instance_variable_set(var_name, prop_override[var_name])
        else
          prop.instance_variable_set(var_name, old_property.instance_variable_get(var_name))
        end
      end
      prop
    end
  end
end

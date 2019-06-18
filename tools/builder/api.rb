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

require 'net/http'
require 'json'
require 'active_support/inflector'

require 'api/product'
require 'api/resource'
require 'api/type'
require 'api/compiler'
require 'api/async'

# All properties are read-only in Magic Modules
# Creating an api.yaml involves a lot of setting values.
# This will create setters for all fields on an api.yaml
# (but only in the context of the linter)
module Api
  class Object
    # Create a setter if the setter doesn't exist
    # Yes, this isn't pretty and I apologize
    def method_missing(method_name, *args)
      matches = /([a-z_]*)=/.match(method_name)
      super unless matches
      create_setter(matches[1])
      method(method_name.to_sym).call(*args)
    end

    def create_setter(variable)
      self.class.define_method("#{variable}=") { |val| instance_variable_set("@#{variable}", val) }
    end
  end
end

TYPES = {
  'string': 'String',
  'boolean': 'Boolean',
  'object': 'Map',
  'integer': 'Integer',
  'number': 'Double',
  'array': 'Array'
}

class HumanApi
  def initialize(discovery, handwritten)
    @discovery = discovery
    @written = written
  end

  def build
    # For each product, inject extra properties.
    @handwritten.products.each { |prod| add_missing_properties(prod.properties) }

    # Inject extra products at end
    missing_products = @discovery.products.select { |x| @handwritten.products.map(&:name).includes?(x.name) }
    @written.products.append(missing_products)
    return @written
  end

  def add_missing_properties(disc_props, hand_props)
    # If nested property, recurse.
    hand_props.select { |p| p.nested_properties? }
              .each do |p|
      matching_discovery_prop = disc_props.select { |d| d.name == p.name }.first
      if p.is_a?(Api::Type::NestedObject)
        add_missing_properties(matching_discovery_prop.properties, p.properties)
      elsif p.is_a?(Api::Type::Array) && p.item_type.is_a?(Api::Type::NestedObject)
        add_missing_properties(matching_discovery_prop.item_type.properties, p.item_type.properties)
      end
    end
    # Inject new properties.
    missing_properties = disc_props.select { |p| !hand_props.any? { |d| d.name == p.name }}
    hand_props.append(missing_properties)
  end
end

class DiscoveryProperty
  attr_reader :schema
  attr_reader :name

  attr_reader :__product

  def initialize(name, schema, product)
    @name = name
    @schema = schema

    @__product = product
  end

  def get_property
    prop = Module.const_get("Api::Type::#{type}").new
    prop.name = @name
    prop.description = @schema.dig('description')
    prop.output = output?
    prop.values = enum if @schema.dig('enum')
    prop.properties = nested if prop.is_a?(Api::Type::NestedObject)
    prop.item_type = array if prop.is_a?(Api::Type::Array)
    prop
  end

  private

  def type
    return "NestedObject" if @schema.dig('$ref')
    return "NestedObject" if @schema.dig('type') == 'object' && @schema.dig('properties')
    return "Enum" if @schema.dig('enum')
    TYPES[@schema.dig('type').to_sym]
  end

  def output?
    (@schema.dig('description') || '').downcase.include?('output only')
  end

  def enum
    @schema.dig('enum').map { |val| val.to_sym }
  end

  def nested
    if @schema.dig('$ref')
      @__product.get_resource(@schema.dig('$ref')).properties
    else
      DiscoveryResource.new(@schema, @__product).properties
    end
  end

  def array
    schema_type = @schema.dig('items', 'type')
    if (!schema_type && @schema.dig('items', '$ref')) || @schema.dig('items', 'properties')
      prop = Api::Type::NestedObject.new
      if @schema.dig('items', '$ref')
        prop.properties = @__product.get_resource(@schema.dig('items', '$ref')).properties
      else
        prop.properties = DiscoveryResource.new(@schema.dig('items'), @__product).properties
      end
      return prop
    end
    return "Api::Type::#{TYPES[schema_type.to_sym]}" if schema_type != 'object'
  end
end

# Holds information about discovery objects
# Two sections: schema (properties) and methods
class DiscoveryResource
  attr_reader :schema

  attr_reader :__product

  def initialize(schema, product)
    @schema = schema
    @__product = product
    @methods = @__product.get_methods_for_resource(@schema.dig('id'))
  end

  def exists?
    !@schema.nil?
  end

  def resource
    res = Api::Resource.new
    res.name = @schema.dig('id')
    res.kind = @schema.dig('properties', 'kind', 'default')
    res.base_url = base_url_format(@methods['list']['path'])
    res.description = @schema.dig('description')
    res.properties = properties
    res
  end

  def properties
    @schema.dig('properties')
           .reject { |k, _| k == 'kind' }
           .map { |k, v| DiscoveryProperty.new(k, v, @__product).get_property }
  end

  private

  def base_url_format(url)
    "projects/#{url.gsub('{', '{{').gsub('}', '}}')}"
  end
end

# Responsible for grabbing Discovery Docs and getting resources from it
class DiscoveryProduct
  attr_reader :results
  attr_reader :doc

  def initialize(url, object)
    @results = send_request(url)
    @object = object
  end

  def get_resources
    @results['schemas'].map do |name, _|
      next if name != @object
      get_resource(name).resource
    end.compact
  end

  def get_resource(resource)
    DiscoveryResource.new(@results['schemas'][resource], self)
  end

  def get_methods_for_resource(resource)
    return if resource.nil?
    @results.dig 'resources', resource.pluralize.camelize(:lower), 'methods'
  end

  def get_product
    product = Api::Product.new
    product.versions = [version]
    product.objects = get_resources
    product
  end

  private

  def send_request(url)
    JSON.parse(Net::HTTP.get(URI(url)))
  end

  def version
    version = Api::Product::Version.new
    version.name = 'ga'
    version.base_url = base_url_format(@results['baseUrl'])
    version.default = true
    version
  end

  def base_url_format(url)
    url.gsub('projects/', '').gsub('{', '{{').gsub('}', '}}')
  end
end

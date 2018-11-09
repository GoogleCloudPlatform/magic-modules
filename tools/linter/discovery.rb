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

TYPES = {
  'string': 'String',
  'boolean': 'Boolean',
  'object': 'NameValues',
  'integer': 'Integer',
  'number': 'Double',
  'array': 'Array'
}

class DiscoveryProperty
  attr_reader :schema
  attr_reader :name

  attr_reader :__product

  def initialize(name, schema, product)
    @name = name
    @schema = schema

    @__product = product
  end

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

  def has_nested_properties?
    if @schema.dig('$ref')
      return true
    elsif @schema.dig('type') == 'object' && @schema.dig('properties')
      return true
    elsif @schema.dig('type') == 'array' && @schema.dig('items', '$ref')
      return true
    elsif @schema.dig('type') == 'array' && @schema.dig('items', 'properties')
      return true
    else
      return false
    end
  end
end

# Holds information about discovery objects
# Two sections: schema (properties) and methods
class DiscoveryResource
  attr_reader :schema
  attr_reader :name

  attr_reader :__product

  def initialize(schema, name, product)
    @schema = schema
    @name = name
    @__product = product

  end

  def exists?
    !@schema.nil?
  end

  def properties
    @schema.dig('properties')
           .reject { |k, _| k == 'kind' }
           .map { |k, v| DiscoveryProperty.new(k, v, @__product) }
  end
end

# Responsible for grabbing Discovery Docs and getting resources from it
class DiscoveryBuilder
  attr_reader :results

  def initialize(url, objects)
    @objects = objects
    @results = send_request(url)
  end

  def resources
    @results['schemas'].map do |name, _|
      next unless @objects.include?(name)
      get_resource(name)
    end.compact
  end

  def get_methods_for_resource(resource)
    return unless @results['resources'][resource.pluralize.camelize(:lower)]
    @results['resources'][resource.pluralize.camelize(:lower)]['methods']
  end

  private

  def send_request(url)
    JSON.parse(Net::HTTP.get(URI(url)))
  end

  def get_resource(resource)
    DiscoveryResource.new(@results['schemas'][resource], resource, self)
  end
end


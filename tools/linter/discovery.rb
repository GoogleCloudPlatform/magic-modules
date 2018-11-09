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

# Holds all Discovery Information about a Api::Type (property)
class DiscoveryProperty
  attr_reader :schema
  attr_reader :name

  attr_reader :__product

  def initialize(name, schema, product)
    @name = name
    @schema = schema

    @__product = product
  end

  def has_nested_properties?
    return !nested_properties.empty?
  end

  def nested_properties
    if @schema.dig('$ref')
      return @__product.get_resource(@schema.dig('$ref')).properties
    elsif @schema.dig('type') == 'object' && @schema.dig('properties')
      return DiscoveryResource.new(@schema, nil, @__product).properties
    elsif @schema.dig('type') == 'array' && @schema.dig('items', '$ref')
      return @__product.get_resource(@schema.dig('items', '$ref')).properties
    elsif @schema.dig('type') == 'array' && @schema.dig('items', 'properties')
      return DiscoveryResource.new(@schema.dig('items'), nil, @__product).properties
    else
      return []
    end
  end
end

# Holds Discovery information about a Resource
# This Resource is usually a Api::Resource,
# although it may be the contents of a NestedObject
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

  def get_resource(resource)
    DiscoveryResource.new(@results['schemas'][resource], resource, self)
  end

  private

  def send_request(url)
    JSON.parse(Net::HTTP.get(URI(url)))
  end
end


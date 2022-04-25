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

# Holds all information about parsing parts of the Discovery docs.
module Discovery
  # Holds all Discovery Information about a Api::Type (property)
  class Property
    attr_reader :schema
    attr_reader :name

    attr_reader :__product

    def initialize(name, schema, product)
      @name = name
      @schema = schema

      @__product = product
    end

    def nested_properties?
      !nested_properties.empty?
    end

    def nested_properties
      if @schema.dig('$ref')
        @__product.get_resource(@schema.dig('$ref')).properties
      elsif @schema.dig('type') == 'object' && @schema.dig('properties')
        Resource.new(@schema, nil, @__product).properties
      elsif @schema.dig('type') == 'array' && @schema.dig('items', '$ref')
        @__product.get_resource(@schema.dig('items', '$ref')).properties
      elsif @schema.dig('type') == 'array' && @schema.dig('items', 'properties')
        Resource.new(@schema.dig('items'), nil, @__product).properties
      else
        []
      end
    end
  end

  # Holds Discovery information about a Resource
  # This Resource is usually a Api::Resource,
  # although it may be the contents of a NestedObject
  class Resource
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
             .map { |k, v| Property.new(k, v, @__product) }
    end
  end

  # Responsible for grabbing Discovery Docs and getting resources from it
  class Builder
    attr_reader :results

    def initialize(doc, objects)
      @product = doc['product']
      @filename = doc['filename']
      @url = doc['url']
      @aliases = doc['aliases'] || {}
      @resources_in_api_yaml = objects
      @results = fetch_discovery_doc(@url)
    end

    def resources
      list_of_resources.map { |name, _| get_resource(name) }
                       .compact
    end

    def get_methods_for_resource(resource)
      return unless @results['resources'][resource.pluralize.camelize(:lower)]

      @results['resources'][resource.pluralize.camelize(:lower)]['methods']
    end

    def get_resource(resource)
      original_resource = resource

      resource = @aliases[resource] if @aliases[resource]
      # Region, Global should resolve to normal.
      if @results['schemas'][resource]
        Resource.new(@results['schemas'][resource], original_resource, self)
      elsif @results['schemas'][resource.sub('Region', '')]
        resource = resource.sub('Region', '')
        Resource.new(@results['schemas'][resource], resource, self)
      elsif @results['schemas'][resource.sub('Global', '')]
        resource = resource.sub('Global', '')
        Resource.new(@results['schemas'][resource], resource, self)
      else
        puts "#{original_resource} from #{@filename} not found in discovery docs - #{@url}"
      end
    end

    private

    def fetch_discovery_doc(url)
      JSON.parse(Net::HTTP.get(URI(url)))
    end

    def list_of_resources
      resources = list_of_resource_keys.map do |k|
        if @results['schemas'][k]
          k
        else
          k.singularize
        end
      end
      resources.map { |k| k.titleize.delete(' ') }
    end

    def list_of_resource_keys
      # There are 2 main ways to format discovery docs, check for the
      # newer format first, then fall back to the old.
      resources = @results.dig('resources', 'projects', 'resources') || @results['resources']

      resources.keys.reject { |o| o.downcase.include?('operation') }
    end
  end
end

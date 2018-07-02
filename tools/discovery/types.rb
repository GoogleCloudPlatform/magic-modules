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

$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '..', '..')
require 'net/http'
require 'json'
require 'erb'
require 'compile/core'

# rubocop:disable Style/MixinUsage
include Compile::Core
# rubocop:enable Style/MixinUsage

# Represents a property on a resource
class Property
  attr_reader :name
  def initialize(name, attributes, json)
    @name = name
    @attributes = attributes
    @json = json
  end

  def type
    return 'NestedObject' if @attributes['$ref']
    @attributes['type'].capitalize
  rescue StandardError
    'NONE'
  end

  def description
    @attributes['description']
  end

  def required
    return 'true' if description.include? '[Required]'
    'false'
  rescue StandardError
    'none'
  end

  def output
    return 'true' if description.match?(/[Oo]utput.[Oo]nly/)
    'false'
  rescue StandardError
    'none'
  end

  def properties
    # Get ref
    @json['schemas'][@attributes['$ref']]['properties'].map do |arr|
      Property.new(arr[0], arr[1], @json)
    end
  end
end

# Resprents a GCP resource
class Resource
  attr_accessor :schema
  attr_reader :properties

  def initialize(name, resource, schema, json)
    @name = name
    @resource = resource
    @schema = schema
    @json = json

    build_properties
  end

  attr_reader :name

  def base_url
    @resource['methods']['list']['path'].gsub('{', '{{').gsub('}', '}}')
  rescue StandardError
    @resource['methods']['get']['path'].gsub('{', '{{').gsub('}', '}}')
  end

  def virtual
    return 'true' if @resource['methods']['insert'].nil?
    'false'
  end

  private

  def build_properties
    @properties = @schema['properties'].map do |arr|
      Property.new(arr[0], arr[1], @json)
    end
  end
end

# Represents a GCP product (a collection of products)
class Product
  attr_reader :resources

  def initialize(results)
    @json = results

    build_resources
  end

  def title
    @json['title']
  end

  def base_url
    @json['baseUrl']
  end

  def version
    @json['version']
  end

  def scope
    @json['auth']['oauth2']['scopes'].keys[0]
  end

  private

  # resources contains a list of objects and their methods
  # schemas contains a list of objects and their properties. This includes
  # NestedObjects
  # We need to build Resource objects from all objects that exist in resources
  # and exist in schemas (otherwise, we'll have a bunch of NestedObjects...)
  # rubocop:disable Metrics/AbcSize
  def build_resources
    @resources = @json['resources'].map do |key, value|
      method = if value['methods']['get']
                 'get'
               elsif value['methods']['list']
                 'list'
               else
                 raise "#{value} does not contain get or list method"
               end
      next unless @json['schemas'][value['methods'][method]['response']['$ref']]
      Resource.new(
        key,
        value,
        @json['schemas'][value['methods'][method]['response']['$ref']],
        @json
      )
    end
  end
  # rubocop:enable Metrics/AbcSize
end

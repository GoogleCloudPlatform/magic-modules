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

require 'api/compiler'
require 'provider/config'
require 'provider/ansible'
require 'provider/terraform/config'
require 'overrides/terraform/resource_override'
require 'overrides/terraform/property_override'

# Takes in a DiscoveryResource + Api::Resource
# Loops through all properties of the DiscoveryResource (at any depth)
# Passes the DiscoveryProperties and their corresponding Api Properties
# to a test block.
class PropertyFetcher
  class << self
    # Runs a block with a matching set of a DiscoveryProperty and Api::Type property
    def fetch_property_pairs(discovery_properties, api_properties, &block)
      run_on_properties(discovery_properties, api_properties, '', &block)
    end

    private

    def run_on_properties(discovery_properties, api_properties, prefix, &block)
      discovery_properties.each do |disc_prop|
        api_props = api_properties&.select { |p| p.name == disc_prop.name }
        raise 'Multiple properties with name' if api_props && api_props.length > 1

        api_prop = api_props&.first
        yield(disc_prop, api_prop, "#{prefix}#{disc_prop.name}")
        if disc_prop.nested_properties?
          run_on_properties(disc_prop.nested_properties, nested_properties_for_api(api_prop),
                            "#{prefix}#{disc_prop.name}.", &block)
        end
      end
    end

    def nested_properties_for_api(api)
      if api&.is_a?(Api::Type::NestedObject)
        api.properties
      elsif api&.is_a?(Api::Type::Array) && api&.item_type&.is_a?(Api::Type::NestedObject)
        api.item_type.properties
      else
        []
      end
    end
  end
end

# Gets a Api::Product from a api.yaml filename
class ApiFetcher
  # Get api from filename
  def self.api_from_file(filename)
    return FakeApi.new(filename) unless File.file?(filename)

    Api::Compiler.new(filename).run
  end

  # Get api from filename and apply overrides from a provider.
  def self.provider_from_file(api_filename, provider_name)
    api = api_from_file(api_filename)
    provider_filename = "#{api_filename.split('/')[0..-2].join('/')}/#{provider_name}.yaml"
    return nil unless File.file?(provider_filename)

    Provider::Config.parse(provider_filename, api, 'ga')
    api
  end
end

# A Fake version of an API object designed to mimic api.yaml when
# none exists
class FakeApi
  attr_reader :objects
  attr_reader :api_name

  def initialize(filename)
    @objects = []
    @api_name = filename.split('/')[1].delete(' ').downcase
  end
end

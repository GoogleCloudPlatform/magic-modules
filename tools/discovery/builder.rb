$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '..', '..')
require 'net/http'
require 'json'
require 'erb'
require 'compile/core'

include Compile::Core

class Resource
  attr_accessor :schema

  def initialize(name, resource, schema)
    @name = name
    @resource = resource
    @schema = schema
  end

  def name
    @name
  end

  def base_url
      @resource['methods']['list']['path']
    rescue
      @resource['methods']['get']['path']
  end
end

class Resources
  # @resources contains the methods.
  # @schemas contains the object model.
  def initialize(resources, schemas)
    @resources = resources
    @schemas = schemas
  end

  # Get all resources that have methods (are GCP resources)
  def resources
    @resources.map do |key, value|
      if @schemas[value['methods']['get']['response']['$ref']]
        Resource.new(
          key,
          value,
          @schemas[value['methods']['get']['response']['$ref']]
        )
      end
    end
  end
end

DISCOVERY_URL = 'https://www.googleapis.com/discovery/v1/apis/compute/v1/rest'
uri = URI(DISCOVERY_URL)
response = Net::HTTP.get(uri)
results = JSON.parse(response)

res = Resources.new(results['resources'], results['schemas']).resources
File.write('output.yaml', lines(compile_file({ results: results,
                                               resources: res }, 'api.yaml.erb')))

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

  # Get all resources that have methods (are GCP resources)
  def build_resources
    @resources = @json['resources'].map do |key, value|
      if @json['schemas'][value['methods']['get']['response']['$ref']]
        Resource.new(
          key,
          value,
          @json['schemas'][value['methods']['get']['response']['$ref']]
        )
      end
    end
  end
end

DISCOVERY_URL = 'https://www.googleapis.com/discovery/v1/apis/compute/v1/rest'
uri = URI(DISCOVERY_URL)
response = Net::HTTP.get(uri)
results = JSON.parse(response)

res = Product.new(results)
File.write('output.yaml', lines(compile_file({product: res}, 'api.yaml.erb')))

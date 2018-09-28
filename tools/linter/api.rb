require 'api/async'
require 'api/compiler'
require 'api/product'
require 'api/resource'
require 'api/type'
require 'google/yaml_validator'

class ProductApi

  attr_reader :api

  def initialize(product_name)
    @api = get_api(product_name)
  end

  def resource(name)
    @api.objects.select { |obj| obj.name == name }.first
  end

  def all_resource_names
    @api.objects.map(&:name)
  end

  private

  def get_api
    Api::Compiler.new("products/#{product_name}/api.yaml").run 
  end
end

require 'net/http'
require 'json'

class Discovery
  attr_reader :results

  def initialize(url)
    @url = url
    @results = send_request(url)
  end

  def resource(name)
    schema = @results['schemas'][name]
    DiscoveryObject.new(schema)
  end

  private

  def send_request(url)
    JSON.parse(Net::HTTP.get(URI(url)))
  end
end

class DiscoveryObject
  attr_reader :schema

  def initialize(schema)
    @schema = schema
  end

  def exists?
    !@schema.nil?
  end
end

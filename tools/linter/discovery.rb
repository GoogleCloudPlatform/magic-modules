require 'net/http'
require 'json'

class Discovery
  attr_reader :results

  def initialize(url)
    @url = url
    @results = send_request(url)
  end

  def resource(name)
    DiscoveryObject.new(
      @results['schemas'][name]
    )
  end

  private

  def send_request(url)
    JSON.parse(Net::HTTP.get(URI(url)))
  end
end

class DiscoveryObject
  attr_reader :schema

  def initialize(schema, methods)
    @schema = schema
  end
end

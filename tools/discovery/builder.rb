require 'net/http'
require 'json'
require 'pp'
require 'erb'

DISCOVERY_URL = 'https://www.googleapis.com/discovery/v1/apis/compute/v1/rest'
uri = URI(DISCOVERY_URL)
response = Net::HTTP.get(uri)
results = JSON.parse(response)

renderer = ERB.new(File.read('api.yaml.erb'), nil, '-%>')
File.write('output.yaml', renderer.result())

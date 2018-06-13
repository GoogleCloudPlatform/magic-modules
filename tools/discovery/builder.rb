$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '..', '..')
require 'net/http'
require 'json'
require 'erb'
require 'compile/core'

include Compile::Core

DISCOVERY_URL = 'https://www.googleapis.com/discovery/v1/apis/compute/v1/rest'
uri = URI(DISCOVERY_URL)
response = Net::HTTP.get(uri)
results = JSON.parse(response)

File.write('output.yaml', lines(compile_file({ results: results }, 'api.yaml.erb')))

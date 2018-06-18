$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '..', '..')
require 'net/http'
require 'json'
require 'erb'
require 'compile/core'
require 'tools/discovery/types'
require 'optparse'

# rubocop:disable Style/MixinUsage
include Compile::Core
# rubocop:enable Style/MixinUsage

options = {
  output: 'output.yaml'
}

OptionParser.new do |parser|
  parser.on('-u', '--url URL', "The discovery URL being parsed") do |v|
    options[:url] = v
  end
  parser.on('-o', '--output FILE', "Output file location") do |v|
    options[:output] = v
  end
end.parse!

raise 'Discovery URL must be specified' unless options[:url]

uri = URI(options[:url])
response = Net::HTTP.get(uri)
results = JSON.parse(response)

res = Product.new(results)
File.write(options[:output], lines(compile_file({ product: res }, 'api.yaml.erb')))

#!/usr/bin/env ruby

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'tools/builder/api'
require 'api/compiler'

raise "Must include a URL, object_name and product" unless ARGV.length == 3
url = ARGV.first
object = ARGV[1]
product = ARGV[2]

discovery = DiscoveryProduct.new(url, object).get_product
new_handwritten = HumanApi.new(discovery, Api::Compiler.new("products/#{product}/api.yaml").run).build
File.write("products/#{product}/api.yaml", YAML.dump(new_handwritten))

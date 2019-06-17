#!/usr/bin/env ruby

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'tools/builder/api'

raise "Must include a URL and an object name" unless ARGV.length == 2
url = ARGV.first
object = ARGV[1]

File.write('api.yaml', YAML.dump(DiscoveryProduct.new(url, object).get_product))

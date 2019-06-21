#!/usr/bin/env ruby

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'tools/builder/api'
require 'api/compiler'
require 'optparse'

options = {}
OptionParser.new do |opts|
  opts.banner = "api.yaml builder: run.rb [options]"

  opts.on("-u", "--url URL", "Discovery Doc URL") do |url|
    options[:url] = url
  end

  opts.on("-o", "--object OBJECT", "The object you want to generate") do |obj|
    options[:obj] = obj
  end

  opts.on("-p", "--product product", "The name of the product you're building") do |prod|
    options[:prod] = prod
  end
end.parse!

raise "Must include a URL, object_name and product" unless options.keys.length == 3

discovery = DiscoveryProduct.new(options[:url], options[:obj]).get_product
new_handwritten = HumanApi.new(discovery, Api::Compiler.new("products/#{options[:product]}/api.yaml").run).build
File.write("api.yaml", YAML.dump(new_handwritten))

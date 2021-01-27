#!/usr/bin/env ruby

# Copyright 2018 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'tools/builder/api'
require 'api/compiler'
require 'optparse'

options = {}
# rubocop:disable Metrics/LineLength
OptionParser.new do |opts|
  opts.banner = 'api.yaml builder: run.rb [options]'

  opts.on('-u', '--url URL', 'Discovery Doc URL') do |url|
    options[:url] = url
  end

  opts.on('-o', '--object OBJECT', 'The objects you want to generate (comma-separated)') do |obj|
    options[:obj] = obj
  end

  opts.on('-p', '--product product', "The name of the product you're building (in products/ format") do |prod|
    options[:prod] = prod
  end
end.parse!
# rubocop:enable Metrics/LineLength

raise 'Must include a URL, object_name and product' unless options.keys.length == 3

discovery = DiscoveryProduct.new(options[:url], options[:obj]).product
handwritten = if File.exist?("#{options[:product]}/api.yaml")
                Api::Compiler.new("#{options[:product]}/api.yaml").run
              end
new_handwritten = HumanApi.new(discovery, handwritten).build
File.write('api.yaml', YAML.dump(new_handwritten))

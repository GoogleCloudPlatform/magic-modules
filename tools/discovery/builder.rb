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
  parser.on('-u', '--url URL', 'The discovery URL being parsed') do |v|
    options[:url] = v
  end
  parser.on('-o', '--output FILE', 'Output file location') do |v|
    options[:output] = v
  end
end.parse!

raise 'Discovery URL must be specified' unless options[:url]

uri = URI(options[:url])
response = Net::HTTP.get(uri)
results = JSON.parse(response)

res = Product.new(results)
File.write(options[:output],
           lines(compile_file({ product: res }, 'api.yaml.erb')))

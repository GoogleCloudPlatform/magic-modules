#!/usr/bin/env ruby
# Copyright 2017 Google Inc.
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

$LOAD_PATH.unshift File.dirname(__FILE__)

# Run from compiler dir so all references are relative to the compiler
# executable. This allows the following command line:
#   puppet-google-compute$ ../puppet-codegen/compiler products/compute $PWD
Dir.chdir(File.dirname(__FILE__))

# Our default timezone is UTC, to avoid local time compromise test code seed
# generation.
ENV['TZ'] = 'UTC'

require 'api/compiler'
require 'optparse'
require 'provider/ansible'
require 'provider/ansible/bundle'
require 'provider/chef'
require 'provider/chef/bundle'
require 'provider/example'
require 'provider/puppet'
require 'provider/puppet/bundle'
require 'provider/terraform'
require 'pp' if ENV['COMPILER_DEBUG']

catalog = nil
output = nil
provider = nil
types_to_generate = []
version = nil

ARGV << '-h' if ARGV.empty?

OptionParser.new do |opt|
  opt.on('-p', '--product PRODUCT', 'Folder with product catalog') do |p|
    catalog = p
  end
  opt.on('-o', '--output OUTPUT', 'Folder for module output') do |o|
    output = o
  end
  opt.on('-e', '--engine ENGINE', 'Technology to build for') do |e|
    provider = "#{e}.yaml"
  end
  opt.on('-t', '--type TYPE[,TYPE...]', Array, 'Types to generate') do |t|
    types_to_generate = t
  end
  opt.on('-v', '--version VERSION', 'API version to generate') do |v|
    version = v
  end
  opt.on('-h', '--help', 'Show this message') do
    puts opt
    exit
  end
end.parse!

raise 'Option -p/--product is a required parameter' if catalog.nil?
raise 'Option -o/--output is a required parameter' if output.nil?
raise 'Option -e/--engine is a required parameter' if provider.nil?

raise "Product '#{catalog}' does not have api.yaml" \
  unless File.exist?(File.join(catalog, 'api.yaml'))
raise "Product '#{catalog}' does not have #{provider} settings" \
  unless File.exist?(File.join(catalog, provider))

raise "Output '#{output}' is not a directory" unless Dir.exist?(output)

Google::LOGGER.info "Compiling '#{catalog}' output to '#{output}'"
Google::LOGGER.info \
  "Generating types: #{types_to_generate.empty? ? 'ALL' : types_to_generate}"

api = Api::Compiler.new(File.join(catalog, 'api.yaml')).run
api.validate
api.version = version
pp api if ENV['COMPILER_DEBUG']

config = Provider::Config.parse(File.join(catalog, provider), api)
pp config if ENV['COMPILER_DEBUG']

provider = config.provider.new(config, api)
provider.generate output, types_to_generate

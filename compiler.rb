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
#   ruby compiler.rb -p products/compute -e ansible -o build/ansible
Dir.chdir(File.dirname(__FILE__))

# Our default timezone is UTC, to avoid local time compromise test code seed
# generation.
ENV['TZ'] = 'UTC'

require 'active_support/inflector'
require 'api/compiler'
require 'google/logger'
require 'optparse'
require 'pathname'
require 'provider/ansible'
require 'provider/ansible_devel'
require 'provider/inspec'
require 'provider/terraform'
require 'provider/terraform_oics'
require 'provider/terraform_object_library'
require 'pp' if ENV['COMPILER_DEBUG']

products_to_generate = nil
all_products = false
yaml_dump = false
output_path = nil
provider_name = nil
force_provider = nil
types_to_generate = []
version = 'ga'
override_dir = nil

ARGV << '-h' if ARGV.empty?
Google::LOGGER.level = Logger::INFO

# rubocop:disable Metrics/BlockLength
OptionParser.new do |opt|
  opt.on('-p', '--product PRODUCT', Array, 'Folder[,Folder...] to generate resources from') do |p|
    products_to_generate = p
  end
  opt.on('-a', '--all', 'Build all products. Cannot be used with --product.') do
    all_products = true
  end
  opt.on('-y', '--yaml-dump', 'Dump the final api.yaml output to a file.') do
    yaml_dump = true
  end
  opt.on('-o', '--output OUTPUT', 'Folder for module output') do |o|
    output_path = o
  end
  opt.on('-e', '--engine ENGINE', 'Provider ("engine") to build') do |e|
    provider_name = e
  end
  opt.on('-f', '--force PROVIDER', 'Force using a non-default provider') do |e|
    force_provider = e
  end
  opt.on('-t', '--type TYPE[,TYPE...]', Array, 'Types to generate') do |t|
    types_to_generate = t
  end
  opt.on('-v', '--version VERSION', 'API version to generate') do |v|
    version = v
  end
  opt.on('-r', '--override OVERRIDE', 'Directory containing api.yaml overrides') do |r|
    override_dir = r
  end
  opt.on('-h', '--help', 'Show this message') do
    puts opt
    exit
  end
  opt.on('-d', '--debug', 'Show all debug logs') do |_debug|
    Google::LOGGER.level = Logger::DEBUG
  end
end.parse!
# rubocop:enable Metrics/BlockLength

# We use ActiveSupport Inflections to perform common string operations like
# going from camelCase/PascalCase -> snake_case -> Title Case.
# In order to not break the world, they've frozen the list of inflections the
# library uses by default.
# Particularly for initialisms, it may need a little help to generate great
# code.
ActiveSupport::Inflector.inflections(:en) do |inflect|
  inflect.acronym 'TPU'
  inflect.acronym 'VPC'
end

raise 'Cannot use -p/--products and -a/--all simultaneously' \
  if products_to_generate && all_products
raise 'Either -p/--products OR -a/--all must be present' \
  if products_to_generate.nil? && !all_products
raise 'Option -o/--output is a required parameter' if output_path.nil?
raise 'Option -e/--engine is a required parameter' if provider_name.nil?

all_product_files = []
Dir['products/**/api.yaml'].each do |file_path|
  all_product_files.push(File.dirname(file_path))
end

if override_dir
  Dir["#{override_dir}/products/**/api.yaml"].each do |file_path|
    product = File.dirname(Pathname.new(file_path).relative_path_from(override_dir))
    all_product_files.push(product) unless all_product_files.include? product
  end
end

products_to_generate = all_product_files if all_products
raise 'No api.yaml files found. Check your provider (-e) name.' if products_to_generate.empty?

start_time = Time.now
Google::LOGGER.info "Generating MM output to '#{output_path}'"
Google::LOGGER.info "Using #{version} version"

# products_for_version entries are a hash of product definitions (:definitions)
# and provider config (:overrides) for the product
products_for_version = []
provider = nil
# rubocop:disable Metrics/BlockLength
all_product_files.each do |product_name|
  product_override_path = ''
  product_override_path = File.join(override_dir, product_name, 'api.yaml') if override_dir
  product_yaml_path = File.join(product_name, 'api.yaml')

  provider_override_path = ''
  provider_override_path = File.join(override_dir, product_name, "#{provider_name}.yaml") \
    if override_dir
  provider_yaml_path = File.join(product_name, "#{provider_name}.yaml")

  unless File.exist?(product_yaml_path) || File.exist?(product_override_path)
    raise "#{product_name} does not contain an api.yaml file"
  end

  if File.exist?(product_override_path)
    result = if File.exist?(product_yaml_path)
               YAML.load_file(product_yaml_path).merge(YAML.load_file(product_override_path))
             else
               YAML.load_file(product_override_path)
             end
    product_yaml = result.to_yaml
  else
    product_yaml = File.read(product_yaml_path)
  end

  unless File.exist?(provider_yaml_path) || File.exist?(provider_override_path)
    Google::LOGGER.info "#{product_name}: Skipped as no #{provider_name}.yaml file exists"
    next
  end

  raise "Output path '#{output_path}' does not exist or is not a directory" \
    unless Dir.exist?(output_path)

  product_api = Api::Compiler.new(product_yaml).run
  product_api.validate
  pp product_api if ENV['COMPILER_DEBUG']

  unless product_api.exists_at_version_or_lower(version)
    Google::LOGGER.info \
      "'#{product_name}' does not have a '#{version}' version, skipping"
    next
  end

  if File.exist?(provider_yaml_path)
    product_api, provider_config, = \
      Provider::Config.parse(provider_yaml_path, product_api, version)
  end
  # Load any dynamic overrides passed in with -r
  if File.exist?(provider_override_path)
    product_api, provider_config, = \
      Provider::Config.parse(provider_override_path, product_api, version, override_dir)
  end

  Google::LOGGER.info "#{product_name}: Compiling provider config"
  pp provider_config if ENV['COMPILER_DEBUG']

  if force_provider.nil?
    provider = provider_config.provider.new(provider_config, product_api, version, start_time)
  else
    override_providers = {
      'oics' => Provider::TerraformOiCS,
      'validator' => Provider::TerraformObjectLibrary,
      'ansible_devel' => Provider::Ansible::Devel
    }

    provider_class = override_providers[force_provider]
    if provider_class.nil?
      raise "Invalid force provider option #{force_provider}." \
        + "\nPossible values #{override_providers} "
    end

    provider = \
      override_providers[force_provider].new(provider_config, product_api, version, start_time)
  end

  # provider_config is mutated by instantiating a provider
  products_for_version.push(definitions: product_api, overrides: provider_config)

  unless products_to_generate.include?(product_name)
    Google::LOGGER.info "#{product_name}: Not specified, skipping generation"
    next
  end

  Google::LOGGER.info \
    "#{product_name}: Generating types: #{types_to_generate.empty? ? 'ALL' : types_to_generate}"
  provider.generate output_path, types_to_generate, product_name, yaml_dump
end

# In order to only copy/compile files once per provider this must be called outside
# of the products loop. This will get called with the provider from the final iteration
# of the loop
provider&.copy_common_files(output_path)
Google::LOGGER.info "Compiling common files for #{provider_name}"
common_compile_file = "provider/#{provider_name}/common~compile.yaml"
provider&.compile_common_files(
  output_path,
  products_for_version.sort_by { |p| p[:definitions].name.downcase },
  common_compile_file
)
if override_dir
  Google::LOGGER.info "Compiling override common files for #{provider_name}"
  common_compile_file = "#{override_dir}/common~compile.yaml"
  provider&.compile_common_files(
    output_path,
    products_for_version.sort_by { |p| p[:definitions].name.downcase },
    common_compile_file,
    override_dir
  )
end
# rubocop:enable Metrics/BlockLength

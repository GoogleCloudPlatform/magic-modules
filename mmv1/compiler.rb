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

# Run from compiler dir so all references are relative to compiler.rb
Dir.chdir(File.dirname(__FILE__))

# Our default timezone is UTC, to avoid local time compromise test code seed
# generation.
ENV['TZ'] = 'UTC'

require 'api/compiler'
require 'openapi_generate/parser'
require 'google/logger'
require 'optparse'
require 'parallel'
require 'pathname'
require 'provider/terraform'
require 'provider/terraform_kcc'
require 'provider/terraform_oics'
require 'provider/terraform_tgc'
require 'provider/terraform_tgc_cai2hcl'

products_to_generate = nil
all_products = false
yaml_dump = false
generate_code = true
generate_docs = true
output_path = nil
provider_name = nil
force_provider = nil
types_to_generate = []
version = 'ga'
override_dir = nil
openapi_generate = false

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
  opt.on('-y', '--yaml-dump', 'Dump the final yaml output to a file.') do
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
  opt.on('-r', '--override OVERRIDE', 'Directory containing yaml overrides') do |r|
    override_dir = r
  end
  opt.on('-h', '--help', 'Show this message') do
    puts opt
    exit
  end
  opt.on('-d', '--debug', 'Show all debug logs') do |_debug|
    Google::LOGGER.level = Logger::DEBUG
  end
  opt.on('-c', '--no-code', 'Do not generate code') do
    generate_code = false
  end
  opt.on('-g', '--no-docs', 'Do not generate documentation') do
    generate_docs = false
  end
  opt.on('--openapi-generate', 'Generate MMv1 YAML from openapi directory (Experimental)') do
    openapi_generate = true
  end
end.parse!
# rubocop:enable Metrics/BlockLength

raise 'Cannot use -p/--products and -a/--all simultaneously' \
  if products_to_generate && all_products
raise 'Either -p/--products OR -a/--all must be present' \
  if products_to_generate.nil? && !all_products
raise 'Option -o/--output is a required parameter' if output_path.nil?
raise 'Option -e/--engine is a required parameter' if provider_name.nil?

if openapi_generate
  # Test write OpenAPI --> YAML
  # This writes to a fake demo product currently. In the future this should
  # produce the entire product folder including product.yaml for a single OpenAPI spec
  OpenAPIGenerate::Parser.new('openapi_generate/openapi/*', 'products').run
  return
end

all_product_files = []
Dir['products/**/product.yaml'].each do |file_path|
  all_product_files.push(File.dirname(file_path))
end

if override_dir
  Google::LOGGER.info "Using override directory '#{override_dir}'"
  Dir["#{override_dir}/products/**/product.yaml"].each do |file_path|
    product = File.dirname(Pathname.new(file_path).relative_path_from(override_dir))
    all_product_files.push(product) unless all_product_files.include? product
  end
end

products_to_generate = all_product_files if all_products
raise 'No product.yaml file found.' if products_to_generate.empty?

start_time = Time.now
Google::LOGGER.info "Generating MM output to '#{output_path}'"
Google::LOGGER.info "Using #{version} version"

allowed_classes = Google::YamlValidator.allowed_classes

# Building compute takes a long time and can't be parallelized within the product
# so lets build it first
all_product_files = all_product_files.sort_by { |product| product == 'products/compute' ? 0 : 1 }

# rubocop:disable Metrics/BlockLength
products_for_version = Parallel.map(all_product_files, in_processes: 8) do |product_name|
  product_override_path = ''
  product_override_path = File.join(override_dir, product_name, 'product.yaml') if override_dir
  product_yaml_path = File.join(product_name, 'product.yaml')

  unless File.exist?(product_yaml_path) || File.exist?(product_override_path)
    raise "#{product_name} does not contain a product.yaml file"
  end

  if File.exist?(product_override_path)
    result = if File.exist?(product_yaml_path)
               YAML.load_file(product_yaml_path, permitted_classes: allowed_classes) \
                   .merge(YAML.load_file(product_override_path, permitted_classes: allowed_classes))
             else
               YAML.load_file(product_override_path, permitted_classes: allowed_classes)
             end
    product_yaml = result.to_yaml
  elsif File.exist?(product_yaml_path)
    product_yaml = File.read(product_yaml_path)
  end

  raise "Output path '#{output_path}' does not exist or is not a directory" \
    unless Dir.exist?(output_path)

  product_api = Api::Compiler.new(product_yaml).run
  product_api.validate
  pp product_api if ENV['COMPILER_DEBUG']

  unless product_api.exists_at_version_or_lower(version)
    Google::LOGGER.info \
      "#{product_name} does not have a '#{version}' version, skipping"
    next
  end

  if File.exist?(product_yaml_path) || File.exist?(product_override_path)
    resources = []
    Dir["#{product_name}/*"].each do |file_path|
      next if File.basename(file_path) == 'product.yaml' \
       || File.extname(file_path) != '.yaml'

      if override_dir
        # Skip if resource will be merged in the override loop
        resource_override_path = File.join(override_dir, file_path)
        next if File.exist?(resource_override_path)
      end
      res_yaml = File.read(file_path)
      resource = Api::Compiler.new(res_yaml).run
      resource.properties = resource.add_labels_related_fields(
        resource.properties_with_excluded, nil
      )
      resource.validate
      resources.push(resource)
    end

    if override_dir
      ovr_prod_dir = File.join(override_dir, product_name)
      Dir["#{ovr_prod_dir}/*"].each do |override_path|
        next if File.basename(override_path) == 'product.yaml' \
        || File.extname(override_path) != '.yaml'

        file_path = File.join(product_name, File.basename(override_path))
        res_yaml = if File.exist?(file_path)
                     YAML.load_file(file_path, permitted_classes: allowed_classes) \
                         .merge(YAML \
                           .load_file(override_path, permitted_classes: allowed_classes)) \
                         .to_yaml
                   else
                     File.read(override_path)
                   end
        unless override_dir.nil?
          # Replace overrides directory if we are running with a provider override
          # This allows providers to reference files in their override path
          res_yaml = res_yaml.gsub('{{override_path}}', override_dir)
        end
        resource = Api::Compiler.new(res_yaml).run
        resource.properties = resource.add_labels_related_fields(
          resource.properties_with_excluded, nil
        )
        resource.validate
        resources.push(resource)
      end
    end
    resources = resources.sort_by(&:name)
    product_api.set_variable(resources, 'objects')
  end

  product_api&.validate

  if force_provider.nil?
    provider = \
      Provider::Terraform.new(product_api, version, start_time)
  else
    override_providers = {
      'oics' => Provider::TerraformOiCS,
      'validator' => Provider::TerraformGoogleConversion,
      'tgc' => Provider::TerraformGoogleConversion,
      'tgc_cai2hcl' => Provider::CaiToTerraformConversion,
      'kcc' => Provider::TerraformKCC
    }

    provider_class = override_providers[force_provider]
    if provider_class.nil?
      raise "Invalid force provider option #{force_provider}." \
        + "\nPossible values #{override_providers} "
    end

    provider = \
      override_providers[force_provider].new(product_api, version, start_time)
  end

  unless products_to_generate.include?(product_name)
    Google::LOGGER.info "#{product_name}: Not specified, skipping generation"
    next { definitions: product_api, provider: provider } # rubocop:disable Style/HashSyntax
  end

  Google::LOGGER.info \
    "#{product_name}: Generating types: #{types_to_generate.empty? ? 'ALL' : types_to_generate}"
  provider.generate(
    output_path,
    types_to_generate,
    product_name,
    yaml_dump,
    generate_code,
    generate_docs
  )

  # we need to preserve a single provider instance to use outside of this loop.
  { definitions: product_api, provider: provider } # rubocop:disable Style/HashSyntax
end

# remove any nil values
products_for_version = products_for_version.compact.sort_by { |p| p[:definitions].name.downcase }

# In order to only copy/compile files once per provider this must be called outside
# of the products loop. This will get called with the provider from the final iteration
# of the loop
final_product = products_for_version.last
provider = final_product[:provider]

provider&.copy_common_files(output_path, generate_code, generate_docs)
Google::LOGGER.info "Compiling common files for #{provider_name}"
common_compile_file = "provider/#{provider_name}/common~compile.yaml"
if generate_code
  provider&.compile_common_files(
    output_path,
    products_for_version,
    common_compile_file
  )

  if override_dir
    Google::LOGGER.info "Compiling override common files for #{provider_name}"
    common_compile_file = "#{override_dir}/common~compile.yaml"
    provider&.compile_common_files(
      output_path,
      products_for_version,
      common_compile_file,
      override_dir
    )
  end
end
# rubocop:enable Metrics/BlockLength

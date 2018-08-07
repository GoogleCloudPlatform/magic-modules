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

require 'json'
require 'provider/chef'
require 'provider/config'
require 'provider/core'

IGNORED_COOKBOOKS = %w[cloud gauth].freeze

module Provider
  # A provider to generate the "bundle" module.
  class ChefBundle < Provider::Core
    # The configuration for the "bundle" module (in chef.yaml)
    class Config < Provider::Config
      attr_reader :manifest

      def provider
        Provider::ChefBundle
      end

      def validate
        check_property :manifest, Provider::ChefBundle::Manifest
      end
    end

    # A manifest for the "bundle" module
    class Manifest < Provider::Chef::Manifest
      attr_reader :products
      attr_accessor :releases
    end

    def generate(output_folder, _types, _version_name)
      # Let's build all the dependencies off of the products we found on our
      # path and has the corresponding provider.yaml file
      @config.manifest.depends.concat(
        products.map do |k, v|
          Provider::Config::Requirements.create(
            "google-#{k.prefix}",
            "< #{next_version(v.manifest.version)}"
          )
        end
      )

      copy_files(output_folder)
      compile_changelog(output_folder)
      compile_files(output_folder)
    end

    def products
      all_products.reject { |k, _v| IGNORED_COOKBOOKS.include?(k.prefix) }
    end

    def product_details
      products.map do |product, config|
        {
          name: product.name,
          prefix: product.prefix,
          objects: (product.objects || []).reject(&:exclude),
          source: config.manifest.source
        }
      end \
      + @config.manifest.products.map do |product|
        {
          name: product.display_name,
          prefix: product.prefix,
          source: product.source
        }
      end
    end

    def releases
      all_products.map { |k, v| { k.prefix[1..-1] => v.manifest.version } }
                  .reduce({}, :merge)
    end

    private

    def release_files
      Dir.glob('products/**/chef.yaml')
    end

    def all_products
      @products ||= begin
        prod_map =
          release_files.map do |product_config|
            product =
              Api::Compiler.new(File.join(File.dirname(product_config),
                                          'api.yaml')).run
            product.validate
            config = Provider::Config.parse(product_config, product)
            config.validate

            [product, config]
          end
        Hash[prod_map.sort_by { |p| p[0].prefix }]
      end
    end

    def next_version(version)
      [Gem::Version.new(version).bump, 0].join('.')
    end
  end
end

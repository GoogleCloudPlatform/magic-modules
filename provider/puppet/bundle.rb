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
require 'provider/puppet'
require 'provider/config'
require 'provider/core'

IGNORED_MODULES = %w[cloud gauth].freeze

module Provider
  # A provider to generate the "bundle" module.
  class PuppetBundle < Provider::Core
    # A manifest for the "bundle" module
    class Manifest < Provider::Puppet::Manifest
      attr_accessor :releases
    end

    # The configuration for the "bundle" module (in puppet.yaml)
    class Config < Provider::Config
      attr_reader :manifest

      def provider
        ::Provider::PuppetBundle
      end
    end

    def generate(output_folder, _types, version_name)
      # Let's build all the dependencies off of the products we found on our
      # path and has the corresponding provider.yaml file
      @config.manifest.releases = releases
      generate_requirements

      # Always include authentication module
      @config.manifest.requires \
        << Provider::Config::Requirements \
           .create('google/gauth', ">= #{auth_ver} < #{next_version(auth_ver)}")

      compile_changelog(output_folder)
      copy_files(output_folder)
      compile_files(output_folder, version_name)
    end

    def products
      all_products.reject { |k, _v| IGNORED_MODULES.include?(k.prefix) }
    end

    def releases
      all_products.map { |k, v| { k.prefix[1..-1] => v.manifest.version } }
                  .reduce({}, :merge)
    end

    private

    def release_files
      Dir.glob('products/**/puppet.yaml')
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

    def auth_ver
      auth_meta = JSON.parse(IO.read('build/puppet/auth/metadata.json'))
      auth_meta['version']
    end

    def next_version(version)
      [Gem::Version.new(version).bump, 0].join('.')
    end

    def generate_requirements
      @config.manifest.requires.concat(
        products.map do |k, _|
          version = @config.manifest.releases[k.prefix[1..-1]]
          Provider::Config::Requirements.create(
            "google/#{k.prefix}", ">= #{version} < #{next_version(version)}"
          )
        end
      ).sort_by!(&:name)
    end
  end
end

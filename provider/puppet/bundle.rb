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

module Provider
  # A provider to generate the "bundle" module.
  class PuppetBundle < Provider::Core
    # A manifest for the "bundle" module
    class Manifest < Puppet::Manifest
      attr_reader :releases

      def validate
        @requires = [] if @requires.nil?
        check_property :releases, Hash
        super
      end
    end

    # The configuration for the "bundle" module (in puppet.yaml)
    class Config < Provider::Config
      attr_reader :manifest

      def provider
        ::Provider::PuppetBundle
      end

      def validate
        check_property :manifest, Provider::PuppetBundle::Manifest
        super
      end
    end

    def generate(output_folder, _types)
      # Let's build all the dependencies off of the products we found on our
      # path and has the corresponding provider.yaml file
      generate_requirements

      # Always include authentication module
      @config.manifest.requires \
        << Provider::Config::Requirements \
           .create('google/gauth', ">= #{auth_ver} < #{next_version(auth_ver)}")

      compile_changelog(output_folder)
      copy_files(output_folder)
      compile_files(output_folder)
    end

    def products
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

    private

    def release_files
      @config.manifest.releases.map { |p, _| "products/#{p}/puppet.yaml" }
             .select { |c| File.exist?(c) }
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

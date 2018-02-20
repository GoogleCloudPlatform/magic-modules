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
require 'provider/ansible'
require 'provider/config'
require 'provider/core'

module Provider
  # A provider to generate the shared files in Ansible.
  # Certain files (mainly the lookup functions) will be a combination
  # of functions from multiple products.
  class AnsibleBundle < Provider::Ansible
    # The configuration for the "bundle" module (in ansible.yaml)
    class Config < Provider::Config
      attr_reader :manifest

      def provider
        Provider::AnsibleBundle
      end

      def validate
        check_property :manifest, Provider::AnsibleBundle::Manifest
      end
    end

    # A manifest for the "bundle" module
    class Manifest < Provider::Ansible::Manifest
      def validate
        @requires = []
        super
      end
    end

    def generate(output_folder, _types)
      compile_files(output_folder)
    end

    def products
      @products ||= begin
        prod_map = Dir['products/**/ansible.yaml']
                   .reject { |f| f.include?('bundle') }
                   .map do |product_config|
                     product = Api::Compiler.new(
                       File.join(File.dirname(product_config), 'api.yaml')
                     ).run
                     product.validate
                     config = Provider::Config.parse(product_config, product)
                     config.validate

                     [product, config]
                   end
        Hash[prod_map.sort_by { |p| p[0].prefix }]
      end
    end
  end
end

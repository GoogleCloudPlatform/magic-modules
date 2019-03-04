# Copyright 2019 Google Inc.
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

CURRENT_ANSIBLE_VERSION = '2.8'

module Provider
  module Ansible
    module VersionAdded
      def build_version_added
        product_name = @api.name.downcase
        versions_file = "products/#{product_name}/ansible_version_added.yaml"
        raise "File not found" unless File.exist?(versions_file)

        versions = if File.exist?(versions_file)
                     YAML.load(File.read(versions_file))
                   else
                     {}
                   end

        struct = {
          facts: {},
          regular: {}
        }

        # Build out paths for regular modules.
        @api.objects.each do |obj|
          resource = {
            version_added: correct_version([:regular, obj.name], versions)
          }

          # Add properties.
          obj.all_user_properties.each do |prop|
            resource[prop.name.to_sym] = property_version(prop, [:regular, obj.name], versions)
          end
          struct[:regular][obj.name.to_sym] = resource

          # Add facts modules from facts datasources.
          struct[:facts][obj.name.to_sym] = {
            version_added: correct_version([:facts, obj.name], versions)
          }
        end

        # Write back to disk.
        File.write("products/#{product_name}/ansible_version_added.yaml", struct.to_yaml)

        struct
      end

      private

      # Builds out property information (with nesting)
      def property_version(prop, path, struct)
        property_hash = {
          version_added: correct_version(path + [prop.name], struct)
        }

        prop.nested_properties.each do |nested_p|
          property_hash[nested_p.name.to_sym] = property_version(nested_p, path + [prop.name], struct)
        end
        property_hash
      end

      def correct_version(path, struct)
        path = path.map(&:to_sym) + [:version_added]
        struct.dig(*path) || CURRENT_ANSIBLE_VERSION
      end
    end
  end
end

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

CURRENT_ANSIBLE_VERSION = '2.8'.freeze

module Provider
  module Ansible
    # All logic involved with altering the version_added yaml file and reading
    # values from it.
    module VersionAdded
      def build_version_added
        product_name = @api.name.downcase
        versions_file = "products/#{product_name}/ansible_version_added.yaml"
        raise 'File not found' unless File.exist?(versions_file)

        versions = if File.exist?(versions_file)
                     YAML.safe_load(File.read(versions_file), [Symbol])
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

      def version_added(object, type = :regular)
        if object.is_a?(Api::Resource)
          correct_version([type, object.name], @version_added)
        else
          path = [type] + build_path(object)
          res_version = correct_version(path[0, 2], @version_added)
          prop_version = correct_version(path, @version_added)
          # We don't need a version added if it matches the resource.
          return nil if res_version == prop_version
          # If our property is the same as the properties above it, we don't
          # need a version added.
          return nil if version_path(path).sort == version_path(path)

          prop_version
        end
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

      # Build out the path of resources/properties that this property exists within.
      def build_path(prop)
        path = []
        while prop
          path << prop if !path.last || (path.last.name != prop.name)
          prop = prop.__parent
        end
        [path.last.__resource.name] + path.map(&:name).reverse
      end

      # Given a path of resources/properties, return the same path, but with versions substituted for names.
      def version_path(path)
        version_path = []
        path.length.times.each do |i|
          version_path << correct_version(path[0, i], @version_added)
        end
        version_path
      end
    end
  end
end

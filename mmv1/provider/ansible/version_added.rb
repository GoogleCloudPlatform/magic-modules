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

CURRENT_ANSIBLE_VERSION = '2.10'.freeze

module Provider
  module Ansible
    # All logic involved with altering the version_added yaml file and reading
    # values from it.
    module VersionAdded
      def build_version_added
        product_name = @api.name.downcase
        versions_file = "products/#{product_name}/ansible_version_added.yaml"

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
        @api.objects.reject(&:exclude).each do |obj|
          next if obj.not_in_version?(@api.version_obj_or_closest('ga'))

          resource = {
            version_added: correct_version([:regular, obj.name], versions)
          }

          # Add properties.
          # Only properties that aren't output-only + excluded should get versions.
          # These are the only properties that become module fields.
          obj.all_user_properties.reject(&:exclude).reject(&:output).each do |prop|
            next if prop.min_version > @api.version_obj_or_closest('ga')

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

      # This fetches a version_added from the config file for a Resource or Property.
      # While the machine-generated config has every property,
      # this function only returns a version_added if it cannot be inferred from
      # elsewhere in the module.
      def version_added(object, type = :regular)
        if object.is_a?(Api::Resource)
          correct_version([type, object.name], @version_added)
        else
          path = [type] + build_path(object)
          res_version = correct_version(path[0, 2], @version_added)
          prop_version = correct_version(path, @version_added)
          # We don't need a version added if it matches the resource.
          return nil if res_version == prop_version

          # If property is the same as the one above it, ignore it.
          return nil if version_path(path).last == version_path(path)[-2]

          prop_version
        end
      end

      private

      # Builds out property information (with nesting)
      def property_version(prop, path, struct)
        property_hash = {
          version_added: correct_version(path + [prop.name], struct)
        }

        # Only properties that aren't output-only + excluded should get versions.
        # These are the only properties that become module fields.
        prop.nested_properties.reject(&:exclude).reject(&:output).each do |nested_p|
          property_hash[nested_p.name.to_sym] = property_version(nested_p,
                                                                 path + [prop.name], struct)
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

      # Given a path of resources/properties, return the same path, but with
      # versions substituted for names.
      def version_path(path)
        version_path = []
        (path.length - 1).times.each do |i|
          version_path << correct_version(path[0, i + 2], @version_added)
        end
        version_path
      end
    end
  end
end

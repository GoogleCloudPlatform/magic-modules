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

require 'provider/config'
require 'provider/core'
require 'provider/ansible/manifest'

module Provider
  # Code generator for Ansible Cookbooks that manage Google Cloud Platform
  # resources.
  class Ansible < Provider::Core
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      def provider
        Provider::Ansible
      end

      def validate
        super
        check_property :manifest, Provider::Ansible::Manifest \
          unless manifest.nil?
      end
    end

    # Takes in a string and returns a multi-line string, where each line
    # is less than max_length characters long and all subsequent lines are
    # indented in by spaces characters
    #
    # Example:
    #   - This is a sentence
    #     that wraps under
    #     the bullet properly
    def bullet_lines(line, spaces)
      # - 2 for "- "
      indented = wrap_field(line, spaces - 2)
      indented = indented.split("\n")
      indented[0] = indented[0].sub(/^../, '- ')
      indented
    end

    # Returns a string representation of the corresponding Python type
    # for a MM type.
    def python_type(type)
      return 'list' if type.is_a? Api::Type::Array
      return 'bool' if type.is_a? Api::Type::Boolean
      return 'int' if type.is_a? Api::Type::Integer
      'str'
    end

    private

    def generate_resource(data)
      target_folder = data[:output_folder]
      FileUtils.mkpath target_folder
      name = Google::StringUtils.underscore(data[:name])
      generate_resource_file data.clone.merge(
        default_template: 'templates/ansible/resource.erb',
        out_file: File.join(target_folder,
                            "lib/ansible/modules/cloud/google/#{name}.py")
      )
    end

    def generate_resource_tests(data) end

    def generate_network_datas(data, object) end

    def generate_base_property(data) end

    def generate_simple_property(type, data) end

    def generate_typed_array(type, data) end

    def emit_nested_object(data) end

    def emit_resourceref_object(data) end
  end
end

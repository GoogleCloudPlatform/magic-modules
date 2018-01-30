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
    INTEGRATION_TEST_DEFAULTS = {
      project: '"{{ gcp_project }}"',
      auth_kind: '"{{ gcp_cred_kind }}"',
      service_account_file: '"{{ gcp_cred_file }}"',
      name: '"{{ resource_name }}"'
    }.freeze

    EXAMPLE_DEFAULTS = {
      name: "testObject",
      project: "testProject",
      auth_kind: "service_account",
      service_account_file: "/tmp/auth.pem"
    }

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
      # All ResourceRefs are dicts with properties.
      return 'dict' if type.is_a? Api::Type::ResourceRef
      return 'list' if type.is_a? Api::Type::Array
      return 'bool' if type.is_a? Api::Type::Boolean
      return 'int' if type.is_a? Api::Type::Integer
      'str'
    end

    # Returns a unicode formatted, quoted string.
    def unicode_string(string)
      return 'Invalid value' if string.nil?
      return "u#{quote_string(string)}" unless string.include? 'u\''
    end

    def self_link_url(resource)
      (product_url, resource_url) = self_link_raw_url(resource)
      full_url = [product_url, resource_url].flatten.join
      # Double {} replaced with single {} to support Python string interpolation
      "\"#{full_url.gsub('{{', '{').gsub('}}', '}')}\""
    end

    def collection_url(resource)
      base_url = resource.base_url.split("\n").map(&:strip).compact
      full_url = [resource.__product.base_url, base_url].flatten.join
      # Double {} replaced with single {} to support Python string interpolation
      "\"#{full_url.gsub('{{', '{').gsub('}}', '}')}\""
    end

    # Returns the name of the module according to Ansible naming standards.
    # Example: gcp_dns_managed_zone
    def module_name(object)
      ["gcp_#{object.__product.prefix[1..-1]}",
       Google::StringUtils.underscore(object.name)].join('_')
    end

    private

    def generate_resource(data)
      prod_name = Google::StringUtils.underscore(data[:object].name)
      path = ["products/#{data[:product_name]}",
              "examples/ansible/#{prod_name}.yaml"].join('/')

      example = compile_template_with_hash(File.open(path).read,
                                           EXAMPLE_DEFAULTS) if File.file?(path)

      target_folder = data[:output_folder]
      FileUtils.mkpath target_folder
      name = module_name(data[:object])
      generate_resource_file data.clone.merge(
        default_template: 'templates/ansible/resource.erb',
        out_file: File.join(target_folder,
                            "lib/ansible/modules/cloud/google/#{name}.py"),
        example: example
      )
    end

    def generate_resource_tests(data)
      prod_name = Google::StringUtils.underscore(data[:object].name)
      path = ["products/#{data[:product_name]}",
              "examples/ansible/#{prod_name}.yaml"].join('/')

      # Unlike other providers, all resources will not be built at once or
      # in close timing to each other (due to external PRs).
      # This means that examples might not be built out for every resource
      # in a GCP product.
      return unless File.file?(path)

      target_folder = data[:output_folder]
      FileUtils.mkpath target_folder

      name = module_name(data[:object])
      generate_resource_file data.clone.merge(
        default_template: 'templates/ansible/example.erb',
        out_file: File.join(target_folder,
                            "test/integration/targets/#{name}/tasks/main.yml"),
        example: compile_template_with_hash(File.open(path).read,
                                            INTEGRATION_TEST_DEFAULTS)
      )
    end

    def compile_template_with_hash(template, hash)
      ERB.new(template).result(OpenStruct.new(hash).instance_eval { binding })
    end

    def generate_network_datas(data, object) end

    def generate_base_property(data) end

    def generate_simple_property(type, data) end

    def generate_typed_array(type, data) end

    def emit_nested_object(data) end

    def emit_resourceref_object(data) end
  end
end

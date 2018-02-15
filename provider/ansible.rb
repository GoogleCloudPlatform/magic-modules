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
require 'provider/ansible/example'

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
        check_optional_property :manifest, Provider::Ansible::Manifest
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
      if type.is_a? Api::Type::ResourceRef
        return 'str' if type.resource_ref.virtual
        return 'dict'
      end
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

    def async_operation_url(resource)
      base_url = resource.__product.base_url
      url = [base_url, resource.async.operation.base_url].join
      "\"#{url.gsub('{{', '{').gsub('}}', '}')}\""
    end

    # Returns the name of the module according to Ansible naming standards.
    # Example: gcp_dns_managed_zone
    def module_name(object)
      ["gcp_#{object.__product.prefix[1..-1]}",
       Google::StringUtils.underscore(object.name)].join('_')
    end

    def build_object_data(object, output_folder)
      # Method is overriden to add Ansible example objects to the data object.
      data = super

      prod_name = Google::StringUtils.underscore(data[:object].name)
      path = ["products/#{data[:product_name]}",
              "examples/ansible/#{prod_name}.yaml"].join('/')

      data.merge(example: (get_example(path) if File.file?(path)))
    end

    # Given a URL and function name, emit a URL.
    # URL functions will have 1-3 parameters.
    # * module will always be included.
    # * extra_data is a dict of extra information.
    # * extra_url will have a URL chunk to be appended after the URL.
    def emit_link(name, url, extra_data = false)
      params = emit_link_var_args(url, extra_data)
      extra = (' + extra' if url.include?('<|extra|>')) || ''
      if extra_data
        [
          "def #{name}(#{params.join(', ')}):",
          indent("if extra_data is None:", 4),
          indent("extra_data = {}", 8),
          indent("url = #{url}#{extra}", 4).gsub('<|extra|>', ''),
          indent("combined = extra_data.copy()", 4),
          indent("combined.update(module.params)", 4),
          indent("return url.format(**combined)", 4),
        ].compact.join("\n")
      else
        url_code = "#{url}.format(**module.params)#{extra}"
        [
          "def #{name}(#{params.join(', ')}):",
          indent("return #{url_code}", 4).gsub('<|extra|>', '')
        ].join("\n")
      end
    end

    def emit_link_var_args(url, extra_data)
      extra_url = url.include?('<|extra|>')
      [
       'module', ("extra_url=''" if extra_url),
       ("extra_data=None" if extra_data)
      ].compact
    end

    # Returns a list of all first-level ResourceRefs that are not virtual
    def nonvirtual_rrefs(object)
      object.all_user_properties
            .select { |prop| prop.is_a? Api::Type::ResourceRef }
            .select { |prop| !prop.resource_ref.virtual }
    end

    # Converts a path in the form a/b/c/d into ['a', 'b', 'c', 'd']
    def path2navigate(path)
      "[#{path.split('/').map { |x| "'#{x}'" }.join(', ')}]"
    end

    private

    def get_example(cfg_file)
      # These examples will have embedded ERB that will be compiled at a later
      # stage.
      ex = Google::YamlValidator.parse(File.read(cfg_file))
      raise "#{cfg_file}(#{ex.class}) is not a Provider::Ansible::Example" \
        unless ex.is_a?(Provider::Ansible::Example)
      ex
    end

    def generate_resource(data)
      target_folder = data[:output_folder]
      FileUtils.mkpath target_folder
      name = module_name(data[:object])
      generate_resource_file data.clone.merge(
        default_template: 'templates/ansible/resource.erb',
        out_file: File.join(target_folder,
                            "lib/ansible/modules/cloud/google/#{name}.py")
      )
    end

    def example_defaults(data)
      obj_name = Google::StringUtils.underscore(data[:object].name)
      path = ["products/#{data[:product_name]}",
              "examples/ansible/#{obj_name}.yaml"].join('/')

      compile_template_with_hash(File.open(path).read, EXAMPLE_DEFAULTS) \
        if File.file?(path)
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
                            "test/integration/targets/#{name}/tasks/main.yml")
      )
    end

    def generate_network_datas(data, object) end

    def generate_base_property(data) end

    def generate_simple_property(type, data) end

    def generate_typed_array(type, data) end

    def emit_nested_object(data) end

    def emit_resourceref_object(data) end
  end
end

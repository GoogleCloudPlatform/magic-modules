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
require 'provider/ansible/config'
require 'provider/ansible/documentation'
require 'provider/ansible/example'
require 'provider/ansible/manifest'
require 'provider/ansible/module'
require 'provider/ansible/property_override'
require 'provider/ansible/request'
require 'provider/ansible/resourceref'
require 'provider/ansible/resource_override'
require 'provider/ansible/property_override'
require 'provider/ansible/facts_override'

module Provider
  module Ansible
    # Code generator for Ansible Cookbooks that manage Google Cloud Platform
    # resources.
    # TODO(alexstephen): Split up class into multiple modules.
    class Core < Provider::Core
      PYTHON_TYPE_FROM_MM_TYPE = {
        'Api::Type::NestedObject' => 'dict',
        'Api::Type::Array' => 'list',
        'Api::Type::Boolean' => 'bool',
        'Api::Type::Integer' => 'int',
        'Api::Type::KeyValuePairs' => 'dict',
        'Provider::Ansible::FilterProp' => 'list'
      }.freeze

      include Provider::Ansible::Documentation
      include Provider::Ansible::Module
      include Provider::Ansible::Request

      def initialize(config, api)
        super(config, api)
        @max_columns = 160
      end

      # Returns a string representation of the corresponding Python type
      # for a MM type.
      def python_type(prop)
        prop = Module.const_get(prop).new('') unless prop.is_a?(Api::Type)
        # All ResourceRefs are dicts with properties.
        if prop.is_a? Api::Type::ResourceRef
          return 'str' if prop.resource_ref.readonly
          return 'dict'
        end
        PYTHON_TYPE_FROM_MM_TYPE.fetch(prop.class.to_s, 'str')
      end

      # Returns a unicode formatted, quoted string.
      def unicode_string(string)
        return 'Invalid value' if string.nil?
        return "u#{quote_string(string)}" unless string.include? 'u\''
      end

      def build_url(url_parts, _extra = false)
        full_url = if url_parts.is_a? Array
                     url_parts.flatten.join
                   else
                     url_parts
                   end

        "\"#{full_url.gsub('{{', '{').gsub('}}', '}')}\""
      end

      # Returns the name of the module according to Ansible naming standards.
      # Example: gcp_dns_managed_zone
      def module_name(object)
        ["gcp_#{object.__product.prefix[1..-1]}",
         object.name.underscore].join('_')
      end

      def build_object_data(object, output_folder, version)
        # Method is overriden to add Ansible example objects to the data object.
        data = super

        prod_name = data[:object].name.underscore
        path = ["products/#{data[:product_name]}",
                "examples/ansible/#{prod_name}.yaml"].join('/')

        data.merge(example: (get_example(path) if File.file?(path)))
      end

      # Given a URL and function name, emit a URL.
      # URL functions will have 1-3 parameters.
      # * module will always be included.
      # * extra_data is a dict of extra information.
      # * extra_url will have a URL chunk to be appended after the URL.
      def emit_link(name, url, object, has_extra_data = false)
        params = emit_link_var_args(url, has_extra_data)
        if rrefs_in_link(url, object)
          url_code = "#{url}.format(**res)"
          [
            "def #{name}(#{params.join(', ')}):",
            indent("res = #{resourceref_hash_for_links(url, object)}", 4),
            indent("return #{url_code}", 4)
          ].join("\n")
        elsif has_extra_data
          [
            "def #{name}(#{params.join(', ')}):",
            indent([
                     'if extra_data is None:',
                     indent('extra_data = {}', 4)
                   ], 4),
            indent("url = #{url}", 4),
            indent([
                     'combined = extra_data.copy()',
                     'combined.update(module.params)',
                     'return url.format(**combined)'
                   ], 4)
          ].compact.join("\n")
        else
          url_code = "#{url}.format(**module.params)"
          [
            "def #{name}(#{params.join(', ')}):",
            indent("return #{url_code}", 4)
          ].join("\n")
        end
      end

      def emit_method(name, args, code, _file_name, _opts = {})
        [
          method_decl(name, args),
          indent(code, 4)
        ].compact.join("\n")
      end

      def rrefs_in_link(link, object)
        props_in_link = link.scan(/{([a-z_]*)}/).flatten
        (object.parameters || []).select do |p|
          props_in_link.include?(p.name.underscore) && \
            p.is_a?(Api::Type::ResourceRef) && !p.resource_ref.readonly
        end.any?
      end

      def resourceref_hash_for_links(link, object)
        props_in_link = link.scan(/{([a-z_]*)}/).flatten
        props = props_in_link.map do |p|
          # Select a resourceref if it exists.
          rref = (object.parameters || []).select do |prop|
            prop.name.underscore == p && \
              prop.is_a?(Api::Type::ResourceRef) && !prop.resource_ref.readonly
          end
          if rref.any?
            [
              "#{quote_string(p)}:",
              "replace_resource_dict(module.params[#{quote_string(p)}],",
              "#{quote_string(rref[0].imports)})"
            ].join(' ')
          else
            "#{quote_string(p)}: module.params[#{quote_string(p)}]"
          end
        end
        ['{', indent_list(props, 4), '}'].join("\n")
      end

      def emit_link_var_args(url, extra_data)
        extra_url = url.include?('<|extra|>')
        [
          'module', ("extra_url=''" if extra_url),
          ('extra_data=None' if extra_data)
        ].compact
      end

      # Returns a list of all first-level ResourceRefs that are not readonly
      def nonreadonly_rrefs(object)
        object.all_resourcerefs
              .reject { |prop| prop.resource_ref.readonly }
      end

      # Converts a path in the form a/b/c/d into ['a', 'b', 'c', 'd']
      def path2navigate(path)
        "[#{path.split('/').map { |x| "'#{x}'" }.join(', ')}]"
      end

      # TODO(alexstephen): Standardize on one version and move to provider/core
      # https://github.com/GoogleCloudPlatform/magic-modules/issues/30
      def wrap_field(field, spaces)
        # field.scan goes from 0 -> avail_columns - 1
        # -1 to account for this
        avail_columns = DEFAULT_FORMAT_OPTIONS[:max_columns] - spaces - 1
        field.scan(/\S.{0,#{avail_columns}}\S(?=\s|$)|\S+/)
      end

      def list_kind(object)
        "#{object.kind}List"
      end

      # Grabs all conflicted properties and returns an array of arrays without
      # any duplicates.
      def conflicted_property_batches(object)
        sets = object.all_user_properties.map do |p|
          if !p.conflicting.empty?
            p.conflicting.map(&:name).map(&:underscore) + [p.name.underscore]
          else
            []
          end
        end
        sets.map(&:sort)
            .uniq
            .reject(&:empty?)
      end

      private

      def get_example(cfg_file)
        # These examples will have embedded ERB that will be compiled at a later
        # stage.
        ex = Google::YamlValidator.parse(File.read(cfg_file))
        raise "#{cfg_file}(#{ex.class}) is not a Provider::Ansible::Example" \
          unless ex.is_a?(Provider::Ansible::Example)
        ex.validate
        ex
      end

      def generate_resource(data)
        add_datasource_info_to_data(data)
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
        obj_name = data[:object].name.underscore
        path = ["products/#{data[:product_name]}",
                "examples/ansible/#{obj_name}.yaml"].join('/')

        compile_file(EXAMPLE_DEFAULTS, path) if File.file?(path)
      end

      def generate_resource_tests(data)
        prod_name = data[:object].name.underscore
        path = ["products/#{data[:product_name]}",
                "examples/ansible/#{prod_name}.yaml"].join('/')

        return unless data[:object].has_tests
        # Unlike other providers, all resources will not be built at once or
        # in close timing to each other (due to external PRs).
        # This means that examples might not be built out for every resource
        # in a GCP product.
        return unless File.file?(path)

        target_folder = data[:output_folder]
        FileUtils.mkpath target_folder

        name = module_name(data[:object])
        generate_resource_file data.clone.merge(
          default_template: 'templates/ansible/integration_test.erb',
          out_file: File.join(target_folder,
                              "test/integration/targets/#{name}/tasks/main.yml")
        )
      end

      def compile_datasource(data)
        target_folder = data[:output_folder]
        FileUtils.mkpath target_folder
        name = "#{module_name(data[:object])}_facts"
        generate_resource_file data.clone.merge(
          default_template: 'templates/ansible/facts.erb',
          out_file: File.join(target_folder,
                              "lib/ansible/modules/cloud/google/#{name}.py")
        )
      end

      def add_datasource_info_to_data(data)
        # We have two sets of overrides - one for regular modules, one for
        # datasources.
        # When building regular modules, we will potentially need some
        # information from the datasource overrides.
        # This method will give the regular module data access to the
        # datasource module overrides.
        name = "@#{data[:object].name}".to_sym
        facts_info = @config&.datasources&.instance_variable_get(name)&.facts
        facts_info ||= Provider::Ansible::FactsOverride.new
        facts_info.validate
        data[:object].instance_variable_set(:@facts, facts_info)
      end
    end
  end
end

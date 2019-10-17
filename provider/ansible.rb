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
require 'provider/ansible/module'
require 'provider/ansible/request'
require 'provider/ansible/resourceref'
require 'provider/ansible/version_added'
require 'provider/ansible/facts_override'
require 'overrides/ansible/resource_override'
require 'overrides/ansible/property_override'

module Provider
  # Ansible Provider module containing helper functions and the Ansible Provider
  # implementation "Core"
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
        'Provider::Ansible::FilterProp' => 'list',
        'Api::Type::Path' => 'path'
      }.freeze

      include Provider::Ansible
      include Provider::Ansible::Documentation
      include Provider::Ansible::Module
      include Provider::Ansible::Request
      include Provider::Ansible::VersionAdded

      # ProductFileTemplate with Ansible specific fields
      class AnsibleProductFileTemplate < Provider::ProductFileTemplate
        # The Ansible example object.
        attr_accessor :example
      end

      def initialize(config, api, version_name, start_time)
        super(config, api, version_name, start_time)
        @version_added = build_version_added
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

      def build_url(url)
        # Return a quoted string, with single pairs of {} brackets and all
        # requested strings are underscored (as they come from the Ansible configs)
        "\"#{url.gsub(/{{\w+}}/) { |param| "{#{param[2..-3].underscore}}" }}\""
      end

      # Returns the name of the module according to Ansible naming standards.
      # Example: gcp_dns_managed_zone
      def module_name(object)
        ["gcp_#{object.__product.api_name}",
         object.name.underscore].join('_')
      end

      def build_object_data(object, output_folder, version)
        # Method is overridden to add Ansible example objects to the data object.
        data = AnsibleProductFileTemplate.file_for_resource(
          output_folder,
          object,
          version,
          @config,
          build_env
        )

        prod_name = data.object.name.underscore
        path = ["products/#{data.product.api_name}",
                "examples/ansible/#{prod_name}.yaml"].join('/')

        data.example = get_example(path) if File.file?(path)
        data
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
        (object.all_user_properties || []).select do |p|
          props_in_link.include?(p.name.underscore) && \
            p.is_a?(Api::Type::ResourceRef) && !p.resource_ref.readonly
        end.any?
      end

      def resourceref_hash_for_links(link, object)
        props_in_link = link.scan(/{([a-z_]*)}/).flatten
        props = props_in_link.map do |p|
          # Select a resourceref if it exists.
          rref = (object.all_user_properties || []).select do |prop|
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

      def list_kind(object)
        "#{object.kind}List"
      end

      # Grabs all conflicting properties and returns an array of arrays without
      # any duplicates.
      # This does not create an optimal list, but it does create a valid list.
      def conflicting_property_batches(object)
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

        ex.provider = self
        ex.validate
        ex
      end

      def generate_resource(data)
        target_folder = data.output_folder
        name = module_name(data.object)
        path = File.join(target_folder,
                         "plugins/modules/#{name}.py")
        data.generate(
          data.object.template || 'templates/ansible/resource.erb',
          path,
          self
        )
      end

      def generate_resource_tests(data)
        prod_name = data.object.name.underscore
        path = ["products/#{data.product.api_name}",
                "examples/ansible/#{prod_name}.yaml"].join('/')

        return if data.object.tests.tests.empty?

        target_folder = data.output_folder

        name = module_name(data.object)

        # Generate the main file with a list of tests.
        path = File.join(target_folder,
                         "tests/integration/targets/#{name}/tasks/main.yml")
        data.generate(
          'templates/ansible/tests_main.erb',
          path,
          self
        )

        # Generate each of the tests individually
        data.object.tests.tests.each do |t|
          path = File.join(target_folder,
                           "tests/integration/targets/#{name}/tasks/#{t.name}.yml")
          data.generate(
            t.path,
            path,
            self
          )
        end

        # Generate 'defaults' file that contains variables.
        path = File.join(target_folder,
                         "tests/integration/targets/#{name}/defaults/main.yml")
        data.generate(
          'templates/ansible/integration_test_variables.erb',
          path,
          self
        )
      end

      def compile_datasource(data)
        target_folder = data.output_folder
        name = module_name(data.object)
        data.generate('templates/ansible/facts.erb',
                      File.join(target_folder,
                                "plugins/modules/#{name}_info.py"),
                      self)
      end

      def generate_objects(output_folder, types)
        # We have two sets of overrides - one for regular modules, one for
        # datasources.
        # When building regular modules, we will potentially need some
        # information from the datasource overrides.
        # This method will give the regular module data access to the
        # datasource module overrides.
        @api.objects.each do |o|
          facts_info = @config&.datasources&.instance_variable_get("@#{o.name}".to_sym)&.facts
          facts_info ||= Provider::Ansible::FactsOverride.new
          facts_info.validate
          o.instance_variable_set(:@facts, facts_info)
        end
        super
      end
    end

    # Returns all URI properties minus those ignored.
    def uri_properties(object, ignored_props = [])
      uri_properties_raw(object)
        .compact
        .map(&:name)
        .reject { |x| ignored_props.include? x }
    end

    # TODO(alexstephen): Update test_constants to use this function.
    # Returns all of the properties that are a part of the self_link or
    # collection URLs
    def uri_properties_raw(object)
      [object.base_url, object.__product.base_url].map do |url|
        parts = url.scan(/\{\{(.*?)\}\}/).flatten
        parts << 'name'
        parts.delete('project')
        parts.map { |pt| object.all_user_properties.select { |p| p.out_name == pt }[0] }
      end.flatten
    end

    # Convert a URL to a regex.
    def regex_url(url)
      url.gsub(/{{[a-z]*}}/, '.*')
    end

    # Generates files on a per-resource basis.
    # All paths are allowed a '%s' where the module name
    # will be added.
    def generate_resource_files(data)
      return unless @config&.files&.resource

      files = @config.files.resource
                     .map { |k, v| [k % module_name(data.object), v] }
                     .to_h

      file_template = ProductFileTemplate.new(
        data.output_folder,
        data.name,
        @api,
        data.version,
        build_env
      )
      compile_file_list(data.output_folder, files, file_template)
    end

    def copy_common_files(output_folder, provider_name = 'ansible')
      super(output_folder, provider_name)
    end

    def module_utils_import_path
      'ansible_collections.google.cloud.plugins.module_utils.gcp_utils'
    end
  end
end

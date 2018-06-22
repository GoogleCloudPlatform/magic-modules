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
require 'provider/ansible/documentation'
require 'provider/ansible/module'
require 'provider/ansible/request'
require 'provider/ansible/resourceref'
require 'provider/ansible/resource_override'
require 'provider/ansible/property_override'
require 'provider/ansible/selflink'

module Provider
  module Ansible
    # Settings for the Ansible provider
    class Config < Provider::Config
      attr_reader :manifest

      def provider
        Provider::Ansible::Core
      end

      def resource_override
        Provider::Ansible::ResourceOverride
      end

      def property_override
        Provider::Ansible::PropertyOverride
      end

      def validate
        super
        check_optional_property :manifest, Provider::Ansible::Manifest
      end
    end

    # Code generator for Ansible Cookbooks that manage Google Cloud Platform
    # resources.
    # TODO(alexstephen): Split up class into multiple modules.
    # rubocop:disable Metrics/ClassLength
    class Core < Provider::Core
      PYTHON_TYPE_FROM_MM_TYPE = {
        'Api::Type::NestedObject' => 'dict',
        'Api::Type::Array' => 'list',
        'Api::Type::Boolean' => 'bool',
        'Api::Type::Integer' => 'int',
        'Api::Type::NameValues' => 'dict'
      }.freeze

      include Provider::Ansible::Documentation
      include Provider::Ansible::Module
      include Provider::Ansible::Request
      include Provider::Ansible::SelfLink

      def initialize(config, api)
        super(config, api)
        @max_columns = 160
      end

      # Returns a string representation of the corresponding Python type
      # for a MM type.
      def python_type(prop)
        prop = Module.const_get(prop).new('') unless prop.is_a?(Api::Type)
        # All ResourceRefs are dicts with properties.
        # We're assuming that all ResourceRefs in the list act the same.
        if prop.is_a? Api::Type::ResourceRef
          return 'str' if prop.resources.first.resource_ref.virtual
          return 'dict'
        end
        PYTHON_TYPE_FROM_MM_TYPE.fetch(prop.class.to_s, 'str')
      end

      # Returns a unicode formatted, quoted string.
      def unicode_string(string)
        return 'Invalid value' if string.nil?
        return "u#{quote_string(string)}" unless string.include? 'u\''
      end

      def self_link_url(resource)
        (product_url, resource_url) = self_link_raw_url(resource)
        full_url = [product_url, resource_url].flatten.join
        # Double {} replaced with single {} to support Python string
        # interpolation
        "\"#{full_url.gsub('{{', '{').gsub('}}', '}')}\""
      end

      def collection_url(resource)
        base_url = resource.base_url.split("\n").map(&:strip).compact
        full_url = [resource.__product.base_url, base_url].flatten.join
        # Double {} replaced with single {} to support Python string
        # interpolation
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

      def build_object_data(object, output_folder, version)
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
      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/AbcSize
      def emit_link(name, url, object, has_extra_data = false)
        params = emit_link_var_args(url, has_extra_data)
        extra = (' + extra_url' if url.include?('<|extra|>')) || ''
        if rrefs_in_link(url, object)
          url_code = "#{url}.format(**res)#{extra}"
          [
            "def #{name}(#{params.join(', ')}):",
            indent("res = #{resourceref_hash_for_links(url, object)}", 4),
            indent("return #{url_code}", 4).gsub('<|extra|>', '')
          ].join("\n")
        elsif has_extra_data
          [
            "def #{name}(#{params.join(', ')}):",
            indent([
                     'if extra_data is None:',
                     indent('extra_data = {}', 4)
                   ], 4),
            indent("url = #{url}#{extra}", 4).gsub('<|extra|>', ''),
            indent([
                     'combined = extra_data.copy()',
                     'combined.update(module.params)',
                     'return url.format(**combined)'
                   ], 4)
          ].compact.join("\n")
        else
          url_code = "#{url}.format(**module.params)#{extra}"
          [
            "def #{name}(#{params.join(', ')}):",
            indent("return #{url_code}", 4).gsub('<|extra|>', '')
          ].join("\n")
        end
      end
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/AbcSize

      def emit_method(name, args, code, _file_name, _opts = {})
        [
          method_decl(name, args),
          indent(code, 4)
        ].compact.join("\n")
      end

      def rrefs_in_link(link, object)
        props_in_link = link.scan(/{([a-z_]*)}/).flatten
        (object.parameters || []).select do |p|
          props_in_link.include?(Google::StringUtils.underscore(p.name)) && \
            p.is_a?(Api::Type::ResourceRef) && !p.resources.first.resource_ref.virtual
        end.any?
      end

      # rubocop:disable Metrics/AbcSize
      def resourceref_hash_for_links(link, object)
        props_in_link = link.scan(/{([a-z_]*)}/).flatten
        props = props_in_link.map do |p|
          # Select a resourceref if it exists.
          rref = (object.parameters || []).select do |prop|
            Google::StringUtils.underscore(prop.name) == p && \
              prop.is_a?(Api::Type::ResourceRef) && !prop.resource_ref.virtual
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
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/AbcSize

      def emit_link_var_args(url, extra_data)
        extra_url = url.include?('<|extra|>')
        [
          'module', ("extra_url=''" if extra_url),
          ('extra_data=None' if extra_data)
        ].compact
      end

      # Returns a list of all first-level ResourceRefs that are not virtual
      def nonvirtual_rrefs(object)
        resourcerefs_for_properties(object.all_user_properties,
                                    object, virtual: 'exclude')
      end

      def virtual_rrefs(object)
        resourcerefs_for_properties(object.all_user_properties,
                                    object, virtual: 'only')
      end

      # Converts a path in the form a/b/c/d into ['a', 'b', 'c', 'd']
      def path2navigate(path)
        "[#{path.split('/').map { |x| "'#{x}'" }.join(', ')}]"
      end

      # Generates a method declaration with function name `name` and args `args`
      # Arguments may have nils and will be ignored.
      def method_decl(name, args)
        "def #{name}(#{args.compact.join(', ')}):"
      end

      # Generates a method call to function name `name` and args `args`
      # Arguments may have nils and will be ignored.
      def method_call(name, args)
        "#{name}(#{args.compact.join(', ')})"
      end

      # TODO(alexstephen): Standardize on one version and move to provider/core
      # https://github.com/GoogleCloudPlatform/magic-modules/issues/30
      def wrap_field(field, spaces)
        # field.scan goes from 0 -> avail_columns - 1
        # -1 to account for this
        avail_columns = DEFAULT_FORMAT_OPTIONS[:max_columns] - spaces - 1
        field.scan(/\S.{0,#{avail_columns}}\S(?=\s|$)|\S+/)
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

        compile_file(EXAMPLE_DEFAULTS, path) if File.file?(path)
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

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      # rubocop:disable Metrics/MethodLength
      # rubocop:disable Metrics/BlockLength
      def resourcerefs_for_properties(props, original_obj, **kwargs)
        rrefs = []
        props.each do |p|
          # We need to recurse on ResourceRefs to get all levels
          # We do not want to recurse on resourcerefs of type self to avoid
          # infinite loop.
          if p.is_a? Api::Type::ResourceRef
            # We want to avoid a circular reference
            # This reference may be the
            # next reference or have some number of refs in between it.
            next if p.resources.first.resource_ref == original_obj
            next if p.resources.first.resource_ref == p.resources.first.__resource
            if p.resources.first.resource_ref.virtual
              next if kwargs[:virtual] == 'exclude'
            else
              next if kwargs[:virtual] == 'only'
            end
            rrefs << p
            rrefs.concat(resourcerefs_for_properties(
                           p.resources[0].resource_ref.required_properties,
                           original_obj,
                           virtual: kwargs[:virtual]
            ))
          elsif p.is_a? Api::Type::NestedObject
            rrefs.concat(resourcerefs_for_properties(p.properties,
                                                          original_obj,
                                                          virtual: kwargs[:virtual]))
          elsif p.is_a? Api::Type::Array
            if p.item_type.is_a? Api::Type::NestedObject
              rrefs.concat(resourcerefs_for_properties(
                             p.item_type.properties,
                             original_obj,
                             virtual: kwargs[:virtual]
              ))
            elsif p.item_type.is_a? Api::Type::ResourceRef
              if p.item_type.resources.first.resource_ref.virtual
                next if kwargs[:virtual] == 'exclude'
              else
                next if kwargs[:virtual] == 'only'
              end
              rrefs << p.item_type
              rrefs.concat(resourcerefs_for_properties(
                             p.item_type.resources[0].resource_ref
                                                     .required_properties,
                             original_obj,
                             virtual: kwargs[:virtual]
              ))
            end
          end
        end
        rrefs.uniq
      end
      # rubocop:enable Metrics/AbcSize
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/PerceivedComplexity
      # rubocop:enable Metrics/BlockLength
    end
    # rubocop:enable Metrics/ClassLength
  end
end

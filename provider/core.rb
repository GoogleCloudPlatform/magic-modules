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

require 'compile/core'
require 'fileutils'
require 'google/extensions'
require 'google/logger'
require 'json'
require 'overrides/runner'
require 'provider/file_template'

module Provider
  # Basic functionality for code generator providers. Provides basic services,
  # such as compiling and including files, formatting data, etc.
  class Core
    include Compile::Core

    def initialize(config, api, start_time)
      @config = config
      @api = api

      # The compiler will error out if a file has been written in this compiler
      # run already. Instead of storing all the modified files in state we'll
      # use the time the file was modified.
      @start_time = start_time
      @py_format_enabled = check_pyformat
      @go_format_enabled = check_goformat
    end

    # This provides the ProductFileTemplate class with access to a provider.
    def provider_binding
      binding
    end

    def check_pyformat
      if system('python3 -m black --help > /dev/null')
        true
      else
        Google::LOGGER.warn 'Either python3 or black is not installed; python ' \
          'code will be poorly formatted and may not pass linter checks.'
        false
      end
    end

    def check_goformat
      if system('which gofmt > /dev/null') && system('which goimports > /dev/null')
        true
      else
        Google::LOGGER.warn 'Either gofmt or goimports is not installed; go ' \
          'code will be poorly formatted and will likely not compile.'
        false
      end
    end

    # Main entry point for the compiler. As this method is simply invoking other
    # generators, it is okay to ignore Rubocop warnings about method size and
    # complexity.
    #
    def generate(output_folder, types, version_name, product_path, dump_yaml)
      generate_objects(output_folder, types, version_name)
      copy_files(output_folder) \
        unless @config.files.nil? || @config.files.copy.nil?
      # Compilation has to be the last step, as some files (e.g.
      # CONTRIBUTING.md) may depend on the list of all files previously copied
      # or compiled.
      # common-compile.yaml is a special file that will get compiled by the last product
      # used in a single invocation of the compiled. It should not contain product-specific
      # information; instead, it should be run-specific such as the version to compile at.
      compile_files(output_folder, version_name) \
        unless @config.files.nil? || @config.files.compile.nil?

      generate_datasources(output_folder, types, version_name) \
        unless @config.datasources.nil?

      generate_operation(output_folder, types, version_name)

      # Write a file with the final version of the api, after overrides
      # have been applied.
      return unless dump_yaml

      raise 'Path to output the final yaml was not specified.' \
        if product_path.nil? || product_path == ''

      File.open("#{product_path}/final_api.yaml", 'w') do |file|
        file.write("# This is a generated file, its contents will be overwritten.\n")
        file.write(YAML.dump(@api))
      end
    end

    def generate_operation(output_folder, types, version_name); end

    def copy_files(output_folder)
      copy_file_list(output_folder, @config.files.copy)
    end

    # version_name is actually used because all of the variables in scope in this method
    # are made available within the templates by the compile call. This means that version_name
    # is exposed to the templating logic and version_name is used in other places in the same
    # way so it needs to be named consistently
    # rubocop:disable Lint/UnusedMethodArgument
    def copy_common_files(output_folder, version_name = 'ga')
      provider_name = self.class.name.split('::').last.downcase
      return unless File.exist?("provider/#{provider_name}/common~copy.yaml")

      Google::LOGGER.info "Copying common files for #{provider_name}"
      files = YAML.safe_load(compile("provider/#{provider_name}/common~copy.yaml"))
      copy_file_list(output_folder, files)
    end
    # rubocop:enable Lint/UnusedMethodArgument

    def copy_file_list(output_folder, files)
      files.map do |target, source|
        Thread.new do
          target_file = File.join(output_folder, target)
          target_dir = File.dirname(target_file)
          Google::LOGGER.debug "Copying #{source} => #{target}"
          FileUtils.mkpath target_dir unless Dir.exist?(target_dir)
          FileUtils.copy_entry source, target_file
        end
      end.map(&:join)
    end

    def compile_files(output_folder, version_name)
      compile_file_list(output_folder, @config.files.compile, version_name)
    end

    def compile_common_files(output_folder, version_name = nil)
      provider_name = self.class.name.split('::').last.downcase
      return unless File.exist?("provider/#{provider_name}/common~compile.yaml")

      Google::LOGGER.info "Compiling common files for #{provider_name}"
      files = YAML.safe_load(compile("provider/#{provider_name}/common~compile.yaml"))
      compile_file_list(output_folder, files, version_name)
    end

    def compile_file_list(output_folder, files, version = nil)
      files.map do |target, source|
        Thread.new do
          Google::LOGGER.debug "Compiling #{source} => #{target}"
          target_file = File.join(output_folder, target)
          ProductFileTemplate.new(
            output_folder,
            target,
            @api,
            version,
            build_env
          ).generate(source, target_file, self)
        end
      end.map(&:join)
    end

    def api_version_setup(version_name)
      version = @api.version_obj_or_closest(version_name)
      @api.set_properties_based_on_version(version)
      version
    end

    def generate_objects(output_folder, types, version_name)
      version = api_version_setup(version_name)
      (@api.objects || []).each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info "Excluding #{object.name} per user request"
        elsif types.empty? && object.exclude
          Google::LOGGER.info "Excluding #{object.name} per API catalog"
        elsif types.empty? && object.not_in_version?(version)
          Google::LOGGER.info "Excluding #{object.name} per API version"
        else
          Google::LOGGER.info "Generating #{object.name}"
          # exclude_if_not_in_version must be called in order to filter out
          # beta properties that are nested within GA resrouces
          object.exclude_if_not_in_version!(version)

          # Make object immutable.
          object.freeze
          object.all_user_properties.each(&:freeze)

          # version_name will differ from version.name if the resource is being
          # generated at its default version instead of the one that was passed
          # in to the compiler. Terraform needs to know which version was passed
          # in so it can name its output directories correctly.
          generate_object object, output_folder, version_name
        end
      end
    end

    def generate_object(object, output_folder, version_name)
      data = build_object_data(object, output_folder, version_name)
      unless object.exclude_resource
        Google::LOGGER.debug "Generating #{object.name} resource"
        generate_resource data.clone
        Google::LOGGER.debug "Generating #{object.name} tests"
        generate_resource_tests data.clone

        generate_resource_files data.clone
      end

      # if iam_policy is not defined or excluded, don't generate it
      return if object.iam_policy.nil? || object.iam_policy.exclude

      Google::LOGGER.debug "Generating #{object.name} IAM policy"
      generate_iam_policy data.clone
    end

    # Generate files at a per-resource basis.
    def generate_resource_files(data) end

    def generate_datasources(output_folder, types, version_name)
      # We need to apply overrides for datasources
      @api = Overrides::Runner.build(@api, @config.datasources,
                                     @config.resource_override,
                                     @config.property_override)
      @api.validate

      version = @api.version_obj_or_closest(version_name)
      @api.set_properties_based_on_version(version)
      @api.objects.each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per user request"
          )
        elsif types.empty? && object.exclude
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per API catalog"
          )
        elsif types.empty? && object.not_in_version?(version)
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per API version"
          )
        else
          generate_datasource object, output_folder, version_name
        end
      end
    end

    def generate_datasource(object, output_folder, version_name)
      data = build_object_data(object, output_folder, version_name)

      compile_datasource data.clone
    end

    def build_object_data(object, output_folder, version)
      ProductFileTemplate.file_for_resource(output_folder, object, version, @config, build_env)
    end

    def build_env
      {
        pyformat_enabled: @py_format_enabled,
        goformat_enabled: @go_format_enabled,
        start_time: @start_time
      }
    end

    # Filter the properties to keep only the ones requiring custom update
    # method and group them by update url & verb.
    def properties_by_custom_update(properties)
      update_props = properties.reject do |p|
        p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP
      end
      update_props.group_by do |p|
        { update_url: p.update_url, update_verb: p.update_verb }
      end
    end

    # Takes a update_url and returns the list of custom updatable properties
    # that can be updated at that URL. This allows flattened objects
    # to determine which parent property in the API should be updated with
    # the contents of the flattened object
    def custom_update_properties_by_url(properties, update_url)
      properties_by_custom_update(properties).select do |k, _|
        k[:update_url] == update_url
      end.first.last
      # .first is to grab the element from the select which returns a list
      # .last is because properties_by_custom_update returns a list of
      # [{update_url}, [properties,...]] and we only need the 2nd part
    end

    def update_url(resource, url_part)
      [resource.__product.base_url, update_uri(resource, url_part)].flatten.join
    end

    def update_uri(resource, url_part)
      return resource.self_link_uri if url_part.nil?

      url_part
    end

    def generate_iam_policy(data) end

    def compile_provider_files(output_folder, products, version) end

    # TODO(nelsonjr): Review all object interfaces and move to private methods
    # that should not be exposed outside the object hierarchy.
    private

    def generate_requires(properties, requires = [])
      requires.concat(properties.collect(&:requires))
    end

    def provider_name
      self.class.name.split('::').last.downcase
    end

    # Determines the copyright year. If the file already exists we'll attempt to
    # recognize the copyright year, and if it finds it will keep it.
    def effective_copyright_year(out_file)
      copyright_mask = /# Copyright (?<year>[0-9-]*) Google Inc./
      if File.exist?(out_file)
        first_line = File.read(out_file).split("\n")
                         .select { |l| copyright_mask.match(l) }
                         .first
        matcher = copyright_mask.match(first_line)
        return matcher[:year] unless matcher.nil?
      end
      Time.now.year
    end
  end
end

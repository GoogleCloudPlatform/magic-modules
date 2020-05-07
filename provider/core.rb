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

    def initialize(config, api, version_name, start_time)
      @config = config
      @api = api

      # @target_version_name is the version specified by MM for this generation
      # run. That's distinct from @version below, which is the best-fit version
      # supported by the product.
      # These values will often match, but if a product supports only GA while
      # MM is ran @ beta, @target_version_name will be at beta and @version will
      # be @ GA.
      # This matters for Terraform, where the primary folder for a provider
      # needs to match the provider name.
      @target_version_name = version_name

      @version = @api.version_obj_or_closest(version_name)
      @api.set_properties_based_on_version(@version)

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

    # Main entry point for generation.
    def generate(output_folder, types, product_path, dump_yaml)
      generate_objects(output_folder, types)
      copy_files(output_folder) \
        unless @config.files.nil? || @config.files.copy.nil?
      # Compilation has to be the last step, as some files (e.g.
      # CONTRIBUTING.md) may depend on the list of all files previously copied
      # or compiled.
      # common-compile.yaml is a special file that will get compiled by the last product
      # used in a single invocation of the compiled. It should not contain product-specific
      # information; instead, it should be run-specific such as the version to compile at.
      compile_product_files(output_folder) \
        unless @config.files.nil? || @config.files.compile.nil?

      generate_datasources(output_folder, types) \
        unless @config.datasources.nil?

      FileUtils.mkpath output_folder unless Dir.exist?(output_folder)
      $pwd = Dir.pwd
      Dir.chdir output_folder
      generate_operation(output_folder, types)
      Dir.chdir $pwd

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

    def generate_operation(output_folder, types); end

    def copy_files(output_folder)
      copy_file_list(output_folder, @config.files.copy)
    end

    def copy_common_files(output_folder, provider_name = nil)
      # version_name is actually used because all of the variables in scope in this method
      # are made available within the templates by the compile call.
      # TODO: remove version_name, use @target_version_name or pass it in expicitly
      # rubocop:disable Lint/UselessAssignment
      version_name = @target_version_name
      # rubocop:enable Lint/UselessAssignment
      provider_name ||= self.class.name.split('::').last.downcase
      return unless File.exist?("provider/#{provider_name}/common~copy.yaml")

      Google::LOGGER.info "Copying common files for #{provider_name}"
      files = YAML.safe_load(compile("provider/#{provider_name}/common~copy.yaml"))
      copy_file_list(output_folder, files)
    end

    def copy_file_list(output_folder, files)
      files.map do |target, source|
        Thread.new do
          target_file = File.join(output_folder, target)
          target_dir = File.dirname(target_file)
          Google::LOGGER.debug "Copying #{source} => #{target}"
          FileUtils.mkpath target_dir unless Dir.exist?(target_dir)

          # If we've modified a file since starting an MM run, it's a reasonable
          # assumption that it was this run that modified it.
          if File.exist?(target_file) && File.mtime(target_file) > @start_time
            raise "#{target_file} was already modified during this run. #{File.mtime(target_file)}"
          end

          FileUtils.copy_entry source, target_file
        end
      end.map(&:join)
    end

    # Compiles files specified within the product
    def compile_product_files(output_folder)
      file_template = ProductFileTemplate.new(
        output_folder,
        nil,
        @api,
        @target_version_name,
        build_env
      )
      compile_file_list(output_folder, @config.files.compile, file_template)
    end

    # Compiles files that are shared at the provider level
    def compile_common_files(
      output_folder,
      products,
      common_compile_file,
      override_path = nil
    )
      return unless File.exist?(common_compile_file)

      files = YAML.safe_load(compile(common_compile_file))
      return unless files

      file_template = ProviderFileTemplate.new(
        output_folder,
        @target_version_name,
        build_env,
        products,
        override_path
      )
      compile_file_list(output_folder, files, file_template)
    end

    def compile_file_list(output_folder, files, file_template)
      FileUtils.mkpath output_folder unless Dir.exist?(output_folder)
      $pwd = Dir.pwd
      Dir.chdir output_folder
      files.map do |target, source|
        Thread.new do
          Google::LOGGER.debug "Compiling #{source} => #{target}"
          file_template.generate(source, target, self)
        end
      end.map(&:join)
      Dir.chdir $pwd
    end

    def generate_objects(output_folder, types)
      (@api.objects || []).each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info "Excluding #{object.name} per user request"
        elsif types.empty? && object.exclude
          Google::LOGGER.info "Excluding #{object.name} per API catalog"
        elsif types.empty? && object.not_in_version?(@version)
          Google::LOGGER.info "Excluding #{object.name} per API version"
        else
          Google::LOGGER.info "Generating #{object.name}"
          # exclude_if_not_in_version must be called in order to filter out
          # beta properties that are nested within GA resources
          object.exclude_if_not_in_version!(@version)

          # Make object immutable.
          object.freeze
          object.all_user_properties.each(&:freeze)

          generate_object object, output_folder, @target_version_name
        end
      end
    end

    def generate_object(object, output_folder, version_name)
      data = build_object_data(object, output_folder, version_name)
      $pwd = Dir.pwd
      unless object.exclude_resource
        FileUtils.mkpath output_folder unless Dir.exist?(output_folder)
        Dir.chdir output_folder
        Google::LOGGER.debug "Generating #{object.name} resource"
        generate_resource data.clone
        Google::LOGGER.debug "Generating #{object.name} tests"
        generate_resource_tests data.clone
        generate_resource_sweepers data.clone
        generate_resource_files data.clone
        Dir.chdir $pwd
      end

      # if iam_policy is not defined or excluded, don't generate it
      return if object.iam_policy.nil? || object.iam_policy.exclude

      FileUtils.mkpath output_folder unless Dir.exist?(output_folder)
      Dir.chdir output_folder
      Google::LOGGER.debug "Generating #{object.name} IAM policy"
      generate_iam_policy data.clone
      Dir.chdir $pwd
    end

    # Generate files at a per-resource basis.
    def generate_resource_files(data) end

    def generate_datasources(output_folder, types)
      # We need to apply overrides for datasources
      @api = Overrides::Runner.build(@api, @config.datasources,
                                     @config.resource_override,
                                     @config.property_override)
      @api.validate

      @api.set_properties_based_on_version(@version)
      @api.objects.each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per user request"
          )
        elsif types.empty? && object.exclude
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per API catalog"
          )
        elsif types.empty? && object.not_in_version?(@version)
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per API version"
          )
        else
          generate_datasource object, output_folder
        end
      end
    end

    def generate_datasource(object, output_folder)
      data = build_object_data(object, output_folder, @target_version_name)

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
    def properties_by_custom_update(properties, behavior = :new)
      update_props = properties.reject do |p|
        p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP
      end

      # TODO(rambleraptor): Add support to Ansible for one-at-a-time updates.
      if behavior == :old
        update_props.group_by do |p|
          { update_url: p.update_url, update_verb: p.update_verb, fingerprint: p.fingerprint_name }
        end
      else
        update_props.group_by do |p|
          {
            update_url: p.update_url,
            update_verb: p.update_verb,
            update_id: p.update_id,
            fingerprint_name: p.fingerprint_name
          }
        end
      end
    end

    # Takes a update_url and returns the list of custom updatable properties
    # that can be updated at that URL. This allows flattened objects
    # to determine which parent property in the API should be updated with
    # the contents of the flattened object
    def custom_update_properties_by_key(properties, key)
      properties_by_custom_update(properties).select do |k, _|
        k[:update_url] == key[:update_url] &&
          k[:update_id] == key[:update_id] &&
          k[:fingerprint_name] == key[:fingerprint_name]
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

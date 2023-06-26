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

    TERRAFORM_PROVIDER_GA = 'github.com/hashicorp/terraform-provider-google'.freeze
    TERRAFORM_PROVIDER_BETA = 'github.com/hashicorp/terraform-provider-google-beta'.freeze
    TERRAFORM_PROVIDER_PRIVATE = 'internal/terraform-next'.freeze
    RESOURCE_DIRECTORY_GA = 'google'.freeze
    RESOURCE_DIRECTORY_BETA = 'google-beta'.freeze
    RESOURCE_DIRECTORY_PRIVATE = 'google-private'.freeze

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
      @go_format_enabled = check_goformat
    end

    # This provides the ProductFileTemplate class with access to a provider.
    def provider_binding
      binding
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
    def generate(output_folder, types, product_path, dump_yaml, generate_code, generate_docs)
      generate_objects(output_folder, types, generate_code, generate_docs)

      # Compilation has to be the last step, as some files (e.g.
      # CONTRIBUTING.md) may depend on the list of all files previously copied
      # or compiled.
      # common-compile.yaml is a special file that will get compiled by the last product
      # used in a single invocation of the compiled. It should not contain product-specific
      # information; instead, it should be run-specific such as the version to compile at.
      compile_product_files(output_folder) \
        unless @config.files.nil? || @config.files.compile.nil?

      FileUtils.mkpath output_folder
      pwd = Dir.pwd
      if generate_code
        Dir.chdir output_folder

        generate_operation(pwd, output_folder, types)
        Dir.chdir pwd
      end

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

    def generate_operation(pwd, output_folder, types); end

    # generate_code and generate_docs are actually used because all of the variables
    # in scope in this method are made available within the templates by the compile call.
    # rubocop:disable Lint/UnusedMethodArgument
    def copy_common_files(output_folder, generate_code, generate_docs, provider_name = nil)
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
    # rubocop:enable Lint/UnusedMethodArgument

    def copy_file_list(output_folder, files)
      files.map do |target, source|
        Thread.new do
          target_file = File.join(output_folder, target)
          target_dir = File.dirname(target_file)
          Google::LOGGER.debug "Copying #{source} => #{target}"
          FileUtils.mkpath target_dir

          # If we've modified a file since starting an MM run, it's a reasonable
          # assumption that it was this run that modified it.
          if File.exist?(target_file) && File.mtime(target_file) > @start_time
            raise "#{target_file} was already modified during this run. #{File.mtime(target_file)}"
          end

          FileUtils.copy_entry source, target_file

          add_hashicorp_copyright_header(output_folder, target) if File.extname(target) == '.go'
          if File.extname(target) == '.go' || File.extname(target) == '.mod'
            replace_import_path(output_folder, target)
          end
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

    def compile_file_list(output_folder, files, file_template, pwd = Dir.pwd)
      FileUtils.mkpath output_folder
      Dir.chdir output_folder
      files.map do |target, source|
        Thread.new do
          Google::LOGGER.debug "Compiling #{source} => #{target}"
          file_template.generate(pwd, source, target, self)

          add_hashicorp_copyright_header(output_folder, target)
          replace_import_path(output_folder, target)
        end
      end.map(&:join)
      Dir.chdir pwd
    end

    def add_hashicorp_copyright_header(output_folder, target)
      unless expected_output_folder?(output_folder)
        Google::LOGGER.info "Unexpected output folder (#{output_folder}) detected " \
                            'when deciding to add HashiCorp copyright headers. ' \
                            'Watch out for unexpected changes to copied files'
      end
      # only add copyright headers when generating TPG and TPGB
      return unless output_folder.end_with?('terraform-provider-google') ||
                    output_folder.end_with?('terraform-provider-google-beta')

      # Prevent adding copyright header to files with paths or names matching the strings below
      # NOTE: these entries need to match the content of the .copywrite.hcl file originally
      #       created in https://github.com/GoogleCloudPlatform/magic-modules/pull/7336
      #       The test-fixtures folder is not included here as it's copied as a whole,
      #       not file by file (see common~copy.yaml)
      ignored_folders = [
        '.release/',
        '.changelog/',
        'examples/',
        'scripts/',
        'META.d/'
      ]
      ignored_files = [
        'go.mod',
        '.goreleaser.yml',
        '.golangci.yml',
        'terraform-registry-manifest.json'
      ]
      should_add_header = true
      ignored_folders.each do |folder|
        # folder will be path leading to file
        next unless target.start_with? folder

        Google::LOGGER.debug 'Not adding HashiCorp copyright headers in ' \
                             "ignored folder #{folder} : #{target}"
        should_add_header = false
      end
      return unless should_add_header

      ignored_files.each do |file|
        # file will be the filename and extension, with no preceding path
        next unless target.end_with? file

        Google::LOGGER.debug 'Not adding HashiCorp copyright headers to ' \
                             "ignored file #{file} : #{target}"
        should_add_header = false
      end
      return unless should_add_header

      Google::LOGGER.debug "Adding HashiCorp copyright header to : #{target}"
      data = File.read("#{output_folder}/#{target}")

      copyright_header = ['Copyright (c) HashiCorp, Inc.', 'SPDX-License-Identifier: MPL-2.0']
      lang = language_from_filename(target)

      # Some file types we don't want to add headers to
      # e.g. .sh where headers are functional
      # Also, this guards against new filetypes being added and triggering build errors
      return unless lang != :unsupported

      # File is not ignored and is appropriate file type to add header to
      header = comment_block(copyright_header, lang)
      File.write("#{output_folder}/#{target}", header)

      File.write("#{output_folder}/#{target}", data, mode: 'a') # append mode
    end

    def expected_output_folder?(output_folder)
      expected_folders = %w[
        terraform-provider-google
        terraform-provider-google-beta
        terraform-next
        terraform-google-conversion
        tfplan2cai
      ]
      folder_name = output_folder.split('/')[-1] # Possible issue with Windows OS
      is_expected = false
      expected_folders.each do |folder|
        next unless folder_name == folder

        is_expected = true
        break
      end
      is_expected
    end

    def replace_import_path(output_folder, target)
      return unless @target_version_name != 'ga'

      # Replace the import pathes in utility files
      case @target_version_name
      when 'beta'
        tpg = TERRAFORM_PROVIDER_BETA
        dir = RESOURCE_DIRECTORY_BETA
      else
        tpg = TERRAFORM_PROVIDER_PRIVATE
        dir = RESOURCE_DIRECTORY_PRIVATE
      end

      data = File.read("#{output_folder}/#{target}")
      data = data.gsub(
        "#{TERRAFORM_PROVIDER_GA}/#{RESOURCE_DIRECTORY_GA}",
        "#{tpg}/#{dir}"
      )
      data = data.gsub(
        "#{TERRAFORM_PROVIDER_GA}/version",
        "#{tpg}/version"
      )

      Google::LOGGER.info "replace_import_path target #{output_folder}/#{target}"
      data = data.gsub(
        "module #{TERRAFORM_PROVIDER_GA}",
        "module #{tpg}"
      )
      File.write("#{output_folder}/#{target}", data)
    end

    def import_path
      case @target_version_name
      when 'ga'
        "#{TERRAFORM_PROVIDER_GA}/#{RESOURCE_DIRECTORY_GA}"
      when 'beta'
        "#{TERRAFORM_PROVIDER_BETA}/#{RESOURCE_DIRECTORY_BETA}"
      else
        "#{TERRAFORM_PROVIDER_PRIVATE}/#{RESOURCE_DIRECTORY_PRIVATE}"
      end
    end

    def generate_objects(output_folder, types, generate_code, generate_docs)
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

          generate_object object, output_folder, @target_version_name, generate_code, generate_docs
        end
      end
    end

    def generate_object(object, output_folder, version_name, generate_code, generate_docs)
      pwd = Dir.pwd
      data = build_object_data(pwd, object, output_folder, version_name)
      unless object.exclude_resource
        FileUtils.mkpath output_folder
        Dir.chdir output_folder
        Google::LOGGER.debug "Generating #{object.name} resource"
        generate_resource(pwd, data.clone, generate_code, generate_docs)
        if generate_code
          Google::LOGGER.debug "Generating #{object.name} tests"
          generate_resource_tests(pwd, data.clone)
          generate_resource_sweepers(pwd, data.clone)
          generate_resource_files(pwd, data.clone)
        end
        Dir.chdir pwd
      end

      # if iam_policy is not defined or excluded, don't generate it
      return if object.iam_policy.nil? || object.iam_policy.exclude

      FileUtils.mkpath output_folder
      Dir.chdir output_folder
      Google::LOGGER.debug "Generating #{object.name} IAM policy"
      generate_iam_policy(pwd, data.clone, generate_code, generate_docs)
      Dir.chdir pwd
    end

    # Generate files at a per-resource basis.
    def generate_resource_files(pwd, data) end

    def build_object_data(_pwd, object, output_folder, version)
      ProductFileTemplate.file_for_resource(output_folder, object, version, @config, build_env)
    end

    def build_env
      {
        goformat_enabled: @go_format_enabled,
        start_time: @start_time
      }
    end

    # used to determine and separate objects that have update methods
    # that target individual fields
    def field_specific_update_methods(properties)
      properties_by_custom_update(properties).length.positive?
    end

    # Filter the properties to keep only the ones requiring custom update
    # method and group them by update url & verb.
    def properties_by_custom_update(properties)
      update_props = properties.reject do |p|
        p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP
      end

      update_props.group_by do |p|
        {
          update_url: p.update_url,
          update_verb: p.update_verb,
          update_id: p.update_id,
          fingerprint_name: p.fingerprint_name
        }
      end
    end

    # Filter the properties to keep only the ones don't have custom update
    # method and group them by update url & verb.
    def properties_without_custom_update(properties)
      properties.select do |p|
        p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP
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

    def generate_iam_policy(pwd, data, generate_code, generate_docs) end

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

    # Adapted from the method used in templating
    # See: mmv1/compile/core.rb
    def comment_block(text, lang)
      case lang
      when :ruby, :python, :yaml, :git, :gemfile
        header = text.map { |t| t&.empty? ? '#' : "# #{t}" }
      when :go
        header = text.map { |t| t&.empty? ? '//' : "// #{t}" }
      else
        raise "Unknown language for comment: #{lang}"
      end

      header_string = header.join("\n")
      "#{header_string}\n" # add trailing newline to returned value
    end

    def language_from_filename(filename)
      extension = filename.split('.')[-1]
      case extension
      when 'go'
        :go
      when 'rb'
        :ruby
      when 'yaml', 'yml'
        :yaml
      else
        :unsupported
      end
    end
  end
end

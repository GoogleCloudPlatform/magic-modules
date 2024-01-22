# Copyright 2020 Google Inc.
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
require 'provider/file_template'
require 'provider/terraform/async'
require 'provider/terraform/import'
require 'provider/terraform/custom_code'
require 'provider/terraform/docs'
require 'provider/terraform/examples'
require 'provider/terraform/sub_template'
require 'google/golang_utils'

module Provider
  # Code generator for Terraform Resources that manage Google Cloud Platform
  # resources.
  class Terraform
    include Compile::Core
    include Provider::Terraform::Import
    include Provider::Terraform::SubTemplate
    include Google::GolangUtils

    attr_accessor :resource_count
    attr_accessor :iam_resource_count
    attr_accessor :resources_for_version

    TERRAFORM_PROVIDER_GA = 'github.com/hashicorp/terraform-provider-google'.freeze
    TERRAFORM_PROVIDER_BETA = 'github.com/hashicorp/terraform-provider-google-beta'.freeze
    TERRAFORM_PROVIDER_PRIVATE = 'internal/terraform-next'.freeze
    RESOURCE_DIRECTORY_GA = 'google'.freeze
    RESOURCE_DIRECTORY_BETA = 'google-beta'.freeze
    RESOURCE_DIRECTORY_PRIVATE = 'google-private'.freeze

    def initialize(api, version_name, start_time)
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

      @resource_count = 0
      @iam_resource_count = 0
      @resources_for_version = []
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

    # Compiles files that are shared at the provider level
    def compile_common_files(
      output_folder,
      products,
      common_compile_file,
      override_path = nil
    )
      return unless File.exist?(common_compile_file)

      generate_resources_for_version(products, @target_version_name)

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
      data = File.read("#{output_folder}/#{target}")

      if data.include? "#{TERRAFORM_PROVIDER_BETA}/#{RESOURCE_DIRECTORY_BETA}"
        raise 'Importing a package from module ' \
              "#{TERRAFORM_PROVIDER_BETA}/#{RESOURCE_DIRECTORY_BETA} " \
              "is not allowed in file #{target.split('/').last}. " \
              'Please import a package from module ' \
              "#{TERRAFORM_PROVIDER_GA}/#{RESOURCE_DIRECTORY_GA}."
      end

      return if @target_version_name == 'ga'

      # Replace the import pathes in utility files
      case @target_version_name
      when 'beta'
        tpg = TERRAFORM_PROVIDER_BETA
        dir = RESOURCE_DIRECTORY_BETA
      else
        tpg = TERRAFORM_PROVIDER_PRIVATE
        dir = RESOURCE_DIRECTORY_PRIVATE
      end

      data = data.gsub(
        "#{TERRAFORM_PROVIDER_GA}/#{RESOURCE_DIRECTORY_GA}",
        "#{tpg}/#{dir}"
      )
      data = data.gsub(
        "#{TERRAFORM_PROVIDER_GA}/version",
        "#{tpg}/version"
      )

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

    # Gets the list of services dependent on the version ga, beta, and private
    # If there are some resources of a servcie is in GA,
    # then this service is in GA. Otherwise, the service is in BETA
    def get_mmv1_services_in_version(products, version)
      services = []
      products.map do |product|
        product_definition = product[:definitions]
        if version == 'ga'
          some_resource_in_ga = false
          product_definition.objects.each do |object|
            break if some_resource_in_ga

            if !object.exclude &&
               !object.not_in_version?(product_definition.version_obj_or_closest(version))
              some_resource_in_ga = true
            end
          end

          services << product[:definitions].name.downcase if some_resource_in_ga
        else
          services << product[:definitions].name.downcase
        end
      end
      services
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
        generate_object_modified object, output_folder, @target_version_name
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
        end
        Dir.chdir pwd
      end

    def generate_object_modified(object, output_folder, version_name)
      pwd = Dir.pwd
      data = build_object_data(pwd, object, output_folder, version_name)
      FileUtils.mkpath output_folder
      Dir.chdir output_folder
      Google::LOGGER.debug "Generating #{object.name} rewrite yaml"
      generate_newyaml(pwd, data.clone)
      Dir.chdir pwd
      end

    def generate_newyaml(pwd, data)
      # @api.api_name is the service folder name
      product_name = @api.api_name
      target_folder = File.join(folder_name(data.version), 'services', product_name)
      FileUtils.mkpath target_folder
      data.generate(pwd,
                    '/templates/terraform/yaml_conversion.erb',
                    "#{target_folder}/#{full_resource_name(data)}.yaml",
                    self)
    end


      # if iam_policy is not defined or excluded, don't generate it
      return if object.iam_policy.nil? || object.iam_policy.exclude

      FileUtils.mkpath output_folder
      Dir.chdir output_folder
      Google::LOGGER.debug "Generating #{object.name} IAM policy"
      generate_iam_policy(pwd, data.clone, generate_code, generate_docs)
      Dir.chdir pwd
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
        p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP ||
          p.is_a?(Api::Type::KeyValueTerraformLabels) ||
          p.is_a?(Api::Type::KeyValueLabels) # effective_labels is used for update
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

    def generating_hashicorp_repo?
      # The default Provider is used to generate TPG and TPGB in HashiCorp-owned repos.
      # The compiler deviates from the default behaviour with a -f flag to produce
      # non-HashiCorp downstreams.
      true
    end

    # ProductFileTemplate with Terraform specific fields
    class TerraformProductFileTemplate < Provider::ProductFileTemplate
      # The async object used for making operations.
      # We assume that all resources share the same async properties.
      attr_accessor :async

      # When generating OiCS examples, we attach the example we're
      # generating to the data object.
      attr_accessor :example

      attr_accessor :resource_name
    end

    # Sorts properties in the order they should appear in the TF schema:
    # Required, Optional, Computed
    def order_properties(properties)
      properties.select(&:required).sort_by(&:name) +
        properties.reject(&:required).reject(&:output).sort_by(&:name) +
        properties.select(&:output).sort_by(&:name)
    end

    def tf_type(property)
      tf_types[property.class]
    end

    # "Namespace" - prefix with product and resource - a property with
    # information from the "object" variable
    def namespace_property_from_object(property, object)
      name = property.name.camelize
      until property.parent.nil?
        property = property.parent
        name = property.name.camelize + name
      end

      "#{property.__resource.__product.api_name.camelize(:lower)}#{object.name}#{name}"
    end

    # Converts between the Magic Modules type of an object and its type in the
    # TF schema
    def tf_types
      {
        Api::Type::Boolean => 'schema.TypeBool',
        Api::Type::Double => 'schema.TypeFloat',
        Api::Type::Integer => 'schema.TypeInt',
        Api::Type::String => 'schema.TypeString',
        # Anonymous string property used in array of strings.
        'Api::Type::String' => 'schema.TypeString',
        Api::Type::Time => 'schema.TypeString',
        Api::Type::Enum => 'schema.TypeString',
        Api::Type::ResourceRef => 'schema.TypeString',
        Api::Type::NestedObject => 'schema.TypeList',
        Api::Type::Array => 'schema.TypeList',
        Api::Type::KeyValuePairs => 'schema.TypeMap',
        Api::Type::KeyValueLabels => 'schema.TypeMap',
        Api::Type::KeyValueTerraformLabels => 'schema.TypeMap',
        Api::Type::KeyValueEffectiveLabels => 'schema.TypeMap',
        Api::Type::KeyValueAnnotations => 'schema.TypeMap',
        Api::Type::Map => 'schema.TypeSet',
        Api::Type::Fingerprint => 'schema.TypeString'
      }
    end

    def updatable?(resource, properties)
      !resource.immutable || !properties.reject { |p| p.update_url.nil? }.empty?
    end

    def force_new?(property, resource)
      (
        (!property.output || property.is_a?(Api::Type::KeyValueEffectiveLabels)) &&
        (property.immutable ||
          (resource.immutable && property.update_url.nil? && property.immutable.nil? &&
            (property.parent.nil? ||
              (force_new?(property.parent, resource) &&
               !(property.parent.flatten_object && property.is_a?(Api::Type::KeyValueLabels))
              )
            )
          )
        )
      ) ||
        (property.is_a?(Api::Type::KeyValueTerraformLabels) &&
          !updatable?(resource, resource.all_user_properties) && !resource.root_labels?
        )
    end

    # Returns tuples of (fieldName, list of update masks) for
    #  top-level updatable fields. Schema path refers to a given Terraform
    # field name (e.g. d.GetChange('fieldName)')
    def get_property_update_masks_groups(properties, mask_prefix: '')
      mask_groups = []
      properties.each do |prop|
        if prop.flatten_object
          mask_groups += get_property_update_masks_groups(
            prop.properties, mask_prefix: "#{prop.api_name}."
          )
        elsif prop.update_mask_fields
          mask_groups << [prop.name.underscore, prop.update_mask_fields]
        else
          mask_groups << [prop.name.underscore, [mask_prefix + prop.api_name]]
        end
      end
      mask_groups
    end

    # Returns an updated path for a given Terraform field path (e.g.
    # 'a_field', 'parent_field.0.child_name'). Returns nil if the property
    # is not included in the resource's properties and removes keys that have
    # been flattened
    # FYI: Fields that have been renamed should use the new name, however, flattened
    # fields still need to be included, ie:
    # flattenedField > newParent > renameMe should be passed to this function as
    # flattened_field.0.new_parent.0.im_renamed
    # TODO(emilymye): Change format of input for
    # exactly_one_of/at_least_one_of/etc to use camelcase, MM properities and
    # convert to snake in this method
    def get_property_schema_path(schema_path, resource)
      nested_props = resource.properties
      prop = nil
      path_tkns = schema_path.split('.0.').map do |pname|
        camel_pname = pname.camelize(:lower)
        prop = nested_props.find { |p| p.name == camel_pname }
        # if we couldn't find it, see if it was renamed at the top level
        prop = nested_props.find { |p| p.name == schema_path } if prop.nil?
        return nil if prop.nil?

        nested_props = prop.nested_properties || []
        prop.flatten_object ? nil : pname.underscore
      end
      if path_tkns.empty? || path_tkns[-1].nil?
        nil
      else
        path_tkns.compact.join('.0.')
      end
    end

    # Transforms a format string with field markers to a regex string with
    # capture groups.
    #
    # For instance,
    #   projects/{{project}}/global/networks/{{name}}
    # is transformed to
    #   projects/(?P<project>[^/]+)/global/networks/(?P<name>[^/]+)
    #
    # Values marked with % are URL-encoded, and will match any number of /'s.
    #
    # Note: ?P indicates a Python-compatible named capture group. Named groups
    # aren't common in JS-based regex flavours, but are in Perl-based ones
    def format2regex(format)
      format
        .gsub(/\{\{%([[:word:]]+)\}\}/, '(?P<\1>.+)')
        .gsub(/\{\{([[:word:]]+)\}\}/, '(?P<\1>[^/]+)')
    end

    # Capitalize the first letter of a property name.
    # E.g. "creationTimestamp" becomes "CreationTimestamp".
    def titlelize_property(property)
      property.name.camelize(:upper)
    end

    # Generates the list of resources, and gets the count of resources and iam resources
    # dependent on the version ga, beta or private.
    # The resource object has the format
    # {
    #    terraform_name:
    #    resource_name:
    #    iam_class_name:
    # }
    # The variable resources_for_version is used to generate resources in file
    # mmv1/third_party/terraform/provider/provider_mmv1_resources.go.erb
    def generate_resources_for_version(products, version)
      products.each do |product|
        product_definition = product[:definitions]
        service = product_definition.name.downcase
        product_definition.objects.each do |object|
          if object.exclude ||
             object.not_in_version?(product_definition.version_obj_or_closest(version))
            next
          end

          @resource_count += 1 unless object&.exclude_resource

          tf_product = (object.__product.legacy_name || product_definition.name).underscore
          terraform_name = object.legacy_name || "google_#{tf_product}_#{object.name.underscore}"

          unless object&.exclude_resource
            resource_name = "#{service}.Resource#{product_definition.name}#{object.name}"
          end

          iam_policy = object&.iam_policy

          @iam_resource_count += 3 unless iam_policy.nil? || iam_policy.exclude

          unless iam_policy.nil? || iam_policy.exclude ||
                 (iam_policy.min_version && iam_policy.min_version < version)
            iam_class_name = "#{service}.#{product_definition.name}#{object.name}"
          end

          @resources_for_version << { terraform_name:, resource_name:, iam_class_name: }
        end
      end

      @resources_for_version = @resources_for_version.compact
    end

    # TODO(nelsonjr): Review all object interfaces and move to private methods
    # that should not be exposed outside the object hierarchy.
    private

    def provider_name
      self.class.name.split('::').last.downcase
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

    # Finds the folder name for a given version of the terraform provider
    def folder_name(version)
      version == 'ga' ? 'google' : "google-#{version}"
    end

    # This function uses the resource.erb template to create one file
    # per resource. The resource.erb template forms the basis of a single
    # GCP Resource on Terraform.
    def generate_resource(pwd, data, generate_code, generate_docs)
      if generate_code
        # @api.api_name is the service folder name
        product_name = @api.api_name
        target_folder = File.join(folder_name(data.version), 'services', product_name)
        FileUtils.mkpath target_folder
        data.generate(pwd,
                      '/templates/terraform/resource.erb',
                      "#{target_folder}/resource_#{full_resource_name(data)}.go",
                      self)
      end

      return unless generate_docs

      generate_documentation(pwd, data)
    end

    def generate_documentation(pwd, data)
      target_folder = data.output_folder
      target_folder = File.join(target_folder, 'website', 'docs', 'r')
      FileUtils.mkpath target_folder
      filepath = File.join(target_folder, "#{full_resource_name(data)}.html.markdown")
      data.generate(pwd, 'templates/terraform/resource.html.markdown.erb', filepath, self)
    end

    def generate_resource_tests(pwd, data)
      return if data.object.examples
                    .reject(&:skip_test)
                    .reject do |e|
                  @api.version_obj_or_closest(data.version) \
                < @api.version_obj_or_closest(e.min_version)
                end
                    .empty?

      product_name = @api.api_name
      target_folder = File.join(folder_name(data.version), 'services', product_name)
      FileUtils.mkpath folder_name(data.version)
      data.generate(
        pwd,
        'templates/terraform/examples/base_configs/test_file.go.erb',
        "#{target_folder}/resource_#{full_resource_name(data)}_generated_test.go",
        self
      )
    end

    def generate_resource_sweepers(pwd, data)
      return if data.object.skip_sweeper ||
                data.object.custom_code.custom_delete ||
                data.object.custom_code.pre_delete ||
                data.object.custom_code.post_delete ||
                data.object.skip_delete

      product_name = @api.api_name
      target_folder = File.join(folder_name(data.version), 'services', product_name)
      file_name =
        "#{target_folder}/resource_#{full_resource_name(data)}_sweeper.go"
      FileUtils.mkpath folder_name(data.version)
      data.generate(pwd,
                    'templates/terraform/sweeper_file.go.erb',
                    file_name,
                    self)
    end

    def generate_operation(pwd, output_folder, _types)
      return if @api.objects.select(&:autogen_async).empty?

      product_name = @api.api_name
      product_name_underscore = @api.name.underscore
      data = build_object_data(pwd, @api.objects.first, output_folder, @target_version_name)

      data.object = @api.objects.select(&:autogen_async).first

      data.async = data.object.async
      target_folder = File.join(folder_name(data.version), 'services', product_name)
      FileUtils.mkpath target_folder
      data.generate(pwd,
                    'templates/terraform/operation.go.erb',
                    "#{target_folder}/#{product_name_underscore}_operation.go",
                    self)
    end

    # Generate the IAM policy for this object. This is used to query and test
    # IAM policies separately from the resource itself
    def generate_iam_policy(pwd, data, generate_code, generate_docs)
      if generate_code \
        && (!data.object.iam_policy.min_version \
        || data.object.iam_policy.min_version >= data.version)
        product_name = @api.api_name
        target_folder = File.join(folder_name(data.version), 'services', product_name)
        FileUtils.mkpath target_folder
        data.generate(pwd,
                      'templates/terraform/iam_policy.go.erb',
                      "#{target_folder}/iam_#{full_resource_name(data)}.go",
                      self)

        # Only generate test if testable examples exist.
        unless data.object.examples.reject(&:skip_test).empty?
          data.generate(
            pwd,
            'templates/terraform/examples/base_configs/iam_test_file.go.erb',
            "#{target_folder}/iam_#{full_resource_name(data)}_generated_test.go",
            self
          )
        end
      end

      return unless generate_docs

      generate_iam_documentation(pwd, data)
    end

    def generate_iam_documentation(pwd, data)
      target_folder = data.output_folder
      resource_doc_folder = File.join(target_folder, 'website', 'docs', 'r')
      datasource_doc_folder = File.join(target_folder, 'website', 'docs', 'd')
      FileUtils.mkpath resource_doc_folder
      filepath =
        File.join(resource_doc_folder, "#{full_resource_name(data)}_iam.html.markdown")

      data.generate(pwd, 'templates/terraform/resource_iam.html.markdown.erb', filepath, self)
      FileUtils.mkpath datasource_doc_folder
      filepath =
        File.join(datasource_doc_folder, "#{full_resource_name(data)}_iam_policy.html.markdown")

      data.generate(pwd, 'templates/terraform/datasource_iam.html.markdown.erb', filepath, self)
    end

    def build_object_data(_pwd, object, output_folder, version)
      TerraformProductFileTemplate.file_for_resource(
        output_folder,
        object,
        version,
        build_env
      )
    end

    def extract_identifiers(url)
      url.scan(/\{\{%?(\w+)\}\}/).flatten
    end

    # Returns the id format of an object, or self_link_uri if none is explicitly defined
    # We prefer the long name of a resource as the id so that users can reference
    # resources in a standard way, and most APIs accept short name, long name or self_link
    def id_format(object)
      object.id_format || object.self_link_uri
    end

    def full_resource_name(data)
      if data.object.legacy_name
        data.object.legacy_name.sub(/^google_/, '')
      else
        name = data.object.filename_override || data.object.name.underscore
        product_name = data.product.legacy_name || data.product.name.underscore
        "#{product_name}_#{name}"
      end
    end

    # Returns the extension for DCL packages for the given version. This is needed
    # as the DCL uses "alpha" for preview resources, while we use "private"
    def dcl_version(version)
      return '' if version == 'ga'
      return '/beta' if version == 'beta'
      return '/alpha' if version == 'private'
    end
  end
end

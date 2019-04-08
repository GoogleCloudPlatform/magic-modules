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
require 'google/hash_utils'
require 'pathname'
require 'json'
require 'overrides/runner'

module Provider
  DEFAULT_FORMAT_OPTIONS = {
    indent: 0,
    start_indent: 0,
    max_columns: 100,
    quiet: false
  }.freeze

  # Responsible for generating a file
  # with a given set of parameters.
  class FileTemplate
    include Compile::Core
    # The name of the resource
    attr_accessor :name
    # The resource itself.
    attr_accessor :object
    # The entire API object.
    attr_accessor :product
    # The API version
    attr_accessor :version
    # The root folder we're outputting to.
    attr_accessor :output_folder
    # The namespace of the product.
    attr_accessor :product_ns
    # The provider-specific configuration.
    attr_accessor :config
    # The provider specific high-level configuration fields.
    attr_accessor :manifest
    # Information about the local environment
    # (which formatters are enabled, start-time)
    attr_accessor :env

    # Ansible stuff.
    # The Ansible example object.
    attr_accessor :example

    # InSpec stuff.
    # Is this a plural resource?
    attr_accessor :plural
    # Should we generate documentation?
    attr_accessor :doc_generation
    # The file name of the attribute
    attr_accessor :attribute_file_name
    attr_accessor :privileged
    attr_accessor :property

    # Terraform stuff.
    # The async object used for making operations.
    # We assume that all resources share the same async properties.
    attr_accessor :async
    attr_accessor :resource_name

    def initialize(options)
      options.each { |k, v| method("#{k}=").call(v) }

      @env = {
        pyformat_enabled: options.dig(:env, :pyformat_enabled),
        goformat_enabled: options.dig(:env, :goformat_enabled),
        start_time: options.dig(:env, :start_time)
      }
    end

    # Given the data object for a file, write that file and verify that it
    # passes these conditions:
    #
    # - The file has not been generated already this run
    # - The file has an autogen exception or an autogen notice defined
    #
    # Once the file's contents are written, set the proper [chmod] mode and
    # format the file with a language-appropriate formatter.
    def generate(template, path, provider)
      folder = File.dirname(path)
      FileUtils.mkpath folder unless Dir.exist?(folder)

      # If we've modified a file since starting an MM run, it's a reasonable
      # assumption that it was this run that modified it.
      if File.exist?(path) && File.mtime(path) > @env[:start_time]
        raise "#{path} was already modified during this run"
      end

      # You're looking at some magic here!
      # This is how variables are made available in templates; we iterate
      # through each key:value pair in this object, and we set them
      # in the scope of the provider.
      #
      # The templates get access to everything in the provider +
      # all of the variables in this object.
      ctx = provider.provider_binding
      instance_variables.each do |name|
        ctx.local_variable_set(name[1..-1], instance_variable_get(name))
      end

      # This variable is used in ansible/resource.erb
      ctx.local_variable_set('file_relative', relative_path(path, @output_folder).to_s)

      Google::LOGGER.debug "Generating #{@name}"
      File.open(path, 'w') { |f| f.puts compile_file(ctx, template) }

      # Files are often generated in parallel.
      # We can use thread-local variables to ensure that autogen checking
      # stays specific to the file each thred represents.
      raise "#{path} missing autogen" unless Thread.current[:autogen]

      old_file_chmod_mode = File.stat(template).mode
      FileUtils.chmod(old_file_chmod_mode, path)

      format_output_file(path)
    end

    private

    def format_output_file(path)
      if path.end_with?('.py') && @env[:pyformat_enabled]
        run_formatter("python3 -m black --line-length 160 -S #{path}")
      elsif path.end_with?('.go') && @env[:goformat_enabled]
        run_formatter("gofmt -w -s #{path}")
        run_formatter("goimports -w #{path}")
      end
    end

    def run_formatter(command)
      output = %x(#{command} 2>&1)
      Google::LOGGER.error output unless $CHILD_STATUS.to_i.zero?
    end

    def relative_path(target, base)
      Pathname.new(target).relative_path_from(Pathname.new(base))
    end
  end

  # Basic functionality for code generator providers. Provides basic services,
  # such as compiling and including files, formatting data, etc.
  class Core
    include Compile::Core

    def initialize(config, api, start_time)
      @config = config
      @api = api
      @max_columns = DEFAULT_FORMAT_OPTIONS[:max_columns]

      # The compiler will error out if a file has been written in this compiler
      # run already. Instead of storing all the modified files in state we'll
      # use the time the file was modified.
      @start_time = start_time
      @py_format_enabled = check_pyformat
      @go_format_enabled = check_goformat
    end

    # This provides the FileTemplate class with access to a provider.
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
      compile_file_list(output_folder, @config.files.compile, version: version_name)
    end

    def compile_common_files(output_folder, version_name = nil)
      provider_name = self.class.name.split('::').last.downcase
      return unless File.exist?("provider/#{provider_name}/common~compile.yaml")

      Google::LOGGER.info "Compiling common files for #{provider_name}"
      files = YAML.safe_load(compile("provider/#{provider_name}/common~compile.yaml"))
      compile_file_list(output_folder, files, version: version_name)
    end

    def compile_file_list(output_folder, files, data = {})
      files.map do |target, source|
        Thread.new do
          Google::LOGGER.debug "Compiling #{source} => #{target}"
          target_file = File.join(output_folder, target)
          manifest = @config.respond_to?(:manifest) ? @config.manifest : {}
          FileTemplate.new({
            name: target,
            product: @api,
            manifest: manifest,
            output_folder: output_folder,
            product_ns: @api.name,
            env: {
              pyformat_enabled: @py_format_enabled,
              goformat_enabled: @go_format_enabled,
              start_time: @start_time
            }
          }.merge(data)).generate(source, target_file, self)
        end
      end.map(&:join)
    end

    def api_version_setup(version_name)
      version = @api.version_obj_or_default(version_name)
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

      generate_resource data.clone
      generate_resource_tests data.clone
    end

    def generate_datasources(output_folder, types, version_name)
      # We need to apply overrides for datasources
      @api = Overrides::Runner.build(@api, @config.datasources,
                                     @config.resource_override,
                                     @config.property_override)
      @api.validate

      version = @api.version_obj_or_default(version_name)
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
      FileTemplate.new(
        name: object.out_name,
        object: object,
        product: object.__product,
        product_ns: object.__product.name,
        output_folder: output_folder,
        version: version,
        config: @config,
        env: {
          pyformat_enabled: @py_format_enabled,
          goformat_enabled: @go_format_enabled,
          start_time: @start_time
        }
      )
    end

    def build_url(url_parts, extra = false)
      (product_url, obj_url) = url_parts
      extra_arg = ''
      extra_arg = ', extra_data' if extra
      ['URI.join(',
       indent([quote_string(product_url) + ',',
               'expand_variables(',
               indent(format_expand_variables(obj_url), 2),
               indent('data' + extra_arg, 2),
               ')'], 2),
       ')'].join("\n")
    end

    # TODO(rileykarson): Rehome this function.
    # For some reason the corresponding quote_string function lives in compile/core.rb
    # and no beside this function.
    def unquote_string(value)
      return value.gsub(/"(.*)"/, '\1') if value.start_with?('"')
      return value.gsub(/'(.*)'/, '\1') if value.start_with?("'")

      value
    end

    def true?(obj)
      obj.to_s.casecmp('true').zero?
    end

    def false?(obj)
      obj.to_s.casecmp('false').zero?
    end

    def emit_link(name, url, emit_self, extra_data = false)
      (params, fn_args) = emit_link_var_args(url, extra_data)
      code = ["def #{emit_self ? 'self.' : ''}#{name}(#{fn_args})",
              indent(url, 2),
              'end']

      if emit_self
        self_code = ['', "def #{name}(#{fn_args})",
                     "  self.class.#{name}(#{params.join(', ')})",
                     'end']
      end

      (code + (self_code || [])).join("\n")
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

    def update_url(resource, url_part)
      return build_url(resource.self_link_url) if url_part.nil?

      [resource.__product.base_url, url_part].flatten.join
    end

    # TODO(nelsonjr): Review all object interfaces and move to private methods
    # that should not be exposed outside the object hierarchy.
    private

    def generate_requires(properties, requires = [])
      requires.concat(properties.collect(&:requires))
    end

    def emit_requires(requires)
      requires.flatten.sort.uniq.map { |r| "require '#{r}'" }.join("\n")
    end

    def emit_link_var_args(url, extra_data)
      params = emit_link_var_args_list(url, extra_data,
                                       %w[data extra extra_data])
      defaults = emit_link_var_args_list(url, extra_data,
                                         [nil, "''", '{}'])
      [params.compact, params.zip(defaults)
                             .reject { |p| p[0].nil? }
                             .map { |p| p[1].nil? ? p[0] : "#{p[0]} = #{p[1]}" }
                             .join(', ')]
    end

    def emit_link_var_args_list(_url, extra_data, args_list)
      [args_list[0],
       (args_list[2] if extra_data)]
    end

    def wrap_field(field, spaces)
      avail_columns = DEFAULT_FORMAT_OPTIONS[:max_columns] - spaces - 5
      indent(field.scan(/\S.{0,#{avail_columns}}\S(?=\s|$)|\S+/), 2)
    end

    def format_section_ruler(size)
      size_pad = (size - size.to_s.length - 4) # < + > + 2 spaces around number.
      return unless size_pad.positive?

      ['<',
       '-' * (size_pad / 2), ' ', size.to_s, ' ', '-' * (size_pad / 2),
       (size_pad.even? ? '' : '-'),
       '>'].join
    end

    def format_box(existing, size)
      result = []
      result << ['+', '-' * size, '+'].join
      result << ['|', format_section_ruler(existing),
                 format_section_ruler(size - existing), '|'].join
      result << yield
      result << ['+', '-' * size, '+'].join
      result.join("\n")
    end

    def format_sources(sources, existing, size)
      format_box(existing, size) do
        sources.map do |source|
          source.split("\n").map do |l|
            right_pad_len = size - existing - l.length
            right_pad = right_pad_len.positive? ? ' ' * right_pad_len : ''
            '|' + '.' * existing + l + right_pad + '|'
          end
        end
      end
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

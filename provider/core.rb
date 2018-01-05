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
require 'dependencies/dependency_graph'
require 'fileutils'
require 'google/logger'
require 'pathname'
require 'provider/properties'
require 'provider/end2end/core'
require 'provider/test_matrix'
require 'provider/test_data/spec_formatter'
require 'provider/test_data/constants'
require 'provider/test_data/property'
require 'provider/test_data/create_data'
require 'provider/test_data/expectations'

module Provider
  DEFAULT_FORMAT_OPTIONS = {
    indent: 0,
    start_indent: 0,
    max_columns: 80,
    quiet: false
  }.freeze

  # Basic functionality for code generator providers. Provides basic services,
  # such as compiling and including files, formatting data, etc.
  class Core
    include Compile::Core
    include Provider::Properties
    include Provider::End2End::Core

    attr_reader :test_data

    def initialize(config, api)
      @config = config
      @api = api
      @property = Provider::TestData::Property.new(self)
      @constants = Provider::TestData::Constants.new(self)
      @data_gen = Provider::TestData::Generator.new
      @create_data = Provider::TestData::CreateData.new(self, @data_gen)
      @prop_data = Provider::TestData::Expectations.new(self, @data_gen)
      @generated = []
      @sourced = []
    end

    # Main entry point for the compiler. As this method is simply invoking other
    # generators, it is okay to ignore Rubocop warnings about method size and
    # complexity.
    #
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def generate(output_folder, types)
      generate_objects(output_folder, types)
      generate_client_functions(output_folder) unless @config.functions.nil?
      copy_files(output_folder) \
        unless @config.files.nil? || @config.files.copy.nil?
      compile_examples(output_folder) unless @config.examples.nil?
      compile_end2end_tests(output_folder) unless @config.examples.nil?
      compile_network_data(output_folder) \
        unless @config.test_data.nil? || @config.test_data.network.nil?
      compile_changelog(output_folder) unless @config.changelog.nil?
      # Compilation has to be the last step, as some files (e.g.
      # CONTRIBUTING.md) may depend on the list of all files previously copied
      # or compiled.
      compile_files(output_folder) \
        unless @config.files.nil? || @config.files.compile.nil?
      apply_file_acls(output_folder) \
        unless @config.files.nil? || @config.files.permissions.nil?
      verify_test_matrixes
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity

    def copy_files(output_folder)
      @config.files.copy.each do |target, source|
        target_file = File.join(output_folder, target)
        target_dir = File.dirname(target_file)
        @sourced << relative_path(target_file, output_folder)
        Google::LOGGER.info "Copying #{source} => #{target}"
        FileUtils.mkpath target_dir unless Dir.exist?(target_dir)
        FileUtils.cp source, target_file
      end
    end

    def compile_files(output_folder)
      compile_file_list(output_folder, @config.files.compile)
    end

    def compile_examples(output_folder)
      compile_file_map(
        output_folder,
        @config.examples,
        lambda do |_object, file|
          ["examples/#{file}",
           "products/#{@api.prefix[1..-1]}/files/examples~#{file}"]
        end
      )
    end

    def compile_network_data(output_folder)
      compile_file_map(
        output_folder,
        @config.test_data.network,
        lambda do |object, file|
          type = Google::StringUtils.underscore(object.name)
          ["spec/data/network/#{object.out_name}/#{file}.yaml",
           "products/#{@api.prefix[1..-1]}/files/spec~#{type}~#{file}.yaml"]
        end
      )
    end

    # Generate the CHANGELOG.md file with the history of the module.
    def compile_changelog(output_folder)
      FileUtils.mkpath output_folder
      generate_file(
        changes: @config.changelog,
        template: 'templates/CHANGELOG.md.erb',
        output_folder: output_folder,
        out_file: File.join(output_folder, 'CHANGELOG.md')
      )
    end

    def apply_file_acls(output_folder)
      @config.files.permissions.each do |perm|
        Google::LOGGER.info "Permission #{perm.path} => #{perm.acl}"
        FileUtils.chmod perm.acl, File.join(output_folder, perm.path)
      end
    end

    def compile_file_map(output_folder, section, mapper)
      create_object_list(section, mapper).each do |o|
        compile_file_list(
          output_folder,
          o
        )
      end
    end

    # Creates an object list by calling a lambda
    # This can be useful for converting a list of config values to something
    # less human-centric.
    def create_object_list(section, mapper)
      @api.objects
          .select { |o| section.key?(o.name) }
          .map do |o|
            Hash[section[o.name].map { |file| mapper.call(o, file) }]
          end
    end

    def list_manual_network_data
      create_object_list(
        @config.test_data.network,
        lambda do |object, file|
          type = Google::StringUtils.underscore(object.name)
          ["spec/data/network/#{object.out_name}/#{file}.yaml",
           "products/#{@api.prefix[1..-1]}/files/spec~#{type}~#{file}.yaml"]
        end
      )
    end

    # rubocop:disable Metrics/MethodLength
    def compile_file_list(output_folder, files, data = {})
      files.each do |target, source|
        Google::LOGGER.info "Compiling #{source} => #{target}"
        target_file = File.join(output_folder, target)
                          .gsub('{{product_name}}', @api.prefix[1..-1])
        generate_file(
          data.clone.merge(
            name: target,
            product: @api,
            object: {},
            config: {},
            scopes: @api.scopes,
            manifest: @config.manifest,
            tests: '',
            template: source,
            generated_files: @generated,
            sourced_files: @sourced,
            compiler: compiler,
            output_folder: output_folder,
            out_file: target_file,
            prop_ns_dir: @api.prefix[1..-1].downcase,
            product_ns: Google::StringUtils.camelize(@api.prefix[1..-1], :upper)
          )
        )
      end
    end
    # rubocop:enable Metrics/MethodLength

    def generate_objects(output_folder, types)
      @api.objects.each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info "Excluding #{object.name} per user request"
        elsif types.empty? && object.exclude
          Google::LOGGER.info "Excluding #{object.name} per API catalog"
        else
          generate_object object, output_folder
        end
      end
    end

    def generate_object(object, output_folder)
      data = build_object_data(object, output_folder)

      generate_resource data
      generate_resource_tests data
      generate_properties data, object.all_user_properties
      generate_network_datas data, object
    end

    # rubocop:disable Metrics/MethodLength
    # Generates all 6 network data files for a object.
    # This includes all combinations of seeds [0-2] and title == / != name
    # Each data file is a YAML file with all properties possible on an object.
    #
    # @config.test_data lists all files that are written by hand and will not
    # be generated.
    #
    # Requires:
    #  object: The Api::Resource used as basis for the network data.
    #  data: A hash with values:
    #    output_folder: root folder for generated module
    def generate_network_datas(data, object)
      target_folder = File.join(data[:output_folder],
                                'spec', 'data', 'network', object.out_name)
      FileUtils.mkpath target_folder

      # Create list of compiled network data
      manual = list_manual_network_data
      3.times.each do |id|
        %w[name title].each do |name|
          out_file = File.join(target_folder, "success#{id + 1}~#{name}.yaml")
          next if manual.include? out_file
          next if true?(Google::HashUtils.navigate(data[:config], %w[manual]))

          generate_network_data data.clone.merge(
            out_file: File.join(target_folder, "success#{id + 1}~#{name}.yaml"),
            id: id,
            title: name,
            object: object
          )
        end
      end
    end
    # rubocop:enable Metrics/MethodLength

    # Generates a single network data file for unit testing.
    # Required values in data:
    #   out_file: path of data file to create
    #   id: a seed value
    #   title: The name of object who is unit tested with this spec file
    #   object: The Api::Resource used as basis for the network data
    def generate_network_data(data)
      formatter = Provider::TestData::SpecFormatter.new(self)

      name = "title#{data[:id]}" if data[:title] == 'title'
      name = "test name##{data[:id]} data" if data[:title] == 'name'
      generate_file data.clone.merge(
        template: 'templates/network_spec.yaml.erb',
        test_data: formatter.generate(data[:object], '', data[:object].kind,
                                      data[:id],
                                      name: name)
      )
    end

    def build_object_data(object, output_folder)
      {
        name: object.out_name,
        object: object,
        config: (@config.objects || {}).select { |o, _v| o == object.name }
                                       .fetch(object.name, {}),
        tests: (@config.tests || {}).select { |o, _v| o == object.name }
                                    .fetch(object.name, {}),
        output_folder: output_folder,
        product_name: object.__product.prefix[1..-1]
      }
    end

    def generate_resource_file(data)
      generate_file(data.clone.merge(
        # Override with provider specific template for this object, if needed
        template: Google::HashUtils.navigate(data[:config], ['template',
                                                             data[:type]],
                                             data[:default_template]),
        product_ns:
          Google::StringUtils.camelize(data[:object].__product.prefix[1..-1],
                                       :upper)
      ))
    end

    def generate_client_functions(output_folder)
      @config.functions.each do |fn|
        info = generate_client_function(output_folder, fn)
        FileUtils.mkpath info[:target_folder]
        generate_file info.clone
      end
    end

    def format_expand_variables(obj_url)
      obj_url = obj_url.split("\n") unless obj_url.is_a?(Array)
      if obj_url.size > 1
        ['[',
         indent_list(obj_url.map { |u| quote_string(u) }, 2),
         '].join,']
      else
        [obj_url.map { |u| quote_string(u) }[0] + ',']
      end
    end

    def build_url(product_url, obj_url, extra = false)
      extra_arg = ''
      extra_arg = ', extra_data' if obj_url.to_s.include?('<|extra|>') || extra
      ['URI.join(',
       indent([quote_string(product_url) + ',',
               'expand_variables(',
               indent(format_expand_variables(obj_url), 2),
               indent('data' + extra_arg, 2),
               ')'], 2),
       ')'].join("\n")
    end

    def async_operation_url(resource)
      build_url(resource.__product.base_url, resource.async.operation.base_url,
                true)
    end

    def collection_url(resource)
      base_url = resource.base_url.split("\n").map(&:strip).compact
      build_url(resource.__product.base_url, base_url)
    end

    def self_link_raw_url(resource)
      base_url = resource.__product.base_url.split("\n").map(&:strip).compact
      if resource.self_link.nil?
        [base_url, [resource.base_url, '{{name}}'].join('/')]
      else
        self_link = resource.self_link.split("\n").map(&:strip).compact
        [base_url, self_link]
      end
    end

    def self_link_url(resource)
      (product_url, resource_url) = self_link_raw_url(resource)
      build_url(product_url, resource_url)
    end

    def extract_variables(template)
      template.scan(/{{[^}]*}}/)
              .map { |v| v.gsub(/{{([^}]*)}}/, '\1') }
              .map(&:to_sym)
    end

    def variable_type(object, var)
      return Api::Type::String::PROJECT if var == :project
      return Api::Type::String::NAME if var == :name
      v = object.all_user_properties.select { |p| p.out_name.to_sym == var ||
                                              p.name.to_sym == var }
                .first
      return v.property if v.is_a?(Api::Type::ResourceRef)
      v
    end

    def quote_string(value)
      raise 'Invalid value' if value.nil?
      if value.include?('#{') || value.include?("'")
        ['"', value, '"'].join
      else
        ["'", value, "'"].join
      end
    end

    # Used to convert a string 'a b c' into a\ b\ c for use in %w[...] form
    def str2warray(value)
      unquote_string(value).gsub(/ /, '\\ ')
    end

    def unquote_string(value)
      return value.gsub(/"(.*)"/, '\1') if value.start_with?('"')
      return value.gsub(/'(.*)'/, '\1') if value.start_with?("'")
      value
    end

    # TODO(alexstephen): Retire in favor of a real code object.
    # No validation is possible on get_code_multiline
    def get_code_multiline(config, node)
      search = node.class <= Array ? node : [node]
      Google::HashUtils.navigate(config, search)
    end

    def lines(code, number = 0)
      return if code.nil? || code.empty?
      code = code.join("\n") if code.is_a?(Array)
      code[-1] = '' while code[-1] == "\n" || code[-1] == "\r"
      "#{code}#{"\n" * (number + 1)}"
    end

    def lines_before(code, number = 0)
      return if code.nil? || code.empty?
      code = code.join("\n") if code.is_a?(Array)
      code[0] = '' while code[0] == "\n" || code[0] == "\r"
      "#{"\n" * (number + 1)}#{code}"
    end

    def true?(obj)
      obj.to_s.casecmp('true').zero?
    end

    def false?(obj)
      obj.to_s.casecmp('false').zero?
    end

    def emit_method(name, args, code, file_name, opts = {})
      method_decl = "def #{name}"
      method_decl << "(#{args.join(', ')})" unless args.empty?
      (rubo_off, rubo_on) = emit_rubo_pair(file_name, name, opts)
      [
        (rubo_off unless rubo_off.empty?),
        method_decl,
        indent(code, 2),
        'end',
        (rubo_on unless rubo_on.empty?)
      ].compact.join("\n")
    end

    def emit_rubo_pair(file_name, name, opts = {})
      [
        emit_rubo_item(file_name, name, :disabled, opts),
        emit_rubo_item(file_name, name, :enabled, opts)
      ]
    end

    def emit_rubo_item(file_name, name, state, opts = {})
      [
        (if opts.key?(:class_name)
           get_rubocop_exceptions(file_name, :function,
                                  [opts[:class_name], name].join('.'), state)
         end),
        get_rubocop_exceptions(file_name, :function, name, state)
      ].compact.flatten
    end

    def emit_rubocop(ctx, pinpoint, name, state)
      get_rubocop_exceptions(ctx.local_variable_get(:file_relative), pinpoint,
                             name, state).join("\n")
    end

    # TODO(nelsonjr): Track usage of exceptions and fail if some are
    # left unused. E.g. we change the function/class name and the setting
    # in the YAML file is now useless.
    def get_rubocop_exceptions(file_name, pinpoint, name, state)
      name = name.flatten.join(' > ') if pinpoint == :test
      flags = get_style_exceptions(file_name, pinpoint, name)
      flags = flags.reverse if state == :enabled
      flags.map do |e|
        "# rubocop:#{state == :enabled ? 'enable' : 'disable'} #{e}"
      end
    end

    def get_style_exceptions(file_name, type, name)
      styles = @config.style
      return [] if styles.nil?
      styles.select { |s| s.name == file_name }
            .map(&:pinpoints)
            .flatten
            .select { |ps| ps.any? { |k, v| k.to_sym == type && v == name } }
            .map { |p| p['exceptions'] }
            .flatten
            .sort
    end

    def emit_link(name, url, emit_self, extra_data = false)
      (params, fn_args) = emit_link_var_args(url, extra_data)
      code = ["def #{emit_self ? 'self.' : ''}#{name}(#{fn_args})",
              indent(url, 2).gsub("'<|extra|>'", 'extra'),
              'end']

      if emit_self
        self_code = ['', "def #{name}(#{fn_args})",
                     "  self.class.#{name}(#{params.join(', ')})",
                     'end']
      end

      (code + (self_code || [])).join("\n")
    end

    # Formats the code and returns the first candidate that fits the alloted
    # column limit.
    def format(sources, indent = 0, start_indent = 0,
               max_columns = DEFAULT_FORMAT_OPTIONS[:max_columns])
      format2(sources, indent: indent,
                       start_indent: start_indent,
                       max_columns: max_columns)
    end

    # TODO(nelsonjr): Make format2 into format and fix all references throughout
    # the code base.
    def format2(sources, overrides = {})
      options = DEFAULT_FORMAT_OPTIONS.merge(overrides)
      output = ''
      avail_columns = options[:max_columns] - options[:start_indent]
      sources.each do |attempt|
        output = indent(attempt, options[:indent])
        return output if format_fits?(output, options[:start_indent],
                                      options[:max_columns])
      end
      unless options[:on_misfit].nil?
        (alt_fit, alt_output) = options[:on_misfit].call(sources, output,
                                                         options, avail_columns)
        return alt_output if alt_fit
      end
      fail_and_log_format_error output, options, avail_columns \
        unless options[:quiet]
    end

    def fail_and_log_format_error(output, options, avail_columns)
      Google::LOGGER.info [
        ["No code option fits in #{avail_columns} columns",
         "w/ #{options[:start_indent]} left indent:"].join(' '),
        format_sources(output.split("\n"), options[:start_indent],
                       options[:max_columns]),
        (unless options[:on_misfit].nil?
           format_sources(alt_output.split("\n"), options[:start_indent],
                          options[:max_columns])
         end)
      ].compact.join("\n")
      raise ArgumentError, "No code fits in #{avail_columns}"
    end

    def format_fits?(output, start_indent,
                     max_columns = DEFAULT_FORMAT_OPTIONS[:max_columns])
      output = output.flatten.join("\n") if output.is_a?(::Array)
      output = output.split("\n") unless output.is_a?(::Array)
      output.select { |l| l.length > (max_columns - start_indent) }.empty?
    end

    def relative_path(target, base)
      Pathname.new(target).relative_path_from(Pathname.new(base))
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

    def emit_link_var_args_list(url, extra_data, args_list)
      [args_list[0],
       (args_list[1] if url.include?('<|extra|>')),
       (args_list[2] if url.include?('<|extra|>') || extra_data)]
    end

    def generate_file(data)
      file_folder = File.dirname(data[:out_file])
      file_relative = relative_path(data[:out_file], data[:output_folder]).to_s
      FileUtils.mkpath file_folder unless Dir.exist?(file_folder)
      @generated << relative_path(data[:out_file], data[:output_folder])
      ctx = binding
      data.each { |name, value| ctx.local_variable_set(name, value) }
      generate_file_write ctx, data
    end

    def generate_file_write(ctx, data)
      Google::LOGGER.info "Generating #{data[:name]} #{data[:type]}"
      write_file data[:out_file], compile_file(ctx, data[:template])
    end

    def compile_file(ctx, source)
      ERB.new(File.read(source), nil, '-%>').result(ctx).split("\n")
    rescue StandardError => e
      puts "Error compiling file: #{source}"
      raise e
    end

    # Write the output to a file. We write one line at a time so tests can
    # reason about what's being written and validate the output.
    def write_file(out_file, output)
      File.open(out_file, 'w') { |f| output.each { |l| f.write("#{l}\n") } }
    end

    def wrap_field(field, spaces)
      avail_columns = DEFAULT_FORMAT_OPTIONS[:max_columns] - spaces - 5
      indent(field.scan(/\S.{0,#{avail_columns}}\S(?=\s|$)|\S+/), 2)
    end

    def verify_test_matrixes
      Provider::TestMatrix::Registry.instance.verify_all
    end

    def format_section_ruler(size)
      size_pad = (size - size.to_s.length - 4) # < + > + 2 spaces around number.
      return unless size_pad > 0
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
            right_pad = right_pad_len > 0 ? ' ' * right_pad_len : ''
            '|' + '.' * existing + l + right_pad + '|'
          end
        end
      end
    end

    def emit_user_agent(product, extra, notes, file_name)
      prov_text = Google::StringUtils.camelize(self.class.name.split('::').last,
                                               :upper)
      prod_text = Google::StringUtils.camelize(product, :upper)
      ua_generator = notes.map { |n| "# #{n}" }.concat(
        [
          "version = '1.0.0'",
          '[',
          indent_list([
            "\"Google#{prov_text}#{prod_text}/\#{version}\"",
            extra
          ].compact, 2),
          "].join(' ')"
        ]
      )
      emit_method('generate_user_agent', [], ua_generator.compact, file_name)
    end

    def provider_name
      self.class.name.split('::').last.downcase
    end

    # Returns true if this module needs access to the saved API response
    # This response is stored in the @fetched variable
    # Requires:
    #   config: The config for an object
    #   object: An Api::Resource object
    def save_api_results?(config, object)
      fetched_props = object.exported_properties.select do |p|
        p.is_a? Api::Type::FetchedExternal
      end
      Google::HashUtils.navigate(config, %w[access_api_results]) || \
        !fetched_props.empty?
    end

    # Generates the documentation for the client side function to be
    # included in the module. Call this function immediately before the function
    # definition and the code generator will use data from api.yaml to build the
    # documentation comment block.
    #
    # rubocop: Method returns a big array. Easier to read a single block
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/MethodLength
    def emit_function_doc(function)
      [
        function.description.strip,
        '',
        'Arguments:',
        indent(function.arguments.map do |arg|
                 [
                   "- #{arg.name}: #{arg.type.split('::').last.downcase}",
                   indent(arg.description.strip.split("\n"), 2)
                 ]
               end, 2),
        (
          unless function.examples.nil?
            [
              '',
              'Examples:',
              indent(function.examples.map { |eg| "- #{eg}" }, 2)
            ]
          end
        ),
        (
          unless function.notes.nil?
            [
              '',
              function.notes.strip
            ]
          end
        )
      ].compact.flatten.join("\n").split("\n").map { |l| "# #{l}".strip }
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize
  end
end

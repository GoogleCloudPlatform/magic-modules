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
require 'google/extensions'
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
    max_columns: 100,
    quiet: false
  }.freeze

  # Basic functionality for code generator providers. Provides basic services,
  # such as compiling and including files, formatting data, etc.
  class Core
    include Compile::Core
    include Provider::Properties
    include Provider::End2End::Core
    include Api::Object::ObjectUtils

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
      @max_columns = DEFAULT_FORMAT_OPTIONS[:max_columns]
    end

    # Main entry point for the compiler. As this method is simply invoking other
    # generators, it is okay to ignore Rubocop warnings about method size and
    # complexity.
    #
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def generate(output_folder, types, version_name)
      version = @api.version_obj_or_default(version_name)
      generate_objects(output_folder, types, version)
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

      generate_datasources(output_folder, types, version) \
        unless @config.datasources.nil?
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
        FileUtils.copy_entry source, target_file
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
          type = object.name.underscore
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
      test_data = @config&.test_data&.network || {}
      create_object_list(
        test_data,
        lambda do |object, file|
          type = object.name.underscore
          ["spec/data/network/#{object.out_name}/#{file}.yaml",
           "products/#{@api.prefix[1..-1]}/files/spec~#{type}~#{file}.yaml"]
        end
      )
    end

    # rubocop:disable Metrics/MethodLength
    # rubocop:disable Metrics/AbcSize
    def compile_file_list(output_folder, files, data = {})
      files.each do |target, source|
        Google::LOGGER.info "Compiling #{source} => #{target}"
        target_file = File.join(output_folder, target)
                          .gsub('{{product_name}}', @api.prefix[1..-1])

        manifest = @config.respond_to?(:manifest) ? @config.manifest : {}
        generate_file(
          data.clone.merge(
            name: target,
            product: @api,
            object: {},
            config: {},
            scopes: @api.scopes,
            manifest: manifest,
            tests: '',
            template: source,
            generated_files: @generated,
            sourced_files: @sourced,
            compiler: compiler,
            output_folder: output_folder,
            out_file: target_file,
            prop_ns_dir: @api.prefix[1..-1].downcase,
            product_ns: @api.prefix[1..-1].camelize(:upper)
          )
        )

        %x(goimports -w #{target_file}) if File.extname(target_file) == '.go'
      end
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize

    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    # rubocop:disable Metrics/AbcSize
    def generate_objects(output_folder, types, version)
      @api.set_properties_based_on_version(version)
      (@api.objects || []).each do |object|
        if !types.empty? && !types.include?(object.name)
          Google::LOGGER.info "Excluding #{object.name} per user request"
        elsif types.empty? && object.exclude
          Google::LOGGER.info "Excluding #{object.name} per API catalog"
        elsif types.empty? && object.exclude_if_not_in_version(version)
          Google::LOGGER.info "Excluding #{object.name} per API version"
        else
          generate_object object, output_folder, version
        end
      end
    end
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity
    # rubocop:enable Metrics/AbcSize

    def generate_object(object, output_folder, version)
      data = build_object_data(object, output_folder, version)

      generate_resource data
      generate_resource_tests data
      generate_properties data, object.all_user_properties
      generate_network_datas data, object
    end

    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def generate_datasources(output_folder, types, version)
      # We need to apply overrides for datasources
      @config.datasources.validate

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
        elsif types.empty? && object.exclude_if_not_in_version(version)
          Google::LOGGER.info(
            "Excluding #{object.name} datasource per API version"
          )
        else
          generate_datasource object, output_folder, version
        end
      end
    end
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity
    # rubocop:enable Metrics/AbcSize

    def generate_datasource(object, output_folder, version)
      data = build_object_data(object, output_folder, version)

      compile_datasource data
    end

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
          next if true?(data[:object].manual)

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

    def build_object_data(object, output_folder, version)
      {
        name: object.out_name,
        object: object,
        config: (@config.objects || {}).select { |o, _v| o == object.name }
                                       .fetch(object.name, {}),
        tests: (@config.tests || {}).select { |o, _v| o == object.name }
                                    .fetch(object.name, {}),
        output_folder: output_folder,
        product_name: object.__product.prefix[1..-1],
        version: version
      }
    end

    def generate_resource_file(data)
      product_ns = if @config.name.nil?
                     data[:object].__product.prefix[1..-1].camelize(:upper)
                   else
                     @config.name
                   end
      generate_file(data.clone.merge(
        # Override with provider specific template for this object, if needed
        template: Google::HashUtils.navigate(data[:config], ['template',
                                                             data[:type]],
                                             data[:default_template]),
        product_ns: product_ns
      ))
    end

    def generate_client_functions(output_folder)
      @config.functions.each do |fn|
        info = generate_client_function(output_folder, fn)
        FileUtils.mkpath info[:target_folder]
        generate_file info.clone
      end
    end

    # rubocop:disable Metrics/AbcSize
    def format_expand_variables(obj_url)
      obj_url = obj_url.split("\n") unless obj_url.is_a?(Array)
      if obj_url.size > 1
        ['[',
         indent_list(obj_url.map { |u| quote_string(u) }, 2),
         '].join,']
      else
        vars = quote_string(obj_url[0])
        vars_parts = obj_url[0].split('/')
        format([
                 [[vars, ','].join],
                 # vars is too big to fit, split in half
                 vars_parts.each_slice((vars_parts.size / 2.0).round).to_a
                 .map.with_index do |p, i|
                   # Use implicit string joining for the first line.
                   quote_string(p.join('/')) + (i.zero? ? ' \\' : ',')
                 end
               ], 0, 8)
      end
    end
    # rubocop:enable Metrics/AbcSize

    def build_url(url_parts, extra = false)
      (product_url, obj_url) = url_parts
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

    def extract_variables(template)
      template.scan(/{{[^}]*}}/)
              .map { |v| v.gsub(/{{([^}]*)}}/, '\1') }
              .map(&:to_sym)
    end

    def variable_type(object, var)
      return Api::Type::String::PROJECT if var == :project
      return Api::Type::String::NAME if var == :name
      v = object.all_user_properties
                .select { |p| p.out_name.to_sym == var || p.name.to_sym == var }
                .first
      return v.property if v.is_a?(Api::Type::ResourceRef)
      v
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

    def true?(obj)
      obj.to_s.casecmp('true').zero?
    end

    def false?(obj)
      obj.to_s.casecmp('false').zero?
    end

    def emit_method(name, args, code, file_name, opts = {})
      (rubo_off, rubo_on) = emit_rubo_pair(file_name, name, opts)
      [
        (rubo_off unless rubo_off.empty?),
        method_decl(name, args),
        indent(code, 2),
        'end',
        (rubo_on unless rubo_on.empty?)
      ].compact.join("\n")
    end

    def method_decl(name, args)
      ["def #{name}", ("(#{args.join(', ')})" unless args.empty?)].compact.join
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
               max_columns = @max_columns)
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

      indent([
               '# rubocop:disable Metrics/LineLength',
               sources.last,
               '# rubocop:enable Metrics/LineLength'
             ], options[:indent])
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

    def check_requires(object, *requires)
      return if object.requires.nil?
      requires_list = object.requires
      missing = requires.flatten.reject { |r| requires_list.any?(r) }
      raise <<~ERROR unless missing.empty?
        Including #{__FILE__} needs the following requires: #{missing}
        Please add them to 'object > requires' section of <provider>.yaml
      ERROR
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
      enforce_file_expectations data[:out_file] do
        Google::LOGGER.info "Generating #{data[:name]} #{data[:type]}"
        write_file data[:out_file], compile_file(ctx, data[:template])
      end
    end

    def enforce_file_expectations(filename)
      @file_expectations = {
        autogen: false
      }
      yield
      raise "#{filename} missing autogen" unless @file_expectations[:autogen]
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

    def emit_user_agent(product, extra, notes, file_name)
      prov_text = self.class.name.split('::').last.camelize(:upper)
      prod_text = product.camelize(:upper)
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

    # Filter the properties to keep only the ones requiring custom update
    # method and group them by update url & verb.
    def properties_by_custom_update(properties)
      update_props = properties.reject do |p|
        p.update_url.nil? || p.update_verb.nil?
      end
      update_props.group_by do |p|
        { update_url: p.update_url, update_verb: p.update_verb }
      end
    end
  end
end

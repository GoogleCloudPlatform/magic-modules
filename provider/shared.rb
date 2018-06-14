# Copyright 2018 Google Inc.
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

module Provider
  # These are functions that are shared across providers.
  # rubocop:disable Metrics/ModuleLength
  module Shared
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
                   .map { |p| quote_string(p.join('/')) + ',' }
               ], 0, 8)
      end
    end
    # rubocop:enable Metrics/AbcSize

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
      build_url(resource.__product.base_url,
                resource.async.operation.base_url,
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
  end
  # rubocop:enable Metrics/ModuleLength
end

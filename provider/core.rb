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

require 'provider/generation'
require 'compile/core'
require 'fileutils'
require 'google/extensions'
require 'google/logger'
require 'google/hash_utils'
require 'pathname'

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
    include Provider::Generation

    def self.generation_steps(*steps)
      @steps = steps
    end

    def initialize(config, api)
      @config = config
      @api = api
      @max_columns = DEFAULT_FORMAT_OPTIONS[:max_columns]
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

    def generate_file(data)
      file_folder = File.dirname(data[:out_file])
      # This variable looks unused, but is used in ansible/resource.erb
      file_relative = relative_path(data[:out_file], data[:output_folder]).to_s
      FileUtils.mkpath file_folder unless Dir.exist?(file_folder)
      ctx = binding
      data.each { |name, value| ctx.local_variable_set(name, value) }
      generate_file_write ctx, data
    end

    def generate_file_write(ctx, data)
      enforce_file_expectations data[:out_file] do
        Google::LOGGER.debug "Generating #{data[:name]} #{data[:type]}"
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

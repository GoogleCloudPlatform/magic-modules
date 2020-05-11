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

require 'binding_of_caller'
require 'erb'
require 'ostruct'

module Compile
  # Unique ID for the Google libraries to be compiled/used by modules
  module Libraries
    NETWORK = 'network'.freeze
  end

  # Helper functions to aid compiling and including files
  module Core
    def compiler
      "#{self.class.name.split('::').last.downcase}-codegen".freeze
    end

    def include(file)
      get_helper_file(file)
    end

    # Compiles a ERB template.
    #
    # This function is used to compile code/docs/tests/etc. by using Embedded
    # Ruby (ERB).
    #
    # This function adjusts the Ruby Bindings* to expose local and instance
    # variables of the caller stack as locally scopes variables in the ERB
    # template context. This non-standard behavior allows for easier writing of
    # templates without having to fetch context from other sources. Any data
    # that the provider writer wants to expose to the template can be made
    # available through the bindings.
    #
    # This allows for exposing instance variables, such as @api and config, or
    # provider or function specific data, such as output_folder. Functions such
    # as Provider::Core::generate_file() or Provider::Core::compile_file() can
    # allow providers to generate such bindings (by passing them in a Hash to
    # the function).
    #
    # This short notation makes the template easier to write. For example in
    # GRPC for you to access a context specific variable you have to call a
    # global function, such as: Contexts.current().get(SOME_GRPC_VAR).
    #
    # Note that if more than 1 compile() function is invoked reentrantly the
    # next run will be adjusted to the same context. So despite the call stack
    # being modified, the view of the variables will not. In the example below
    # the text of the 'file2.erb' template will print properly the name of the
    # product from the @api instance variable of the provider object N stack
    # traces behind.
    #
    #   Call stack:                    Context:
    #
    #   provider::some_func
    #                                  instance variable: @api
    #                                  local variable: object
    #                                  local variable: output_file
    #   compile('file1.erb')
    #
    #     ...
    #     <%= compile('TEST.md.erb') -%>
    #     ...
    #
    #   compile('TEST.md.erb')
    #
    #     ...
    #     We're compiling object '<%= object.name -%>' for '<%= @api.name -%>'.
    #     This file is named '<%= output_file -%>'.
    #     ...
    #
    # It would generate an output like:
    #
    #   We're compiling object 'Disk' for 'Google Compute Engine'.
    #   This file is named 'build/provider/compute/TEST.md'.
    #
    # If you want to avoid the load/save you can use one helper functions, such
    # as Provider::Core.compile_file(), which take a hash with overrides. For
    # example:
    #
    # parent.erb:
    #   <% indent_spaces = 2 -%>
    #   Parent: <%= indent_spaces -%>
    #   <%=
    #     compile_file({ indent_spaces: indent_spaces + 2 },
    #                  'child.erb')
    #   -%>
    #   Parent: <%= indent_spaces -%>
    #
    # child.erb:
    #   Child: <%= indent_spaces -%>
    #
    # This would produce something like:
    #
    #   Parent: 2
    #   Child: 4
    #   Parent 2
    #
    # This function is reentrant.
    #
    # * Binding: Objects of class Binding encapsulate the execution context at
    # some particular place in the code and retain this context for future use.
    # The variables, methods, value of self, and possibly an iterator block that
    # can be accessed in this context are all retained. Binding objects can be
    # created using Kernel#binding, and are made available to the callback of
    # Kernel#set_trace_func. -- https://ruby-doc.org/core-2.2.0/Binding.html
    #
    # WARNING: Do *NOT* change caller_frame = 1 value unless you really know
    # what you're doing. It is reserved for future use.
    def compile(file, caller_frame = 1)
      ctx = binding.of_caller(caller_frame)
      has_erbout = ctx.local_variables.include?(:_erbout)
      content = ctx.local_variable_get(:_erbout) if has_erbout # save code
      ctx.local_variable_set(:compiler, compiler)
      Google::LOGGER.debug "Compiling #{file}"
      input = ERB.new get_helper_file(file), trim_mode: '->'
      compiled = input.result(ctx)
      ctx.local_variable_set(:_erbout, content) if has_erbout # restore code
      compiled
    rescue StandardError
      Google::LOGGER.fatal "Error compiling #{file}"
      raise
    end

    def ansible_style_yaml(obj, options = {})
      if obj.is_a?(::Hash)
        obj.reject { |_, v| v.nil? }.to_yaml(options).sub("---\n", '')
      else
        obj.to_yaml(options).sub("---\n", '')
      end
    end

    # Compiles a ERB template from a file.
    #
    # Arguments:
    #
    # - ctx: A binding (or hash) that provides the context to expose to the ERB
    #        environment
    # - source: A path to a file in ERB format.
    #
    # Refer to Compile::Core.compile for full details about the compilation
    # process.
    def compile_file(ctx, source)
      $pwd ||= Dir.pwd
      compile_string(ctx, File.read($pwd + '/' + source))
    rescue StandardError => e
      puts "Error compiling file: #{source}"
      raise e
    end

    def indent(text, spaces, filler = ' ')
      indent_array(text, spaces, filler).join("\n")
    end

    def indent_list(text, spaces, last_line_comma = false, filler = ' ')
      if last_line_comma
        [indent_array(text, spaces, filler).join(",\n"), ','].join
      else
        indent_array(text, spaces, filler).join(",\n")
      end
    end

    def indent_array(text, spaces, filler = ' ')
      return [] if text.nil?

      lines = text.class <= Array ? text : text.split("\n")
      lines.map do |line|
        if line.class <= Array
          indent(line, spaces, filler)
        elsif line.include?("\n")
          indent(line.split("\n"), spaces, filler)
        elsif line.strip.empty?
          ''
        else
          ' ' * spaces + line.gsub(/\n/, "\n" + ' ' * spaces)
        end
      end
    end

    def quote_string(value)
      raise 'Invalid value' if value.nil?

      if value.include?('#{') || value.include?("'")
        ['"', value, '"'].join
      else
        ["'", value, "'"].join
      end
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

    # Compiles an ERB template using the data from a key-value pair.
    # The key-value pair may be a Hash or a Binding
    def compile_string(ctx, source)
      if ctx.is_a? Binding
        ERB.new(source, trim_mode: '->').result(ctx).split("\n")
      elsif ctx.is_a? Hash
        ERB.new(source, trim_mode: '->').result(
          OpenStruct.new(ctx).instance_eval { binding.of_caller(1) }
        ).split("\n")
      else
        raise TypeError, "#{ctx.class} is not a valid type for compilation"
      end
    end

    def autogen_notice(lang)
      Thread.current[:autogen] = true
      comment_block(compile('templates/autogen_notice.erb').split("\n"), lang)
    end

    def autogen_exception
      Thread.current[:autogen] = true
    end

    def comment_block(text, lang)
      case lang
      when :ruby, :python, :yaml, :git, :gemfile
        text.map { |t| t&.empty? ? '#' : "# #{t}" }
      when :go
        text.map { |t| t&.empty? ? '//' : "// #{t}" }
      when :html, :markdown
        [
          '<!--',
          indent(text, 2),
          '-->'
        ]
      else
        raise "Unknown language for comment: #{lang}"
      end
    end

    private

    def autogen_notice_contrib
      ['Please read more about how to change this file in README.md and',
       'CONTRIBUTING.md located at the root of this package.']
    end

    def get_helper_file(file, remove_copyright_notice = true)
      $pwd ||= Dir.pwd
      content = IO.read($pwd + '/' + file)
      remove_copyright_notice ? strip_copyright_notice(content) : content
    end

    def strip_copyright_notice(content, comment_marker = '#')
      lines = content.split("\n")
      return content unless lines[0].include?('Copyright 20')

      lines = lines.drop(1) while lines[0].start_with?(comment_marker)
      lines = lines.drop(1) while lines[0].strip.empty?
      lines.join("\n")
    end
  end
end

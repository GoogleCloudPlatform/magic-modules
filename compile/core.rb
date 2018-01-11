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

    def compile(file, caller_frame = 1)
      ctx = binding.of_caller(caller_frame)
      has_erbout = ctx.local_variables.include?(:_erbout)
      content = ctx.local_variable_get(:_erbout) if has_erbout # save code
      ctx.local_variable_set(:compiler, compiler)
      Google::LOGGER.info "Compiling #{file}"
      input = ERB.new get_helper_file(file, true), nil, '-%>'
      compiled = input.result(ctx)
      ctx.local_variable_set(:_erbout, content) if has_erbout # restore code
      compiled
    end

    def compile_if(config, node)
      file = Google::HashUtils.navigate(config, node)
      compile(file, 2) unless file.nil?
    end

    def indent(text, spaces)
      indent_array(text, spaces).join("\n")
    end

    def indent_list(text, spaces, last_line_comma = false)
      if last_line_comma
        [indent_array(text, spaces).join(",\n"), ','].join
      else
        indent_array(text, spaces).join(",\n")
      end
    end

    def indent_array(text, spaces)
      return [] if text.nil?
      lines = text.class <= Array ? text : text.split("\n")
      lines.map do |line|
        if line.class <= Array
          indent(line, spaces)
        elsif line.include?("\n")
          indent(line.split("\n"), spaces)
        elsif line.strip.empty?
          ''
        else
          ' ' * spaces + line.gsub(/\n/, "\n" + ' ' * spaces)
        end
      end
    end

    # Includes a require clause and schedules library to be copied, potentially
    # with its dependencies & tests.
    def emit_google_lib(ctx, lib, file)
      product_ns = ctx.local_variable_get(:object).__product.prefix[1..-1]

      files = case lib
              when Libraries::NETWORK
                google_lib_network file, product_ns
              else
                google_lib_basic lib, file, product_ns
              end

      emit_libraries(ctx.local_variable_get(:output_folder), product_ns, files)

      File.join(*%w[google].concat([product_ns, lib, file]))
    end

    private

    def get_helper_file(file, remove_copyright_notice = true)
      content = IO.read(file)
      remove_copyright_notice ? strip_copyright_notice(content) : content
    end

    def strip_copyright_notice(content, comment_marker = '#')
      lines = content.split("\n")
      return content unless lines[0].include?('Copyright 20')
      lines = lines.drop(1) while lines[0].start_with?(comment_marker)
      lines = lines.drop(1) while lines[0].strip.empty?
      lines.join("\n")
    end

    def emit_libraries(output_folder, product_name, files)
      product_ns = product_name.downcase
      compile_file_list(output_folder, files,
                        product_ns: Google::StringUtils.camelize(product_name,
                                                                 :upper),
                        product_ns_dir: product_ns)
    end

    def google_lib_basic_files(file, product_ns, *google_root)
      {
        File.join(*google_root, product_ns, folder, "#{file}.rb") =>
          File.join('templates', folder, "#{file}.rb.erb"),
        File.join('spec', "#{folder}_#{file}_spec.rb") =>
          File.join('templates', folder, "#{file}_spec.rb.erb")
      }
    end

    def google_lib_network_files(file, product_ns, *google_root)
      files = {}
      folder = Libraries::NETWORK
      ['base', file].each do |f|
        files[File.join(*google_root, product_ns, folder, "#{f}.rb")] =
          File.join('templates', folder, "#{f}.rb.erb")
        base_spec = File.join('templates', folder, "#{f}_spec.rb.erb")
        files[File.join('spec', "#{folder}_#{f}_spec.rb")] = base_spec \
          if f != 'base' || File.exist?(base_spec)
      end
      %w[network_blocker network_blocker_spec].each do |f|
        files[File.join('spec', "#{f}.rb")] = File.join('templates', folder,
                                                        "#{f}.rb.erb")
      end
      files
    end
  end
end

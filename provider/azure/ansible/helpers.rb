require 'api/azure/type'

module Provider
  module Azure
    module Ansible
      module Helpers
        def is_resource_name?(property)
          property.parent.nil? && property.name == 'name'
        end

        def is_tags?(property)
          property.is_a? Api::Azure::Type::Tags
        end

        def is_tags_defined?(object)
          object.all_user_properties.any?{|p| is_tags?(p)}
        end

        def get_tags_property(object)
          object.all_user_properties.find{|p| is_tags?(p)}
        end

        def is_location?(property)
          property.parent.nil? && property.is_a?(Api::Azure::Type::Location)
        end

        def is_location_defined?(object)
          object.all_user_properties.any?{|p| is_location?(p)}
        end

        def is_resource_group?(property)
          property.parent.nil? && property.is_a?(Api::Azure::Type::ResourceGroupName)
        end

        def always_has_value?(property)
          property.required || !property.default_value.nil?
        end

        def word_wrap_for_yaml(lines, width = 160)
          wrapped = Array.new
          lines.each do |line|
            quoted = false
            first_line = true
            while line.length > width
              # Calculate leading spaces for the following lines
              striped = line.lstrip
              spaces = line.length - striped.length
              if first_line
                spaces += 2
                first_line = false
              end

              # Quote the whole line using quotation mark if not quoted
              quoted = true unless striped.start_with? '- '
              unless quoted
                line = line[0..spaces - 1] + '"' + line[spaces..-1] + '"' if line[spaces] != '"'
                quoted = true
              end

              # Find the last possible word-break character
              wb_index = find_word_break_index(line, /[ \t@=,;]/, width - spaces)
              wb_index = find_word_break_index(line, /[ \t@=,;\/]/, width - spaces) if wb_index.nil? || wb_index <= spaces
              break if wb_index.nil? || wb_index <= spaces

              # Break this line into two
              wb_char = line[wb_index]
              cur_line = line[0..wb_index - 1]
              cur_line += wb_char unless wb_char == ' '
              line = ' ' * spaces + line[wb_index + 1..-1]
              wrapped << cur_line
            end
            wrapped << line
          end
          wrapped
        end

        private

        def find_word_break_index(line, wb_chars, width)
          wb_index = line.rindex(wb_chars)
          while !wb_index.nil? && wb_index > width
            wb_index = line.rindex(wb_chars, wb_index - 1)
          end
          wb_index
        end
      end
    end
  end
end

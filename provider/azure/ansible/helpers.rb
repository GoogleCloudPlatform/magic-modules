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

        def word_wrap_for_yaml(lines, width = 150)
          wrapped = Array.new
          lines.each do |line|
            quoted = false
            while line.length > width
              # Calculate leading spaces for the following lines
              striped = line.lstrip
              spaces = line.length - striped.length
              spaces += 2 if striped.start_with? '- '

              # Quote the whole line using quotation mark if not quoted
              unless quoted
                line = line[0..spaces - 1] + '"' + line[spaces..-1] + '"' if line[spaces] != '"'
                quoted = true
              end

              # Find the last possible word-break character
              wb_index = -1
              wb_index_try = line.index(/[ \t@=,;]/)
              while !wb_index_try.nil? && wb_index_try < width
                wb_index = wb_index_try
                wb_index_try = line.index(/[ \t@=,;]/, wb_index_try + 1)
              end
              break if wb_index == -1

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
      end
    end
  end
end

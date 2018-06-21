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

require 'json'
require 'time'
require 'zlib'

module Provider
  module TestData
    # Arrays are of size 3 to ensure that only 3 network response YAML files
    # are needed for ResourceRefs
    MAX_ARRAY_SIZE = 3

    # A helper class to format data for testing purposes
    # rubocop:disable Metrics/ClassLength
    class Generator
      def comparator(property)
        generator = comparators.select do |type, _|
          property.class <= type
        end.values
        raise "Unknown property type: #{property.class}" if generator.empty?
        generator[0]
      end

      def value(for_type, property, seed)
        return property.default_value if property.default_value
        if for_type == Api::Type::Array
          for_type = [Api::Type::Array, property.item_type_class]
        end
        raise "Unknown property type: #{for_type} @ #{property}" \
          unless values.key?(for_type)
        values[for_type].call(property, seed)
      end

      # NameValues and Arrays require a size.
      # This function returns the size of a property of arbitrary size.
      # Use inside_array to manually specify that this object is being created
      # inside of an array (for resourceref counting purposes)
      def object_size(prop, seed, inside_array = false)
        size = (2 + integer_value(prop, seed) % 4)
        inside_array ||= prop.is_a?(Api::Type::Array)

        # This Nested Object may contain a resourceref.
        # If so, we need to ensure that there are only 3 objects in the array.
        # 3 is chosen because only 3 network response YAML files are written
        # per object.
        alt_size = 1 + (size % MAX_ARRAY_SIZE)

        # Array of ResourceRefs
        return alt_size if inside_array &&
                           prop.item_type.is_a?(Api::Type::ResourceRef)

        # Array of NestedObjects with ResourceRefs
        return alt_size if inside_array &&
                           prop.item_type.is_a?(Api::Type::NestedObject) &&
                           contains_resourcerefs?(prop.item_type)
        size
      end

      private

      def values
        {
          Api::Type::Boolean => method(:boolean_value),
          Api::Type::Constant => method(:constant_value),
          Api::Type::Double => method(:double_value),
          Api::Type::Enum => method(:enum_value),
          Api::Type::Integer => method(:integer_value),
          Api::Type::SelfLink => method(:selflink_value),
          Api::Type::FetchedExternal => method(:string_value),
          Api::Type::String => method(:string_value),
          Api::Type::Time => method(:time_value),
          Api::Type::Array::STRING_ARRAY_TYPE => method(:array_string),
          Api::Type::Array::NESTED_ARRAY_TYPE => method(:array_nested_cb),
          Api::Type::Array::RREF_ARRAY_TYPE => method(:array_rref_cb),
          Api::Type::NameValues => method(:name_values),
          Api::Type::ResourceRef => method(:resource_value),
          Api::Type::NestedObject => method(:nested_value)
        }
      end
      # rubocop:enable Metrics/MethodLength

      def comparators
        {
          Api::Type::Array => 'match_array',
          Api::Type::Boolean => 'is',
          Api::Type::Double => 'eq',
          Api::Type::Enum => 'eq',
          Api::Type::Integer => 'eq',
          Api::Type::NameValues => 'eq',
          Api::Type::NestedObject => 'eq',
          Api::Type::ResourceRef => 'eq',
          Api::Type::String => 'eq',
          Api::Type::Time => 'eq'
        }
      end

      def nested_value(prop, seed)
        Hash[prop.properties.map do |p|
          [p.out_name, if p.type == Api::Type::Integer
                         calc_integer_value(p, seed)
                       else
                         value(p.class, p, seed)
                       end]
        end]
      end

      def boolean_value(_prop, seed)
        (seed % 2).zero?
      end

      def constant_value(prop, seed)
        string_value(prop, seed).tr(' ', '').upcase
      end

      def enum_value(prop, seed)
        prop.values[seed % prop.values.length]
      end

      def calc_double_value(prop, seed)
        Zlib.crc32(prop.name).to_i * (seed + 1) * 0.67
      end

      def double_value(prop, seed)
        value = calc_double_value(prop, seed)
        [(value / 100).to_i, '.', (value % 100).to_i].join.to_f
      end

      def integer_value(prop, seed)
        calc_double_value(prop, seed).to_i
      end

      def string_value(prop, seed)
        "test #{prop.out_name}##{seed} data"
      end

      def selflink_value(prop, seed)
        name = Google::StringUtils.underscore(prop.resource.name)
        "selflink(resource(#{name},#{seed}))"
      end

      def time_value(prop, seed)
        Time.at(integer_value(prop, seed))
      end

      def resource_value(prop, seed)
        # Always use the first in the list for testing purposes.
        name = Google::StringUtils.underscore(prop.resources[0].resource_ref.name)
        "'resource(#{name},#{seed})'"
      end

      def name_values(prop, seed)
        size = object_size(prop, seed)
        Hash[(1..size).map do |i|
               [string_value(prop, seed + i),
                if i.even?
                  integer_value(prop, seed + i)
                else
                  string_value(prop, seed + i)
                end]
             end]
      end

      # Returns a callback to process string values
      def array_string(prop, seed)
        size = object_size(prop, seed)
        start = 'a'.ord + (integer_value(prop, seed) - 3) % 26
        start = 'z'.ord - size if start + size > 'z'.ord
        (1..size).map { |i| (start + i).chr * 2 }
      end

      def array_nested_cb(prop, seed)
        lambda do |&block|
          size = object_size(prop, seed)
          (1..size).map { |index| block.call index }
        end
      end

      def array_rref_cb(prop, seed)
        lambda do |hash|
          size = object_size(prop, seed)
          (0..size - 1).map do |index|
            if hash[:exported_values]
              # Return the exported value.
              imports = prop.item_type.imports.downcase
              resource = Google::StringUtils.underscore(prop.item_type.resource)
              "#{imports}(resource(#{resource},#{index}))"
            else
              resource_value(prop.item_type, index)
            end
          end
        end
      end

      # Returns true if a NestedObject property contains a resourceref
      # rubocop:disable Metrics/CyclomaticComplexity
      # rubocop:disable Metrics/PerceivedComplexity
      def contains_resourcerefs?(prop)
        return false unless prop.is_a? Api::Type::NestedObject
        prop.properties.each do |p|
          return true if p.is_a? Api::Type::ResourceRef

          if p.is_a? Api::Type::NestedObject
            return true if contains_resourcerefs?(p)
          elsif p.is_a? Api::Type::Array
            return true if p.item_type == 'Api::Type::ResourceRef'
            return true if contains_resourcerefs?(p.item_type)
          end
        end
        false
      end
      # rubocop:enable Metrics/CyclomaticComplexity
      # rubocop:enable Metrics/PerceivedComplexity
    end
    # rubocop:enable Metrics/ClassLength
  end
end

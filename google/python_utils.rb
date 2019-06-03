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

module Google
  # Functions to deal with the Python language.
  module PythonUtils
    # Represents a Unicode value.
    class UnicodeString
      attr_reader :value
      def initialize(val)
        @value = val
      end

      def to_s
        @value
      end
    end

    # Represents a line of Python code
    class PythonCode
      attr_reader :value
      def initialize(val)
        @value = val
      end

      def to_s
        @value
      end
    end
    # Prints literal in Python
    #
    # For instance, an int is printed as-is and a string is quoted:
    #  some_key: 80
    #  other_key: "foo"
    #
    # When a yaml file is parsed, strings specified with quotes or without
    # quotes becomes a ruby string without quotes unless you explicitly set
    # quotes in the string like "\"foo\"" which is not a pattern we want to
    # see in our yaml config files.
<<<<<<< HEAD
    def python_literal(value)
      if value.is_a?(String)
        "'#{value}'"
      elsif value.is_a?(Symbol)
        "'#{value.to_s.underscore}'"
=======
    # options:
    # use_hash_brackets - use {} instead of dict() notation.
    def python_literal(value, **opts)
      if value.is_a?(String) || value.is_a?(Symbol)
        "'#{value}'"
      elsif value.is_a?(PythonCode)
        value
>>>>>>> master
      elsif value.is_a?(Numeric)
        value.to_s
      elsif value.is_a?(Hash) && opts[:use_hash_brackets]
        hash_format(value)
      elsif value.is_a?(Hash)
        "dict(#{value.map { |k, v| "#{k}=#{python_literal(v)}" if v }.compact.join(', ')})"
      elsif value.is_a?(Array)
        values = value.map { |x| python_literal(x) }
        array_format(values)
      elsif value == true
        'True'
      elsif value == false
        'False'
      else
        raise "Unsupported Python literal #{value}"
      end
    end

    # Generates a method declaration with function name `name` and args `args`
    # Arguments may have nils and will be ignored.
    def method_decl(name, args)
      "def #{name}(#{args.compact.join(', ')}):"
    end

    # Generates a method call to function name `name` and args `args`
    # Arguments may have nils and will be ignored.
    def method_call(name, args)
      "#{name}(#{args.compact.join(', ')})"
    end

<<<<<<< HEAD
    def python_variable_name(property, sdk_op_def)
      sdk_ref = get_applicable_reference(property.azure_sdk_references, sdk_op_def.request)
      return property.out_name.underscore if sdk_ref.nil?
      python_var = get_sdk_typedef_by_references(property.azure_sdk_references, sdk_op_def.request).python_variable_name
      return property.out_name.underscore if python_var.nil?
      python_var
=======
    private

    def array_format(values)
      '[' + values.join(', ') + ']'
    end

    def hash_format(value)
      hash_vals = value.map do |k, v|
        next if v.nil?

        if k.is_a?(UnicodeString)
          "u'#{k}': #{python_literal(v)}"
        else
          "'#{k}': #{python_literal(v)}"
        end
      end.compact
      "{ #{hash_vals.join(',')} }"
>>>>>>> master
    end
  end
end

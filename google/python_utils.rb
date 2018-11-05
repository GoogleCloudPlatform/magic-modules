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
    def python_literal(value, spaces=0)
      if value.is_a?(String) || value.is_a?(Symbol)
        "'#{value}'"
      elsif value.is_a?(Numeric)
        value.to_s
      elsif value.is_a?(Array)
        format([
          ["[#{value.map { |x| python_literal(x) }.join(', ')}]"],
          ["[#{value.map { |x| python_literal(x) }.join("\n") }]"],
        ], 0, spaces)
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
  end
end

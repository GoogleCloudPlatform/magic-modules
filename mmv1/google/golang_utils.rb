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
  # Functions to deal with the Go language.
  module GolangUtils
    # Prints literal in go.
    #
    # For instance, an int is printed as-is and a string is quoted:
    #  some_key: 80
    #  other_key: "foo"
    #
    # When a yaml file is parsed, strings specified with quotes or without
    # quotes becomes a ruby string without quotes unless you explicitly set
    # quotes in the string like "\"foo\"" which is not a pattern we want to
    # see in our yaml config files.
    def go_literal(value)
      if value.is_a?(String) || value.is_a?(Symbol)
        "\"#{value}\""
      elsif value.is_a?(Numeric)
        value.to_s
      elsif value.is_a?(Array) && value.all? { |v| v.is_a?(String) || v.is_a?(Symbol) }
        "[]string{#{value.map(&method(:go_literal)).join(', ')}}"
      elsif value.is_a?(TrueClass) || value.is_a?(FalseClass)
        value.to_s
      else
        raise "Unsupported go literal #{value}"
      end
    end
  end
end

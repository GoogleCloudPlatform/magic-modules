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
  # Functions to deal with the Ruby language.
  module RubyUtils
    # Prints literal in Ruby
    #
    # For instance, an int is printed as-is and a string is quoted:
    #  some_key: 80
    #  other_key: "foo"
    #
    # When a yaml file is parsed, strings specified with quotes or without
    # quotes becomes a ruby string without quotes unless you explicitly set
    # quotes in the string like "\"foo\"" which is not a pattern we want to
    # see in our yaml config files.
    def ruby_literal(value)
      if value.is_a?(String) || value.is_a?(Symbol)
        "'#{value}'"
      elsif value.is_a?(Numeric)
        value.to_s
      else
        raise "Unsupported Ruby literal #{value}"
      end
    end

    def method_decl(name, args)
      ["def #{name}", ("(#{args.compact.join(', ')})" unless args.empty?)].compact.join
    end
    def method_call(name, args, indent = 0)
      args = args.compact
      format([
               # All on one line.
               [
                 [name.to_s, ("(#{args.join(', ')})" unless args.empty?)].compact.join
               ],
               # All but first on one line.
               [
                 [name.to_s, ("(#{args[0..-1].join(', ')}" unless args.empty?)].compact.join,
                 "#{indent(args.last, indent + name.length + 2)})"
               ],
               # All on separate lines.
               [
                 "#{name}(#{args.first},",
                 indent_list(args.slice(1..-2), indent + name.length + 1, true),
                 indent("#{args.last})", indent + name.length + 1)
               ]
             ], 0, indent)
    end
    # rubocop:enable Metrics/AbcSize
  end
end

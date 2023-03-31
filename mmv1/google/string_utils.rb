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
  # Helper class to process and mutate strings.
  class StringUtils
    # Converts string from camel case to underscore
    def self.underscore(source)
      source.gsub(/::/, '/')
            .gsub(/([A-Z]+)([A-Z][a-z])/, '\1_\2')
            .gsub(/([a-z\d])([A-Z])/, '\1_\2')
            .tr('-', '_')
            .tr('.', '_')
            .downcase
    end

    # Converts from PascalCase to Space Separated
    def self.space_separated(source)
      tmp = source.gsub(/([A-Z]+)([A-Z][a-z])/, '\1 \2')
                  .gsub(/([a-z\d])([A-Z])/, '\1 \2')
                  .downcase
      tmp[0].upcase.concat(tmp[1..])
    end

    # rubocop:disable Style/SafeNavigation # support Ruby < 2.3.0
    def self.symbolize(key)
      key.to_sym unless key.nil?
    end
    # rubocop:enable Style/SafeNavigation

    # Returns all the characters up until the period (.) or returns text
    # unchanged if there is no period.
    def self.first_sentence(text)
      period_pos = text.index(/[.?!]/)
      return text if period_pos.nil?

      text[0, period_pos + 1]
    end

    # Returns the plural form of a word
    def self.plural(source)
      # policies -> policies
      # indices -> indices
      return source if source.end_with?('ies') || source.end_with?('es')

      # index -> indices
      return "#{source.gsub(/ex$/, '')}ices" if source.end_with?('ex')

      # mesh -> meshes
      return "#{source}es" if source.end_with?('esh')

      # key -> keys
      # gateway -> gateways
      return "#{source}s" if source.end_with?('ey') || source.end_with?('ay')

      # policy -> policies
      return "#{source.gsub(/y$/, '')}ies" if source.end_with?('y')

      "#{source}s"
    end
  end
end

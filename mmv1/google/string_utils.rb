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

    # Converts a string to space-separated capitalized words
    def self.title(source)
      ss = space_separated(source)
      ss.gsub(/\b(?<!\w['’`()])[a-z]/, &:capitalize)
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

    # ActiveSupport::Inflector code
    # Converts strings to UpperCamelCase.
    # If the +uppercase_first_letter+ parameter is set to false, then produces
    # lowerCamelCase.
    #
    # Also converts '/' to '::' which is useful for converting
    # paths to namespaces.
    #
    #   camelize('active_model')                # => "ActiveModel"
    #   camelize('active_model', false)         # => "activeModel"
    #   camelize('active_model/errors')         # => "ActiveModel::Errors"
    #   camelize('active_model/errors', false)  # => "activeModel::Errors"
    #
    # As a rule of thumb you can think of +camelize+ as the inverse of
    # #underscore, though there are cases where that does not hold:
    #
    #   camelize(underscore('SSLError'))        # => "SslError"
    def self.camelize(term, uppercase_first_letter = true)
      # patched in to this fn
      define_acronym_regex_patterns
      inflections = {"tpu": "TPU", "vpc": "VPC"}

      string = term.to_s
      # String#camelize takes a symbol (:upper or :lower), so here we also support :lower to keep the methods consistent.
      if !uppercase_first_letter || uppercase_first_letter == :lower
        string = string.sub(@acronyms_camelize_regex) { |match| match.downcase! || match }
      else
        string = string.sub(/^[a-z\d]*/) { |match| inflections[match] || match.capitalize! || match }
      end
      string.gsub!(/(?:_|(\/))([a-z\d]*)/i) do
        word = $2
        substituted = inflections[word] || word.capitalize! || word
        $1 ? "::#{substituted}" : substituted
      end
      string
    end

    def self.define_acronym_regex_patterns
      @acronyms = {}
      @acronym_regex             = @acronyms.empty? ? /(?=a)b/ : /#{@acronyms.values.join("|")}/
      @acronyms_camelize_regex   = /^(?:#{@acronym_regex}(?=\b|[A-Z_])|\w)/
      @acronyms_underscore_regex = /(?:(?<=([A-Za-z\d]))|\b)(#{@acronym_regex})(?=\b|[^a-z])/
    end
  end
end

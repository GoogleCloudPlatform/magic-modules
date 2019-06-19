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

require 'compile/core'
require 'provider/config'
require 'provider/core'
require 'redcarpet'

module Provider
  module Ansible
    # Responsible for building out YAML documentation
    # from Markdown
    # This is primarily done for parsing descriptions.
    module Markdown
      DELIMITER = '$$$'.freeze

      def description(text)
        # Description texts should have leading + trailing spaces
        # removed so that they are not mistaken for code blocks by the
        # Markdown parser.
        Redcarpet::Markdown.new(AnsibleDescriptionRender)
                           .render(text.strip.squeeze("\n"))
                           .split(DELIMITER)
                           .map(&:strip)
      end

      # This is a rendering class that takes in
      # a Markdown description and returns an
      # array of strings. This is used exclusively for
      # description documentation.
      #
      # Redcarpet will return a String (because that's the expectation of markdown).
      # Ansible wants an array of strings, so this class will return a single string
      # with the '$$$' character denoting where the string should be split.
      class AnsibleDescriptionRender < Redcarpet::Render::Base
        LIST_DELIMITER = '%%%'.freeze
        # Returns a paragraph with delimiters showing where it should be split.
        def paragraph(text)
          text.split(".\n").map do |paragraph|
            paragraph += '.' unless paragraph.end_with?('.')
            paragraph = format_url(paragraph)
            paragraph.tr("\n", ' ').squeeze(' ')
          end.join(DELIMITER)
        end

        def codespan(code)
          "\"#{code}\""
        end

        def normal_text(text)
          text
        end

        def link(link, _title, content)
          if content
            "L(#{content},#{link})"
          else
            "U(#{link})"
          end
        end

        def list(content, _list_type)
          content.split(LIST_DELIMITER).join(', ')
        end

        # List items come first. We have to place special delimiters
        # because all of the list strings are joined together before
        # list() is called.
        def list_item(text, _list_type)
          "#{text.sub("\n", '')}#{LIST_DELIMITER}"
        end

        private

        # Find URLs and surround with U()
        # If there's a period at the end of the URL, make sure the
        # period is outside of the ()
        def format_url(paragraph)
          paragraph.gsub(%r{
            https?:\/\/(?:www\.|(?!www))[a-zA-Z0-9]
            [a-zA-Z0-9-]+[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]+
            [a-zA-Z0-9]\.[^\s]{2,}|https?:\/\/(?:www\.|(?!www))
            [a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9]\.[^\s]{2,}
          }x, 'U(\\0)').gsub('.)', ').')
        end
      end
    end
  end
end

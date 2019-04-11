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

module Google
  # Functions to deal with the HCL language.
  module HclUtils
    include Compile::Core

    # Takes a structure of the following format and converts to HCL.
    # "resource" => [{
    #   "google_compute_autoscalar" => [{
    #      "my-autoscalar" => {
    #        "field" => true
    #      }
    #   }]
    # }]
    def hcl(dictionary)
      raise 'Only accepts dictionary' unless dictionary.is_a?(Hash)
      raise 'Only accepts dictinonaries with one key' \
        unless dictionary.keys.length == 1

      hcl_type = dictionary.keys.first
      gcp_type = dictionary[hcl_type].first.keys.first
      name = dictionary[hcl_type].first[gcp_type].first.keys.first

      values = dictionary[hcl_type].first[gcp_type].first[name].first

      # Find longest name to understand how spacing should work.
      longest = values.keys.max_by(&:length).length
      values = values.map { |k, v| "#{k}#{' ' * (longest - k.length)} = #{hcl_literal(v)}" }

      [
        "#{hcl_type} \"#{gcp_type}\" \"#{name}\" {",
        values.map { |k| "  #{k}" },
        '}'
      ].flatten.join("\n")
    end

    def hcl_literal(literal)
      if literal.is_a?(String)
        "\"#{literal}\""
      elsif literal.is_a?(Array)
        "[#{literal.map { |v| hcl_literal(v) }.join(',')}]"
      else
        raise "HCL type: #{literal.class} not supported"
      end
    end
  end
end

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

module Provider
  module Ansible
    # Responsible for building out the resource_to_request and
    # request_from_hash methods.
    module Request
      # Takes in a list of properties and outputs a python hash that takes
      # in a module and outputs a formatted JSON request.
      def request_properties(properties, indent = 4)
        indent_list(properties.map { |prop| request_property(prop) },
                    indent)
      end

      def response_properties(properties, indent = 8)
        indent_list(properties.map { |prop| response_property(prop) },
                    indent)
      end

      private

      def request_property(prop)
        [
         "#{unicode_string(prop.field_name)}:",
         "module.params.get(#{quote_string(prop.out_name)})"
        ].join(' ')
      end

      def response_property(prop)
        [
         "#{unicode_string(prop.field_name)}:",
         "response.get(#{unicode_string(prop.name)})"
        ].join(' ')
      end
    end
  end
end

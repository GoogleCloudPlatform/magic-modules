# Copyright 2019 Google Inc.
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

require 'api/object'
require 'google/string_utils'

module Api
  # An object available in the product
  class Resource < Api::Object::Named
    # Represents a response from the API that returns a list of objects.
    class ResponseList < Api::Object
      attr_reader :kind
      attr_reader :items

      def validate
        super

        check :items, default: 'items', type: ::String, required: true
        check :kind, type: ::String
      end

      def kind?
        !@kind.nil?
      end
    end
  end
end

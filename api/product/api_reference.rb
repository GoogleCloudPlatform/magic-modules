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

module Api
  class Product < Api::Object::Named
    # Represents any APIs that are required to be enabled to use this product
    class ApiReference < Api::Object
      attr_reader :name
      attr_reader :url

      def validate
        super
        check :name, type: String, required: true
        check :url, type: String, required: true
      end
    end
  end
end

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
    # Query information for finding resource nested in an returned API object
    # i.e. fine-grained resources
    class NestedQuery < Api::Object
      # A list of keys to traverse in order.
      # i.e. backendBucket --> cdnPolicy.signedUrlKeyNames
      # should be ["cdnPolicy", "signedUrlKeyNames"]
      attr_reader :keys

      # If true, we expect the the nested list to be
      # a list of IDs for the nested resource, rather
      # than a list of nested resource objects
      # i.e. backendBucket.cdnPolicy.signedUrlKeyNames is a list of key names
      # rather than a list of the actual key objects
      attr_reader :is_list_of_ids

      # This is used by Ansible, but may not be necessary.
      attr_reader :kind

      def validate
        super

        check :keys, type: Array, item_type: String, required: true
        check :is_list_of_ids, type: :boolean, default: false

        check :kind, type: String
      end
    end
  end
end

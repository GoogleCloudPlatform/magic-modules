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
    # Metadata for resources that are nested within a parent resource, as
    # a list of resources or single object within the parent.
    # e.g. Fine-grained resources
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

      # If true, the resource is created/updated/deleted by patching
      # the parent resource and appropriate encoders/update_encoders/pre_delete
      # custom code will be included automatically. Only use if parent resource
      # does not have a separate endpoint (set as create/delete/update_urls)
      # for updating this resource.
      # The resulting encoded data will be mapped as
      # {
      #  keys[-1] : list_of_objects
      # }
      attr_reader :modify_by_patch

      # Nested resources generally don't have a kind field.
      # This is used as a (potentially unnecessary) placeholder by Ansible
      attr_reader :kind

      def validate
        super

        check :keys, type: Array, item_type: String, required: true
        check :is_list_of_ids, type: :boolean, default: false
        check :modify_by_patch, type: :boolean, default: false

        check :kind, type: String
      end
    end
  end
end

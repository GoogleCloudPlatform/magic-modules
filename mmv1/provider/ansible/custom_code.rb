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

require 'google/python_utils'

module Provider
  module Ansible
    # Holds all of the custom code injection points in Ansible.
    class CustomCode < Api::Object
      attr_reader :create
      attr_reader :delete
      attr_reader :update

      # Code that is run before every action (create, delete, update, fetch)
      attr_reader :pre_action

      # Code that is run after a create.
      attr_reader :post_create

      # Code that is run after every action (create, delete, update, fetch)
      attr_reader :post_action

      # Custom function that takes in a requests 'Response' object
      # and returns a JSON body or errors out.
      attr_reader :return_if_object

      # Says if a custom unwrap_resource function is being used.
      attr_reader :unwrap_resource

      # Says if different resource_to_request body sent for create calls.
      attr_reader :custom_create_resource

      # Says if different resource_to_request body sent for update calls.
      attr_reader :custom_update_resource

      # Custom function to get async url.
      attr_reader :custom_async_function

      def validate
        check :create, type: String
        check :custom_create_resource, type: :boolean
        check :custom_update_resource, type: :boolean
        check :delete, type: String
        check :pre_action, type: String
        check :post_action, type: String
        check :post_create, type: String
        check :return_if_object, type: String
        check :update, type: String
        check :unwrap_resource, type: :boolean
        check :custom_async_function, type: String
      end
    end
  end
end

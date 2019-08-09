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
    # Information about the IAM policy for this resource
    # Several GCP resources have IAM policies that are scoped to
    # and accessed via their parent resource
    # See: https://cloud.google.com/iam/docs/overview
    class IamPolicy < Api::Object
      # boolean of if this binding should be generated
      attr_reader :exclude

      # Character that separates resource identifier from method call in URL
      # For example, PubSub subscription uses {resource}:getIamPolicy
      # While Compute subnetwork uses {resource}/getIamPolicy
      attr_reader :method_name_separator

      # boolean of if a parent resource is required to apply this IAM policy to
      # Some IAM policies are applied across multiple resources rather than a
      # specific resource
      attr_reader :parent_resource_required

      # The terraform type of the parent resource if it is not the same as the
      # IAM resource. The IAP product needs these as its IAM policies refer
      # to compute resources
      attr_reader :parent_resource_type

      # Some resources allow retrieving the IAM policy with GET requests,
      # others expect POST requests
      attr_reader :fetch_iam_policy_verb

      # Certain resources allow different sets of roles to be set with IAM policies
      # This is a role that is acceptable for the given IAM policy resource for use in tests
      attr_reader :allowed_iam_role

      # The code that is rendered within qualify{{resource.name}}Url. This is used for IAP
      # that supports different URL endpoints based values set in the config
      attr_reader :custom_url_qualifier

      # The code that is rendered within qualify{{resource.name}}Url. This is used for IAP
      # that supports different URL endpoints based values set in the config
      attr_reader :custom_id_function

      # Allows for optional properties to be specified that are allowed in the IAM policy
      # these are needed for IAP that switches URLs based on the presence of these properties
      attr_reader :optional_properties

      def validate
        super

        check :exclude, type: :boolean, default: false
        check :method_name_separator, type: String, default: '/'
        check :parent_resource_required, type: :boolean, default: true
        check :parent_resource_type, type: String
        check :fetch_iam_policy_verb, type: String, default: 'GET'
        check :allowed_iam_role, type: String, default: 'roles/editor'
        check :custom_url_qualifier, type: String
        check :custom_id_function, type: String
        check :optional_properties, type: Array, item_type: Api::Type, default: []
      end
    end
  end
end

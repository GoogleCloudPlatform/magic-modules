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

      # Certain resources need an attribute other than "id" from their parent resource
      # Especially when a parent is not the same type as the IAM resource
      attr_reader :parent_resource_attribute

      # If the IAM resource test needs a new project to be created, this is the name of the project
      attr_reader :test_project_name

      # Resource name may need a custom diff suppress function. Default is to use
      # compareSelfLinkOrResourceName
      attr_reader :custom_diff_suppress

      # Some resources (IAP) use fields named differently from the parent resource.
      # We need to use the parent's attributes to create an IAM policy, but they may not be
      # named as the IAM IAM resource expects.
      # This allows us to specify a file (relative to MM root) containing a partial terraform
      # config with the test/example attributes of the IAM resource.
      attr_reader :example_config_body

      # How the API supports IAM conditions
      attr_reader :iam_conditions_request_type

      def validate
        super

        check :exclude, type: :boolean, default: false
        check :method_name_separator, type: String, default: '/'
        check :parent_resource_type, type: String
        check :fetch_iam_policy_verb, type: Symbol, default: :GET, allowed: %i[GET POST]
        check :allowed_iam_role, type: String, default: 'roles/viewer'
        check :parent_resource_attribute, type: String, default: 'id'
        check :test_project_name, type: String
        check :iam_conditions_request_type, type: Symbol, allowed: %i[REQUEST_BODY QUERY_PARAM]
        check(
          :example_config_body,
          type: String, default: 'templates/terraform/iam/iam_attributes.tf.erb'
        )
      end
    end
  end
end

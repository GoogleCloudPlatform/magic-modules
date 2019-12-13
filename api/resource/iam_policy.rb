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

      # Last part of URL for fetching IAM policy.
      attr_reader :fetch_iam_policy_method

      # Some resources allow setting the IAM policy with POST requests,
      # others expect PUT requests
      attr_reader :set_iam_policy_verb

      # Last part of URL for setting IAM policy.
      attr_reader :set_iam_policy_method

      # Whether the policy JSON is contained inside of a 'policy' object.
      attr_reader :wrapped_policy_obj

      # Certain resources allow different sets of roles to be set with IAM policies
      # This is a role that is acceptable for the given IAM policy resource for use in tests
      attr_reader :allowed_iam_role

      # This is a role that grants create/read/delete for the parent resource for use in tests.
      # If set, the test runner will receive a binding to this role in _policy tests in order to
      # avoid getting locked out of the resource.
      attr_reader :admin_iam_role

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

      # Allows us to override the base_url of the resource. This is required for Cloud Run as the
      # IAM resources use an entirely different base URL from the actual resource
      attr_reader :base_url

      # Allows us to override the import format of the resource. Useful for Cloud Run where we need
      # variables that are outside of the base_url qualifiers.
      attr_reader :import_format

      # This code replaces the portion of code that manipulates the import format after getting
      # the import id qualifiers.  It's useful in checking the the qualifiers were parsed correctly
      # in cases that they may not have been (for example, if the name of the resource has a forward
      # slash in it)
      attr_reader :custom_import

      def validate
        super

        check :exclude, type: :boolean, default: false
        check :method_name_separator, type: String, default: '/'
        check :parent_resource_type, type: String
        check :fetch_iam_policy_verb, type: Symbol, default: :GET, allowed: %i[GET POST]
        check :fetch_iam_policy_method, type: String, default: 'getIamPolicy'
        check :set_iam_policy_verb, type: Symbol, default: :POST, allowed: %i[POST PUT]
        check :set_iam_policy_method, type: String, default: 'setIamPolicy'
        check :wrapped_policy_obj, type: :boolean, default: true
        check :allowed_iam_role, type: String, default: 'roles/viewer'
        check :admin_iam_role, type: String
        check :parent_resource_attribute, type: String, default: 'id'
        check :test_project_name, type: String
        check :iam_conditions_request_type, type: Symbol, allowed: %i[REQUEST_BODY QUERY_PARAM]
        check :base_url, type: String
        check :import_format, type: Array, item_type: String
        check(
          :example_config_body,
          type: String, default: 'templates/terraform/iam/iam_attributes.tf.erb'
        )
      end
    end
  end
end

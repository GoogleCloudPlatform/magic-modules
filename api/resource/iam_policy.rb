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

      def validate
        super

        check :exclude, type: :boolean, default: false
        check :method_name_separator, type: String, default: '/'
      end
    end
  end
end

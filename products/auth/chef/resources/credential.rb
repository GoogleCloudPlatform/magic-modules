# Copyright 2018 Google Inc.
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

require 'chef/resource'
require_relative 'google/authorization'

module Google
  module Auth
    # Chef Resource to authenticate to GCP
    class Credential < Chef::Resource::LWRPBase
      resource_name :gauth_credential

      default_action :nothing

      property :name, String, identity: true, desired_state: false
      property :path, String, desired_state: false
      property :scopes, Array, desired_state: false,
                               default: ['https://www.googleapis.com/auth/compute']
      property :__auth, ::Google::Authorization, desired_state: false

      action :serviceaccount do
        if new_resource.path.nil?
          raise ["Missing 'path' parameter in",
                 "gauth_credential[#{new_resource.name}]"].join(' ')
        end

        if new_resource.scopes.nil?
          raise ["Missing 'scopes' parameter in",
                 "gauth_credential[#{new_resource.name}]"].join(' ')
        end

        # TODO: How do define a private property, or better, how to store a
        # variable only for this instance
        new_resource.__auth ::Google::Authorization.new.for!(
          new_resource.scopes
        ).from_service_account_json!(
          new_resource.path
        )
      end

      action :defaultuseraccount do
        __auth ::Google::Authorization.new.from_user_credential!
      end

      action :nothing do
        raise 'An action for a provider is required: service_account'
      end

      def authorization
        raise "Failed to authenticate gauth_credential[#{name}]" if __auth.nil?
        __auth
      end
    end
  end
end

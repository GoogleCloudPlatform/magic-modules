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

require 'google/sql/network/post'

module Google
  module Sql
    module Api
      # A helper class to provide access to (some) Google Cloud SQL API.
      class User
        def initialize(name, instance, project, cred)
          @name = name
          @instance = instance
          @project = project
          @cred = cred
        end

        # TODO(ody): This function is the same as gsql_user if you ignore
        # idempotency. Is there a way to create the resource and call create on
        # it?
        def passwd(host, password)
          # TODO(ody): If not the above this is resource_to_request. Can we call
          # from the provider?
          request = { name: @name, host: host, password: password }.to_json
          post_request = ::Google::Sql::Network::Post.new(
            gsql_user_collection, @cred, 'application/json', request
          )
          response = JSON.parse(post_request.send.body)
          raise Puppet::Error, response['error']['errors'][0]['message'] \
            if response['error']

          # TODO(nelsonjr): Make this function wait for the operation to
          # complete
        end

        private

        def gsql_user_collection
          Puppet::Type.type(:gsql_user).provider(:google).collection(
            instance: @instance, project: @project
          )
        end
      end
    end
  end
end

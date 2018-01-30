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
      class Instance
        def initialize(name, project, cred)
          @name = name
          @project = project
          @cred = cred
        end

        # TODO(nelsonjr): Make this function wait for the operation to complete
        def clone(target)
          clone_request = ::Google::Sql::Network::Post.new(
            gsql_instance_clone, @cred, 'application/json', {
              cloneContext: { kind: 'sql#cloneContext',
                              destinationInstanceName: target }
            }.to_json
          )
          response = JSON.parse(clone_request.send.body)
          raise Puppet::Error, response['error']['errors'][0]['message'] \
             if response['error']
        end

        private

        def gsql_instance_clone
          URI.parse(
            format(
              '%<self_link>s/%<method>s',
              self_link: Puppet::Type.type(:gsql_instance).provider(:google)
                                     .self_link(name: @name, project: @project),
              method: 'clone'
            )
          )
        end
      end
    end
  end
end

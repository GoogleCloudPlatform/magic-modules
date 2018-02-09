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

require 'google/compute/network/post'

module Google
  module Compute
    module Api
      # A helper class to provide access to (some) Google Compute Engine API.
      class Instance
        def initialize(name, zone, project, cred)
          @name = name
          @zone = zone
          @project = project
          @cred = cred
        end

        # TODO(nelsonjr): Make this function wait for the operation to complete
        def reset
          reset_request = ::Google::Compute::Network::Post.new(
            gcompute_instance_reset, @cred, 'application/json', {}.to_json
          )
          response = JSON.parse(reset_request.send.body)
          raise Puppet::Error, response['error']['errors'][0]['message'] \
             if response['error']
        end

        private

        def gcompute_instance_reset
          URI.parse(
            format(
              '%<self_link>s/%<method>s',
              self_link: Puppet::Type.type(:gcompute_instance).provider(:google)
                                     .self_link(name: @name, zone: @zone,
                                                project: @project),
              method: 'reset'
            )
          )
        end
      end
    end
  end
end

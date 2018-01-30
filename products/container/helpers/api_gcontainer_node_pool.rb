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

require 'google/container/network/post'

module Google
  module Container
    module Api
      # A helper class to provide access to (some) Google Container Engine API.
      class NodePool
        def initialize(name, cluster, zone, project, cred)
          @name = name
          @cluster = cluster
          @zone = zone
          @project = project
          @cred = cred
        end

        # TODO(nelsonjr): Implement this as gcontainer_node_pool { size }
        # TODO(nelsonjr): Make this function wait for the operation to complete
        # TODO(nelsonjr): Add error checking on response on this task
        #                 (ditto on all Bolt tasks)
        def resize(size)
          resize_request = ::Google::Container::Network::Post.new(
            gcontainer_node_pool_resize, @cred, 'application/json',
            { 'nodeCount' => size }.to_json
          )
          response = JSON.parse(resize_request.send.body)
          raise Puppet::Error, response['error']['errors'][0]['message'] \
             if response['error']
        end

        private

        def gcontainer_node_pool_resize
          URI.parse(
            format(
              '%<self_link>s/%<method>s',
              self_link: Puppet::Type.type(:gcontainer_node_pool)
                                     .provider(:google)
                                     .self_link(name: @name, zone: @zone,
                                                project: @project,
                                                cluster: @cluster),
              method: 'setSize'
            )
          )
        end
      end
    end
  end
end

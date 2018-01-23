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

require 'google/pubsub/network/post'

module Google
  module Pubsub
    module Api
      # A helper class to provide access to (some) Google Container Engine API.
      class Topic
        def initialize(topic, project, cred)
          @topic = topic
          @project = project
          @cred = cred
        end

        # TODO(nelsonjr): Implement this as gcontainer_node_pool { size }
        # TODO(nelsonjr): Make this function wait for the operation to complete
        # TODO(nelsonjr): Add error checking on response on this task
        #                 (ditto on all Bolt tasks)
        def publish(message, attributes)
          message = {
            messages: [{
              attributes: attributes,
              data: Base64.encode64(message).strip
            }]
          }

          request = ::Google::Pubsub::Network::Post.new(
            gtopic_publish, @cred, 'application/json',
            message.to_json
          )

          response = JSON.parse(request.send.body)
          raise Puppet::Error, response['error']['message'] if response['error']
          response
        end

        private

        def gtopic_publish
          URI.parse(
            format(
              '%s:%s',
              Puppet::Type.type(:gpubsub_topic).provider(:google)
                          .self_link(name: @topic, project: @project,
                                     cluster: @cluster),
              'publish'
            )
          )
        end
      end
    end
  end
end

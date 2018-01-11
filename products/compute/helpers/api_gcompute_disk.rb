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
      class Disk
        def initialize(name, zone, project, cred)
          @name = name
          @zone = zone
          @project = project
          @cred = cred
        end

        # TODO(nelsonjr): Implement this as gcompute_disk_snapshot { }
        # TODO(nelsonjr): Make this function wait for the operation to complete
        def snapshot(target, properties = {})
          snapshot_request = ::Google::Compute::Network::Post.new(
            gcompute_disk_snapshot, @cred, 'application/json',
            # Ordering of 'kind' must be moved so that testing
            # expectations do not fail.
            { kind: properties[:kind], name: target }.merge(properties).to_json
          )
          response = JSON.parse(snapshot_request.send.body)
          raise Puppet::Error, response['error']['errors'][0]['message'] \
             if response['error']
        end

        private

        def gcompute_disk_snapshot
          URI.parse(
            format(
              '%s/%s',
              Puppet::Type.type(:gcompute_disk).provider(:google).self_link(
                name: @name, zone: @zone, project: @project
              ), 'createSnapshot'
            )
          )
        end
      end
    end
  end
end

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

module Google
  # A helper class to aggregate all resources when multiple providers produce
  # them at the same time. This singleton is used to collect different
  # "flavors" of resources, managed by different providers, when the client
  # does not care which provider created it.
  #
  # For example when authenticating requests, as long as the resource honors
  # the "authenticate" API which provider, or which parameters it needed to
  # create itself matters little to the consumer of the authenticator.
  class ObjectStore
    include Singleton

    attr_reader :resources

    def initialize
      @resources = {}
    end

    # Adds an instance of the resource to the global collection.
    def add(type, resource)
      Puppet.debug "Registering resource #{resource}"
      @resources[type] = [] if @resources[type].nil?
      @resources[type] << resource
    end

    def [](type)
      if @resources[type].nil?
        []
      else
        @resources[type]
      end
    end
  end
end

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

require 'api/object'

module Provider
  # Override to an Api::Resource in api.yaml
  class ResourceOverride < ::Hash
    # Used for testing.
    def initialize(hash)
      hash.each { |k, v| self[k] = v }
    end

    def [](key)
      if key.to_s[0] == '@'
        dig key.to_s[1..-1]
      else
        dig key.to_s
      end
    end
  end
end

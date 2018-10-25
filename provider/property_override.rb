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
require 'api/type'

module Provider
  # Override a resource property (Api::Type) in api.yaml
  # TODO(rosbo): Shared common logic with ResourceOverride via a base class.
  class PropertyOverride < Api::Object
    # Used for testing.
    def initialize(hash)
      hash.each { |k, v| instance_variable_set("@#{k}", v) }
    end

    def [](key)
      if key[0] == '@'
        instance_variable_get(key)
      else
        instance_variable_get("@#{key}")
      end
    end
  end
end

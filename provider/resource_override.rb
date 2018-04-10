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
  class ResourceOverride < Api::Object
    attr_reader :description

    def consume_api(api)
      @__api = api
    end

    def validate
      super

      check_optional_property :description, String
    end

    # Apply this override to the given instance of Api::Resource
    def apply(api_resource)
      extend_string api_resource, :description, @description
    end

    # Replace the `object_key` instance variable on `object` by the
    # `override_val`. If `override_val` includes the tag '{{<object_key>}}',
    # this tag will be substituted by the object value.
    def extend_string(object, object_key, override_val)
      return if override_val.nil?

      object_val = object.send object_key
      new_val = override_val.gsub "{{#{object_key}}}", object_val

      object.instance_variable_set("@#{object_key}", new_val)
    end
  end
end

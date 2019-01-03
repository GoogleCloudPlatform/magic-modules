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
require 'provider/abstract_core'
require 'provider/property_override'

module Provider
  module Ansible
    # Collection of fields allowed in the PropertyOverride section for
    # Ansible. All fields should be `attr_reader :<property>`
    module OverrideFields
      def self.attributes
        [
          :aliases,
          :version_added
        ]
      end

      attr_reader(*self.attributes)
    end

    module PropertyOverrideSharedCode
      def validate
        super

        check_optional_property :aliases, ::Array
        check_optional_property :version_added, ::String
      end
    end

    # Ansible-specific overrides to api.yaml.
    class PropertyOverride < Provider::PropertyOverride
      include OverrideFields
      include PropertyOverrideSharedCode

      private

      def overriden
        Provider::Ansible::OverrideFields
      end
    end
  end
end

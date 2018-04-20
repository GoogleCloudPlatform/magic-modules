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
  class Puppet < Provider::Core
    # Collection of fields allowed in the PropertyOverride section for
    # Puppet. All fields should be `attr_reader :<property>`
    module OverrideFields
      attr_reader :validation # Adds a validate block to type
    end

    # Puppet-specific overrides to api.yaml.
    class PropertyOverride < Provider::PropertyOverride
      include OverrideFields

      def validate
        super

        check_optional_property :validation, String
      end

      private

      def overriden
        Provider::Puppet::OverrideFields
      end
    end
  end
end

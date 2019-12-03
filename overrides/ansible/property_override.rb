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
require 'overrides/resources'

module Overrides
  module Ansible
    # Ansible-specific overrides to api.yaml.
    class PropertyOverride < Overrides::PropertyOverride
      # Collection of fields allowed in the PropertyOverride section for
      # Ansible. All fields should be `attr_reader :<property>`
      def self.attributes
        %i[
          aliases
          contain_extra_docs
        ]
      end

      attr_reader(*attributes)

      def validate
        super

        check :aliases, type: ::Array, item_type: ::String
        check :contain_extra_docs, type: :boolean, default: true
      end
    end
  end
end

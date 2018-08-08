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

module Provider
  module Ansible
    # Ansible specific properties to be added to Api::Resource
    class FactsOverride < Api::Object
      attr_reader :list_key
      attr_reader :has_filters

      def validate
        super

        default_value_property :list_key, 'items'
        default_value_property :has_filters, true

        check_property :list_key, ::String
        check_property :has_filters, :boolean
      end
    end
  end
end

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

require 'provider/resource_override'

module Provider
  module Ansible
    # Ansible specific properties to be added to Api::Resource
    module OverrideProperties
      attr_reader :access_api_results
    end

    # Product specific overriden properties for Ansible
    class ResourceOverride < Provider::ResourceOverride
      include OverrideProperties

      def validate
        super

        default_value_property :access_api_results, false

        check_property :access_api_results, :boolean
      end

      private

      def overriden
        Provider::Ansible::OverrideProperties
      end
    end
  end
end

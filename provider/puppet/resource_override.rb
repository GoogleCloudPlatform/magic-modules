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

require 'provider/core'
require 'provider/resource_override'

module Provider
  class Puppet < Provider::Core
    # Puppet specific properties to be added to Api::Resource
    module OverrideProperties
      attr_reader :handlers
      attr_reader :provider_helpers
      attr_reader :requires
    end

    # Custom Puppet code to handle type convergence operations
    class Handlers < Api::Object
      attr_reader :create
      attr_reader :delete
      attr_reader :flush
      attr_reader :resource_to_request_patch

      def validate
        super

        check_optional_property :create, String
        check_optional_property :delete, String
        check_optional_property :flush, String
        check_optional_property :resource_to_request_patch, String
      end
    end

    # Product specific overriden properties for Puppet
    class ResourceOverride < Provider::ResourceOverride
      include OverrideProperties

      def validate
        @provider_helpers ||= []

        super

        check_optional_property :access_api_results, :boolean
        check_optional_property :handlers, Provider::Puppet::Handlers

        check_property_list :provider_helpers, String
        check_property_list :requires, String
      end

      private

      def overriden
        Provider::Puppet::OverrideProperties
      end
    end
  end
end

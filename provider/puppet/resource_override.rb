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
      attr_reader :access_api_results
      attr_reader :handlers
      attr_reader :provider_helpers
      attr_reader :requires
      attr_reader :resource_to_request
      attr_reader :return_if_object
      attr_reader :unwrap_resource
      attr_reader :custom_create_resource
      attr_reader :custom_update_resource
    end

    # Custom Puppet code to handle type convergence operations
    class Handlers < Api::Object
      attr_reader :collection # A custom collection function to use
      attr_reader :create
      attr_reader :delete
      attr_reader :flush
      attr_reader :prefetch
      attr_reader :request_to_query
      attr_reader :resource_to_request_patch
      attr_reader :return_if_object
      attr_reader :self_link # A custom self_link function to use

      def validate
        super

        check_optional_property :create, String
        check_optional_property :delete, String
        check_optional_property :flush, String
        check_optional_property :prefetch, String
        check_optional_property :request_to_query, String
        check_optional_property :resource_to_request_patch, String
      end
    end

    # Product specific overriden properties for Puppet
    class ResourceOverride < Provider::ResourceOverride
      include OverrideProperties

      def validate
        assign_defaults

        super

        check_property :access_api_results, :boolean
        check_property :custom_create_resource, :boolean
        check_property :custom_update_resource, :boolean
        check_optional_property :handlers, Provider::Puppet::Handlers
        check_optional_property :requires, Array
        check_property :resource_to_request, :boolean
        check_property :return_if_object, :boolean
        check_property :unwrap_resource, :boolean

        check_property_list :provider_helpers, String
        check_property_list :requires, String
      end

      private

      def assign_defaults
        default_value_property :access_api_results, false
        default_value_property :custom_create_resource, false
        default_value_property :custom_update_resource, false
        default_value_property :manual, false
        default_value_property :provider_helpers, []
        default_value_property :resource_to_request, true
        default_value_property :return_if_object, true
        default_value_property :unwrap_resource, true
      end

      def overriden
        Provider::Puppet::OverrideProperties
      end
    end
  end
end

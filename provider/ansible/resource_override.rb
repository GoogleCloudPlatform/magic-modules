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
require 'provider/ansible/facts_override'

module Provider
  module Ansible
    # Ansible specific properties to be added to Api::Resource
    module OverrideProperties
      attr_reader :access_api_results
      attr_reader :collection
      attr_reader :custom_create_resource
      attr_reader :custom_update_resource
      attr_reader :create
      attr_reader :delete
      attr_reader :has_tests
      attr_reader :hidden
      attr_reader :imports
      attr_reader :post_create
      attr_reader :post_action
      attr_reader :provider_helpers
      attr_reader :return_if_object
      attr_reader :template
      attr_reader :unwrap_resource
      attr_reader :update
      attr_reader :version_added

      attr_reader :facts
    end

    # Product specific overriden properties for Ansible
    class ResourceOverride < Provider::ResourceOverride
      include OverrideProperties
      def validate
        super

        default_value_property :access_api_results, false
        default_value_property :custom_create_resource, false
        default_value_property :custom_update_resource, false
        default_value_property :exclude, false
        default_value_property :has_tests, true
        default_value_property :imports, []
        default_value_property :provider_helpers, []
        default_value_property :unwrap_resource, false

        check_property :access_api_results, :boolean
        check_optional_property :collection, ::String
        check_property :custom_create_resource, :boolean
        check_property :custom_update_resource, :boolean
        check_optional_property :create, ::String
        check_optional_property :delete, ::String
        check_property :has_tests, :boolean
        check_optional_property :hidden, ::Array
        check_property :imports, ::Array
        check_optional_property :post_create, ::String
        check_optional_property :post_action, ::String
        check_property :provider_helpers, ::Array
        check_optional_property :return_if_object, ::String
        check_optional_property :template, ::String
        check_optional_property :update, ::String
        check_optional_property :unwrap_resource, :boolean
        check_optional_property :version_added, ::String

        @facts ||= FactsOverride.new
        check_property :facts, FactsOverride
      end

      private

      def overriden
        Provider::Ansible::OverrideProperties
      end
    end
  end
end

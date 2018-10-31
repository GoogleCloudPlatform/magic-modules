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

require 'provider/overrides/resources'
require 'provider/ansible/facts_override'

module Provider
  module Ansible
    # Product specific overriden properties for Ansible
    class ResourceOverride < Provider::Overrides::ResourceOverride
      def self.attributes
        %i[
          access_api_results
          collection
          custom_create_resource
          custom_update_resource
          create
          delete
          editable
          has_tests
          hidden
          imports
          post_create
          post_action
          provider_helpers
          return_if_object
          unwrap_resource
          update
          version_added

          facts
        ]
      end

      attr_reader(*attributes)

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/MethodLength
      def validate
        super

        default_value_property :access_api_results, false
        default_value_property :custom_create_resource, false
        default_value_property :custom_update_resource, false
        default_value_property :exclude, false
        default_value_property :editable, true
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
        check_property :editable, :boolean
        check_property :has_tests, :boolean
        check_optional_property :hidden, ::Array
        check_property :imports, ::Array
        check_optional_property :post_create, ::String
        check_optional_property :post_action, ::String
        check_property :provider_helpers, ::Array
        check_optional_property :return_if_object, ::String
        check_optional_property :update, ::String
        check_optional_property :unwrap_resource, :boolean
        check_optional_property :version_added, ::String

        @facts ||= FactsOverride.new
        check_property :facts, FactsOverride
      end
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/AbcSize
    end
  end
end

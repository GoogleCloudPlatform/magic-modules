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
require 'provider/overrides/runner'

module Provider
  module Ansible
    # Contains alternate tests for verifying resource existence
    # using facts modules.
    # Contains a test to verify that a resource does exist and does not.
    # These tests may be the same, or they may differ.
    class AnsibleFactsTestInformation < Api::Object
      attr_reader :exists
      attr_reader :does_not_exist
      def validate
        super
        check_optional_property :exists, ::String
        check_optional_property :does_not_exist, ::String
      end
    end

    # Ansible specific properties to be added to Api::Resource
    class FactsOverride < Api::Object
      attr_reader :has_filters
      attr_reader :filter
      attr_reader :query_options
      attr_reader :filter_api_param
      attr_reader :test

      def validate
        super
        default_value_property :has_filters, true
        default_value_property :filter, FilterProp.new
        default_value_property :query_options, true
        default_value_property :filter_api_param, 'filter'

        check_property :has_filters, :boolean
        check_property :filter, Api::Object
        check_property :query_options, :boolean
        check_property :filter_api_param, ::String
        check_optional_property :test, AnsibleFactsTestInformation

        # We have to apply the property overrides and validate
        # the filtering property
        @filter = Provider::Overrides::Runner.build_single_property(
          @filter, {}, Provider::Overrides::Ansible::PropertyOverride
        )
      end
    end
    # This is a property exclusive to Ansible filters.
    # This is the default property used for filter information on Ansible.
    # By using Api::Types, we get more flexibility and a lot for free.
    class FilterProp < Api::Type::Array
      def validate
        @item_type ||= 'Api::Type::String'
        # GCE (and some others) uses the 'filters' property by default.
        # By default, assume that these are for GCE.
        @name ||= 'filters'
        @description ||= <<-STRING
        A list of filter value pairs. Available filters are listed here
                U(https://cloud.google.com/sdk/gcloud/reference/topic/filters).
                Each additional filter in the list will act be added as an AND condition
                (filter1 and filter2)
        STRING
        true
      end

      def gce?
        @name == 'filters'
      end
    end
  end
end

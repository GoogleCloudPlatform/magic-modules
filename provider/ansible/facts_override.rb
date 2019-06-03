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
require 'overrides/runner'

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
        check :exists, type: ::String
        check :does_not_exist, type: ::String
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
<<<<<<< HEAD
        default_value_property :has_filters, false
        default_value_property :filter, FilterProp.new
        default_value_property :query_options, true
        default_value_property :filter_api_param, 'filter'
=======
        check :has_filters, type: :boolean, default: true
        check :filter, type: Api::Object, default: FilterProp.new
        check :query_options, type: :boolean, default: true
        check :filter_api_param, type: ::String, default: 'filter'
        check :test, type: AnsibleFactsTestInformation
>>>>>>> master

        # We have to apply the property overrides and validate
        # the filtering property
        @filter = Overrides::Runner.build_single_property(
          @filter, {}, Overrides::Ansible::PropertyOverride
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
                https://cloud.google.com/sdk/gcloud/reference/topic/filters.
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

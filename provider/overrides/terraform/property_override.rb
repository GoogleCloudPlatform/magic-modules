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
require 'provider/overrides/resources'
require 'provider/terraform/property_override'

module Provider
  module Overrides
    module Terraform
      # Terraform-specific overrides to api.yaml.
      class PropertyOverride < Provider::Overrides::PropertyOverride
        # Collection of fields allowed in the PropertyOverride section for
        # Terraform. All fields should be `attr_reader :<property>`
        def self.attributes
          Provider::Terraform::OverrideFields.attributes
        end

        attr_reader(*attributes)

        # Used to allow us to easily access these values in `apply`
        # without resorting to "instance_variable_get"
        attr_reader :description

        def validate
          super

          # Ensures boolean values are set to false if nil
          @sensitive ||= false
          @is_set ||= false
          @unordered_list ||= false
          @default_from_api ||= false
          @conflicts_with ||= []

          check_property :sensitive, :boolean
          check_property :is_set, :boolean
          check_property :default_from_api, :boolean
          check_property_list :conflicts_with, ::String

          check_optional_property :diff_suppress_func, String
          check_optional_property :state_func, String
          check_optional_property :validation, Provider::Terraform::Validation
          check_optional_property :set_hash_func, String

          check_optional_property :update_statement, String
          check_optional_property :custom_flatten, String
          check_optional_property :custom_expand, String
        end
        # rubocop:enable Metrics/MethodLength

        # rubocop:disable Metrics/CyclomaticComplexity
        def apply(api_property)
          unless description.nil?
            @description = format_string(:description, @description,
                                         api_property.description)
          end

          unless api_property.is_a?(Api::Type::Array) ||
                 api_property.is_a?(Api::Type::Map)
            if @is_set
              raise 'Set can only be specified for Api::Type::Array ' \
                    'or Api::Type::NameValues<String, NestedObject>. ' \
                    "Type is #{api_property.class} for property "\
                    "'#{api_property.name}'"
            end
          end

          raise "'default_value' and 'default_from_api' cannot be both set"  \
            if default_from_api && !api_property.default_value.nil?

          super
        end
        # rubocop:enable Metrics/CyclomaticComplexity

        private

        # Formats the string and potentially uses its old value as part of the new
        # value. The marker should be in the form `{{name}}` where `name` is the
        # field being formatted.
        #
        # Note: This function only supports the variable with the same name as the
        # property being updated.
        def format_string(name, mask, current_value)
          mask.gsub "{{#{name.id2name}}}", current_value
        end
      end
    end
  end
end

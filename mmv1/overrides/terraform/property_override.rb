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

module Provider
  class Terraform
    # Support for schema ValidateFunc functionality.
    class Validation < Api::Object
      # Ensures the value matches this regex
      attr_reader :regex
      attr_reader :function

      def validate
        super

        check :regex, type: String
        check :function, type: String
      end
    end
  end
end

module Overrides
  module Terraform
    # Terraform-specific overrides to api.yaml.
    class PropertyOverride < Overrides::PropertyOverride
      # Collection of fields allowed in the PropertyOverride section for
      # Terraform. All fields should be `attr_reader :<property>`
      def self.attributes
        []
      end

      attr_reader(*attributes)

      # Used to allow us to easily access these values in `apply`
      # without resorting to "instance_variable_get"
      attr_reader :description

      def validate
        super
      end

      def apply(api_property)
        unless description.nil?
          @description = format_string(:description, @description,
                                       api_property.description)
        end

        if @flatten_object && !api_property.is_a?(Api::Type::NestedObject)
          raise 'Only NestedObjects can be flattened with flatten_object. Type'\
            " is #{api_property.class} for property #{api_property.name}"
        end

        unless api_property.is_a?(Api::Type::Array) ||
               api_property.is_a?(Api::Type::Map)
          if @is_set
            raise 'Set can only be specified for Api::Type::Array ' \
                  'or Api::Type:Map. ' \
                  "Type is #{api_property.class} for property "\
                  "'#{api_property.name}'"
          end

          if @schema_mode_config_attr && !@default_from_api
            raise 'default_from_api must be true on a nested block to set' \
                  'schema_mode_config_attr'
          end
        end

        super
      end

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

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
require 'provider/abstract_core'
require 'provider/property_override'

module Provider
  class Terraform < Provider::AbstractCore
    # Collection of fields allowed in the PropertyOverride section for
    # Terraform. All fields should be `attr_reader :<property>`
    module OverrideFields
      attr_reader :exclude
      # Adds a DiffSuppressFunc to the schema field with the given name
      attr_reader :diff_suppress_func

      attr_reader :validation
    end

    # Adds a ValidateFunc to the schema field.
    class Validation < Api::Object
      # Ensures the value matches this regex
      attr_reader :regex

      # TODO(rosbo): Add support for complex validation (i.e. custom go code)

      def validate
        super

        check_optional_property :regex, String
      end
    end

    # Terraform-specific overrides to api.yaml.
    class PropertyOverride < Provider::PropertyOverride
      include OverrideFields

      def validate
        super

        @exclude ||= false # sets to false if nil

        check_property :exclude, :boolean

        check_optional_property :diff_suppress_func, String
        check_optional_property :validation, Provider::Terraform::Validation
      end

      def apply(property)
        unless description.nil?
          @description = format_string(:description, @description,
                                       property.description)
        end

        super
      end

      private

      def overriden
        Provider::Terraform::OverrideFields
      end

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

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

require 'provider/terraform/custom_code'
require 'provider/overrides/resources'
require 'provider/terraform/resource_override'

module Provider
  module Overrides
    module Terraform
      # A class to control overridden properties on terraform.yaml in lieu of
      # values from api.yaml.
      class ResourceOverride < Provider::Overrides::ResourceOverride
        attr_reader :description

        def self.attributes
          Provider::Terraform::OverrideProperties.attributes
        end

        attr_reader(*attributes)

        def validate
          super

          @id_format ||= '{{name}}'
          @import_format ||= []
          @custom_code ||= Provider::Terraform::CustomCode.new
          @docs ||= Provider::Terraform::Docs.new
          @example ||= []

          check_property :id_format, String

          check_optional_property :examples, String
          check_optional_property_list :example, Provider::Terraform::Examples

          check_optional_property :custom_code, Provider::Terraform::CustomCode
          check_optional_property :docs, Provider::Terraform::Docs
          check_property :import_format, Array
          check_property_list :import_format, String
        end

        def apply(resource)
          unless description.nil?
            @description = format_string(:description, @description,
                                         resource.description)
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
end

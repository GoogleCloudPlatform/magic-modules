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

require 'provider/abstract_core'
require 'provider/resource_override'

module Provider
  class Terraform < Provider::AbstractCore
    # Collection of properties allowed in the ResourceOverride section for
    # Terraform. All properties should be `attr_reader :<property>`
    module OverrideProperties
    end

    # A class to control overridden properties on terraform.yaml in lieu of
    # values from api.yaml.
    class ResourceOverride < Provider::ResourceOverride
      include OverrideProperties

      def apply(resource)
        unless description.nil?
          @description = format_string(:description, @description,
                                       resource.description)
        end

        super
      end

      private

      def overriden
        Provider::Terraform::OverrideProperties
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

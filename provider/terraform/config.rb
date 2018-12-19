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
require 'provider/config'

module Provider
  class Terraform < Provider::AbstractCore
    # Settings for the provider
    class Config < Provider::Config
      def provider
        Provider::Terraform
      end

      def resource_override
        Provider::Terraform::ResourceOverride
      end

      def property_override
        Provider::Terraform::PropertyOverride
      end

      # These two methods are for the new set of overrides.
      # They'll replace `resource_override` and `property_override`
      # when the old overrides are deprecated.
      def new_resource_override
        Provider::Overrides::Terraform::ResourceOverride
      end

      def new_property_override
        Provider::Overrides::Terraform::PropertyOverride
      end
    end
  end
end

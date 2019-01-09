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

require 'provider/config'
require 'provider/core'

module Provider
  module Ansible
    # Settings for the Ansible provider
    class Config < Provider::Config
      attr_reader :manifest

      def provider
        Provider::Ansible::Core
      end

      def resource_override
        Provider::Ansible::ResourceOverride
      end

      def property_override
        Provider::Ansible::PropertyOverride
      end

      # These two methods are for the new set of overrides.
      # They'll replace `resource_override` and `property_override`
      # when the old overrides are deprecated.
      def new_resource_override
        Provider::Overrides::Ansible::ResourceOverride
      end

      def new_property_override
        Provider::Overrides::Ansible::PropertyOverride
      end

      def validate
        super
        check_optional_property :manifest, Provider::Ansible::Manifest
      end
    end
  end
end

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

module Provider
  class Chef < Provider::Core
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      attr_reader :operating_systems
      # TODO(alexstephen): Convert this to a regular function generator like Puppet.
      attr_reader :functions

      def provider
        Provider::Chef
      end

      def resource_override
        Provider::Chef::ResourceOverride
      end

      def property_override
        Provider::Chef::PropertyOverride
      end

      def validate
        super
        check_optional_property :manifest, Provider::Chef::Manifest
        check_property_list \
          :operating_systems, Provider::Config::OperatingSystem
      end
    end
  end
end
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
  class Puppet < Provider::Core
    # Settings for the provider
    class Config < Provider::Config
      attr_reader :manifest
      attr_reader :functions
      attr_reader :bolt_tasks

      def provider
        Provider::Puppet
      end

      def resource_override
        Provider::Puppet::ResourceOverride
      end

      def property_override
        Provider::Puppet::PropertyOverride
      end

      def validate
        super

        check_optional_property :manifest, Provider::Puppet::Manifest
        check_property_list :functions, Provider::Config::Function
        check_property_list :bolt_tasks, Provider::Puppet::BoltTask
      end
    end
  end
end

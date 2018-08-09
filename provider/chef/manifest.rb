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

require 'provider/core'

module Provider
  class Chef < Provider::Core
    # Metadata for manifest.json
    class Manifest < Api::Object
      attr_reader :additional_info
      attr_reader :depends
      attr_reader :description
      attr_reader :issues
      attr_reader :operating_systems
      attr_reader :requires
      attr_reader :source
      attr_reader :summary
      attr_reader :version

      def validate
        check_optional_property :additional_info, Array
        check_optional_property :depends, Array
        check_property :description, String
        check_property :issues, String
        check_property :operating_systems, Array
        check_property :source, String
        check_property :summary, String
        check_property :version, String
        check_property_list :depends, Provider::Config::Requirements
        check_property_list \
          :operating_systems, Provider::Config::OperatingSystem
        super
      end
    end
  end
end

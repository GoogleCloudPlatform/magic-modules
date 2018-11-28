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
  module Ansible
    # Metadata for manifest.json
    class Manifest < Api::Object
      attr_reader :metadata_version
      attr_reader :status
      attr_reader :supported_by
      attr_reader :requirements
      attr_reader :version_added
      attr_reader :author

      def validate
        check_property :metadata_version, String
        check_property :status, Array
        check_property_list :status, String
        check_property :supported_by, String
        check_property :requirements, Array
        check_property_list :requirements, String
        check_property :version_added, String
        check_property :author, String
      end

      # Get value from config and fallback to manifest.
      def get(value, object)
        return object.instance_variable_get("@#{value}".to_sym) \
          unless object.instance_variable_get("@#{value}".to_sym).nil?

        instance_variable_get("@#{value}".to_sym)
      end
    end
  end
end

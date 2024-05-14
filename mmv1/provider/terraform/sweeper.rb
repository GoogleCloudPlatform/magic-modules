# Copyright 2023 Google Inc.
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

require 'uri'
require 'api/object'
require 'compile/core'
require 'google/golang_utils'

module Provider
  class Terraform
    # Support for schema ValidateFunc functionality.
    class Sweeper < Google::YamlValidator
      # The field checked by sweeper to determine
      # eligibility for deletion for generated resources
      attr_reader :sweepable_identifier_field

      def validate
        super

        check :sweepable_identifier_field, type: String
      end
    end
  end
end

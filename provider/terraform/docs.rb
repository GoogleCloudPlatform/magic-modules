# Copyright 2019 Google Inc.
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
require 'provider/abstract_core'

module Provider
  class Terraform < Provider::AbstractCore
    # Inserts custom strings into terraform resource docs.
    class Docs < Api::Object
      # All these values should be strings, which will be inserted
      # directly into the terraform resource documentation.  The
      # strings should _not_ be the names of template files
      # (This should be reconsidered if we find ourselves repeating
      # any string more than ones), but rather the actual text
      # (including markdown) which needs to be injected into the
      # template.
      # The text will be injected at the bottom of the specified
      # section.
      attr_reader :warning
      attr_reader :required_properties
      attr_reader :optional_properties
      attr_reader :attributes

      def validate
        super
        check :warning, type: String
        check :required_properties, type: String
        check :optional_properties, type: String
        check :attributes, type: String
      end
    end
  end
end

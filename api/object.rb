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

require 'google/extensions'
require 'google/logger'
require 'google/yaml_validator'

module Api
  # Represents a base object
  class Object < Google::YamlValidator
    # Represents an object that has a (mandatory) name
    class Named < Api::Object
      # The list of properties (attr_reader) that can be overridden in
      # <provider>.yaml.
      module Properties
        attr_reader :name
      end

      include Properties

      # original value of :name before the provider override happens
      # same as :name if not overridden in provider
      attr_reader :api_name

      def validate
        super
        check :name, type: String, required: true
        check :api_name, type: String, default: @name
      end
    end

    def out_name
      @name.underscore
    end
  end
end

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

require 'google/logger'
require 'google/yaml_validator'

module Api
  # Repesents a base object
  class Object < Google::YamlValidator
    # Basic functions for defining classes that explicitly marks a property as
    # not in use, such as Provider::Config::TestData::NONE
    module MissingObject
      attr_reader :reason

      def validate
        return if @validated
        check_property :reason, String
        # Now that we verified a reason was provided, delete it so it does not
        # end in the object mapping (and eventually fail as the type "reason"
        # will likely never exist)
        remove_instance_variable('@reason')
        super
        @validated = true
      end
    end

    # Represents an object that has a (mandatory) name
    class Named < Api::Object
      attr_reader :name

      def validate
        super
        check_property :name, String
      end
    end

    def out_name
      Google::StringUtils.underscore(@name)
    end
  end
end

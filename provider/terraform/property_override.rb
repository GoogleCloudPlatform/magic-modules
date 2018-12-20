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

require 'api/object'
require 'provider/abstract_core'
require 'provider/property_override'

module Provider
  class Terraform < Provider::AbstractCore
    # Collection of fields allowed in the PropertyOverride section for
    # Terraform. All fields should be `attr_reader :<property>`
    module OverrideFields
      def self.attributes
        [
          :diff_suppress_func, # Adds a DiffSuppressFunc to the schema
          :state_func, # Adds a StateFunc to the schema
          :sensitive, # Adds `Sensitive: true` to the schema
          # Does not set this value to the returned API value.  Useful for fields
          # like secrets where the returned API value is not helpful.
          :ignore_read,
          :validation, # Adds a ValidateFunc to the schema
          # Indicates that this is an Array that should have Set diff semantics.
          :unordered_list,

          :is_set, # Uses a Set instead of an Array
          # Optional function to determine the unique ID of an item in the set
          # If not specified, schema.HashString (when elements are string) or
          # schema.HashSchema are used.
          :set_hash_func,

          # if true, then we get the default value from the Google API if no value
          # is set in the terraform configuration for this field.
          # It translates to setting the field to Computed & Optional in the schema.
          :default_from_api,

          # Names of attributes that can't be set alongside this one
          :conflicts_with,

          # ===========
          # Custom code
          # ===========
          # All custom code attributes are string-typed.  The string should
          # be the name of a template file which will be compiled in the
          # specified / described place.
          #
          # Property Updates are used when a resource is updateable but
          # resource.input is true.  In this case, only individual
          # properties can be updated.  The value of this attribute should
          # be the path to a template which will be compiled. This code is placed
          # *inline* in the obj := { ... } definition - it is not a custom
          # function, it is a custom statement.  Note that this cannot
          # be used for nested properties, as they are not present in the
          # obj := {...} statement.  This statement template receives `property`
          # and `prefix` to aid in code reuse.
          :update_statement,
          # A custom flattener replaces the default flattener for an attribute.
          # It is called as part of Read.  It can return an object of any
          # type, and may sometimes need to return an object with non-interface{}
          # type so that the d.Set() call will succeed, so the function
          # header *is* a part of the custom code template.  To help with
          # creating the function header, `property` and `prefix` are available,
          # just as they are in the standard flattener template.
          :custom_flatten,
          # A custom expander replaces the default expander for an attribute.
          # It is called as part of Create, and as part of Update if
          # object.input is false.  It can return an object of any type,
          # so the function header *is* part of the custom code template.
          # As with flatten, `property` and `prefix` are available.
          :custom_expand
        ]
      end

      attr_reader(*attributes)
    end

    # Support for schema ValidateFunc functionality.
    class Validation < Api::Object
      # Ensures the value matches this regex
      attr_reader :regex
      attr_reader :function

      def validate
        super

        check_optional_property :regex, String
        check_optional_property :function, String
      end
    end

    # Shared code between new property overrides and old.
    module PropertyOverrideSharedCode
      def validate
        super

        # Ensures boolean values are set to false if nil
        @sensitive ||= false
        @is_set ||= false
        @unordered_list ||= false
        @default_from_api ||= false

        check_property :sensitive, :boolean
        check_property :is_set, :boolean
        check_property :default_from_api, :boolean

        check_optional_property :diff_suppress_func, String
        check_optional_property :state_func, String
        check_optional_property :validation, Provider::Terraform::Validation
        check_optional_property :set_hash_func, String

        check_optional_property :update_statement, String
        check_optional_property :custom_flatten, String
        check_optional_property :custom_expand, String

        raise "'default_value' and 'default_from_api' cannot be both set"  \
          if @default_from_api && !@default_value.nil?
      end

      def apply(api_property)
        unless description.nil?
          @description = format_string(:description, @description,
                                       api_property.description)
        end

        unless api_property.is_a?(Api::Type::Array) ||
               api_property.is_a?(Api::Type::Map)
          if @is_set
            raise 'Set can only be specified for Api::Type::Array ' \
                  'or Api::Type:Map. ' \
                  "Type is #{api_property.class} for property "\
                  "'#{api_property.name}'"
          end
        end

        super
      end

      private

      # Formats the string and potentially uses its old value as part of the new
      # value. The marker should be in the form `{{name}}` where `name` is the
      # field being formatted.
      #
      # Note: This function only supports the variable with the same name as the
      # property being updated.
      def format_string(name, mask, current_value)
        mask.gsub "{{#{name.id2name}}}", current_value
      end
    end

    # Terraform-specific overrides to api.yaml.
    class PropertyOverride < Provider::PropertyOverride
      include OverrideFields
      include PropertyOverrideSharedCode

      private

      def overriden
        Provider::Terraform::OverrideFields
      end
    end
  end
end

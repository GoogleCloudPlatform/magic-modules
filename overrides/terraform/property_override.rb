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
require 'overrides/resources'

module Provider
  class Terraform
    # Support for schema ValidateFunc functionality.
    class Validation < Api::Object
      # Ensures the value matches this regex
      attr_reader :regex
      attr_reader :function

      def validate
        super

        check :regex, type: String
        check :function, type: String
      end
    end
  end
end

module Overrides
  module Terraform
    # Terraform-specific overrides to api.yaml.
    class PropertyOverride < Overrides::PropertyOverride
      # Collection of fields allowed in the PropertyOverride section for
      # Terraform. All fields should be `attr_reader :<property>`
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

          # https://github.com/hashicorp/terraform/pull/20837
          # Apply a ConfigMode of SchemaConfigModeAttr to the field.
          # This should be avoided for new fields, and only used with old ones.
          :schema_config_mode_attr,

          # Names of attributes that can't be set alongside this one
          :conflicts_with,

          # Names of attributes that at least one of must be set
          :at_least_one_of,

          # Names of attributes that exactly one of must be set
          :exactly_one_of,

          # Names of fields that should be included in the updateMask.
          :update_mask_fields,

          # For a TypeMap, the expander function to call on the key.
          # Defaults to expandString.
          :key_expander,

          # For a TypeMap, the DSF to apply to the key.
          :key_diff_suppress_func,

          # ====================
          # Schema Modifications
          # ====================
          # Schema modifications change the schema of a resource in some
          # fundamental way. They're not very portable, and will be hard to
          # generate so we should limit their use. Generally, if you're not
          # converting existing Terraform resources, these shouldn't be used.
          #
          # With great power comes great responsibility.

          # Flattens a NestedObject by removing that field from the Terraform
          # schema but will preserve it in the JSON sent/retrieved from the API
          #
          # EX: a API schema where fields are nested (eg: `one.two.three`) and we
          # desire the properties of the deepest nested object (eg: `three`) to
          # become top level properties in the Terraform schema. By overriding
          # the properties `one` and `one.two` and setting flatten_object then
          # all the properties in `three` will be at the root of the TF schema.
          #
          # We need this for cases where a field inside a nested object has a
          # default, if we can't spend a breaking change to fix a misshapen
          # field, or if the UX is _much_ better otherwise.
          #
          # WARN: only fully flattened properties are currently supported. In the
          # example above you could not flatten `one.two` without also flattening
          # all of it's parents such as `one`
          :flatten_object,

          # ===========
          # Custom code
          # ===========
          # All custom code attributes are string-typed.  The string should
          # be the name of a template file which will be compiled in the
          # specified / described place.

          # A custom expander replaces the default expander for an attribute.
          # It is called as part of Create, and as part of Update if
          # object.input is false.  It can return an object of any type,
          # so the function header *is* part of the custom code template.
          # As with flatten, `property` and `prefix` are available.
          :custom_expand,

          # A custom flattener replaces the default flattener for an attribute.
          # It is called as part of Read.  It can return an object of any
          # type, and may sometimes need to return an object with non-interface{}
          # type so that the d.Set() call will succeed, so the function
          # header *is* a part of the custom code template.  To help with
          # creating the function header, `property` and `prefix` are available,
          # just as they are in the standard flattener template.
          :custom_flatten
        ]
      end

      attr_reader(*attributes)

      # Used to allow us to easily access these values in `apply`
      # without resorting to "instance_variable_get"
      attr_reader :description

      def validate
        super

        check :sensitive, type: :boolean, default: false
        check :is_set, type: :boolean, default: false
        check :default_from_api, type: :boolean, default: false
        check :unordered_list, type: :boolean, default: false
        check :schema_config_mode_attr, type: :boolean, default: false

        # technically set as a default everywhere, but only maps will use this.
        check :key_expander, type: String, default: 'expandString'
        check :key_diff_suppress_func, type: String

        check :diff_suppress_func, type: String
        check :state_func, type: String
        check :validation, type: Provider::Terraform::Validation
        check :set_hash_func, type: String

        check :custom_flatten, type: String
        check :custom_expand, type: String

        raise "'default_value' and 'default_from_api' cannot be both set"  \
          if @default_from_api && !@default_value.nil?
      end

      def apply(api_property)
        unless description.nil?
          @description = format_string(:description, @description,
                                       api_property.description)
        end

        if @flatten_object && !api_property.is_a?(Api::Type::NestedObject)
          raise 'Only NestedObjects can be flattened with flatten_object. Type'\
            " is #{api_property.class} for property #{api_property.name}"
        end

        unless api_property.is_a?(Api::Type::Array) ||
               api_property.is_a?(Api::Type::Map)
          if @is_set
            raise 'Set can only be specified for Api::Type::Array ' \
                  'or Api::Type:Map. ' \
                  "Type is #{api_property.class} for property "\
                  "'#{api_property.name}'"
          end

          if @schema_mode_config_attr && !@default_from_api
            raise 'default_from_api must be true on a nested block to set' \
                  'schema_mode_config_attr'
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
  end
end

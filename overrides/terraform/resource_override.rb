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

require 'overrides/resources'
require 'provider/terraform/custom_code'
require 'provider/terraform/docs'
require 'provider/terraform/examples'
require 'provider/terraform/virtual_fields'

module Overrides
  module Terraform
    # A class to control overridden properties on terraform.yaml in lieu of
    # values from api.yaml.
    class ResourceOverride < Overrides::ResourceOverride
      def self.attributes
        [
          # If non-empty, overrides the full filename prefix
          # i.e. google/resource_product_{{resource_filename_override}}.go
          # i.e. google/resource_product_{{resource_filename_override}}_test.go
          # Note this doesn't override the actual resource name
          # use :legacy_name instead.
          :filename_override,

          # If non-empty, overrides the full given resource name.
          # i.e. 'google_project' for resourcemanager.Project
          # Use Provider::Terraform::Config.legacy_name to override just
          # product name.
          :legacy_name,

          # The Terraform resource id format used when calling #setId(...).
          # For instance, `{{name}}` means the id will be the resource name.
          :id_format,
          :import_format,
          :custom_code,
          :docs,

          # Lock name for a mutex to prevent concurrent API calls for a given
          # resource.
          :mutex,

          # Examples in documentation. Backed by generated tests, and have
          # corresponding OiCS walkthroughs.
          :examples,

          # Virtual fields on the Terraform resource. Usage and differences from url_param_only
          # are documented in provider/terraform/virtual_fields.rb
          :virtual_fields,

          # TODO(alexstephen): Deprecate once all resources using autogen async.
          :autogen_async,

          # If true, resource is not importable
          :exclude_import,

          # If true, exclude resource from Terraform Validator
          # (i.e. terraform-provider-conversion)
          :exclude_validator,

          :timeouts,

          # An array of function names that determine whether an error is retryable.
          :error_retry_predicates,

          :schema_version,

          # If true, skip sweeper generation for this resource
          :skip_sweeper,

          # Set to true for resources that are unable to be deleted, such as KMS keyrings or project
          # level resources such as firebase project
          :skip_delete,

          # This enables resources that get their project via a reference to a different resource
          # instead of a project field to use User Project Overrides
          :supports_indirect_user_project_override,

          # Function to transform a read error so that handleNotFound recognises
          # it as a 404. This should be added as a handwritten fn that takes in
          # an error and returns one.
          :read_error_transform
        ]
      end

      attr_reader(*attributes)
      attr_reader :description

      def validate
        super

        @examples ||= []

        check :filename_override, type: String
        check :legacy_name, type: String
        check :id_format, type: String
        check :examples, item_type: Provider::Terraform::Examples, type: Array, default: []
        check :virtual_fields,
              item_type: Api::Type,
              type: Array,
              default: []

        check :custom_code, type: Provider::Terraform::CustomCode,
                            default: Provider::Terraform::CustomCode.new
        check :docs, type: Provider::Terraform::Docs, default: Provider::Terraform::Docs.new
        check :import_format, type: Array, item_type: String, default: []
        check :autogen_async, type: :boolean, default: false
        check :exclude_import, type: :boolean, default: false

        check :timeouts, type: Api::Timeouts
        check :error_retry_predicates, type: Array, item_type: String
        check :schema_version, type: Integer
        check :skip_sweeper, type: :boolean, default: false
        check :skip_delete, type: :boolean, default: false
        check :supports_indirect_user_project_override, type: :boolean, default: false
        check :read_error_transform, type: String
      end

      def apply(resource)
        unless description.nil?
          @description = format_string(:description, @description,
                                       resource.description)
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

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
require 'compile/core'
require 'google/golang_utils'
require 'provider/abstract_core'
require 'provider/property_override'

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
        check_optional_property :warning, String
        check_optional_property :required_properties, String
        check_optional_property :optional_properties, String
        check_optional_property :attributes, String
      end
    end

    # Generates configs to be shown as examples in docs and outputted as tests
    # from a shared template
    class Examples < Api::Object
      include Compile::Core
      include Google::GolangUtils

      # The name of the example in lower snake_case.
      # Generally takes the form of the resource name followed by some detail
      # about the specific test. For example, "address_with_subnetwork".
      # The template for the example is expected at the path
      # "templates/terraform/examples/{{name}}.tf.erb"
      attr_reader :name

      # The id of the "primary" resource in an example. Used in import tests.
      # This is the value that will appear in the Terraform config url. For
      # example:
      # resource "google_compute_address" {{primary_resource_id}} {
      #   ...
      # }
      attr_reader :primary_resource_id

      # vars is a Hash from template variable names to output variable names
      attr_reader :vars

      # Extra properties to ignore read on during import.
      # These properties will likely be custom code.
      attr_reader :ignore_read_extra

      def config_documentation
        body = lines(compile_file(
                       {
                         vars: vars,
                         primary_resource_id: primary_resource_id
                       },
                       "templates/terraform/examples/#{name}.tf.erb"
        ))
        lines(compile_file(
                { content: body },
                'templates/terraform/examples/base_configs/documentation.tf.erb'
        ))
      end

      def config_test
        @vars ||= []
        body = lines(compile_file(
                       {
                         vars: vars.map { |k, str| [k, "#{str}-%s"] }.to_h,
                         primary_resource_id: primary_resource_id
                       },
                       "templates/terraform/examples/#{name}.tf.erb"
        ))

        body = substitute_test_paths body

        lines(compile_file(
                {
                  content: body,
                  count: vars.length
                },
                'templates/terraform/examples/base_configs/test_body.go.erb'
        ))
      end

      def config_example
        lines(compile_file(
                {
                  vars: vars.map { |k, str| [k, "#{str}-${local.name_suffix}"] }.to_h,
                  primary_resource_id: primary_resource_id
                },
                "templates/terraform/examples/#{name}.tf.erb"
        ))
      end

      def substitute_test_paths(config)
        config = config.gsub('path/to/private.key', 'test-fixtures/ssl_cert/test.key')
        config.gsub('path/to/certificate.crt', 'test-fixtures/ssl_cert/test.crt')
      end

      def validate
        super
        @ignore_read_extra ||= []

        check_property :name, String
        check_property :primary_resource_id, String

        check_optional_property :vars, Hash
        check_optional_property_list :ignore_read_extra, String
      end
    end

    # Inserts custom code into terraform resources.
    class CustomCode < Api::Object
      # Collection of fields allowed in the CustomCode section for
      # Terraform.

      # All custom code attributes are string-typed.  The string should
      # be the name of a template file which will be compiled in the
      # specified / described place.
      #
      # ======================
      # schema.Resource stuff
      # ======================
      # Extra Schema Entries go below all other schema entries in the
      # resource's Resource.Schema map.  They should be formatted as
      # entries in the map, e.g. `"foo": &schema.Schema{ ... },`.
      attr_reader :extra_schema_entry
      # Resource definition code is inserted below everything else
      # in the resource's Resource {...} definition.  This may be useful
      # for things like a MigrateState / SchemaVersion pair.
      # This is likely to be used rarely and may be removed if all its
      # use cases are covered in other ways.
      attr_reader :resource_definition
      # ====================
      # Encoders & Decoders
      # ====================
      # The encoders are functions which take the `obj` map after it
      # has been assembled in either "Create" or "Update" and mutate it
      # before it is sent to the server.  There are lots of reasons you
      # might want to use these - any differences between local schema
      # and remote schema will be placed here.
      # Because the call signature of this function cannot be changed,
      # the template will place the function header and closing } for
      # you, and your custom code template should *not* include them.
      attr_reader :encoder
      # The update encoder is the encoder used in Update - if one is
      # not provided, the regular encoder is used.  If neither is
      # provided, of course, neither is used.  Similarly, the custom
      # code should *not* include the function header or closing }.
      # Update encoders are only used if object.input is false,
      # because when object.input is true, only individual fields
      # can be updated - in that case, use a custom expander.
      attr_reader :update_encoder
      # The decoder is the opposite of the encoder - it's called
      # after the Read succeeds, rather than before Create / Update
      # are called.  Like with encoders, the decoder should not
      # include the function header or closing }.
      attr_reader :decoder

      # =====================
      # Simple customizations
      # =====================
      # Constants go above everything else in the file, and include
      # things like methods that will be referred to by name elsewhere
      # (e.g. "fooBarDiffSuppress") and regexes that are necessarily
      # exported (e.g. "fooBarValidationRegex").
      attr_reader :constants
      # This code is run after the Create call succeeds.  It's placed
      # in the Create function directly without modification.
      attr_reader :post_create
      # This code is run before the Update call happens.  It's placed
      # in the Update function, just after the encoder call, before
      # the Update call.  Just like the encoder, it is only used if
      # object.input is false.
      attr_reader :pre_update
      # This code is run after the Update call happens.  It's placed
      # in the Update function, just after the call succeeds.
      # Just like the encoder, it is only used if object.input is
      # false.
      attr_reader :post_update
      # This code is run just before the Delete call happens.  It's
      # useful to prepare an object for deletion, e.g. by detaching
      # a disk before deleting it.
      attr_reader :pre_delete
      # This code replaces the entire import method.  Since the import
      # method's function header can't be changed, the template
      # inserts that for you - do not include it in your custom code.
      attr_reader :custom_import
      # This code is run just after the import method succeeds - it
      # is useful for parsing attributes that are necessary for
      # the Read() method to succeed.
      attr_reader :post_import

      def validate
        super

        check_optional_property :extra_schema_entry, String
        check_optional_property :resource_definition, String
        check_optional_property :encoder, String
        check_optional_property :update_encoder, String
        check_optional_property :decoder, String
        check_optional_property :constants, String
        check_optional_property :post_create, String
        check_optional_property :pre_update, String
        check_optional_property :post_update, String
        check_optional_property :pre_delete, String
        check_optional_property :custom_import, String
        check_optional_property :post_import, String
      end
    end
  end
end

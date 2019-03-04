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

      # vars is a Hash from template variable names to output variable names.
      # It will use the provided value as a prefix for generated tests, and
      # insert it into the docs verbatim.
      attr_reader :vars
      # Some variables need to hold special values during tests, and cannot
      # be inferred by Open in Cloud Shell.  For instance, org_id
      # needs to be the correct value during integration tests, or else
      # org tests cannot pass. Other examples include an existing project_id,
      # a zone, a service account name, etc.
      #
      # test_env_vars is a Hash from template variable names to one of the
      # following symbols:
      #  - :PROJECT_NAME
      #  - :CREDENTIALS
      #  - :REGION
      #  - :ORG_ID
      #  - :ORG_TARGET
      #  - :BILLING_ACCT
      #  - :SERVICE_ACCT
      # This list corresponds to the `get*FromEnv` methods in provider_test.go.
      attr_reader :test_env_vars

      # the version of the example. Note that _all features_ used in an example
      # must be set to the example min version.
      attr_reader :min_version

      # Extra properties to ignore read on during import.
      # These properties will likely be custom code.
      attr_reader :ignore_read_extra

      # Whether to skip generating tests for this resource
      attr_reader :skip_test

      def config_documentation
        docs_defaults = {
          PROJECT_NAME: 'my-project-name',
          CREDENTIALS: 'my/credentials/filename.json',
          REGION: 'us-west1',
          ORG_ID: '123456789',
          ORG_TARGET: '123456789',
          BILLING_ACCT: '000000-0000000-0000000-000000',
          SERVICE_ACCT: 'emailAddress:my@service-account.com'
        }
        @vars ||= {}
        @test_env_vars ||= {}
        body = lines(compile_file(
                       {
                         vars: vars,
                         test_env_vars: test_env_vars.map { |k, v| [k, docs_defaults[v]] }.to_h,
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
        @vars ||= {}
        @test_env_vars ||= {}
        body = lines(compile_file(
                       {
                         vars: vars.map { |k, str| [k, "#{str}-%{random_suffix}"] }.to_h,
                         test_env_vars: test_env_vars.map { |k, _| [k, "%{#{k}}"] }.to_h,
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
        @vars ||= []
        # Examples with test_env_vars are skipped elsewhere
        body = lines(compile_file(
                       {
                         vars: vars.map { |k, str| [k, "#{str}-${local.name_suffix}"] }.to_h,
                         primary_resource_id: primary_resource_id
                       },
                       "templates/terraform/examples/#{name}.tf.erb"
                     ))

        substitute_example_paths body
      end

      def oics_link
        hash = {
          cloudshell_git_repo: 'https://github.com/terraform-google-modules/docs-examples.git',
          cloudshell_working_dir: @name,
          cloudshell_image: 'gcr.io/graphite-cloud-shell-images/terraform:latest',
          open_in_editor: 'main.tf',
          cloudshell_print: './motd',
          cloudshell_tutorial: './tutorial.md'
        }
        URI::HTTPS.build(
          host: 'console.cloud.google.com',
          path: '/cloudshell/open',
          query: URI.encode_www_form(hash)
        )
      end

      def substitute_test_paths(config)
        config = config.gsub('../static/img/header-logo.png', 'test-fixtures/header-logo.png')
        config = config.gsub('path/to/private.key', 'test-fixtures/ssl_cert/test.key')
        config.gsub('path/to/certificate.crt', 'test-fixtures/ssl_cert/test.crt')
      end

      def substitute_example_paths(config)
        config = config.gsub('../static/img/header-logo.png', '../static/header-logo.png')
        config = config.gsub('path/to/private.key', '../static/ssl_cert/test.key')
        config.gsub('path/to/certificate.crt', '../static/ssl_cert/test.crt')
      end

      def validate
        super
        check :name, type: String, required: true
        check :primary_resource_id, type: String
        check :min_version, type: String
        check :vars, type: Hash
        check :test_env_vars, type: Hash
        check :ignore_read_extra, type: Array, item_type: String, default: []
        check :skip_test, type: TrueClass
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

        check :extra_schema_entry, type: String
        check :resource_definition, type: String
        check :encoder, type: String
        check :update_encoder, type: String
        check :decoder, type: String
        check :constants, type: String
        check :post_create, type: String
        check :pre_update, type: String
        check :post_update, type: String
        check :pre_delete, type: String
        check :custom_import, type: String
        check :post_import, type: String
      end
    end
  end
end

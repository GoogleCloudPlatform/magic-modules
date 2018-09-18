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

require 'provider/abstract_core'

module Provider
  class Terraform < Provider::AbstractCore
    # This module converts a generic (api.yaml) example config definition into
    # the shape of a Terraform config.
    module Examples
      # Params:
      #   - a generic example config
      # Return a list of Hash's with:
      # :example -> a valid TF config wrapped in a markdown code block
      def config_documentation(docs_examples)
        new_examples = []
        docs_examples&.each do |example|
          config = generic_example_to_config(example)
          example.vars.each do |k, v|
            config = config.gsub "{{#{k}}}", v
          end

          new_example = lines(compile_file(
                                { content: config },
                                'templates/terraform/examples/base_configs/documentation.tf.erb'
          ))
          new_examples << {
            example: new_example
          }
        end
        new_examples
      end

      # Params:
      #   - a generic example config
      # Return a list of Hash's with:
      # :example -> the method body of a _test file tf config
      # :primary_resource_id -> the tf uri name of the resource being tested
      # :name -> the name of a test case
      def config_test(docs_examples)
        new_examples = []
        docs_examples&.each do |example|
          config = generic_example_to_config(example)
          test_vars = example.vars.map { |k, str| [k, "#{str}-%s"] }.to_h
          test_vars.each do |k, v|
            config = config.gsub "{{#{k}}}", v
          end
          #

          new_example = lines(compile_file(
                                {
                                  content: config,
                                  count: example.vars.length
                                },
                                'templates/terraform/examples/base_configs/test_body.go.erb'
          ))
          new_examples << {
            example: new_example,
            name: example.name,
            primary_resource_id: example.primary_resource.name
          }
        end
        new_examples
      end

      private

      # Turns a single Api::Resource::Example into a tf config
      def generic_example_to_config(example)
        lines(compile_file(
                terraformify_example(example),
                'templates/terraform/examples/base_configs/config.tf.erb'
        ))
      end

      # Turns an Api::Resource::Example into a similarly shaped Hash with
      # values in the shape they need to be for Terraform.
      def terraformify_example(example)
        {
          name: example.name,
          vars: example.vars,
          primary_resource: terraformify_resource(example.primary_resource),
          resources: example.resources&.map { |r| terraformify_resource(r) }
        }
      end

      # Turns an Api::Resource::Example::Resource into a similarly shaped Hash
      # with values coerced in the shape they need to be for Terraform.
      def terraformify_resource(resource)
        {
          name: resource.name,
          type: generic_type_to_terraform_type(resource.type),
          properties: Hash[resource.properties&.map { |k, v| terraformify_property(k, v) }]
        }
      end

      # Turns a property key, value into the right shape for Terraform.
      # Keys (field_name) are underscored properly
      #
      # Values starting with @ are turned into interpolation
      # For example:
      # "@gcompute/Subnetwork/default/selfLink"
      # "${google_compute_network.vpc_network.self_link}"
      def terraformify_property(field_name, value)
        if value.start_with?('@')
          parts = value.split('/')

          # TODO(rileykarson): actually infer product here
          product = 'compute'
          resource = parts[1].underscore
          uri = parts[2]
          field = parts[3].underscore

          value = "${google_#{product}_#{resource}.#{uri}.#{field}}"
        end
        [field_name.underscore, value]
      end

      # Turns a resource type into it's name in Terraform
      # For example:
      # "gcompute/Address"
      # "compute_address"
      def generic_type_to_terraform_type(type)
        parts = type.split('/')

        # TODO(rileykarson): actually infer product here
        product = 'compute'
        resource = parts[1].underscore
        "#{product}_#{resource}"
      end
    end
  end
end

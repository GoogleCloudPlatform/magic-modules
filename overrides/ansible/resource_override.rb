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
require 'provider/ansible/facts_override'
require 'provider/ansible/custom_code'
require 'provider/ansible/tests'

module Overrides
  module Ansible
    # Allows overriding snowflake transport requests
    # Each one of these takes in a filename or a function name.
    # If filename, each file should contain one Python file.
    class Transport < Api::Object
      include Compile::Core

      attr_reader :encoders
      attr_reader :decoders
      attr_reader :remove_nones_post_encoder

      def validate
        super
        check :encoders, type: ::Array, default: [], item_type: ::String
        check :decoders, type: ::Array, default: [], item_type: ::String
        check :remove_nones_post_encoder, type: :boolean, default: true
      end

      def encoder_functions
        function_names(@encoders)
      end

      def decoder_functions
        function_names(@decoders)
      end

      private

      # Given an array of files that contain a single function and function names,
      # return the list of function names
      def function_names(array)
        array.map do |e|
          if File.file?(e)
            compile(e).match(/def ([a-zA-Z_]*)\([a-zA-Z_]*, [a-zA-Z_]*\)/)[1]
          else
            e
          end
        end
      end
    end

    # A class to control overridden properties on ansible.yaml in lieu of
    # values from api.yaml.
    class ResourceOverride < Overrides::ResourceOverride
      def self.attributes
        %i[
          access_api_results
          collection
          custom_code
          hidden
          imports
          notes
          provider_helpers
          template
          transport
          unwrap_resource

          tests

          facts
        ]
      end

      attr_reader(*attributes)

      def validate
        super

        @exclude ||= false

        check :access_api_results, type: :boolean, default: false
        check :collection, type: ::String
        check :custom_code, type: Provider::Ansible::CustomCode,
                            default: Provider::Ansible::CustomCode.new
        check :hidden, type: ::Array, item_type: String, default: []
        check :imports, type: ::Array, default: [], item_type: String
        check :notes, type: ::Array, item_type: String
        check :provider_helpers, type: ::Array, default: [], item_type: String
        check :return_if_object, type: ::String
        check :transport, type: Transport, default: Transport.new
        check :template, type: ::String
        check :update, type: ::String
        check :unwrap_resource, type: :boolean, default: false

        check :tests, type: Provider::Ansible::Tests,
                      default: Provider::Ansible::Tests.new

        check :facts, type: Provider::Ansible::FactsOverride,
                      default: Provider::Ansible::FactsOverride.new
      end
    end
  end
end

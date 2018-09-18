# Copyright 2018 Google Inc.
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

require 'api/resource/example/resource'

module Api
  class Resource < Api::Object::Named
    # Source for configs to be shown as examples in docs and outputted as tests
    # from a shared template
    class Example < Api::Object
      include Compile::Core

      # The name of the example in lower snake_case.
      # Generally takes the form of the resource name followed by some detail
      # about the specific test. For example, "address_with_subnetwork".
      attr_reader :name

      # vars is a Hash from template variable names to output variable names
      attr_reader :vars

      # The main resource block in an example config.
      attr_reader :primary_resource

      # Every supporting resource in an example config.
      attr_reader :resources

      def validate
        super
        @resources ||= []

        check_property :name, String
        check_property :vars, Hash
        check_property :primary_resource, Resource::Example::Resource
        check_property :resources, Array # Resource::Example::Resource
      end
    end
  end
end

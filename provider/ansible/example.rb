# frozen_string_literal: true

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
require 'provider/config'
require 'provider/core'
require 'provider/ansible/manifest'

module Provider
  class Ansible
    INTEGRATION_TEST_DEFAULTS = {
      project: '"{{ gcp_project }}"',
      auth_kind: '"{{ gcp_cred_kind }}"',
      service_account_file: '"{{ gcp_cred_file }}"',
      name: '"{{ resource_name }}"'
    }.freeze

    EXAMPLE_DEFAULTS = {
      name: 'testObject',
      project: 'testProject',
      auth_kind: 'service_account',
      service_account_file: '/tmp/auth.pem'
    }.freeze

    # Class responsible for holding a single Ansible task. This task may create
    # a GCP resource or create a dependent GCP resource.
    class Task < Api::Object
      include Compile::Core
      attr_reader :name
      attr_reader :code

      def validate
        super
        check_property :name, String
        check_property :code, String
      end

      def build_test(state, noop = false)
        build_task(state, INTEGRATION_TEST_DEFAULTS, noop)
      end

      def build_example(state)
        build_task(state, EXAMPLE_DEFAULTS)
      end

      private

      def build_task(state, hash, noop = false)
        verb = 'create'
        verb = 'delete' if state == 'absent'

        again = ''
        again = ' again' if noop
        [
          "- name: #{verb} a #{object_name_from_module_name(@name)}#{again}",
          "  #{@name}:",
          indent(compile_template_with_hash(@code, hash), 4).to_s,
          "    state: #{state}"
        ].join("\n")
      end

      def compile_template_with_hash(template, hash)
        ERB.new(template).result(OpenStruct.new(hash).instance_eval { binding })
      end

      def object_name_from_module_name(mod_name)
        product_name = mod_name.match(/gcp_[a-z]*_(.*)/).captures[0]
        product_name.tr('_', ' ')
      end
    end

    # Class responsible for holding all information relating to Ansible
    # examples.
    class Example < Api::Object
      attr_reader :task
      attr_reader :dependencies

      def validate
        super
        check_property :task, Task
        check_property_list :dependencies, @dependencies, Task \
          unless @dependencies.nil?
      end
    end
  end
end

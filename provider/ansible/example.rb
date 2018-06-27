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
  module Ansible
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

    # Holds all information necessary to run gcloud and verify the creation
    # of a resource.
    #
    # command         - The gcloud command being run
    # failed_name     - The name of the resource being created
    #                   (defaults to {{resource_name}})
    #                   Used for verifying deletion.
    # failed_verifier - A custom Jinja2 verifier to test if gcloud
    #                   command post-deletion worked as intended.
    class Verifier < Api::Object
      include Compile::Core
      attr_reader :command
      attr_reader :failure

      def validate
        @failure ||= FailureCondition.new

        check_property :command, String
        check_property :failure, FailureCondition
      end

      # rubocop:disable Metrics/AbcSize
      # rubocop:disable Metrics/MethodLength
      def build_task(state, object)
        raise 'State must be present or absent' \
          unless %w[present absent].include? state

        obj_name = Google::StringUtils.underscore(object.name)
        verb = verbs[state.to_sym]
        check_on_failure = @failure.enabled && state == 'absent'
        status = state == 'present' ? 0 : 1
        [
          "- name: verify that #{obj_name} was #{verb}",
          indent([
            'shell: |',
            # Ansible doesn't like multiline shell commands in playbooks.
            indent(@command.tr("\n", ''), 2),
            'register: results',
            ('failed_when: results.rc == 0' if state == 'absent')
          ].compact, 2),

          '- name: verify that command succeeded',
          indent([
                   'assert:',
                   indent([
                            'that:',
                            indent([
                              "- results.rc == #{status}",
                              ("- #{@failure.test}" if check_on_failure)
                            ].compact, 2)
                          ], 2)
                 ], 2)
        ].compact
      end
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/AbcSize

      private

      def verbs
        {
          present: 'created',
          absent: 'deleted'
        }
      end
    end

    # A gcloud command failing is not enough to verify that a resource does not
    # exist
    # Stderr should be checked to verify that the resource actually does not
    # exist.
    # @name - the name of the resource we're looking for
    # @error - the full line in stderr that's being looked for.
    # @test - the full test to verify resource does not exist.
    class FailureCondition < Api::Object
      attr_reader :enabled
      attr_reader :name
      attr_reader :test

      def validate
        check_optional_property :name, ::String
        check_optional_property :enabled, [TrueClass, FalseClass]
        check_optional_property :test, ::String

        @enabled ||= true
        @name ||= '{{ resource_name }}'
        @error ||= "#{@name} was not found."
        @test ||= "\"\\\"#{@error.strip}\\\" in results.stderr\""
      end
    end

    # GCE gcloud commands follow a relatively standard pattern.
    class ComputeFailureCondition < FailureCondition
      attr_reader :region
      attr_reader :type

      def validate
        raise 'Region must be slash delineated (e.g. regions/us-west1)' \
          unless @region == 'global' || @region.match?(%r{.*\/.*})

        check_optional_property :type, ::String

        @name ||= '{{ resource_name }}'
        @error = [
          "'projects/{{ gcp_project }}/#{@region}/#{@type}/#{@name}\'",
          'was not found'
        ].join(' ')
        super
      end
    end

    # A PERMISSION_DENIED is returned if the service account does not exist.
    class IamServiceAccountFailureCondition < FailureCondition
      def validate
        @error = [
          'PERMISSION_DENIED: Permission iam.serviceAccounts.get',
          'is required to perform this operation on service account'
        ].join(' ')
        super
      end
    end

    # Class responsible for holding a single Ansible task. This task may create
    # a GCP resource or create a dependent GCP resource.
    class Task < Api::Object
      include Compile::Core
      attr_reader :name
      attr_reader :code
      attr_reader :scopes
      attr_reader :register

      def validate
        super
        check_property :name, String
        check_property :code, String
        check_optional_property_list :scopes, ::String
      end

      def build_test(state, object, noop = false)
        build_task(state, INTEGRATION_TEST_DEFAULTS, object, noop)
      end

      def build_example(state, object)
        build_task(state, EXAMPLE_DEFAULTS, object)
      end

      def verbs
        {
          present: 'create',
          absent: 'delete'
        }
      end

      private

      # rubocop:disable Metrics/CyclomaticComplexity
      def build_task(state, hash, object, noop = false)
        verb = verbs[state.to_sym]

        again = ''
        again = ' that already exists' if noop && state == 'present'
        again = ' that does not exist' if noop && state == 'absent'
        scopes = (@scopes || object.__product.scopes).map { |x| "- #{x}" }
        [
          "- name: #{verb} a #{object_name_from_module_name(@name)}#{again}",
          indent([
            "#{@name}:",
            indent([
                     compile_string(hash, @code),
                     'scopes:',
                     indent(lines(scopes), 2),
                     "state: #{state}"
                   ], 4),
            ("register: #{@register}" unless @register.nil?)
          ].compact, 2)
        ]
      end
      # rubocop:enable Metrics/MethodLength
      # rubocop:enable Metrics/CyclomaticComplexity

      def object_name_from_module_name(mod_name)
        product_name = mod_name.match(/gcp_[a-z]*_(.*)/).captures[0]
        product_name.tr('_', ' ')
      end

      def dependency_name(dependency, resource)
        quote_string("#{dependency.downcase}-#{resource.downcase}")
      end
    end

    # Class responsible for holding all information relating to Ansible
    # examples.
    class Example < Api::Object
      attr_reader :task
      attr_reader :verifier
      attr_reader :dependencies

      def validate
        super
        check_property :task, Task
        check_optional_property :verifier, Verifier
        check_optional_property_list :dependencies, Task
      end
    end
  end
end

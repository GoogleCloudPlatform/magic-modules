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
      project: '{{ gcp_project }}',
      auth_kind: '{{ gcp_cred_kind }}',
      service_account_file: '{{ gcp_cred_file }}',
      name: '{{ resource_name }}'
    }.freeze

    EXAMPLE_DEFAULTS = {
      name: 'test_object',
      project: 'test_project',
      auth_kind: 'serviceaccount',
      service_account_file: '/tmp/auth.pem'
    }.freeze

    # Examples are used to generate the EXAMPLES block of Ansible documentation
    # and the integration tests.
    # Integration tests are a series of YAML tasks (standalone actions).
    # Integration tests are broken into three parts:
    # * a list of dependency tasks that should be run.
    # * a 'task' that is being tested (also used for EXAMPLES block)
    # * a verifier that will verify cloud status
    class Example < Api::Object
      attr_reader :task
      attr_reader :verifier
      attr_reader :dependencies
      attr_reader :facts

      def validate
        super
        default_value_property :facts, FactsTask.new
        default_value_property :verifier, FactsVerifier.new

        check_property :task, Task
        check_optional_property :verifier, Verifier
        check_optional_property_list :dependencies, Task
        check_optional_property :facts, Task

        @facts&.set_variable(self, :__example)
        @verifier.set_variable(self, :__example) if @verifier.respond_to?(:__example)
      end
    end

    # A Task represents a single Ansible action. This action is represented
    # as a standalone YAML block.
    class Task < Api::Object
      include Compile::Core
      attr_reader :name
      attr_reader :code
      attr_reader :scopes
      attr_reader :register

      def validate
        super
        check_property :name, String
        check_property :code, Hash
        check_optional_property_list :scopes, ::String
      end

      def build_test(state, object, noop = false)
        build_task(state, INTEGRATION_TEST_DEFAULTS, object, noop)
      end

      def build_example(state, object)
        build_task(state, EXAMPLE_DEFAULTS, object)
      end

      private

      # All of the arguments are used inside the ERB file, so we need
      # to disable rubocop complaining about unused methods
      # rubocop:disable Lint/UnusedMethodArgument
      def build_task(state, hash, object, noop = false)
        compile 'templates/ansible/tasks/task.yaml.erb'
      end
      # rubocop:enable Lint/UnusedMethodArgument

      def object_name_from_module_name(mod_name)
        product_name = mod_name.match(/gcp_[a-z]*_(.*)/).captures.first
        product_name.tr('_', ' ')
      end

      def dependency_name(dependency, resource)
        "#{dependency.downcase}-#{resource.downcase}"
      end

      def verbs
        {
          present: 'create',
          absent: 'delete'
        }
      end
    end

    # Verifiers verify that the Ansible modules actually created changes
    # in the cloud.
    # A Verifier has:
    # * A bash command.
    # * A failure check. If the bash command fails, that may not be enough
    #   to verify that the cloud status is correct.
    class Verifier < Api::Object
      include Compile::Core
      attr_reader :command
      attr_reader :failure

      def validate
        @failure ||= FailureCondition.new

        check_property :command, String
        check_property :failure, FailureCondition
      end

      # All of the arguments are used inside the ERB file, so we need
      # to disable rubocop complaining about unused methods
      # rubocop:disable Lint/UnusedMethodArgument
      def build_task(state, object)
        raise 'State must be present or absent' \
          unless %w[present absent].include? state

        compile 'templates/ansible/verifiers/bash.yaml.erb'
      end
      # rubocop:enable Lint/UnusedMethodArgument

      private

      def verbs
        {
          present: 'created',
          absent: 'deleted',
          facts: 'verify'
        }
      end
    end

    # A Verifier that doesn't build anything.
    class NoVerifier < Verifier
      attr_reader :reason
      def validate() end

      def build_task(_state, _object)
        ''
      end
    end

    # A Task that doesn't build anything.
    class NoTask < Task
      attr_reader :reason
      def validate() end

      def build_task(_state, _hash, _object, _noop = false)
        ''
      end
    end

    # Holds all information necessary to run a facts module and verify the
    # creation / deletion of a resource.
    # FactsVerifiers are verifiers in the sense that they verify GCP status.
    # They do not do this with bash commands, but with a Ansible facts module.
    # This verifier will look + an act a lot like a Task.
    class FactsVerifier < Verifier
      # Ruby YAML requires at least one value to create the object.
      attr_reader :noop

      attr_reader :__example
      include Compile::Core

      def validate
        @failure ||= FailureCondition.new
      end

      def build_task(_state, object)
        @parameters = build_parameters(object)
        compile 'templates/ansible/verifiers/facts.yaml.erb'
      end

      private

      def verbs
        {
          present: 'created',
          absent: 'deleted'
        }
      end

      def build_parameters(object)
        sample_code = @__example.task.code
        ignored_props = %w[project name]

        url_parts = object.uri_properties
                          .map(&:name)
                          .reject { |x| ignored_props.include? x }
        # Grab all code values for parameters
        object.all_user_properties
              .map(&:name)
              .select { |para| url_parts.include? para }
              .map { |para| { para => sample_code[para] } }
              .reduce({}, :merge)
      end

      def name_parameter
        compile_string(INTEGRATION_TEST_DEFAULTS, @__example.task.code['name']).join
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

    # Grpc gcloud commands seem to follow a similar pattern
    class GrpcFailureCondition < FailureCondition
      attr_reader :single
      attr_reader :plural

      def validate
        check_optional_property :single, ::String
        check_optional_property :plural, ::String

        @name ||= '{{ resource_name }}'
        @error = [
          "#{single.capitalize} not found:",
          "projects/{{ gcp_project }}/#{plural}/#{@name}"
        ].join(' ')
        super
      end
    end

    # A task for Ansible Facts.
    # Uses information from a traditional Ansible task.
    class FactsTask < Task
      # Ruby YAML requires at least one value to create the object.
      attr_reader :noop

      attr_reader :__example

      def validate; end

      def build_test(state, object, noop = false)
        @code = build_code(object, INTEGRATION_TEST_DEFAULTS)
        @name = ["gcp_#{object.__product.prefix[1..-1]}",
                 object.name.underscore,
                 'facts'].join('_')
        super(state, object, noop)
      end

      def build_example(state, object)
        @code = build_code(object, EXAMPLE_DEFAULTS)
        @name = ["gcp_#{object.__product.prefix[1..-1]}",
                 object.name.underscore,
                 'facts'].join('_')
        super(state, object)
      end

      private

      def build_code(object, hash)
        sample_code = @__example.task.code
        ignored_props = %w[project name]

        url_parts = object.uri_properties
                          .map(&:name)
                          .reject { |x| ignored_props.include? x }
        # Grab all code values for parameters
        code = object.all_user_properties
                     .map(&:name)
                     .select { |para| url_parts.include? para }
                     .map { |para| { para => sample_code[para] } }
                     .reduce({}, :merge)

        code['filters'] = ["name = #{hash[:name]}"]
        hash.each { |k, v| code[k.to_s] = v unless k == :name }
        code
      end
    end
  end
end

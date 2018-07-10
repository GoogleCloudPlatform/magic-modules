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

require 'google/string_utils'
require 'test_template'

module Puppet
  # A template for a test that follows the standard test workflow:
  #
  #   @pre
  #   - cleanup
  #   - create
  #   - create {again} : idempotency
  #   @change
  #   - delete
  #   - delete {again} : idempotency
  #   @post
  #
  # If 'pre' is defined it is attached to the beginning of the test, in case
  # some tests require non-standard startup.
  #
  # If 'post' is defined it is attached to the end of the test, in case some
  # tests require non-standard cleanup.
  #
  # If 'change' is defined it is attached after the create phases. Sometimes,
  # tests involve changing a resource before deleting it.
  class StandardTest < TestTemplate
    private

    # rubocop:disable Metrics/MethodLength
    def build_plan
      file = Google::StringUtils.underscore(@name)
      delete_name = @delete || "delete_#{file}.pp"
      create_name = @create || "#{file}.pp"
      @vars = {
        'name' => @name,
        'phases' => (@pre || []).concat(
          [
            {
              'name' => 'cleanup',
              'apply' => [{ 'run' => delete_name, 'exits' => [0, 2] }]
            },
            {
              'name' => 'cleanup{again}',
              'apply' => [{ 'run' => delete_name, 'exits' => 0 }]
            },
            {
              'name' => 'create',
              'apply' => [{ 'run' => create_name, 'exits' => 2 }]
            },
            {
              'name' => 'create{again}',
              'apply' => [{ 'run' => create_name, 'exits' => 0 }]
            }
          ]
        )
      }

      @vars['phases'] = @vars['phases'].concat(@change || [])
      @vars['phases'] = @vars['phases'].concat(
        [
          {
            'name' => 'delete',
            'apply' => [{ 'run' => delete_name, 'exits' => 2 }]
          },
          {
            'name' => 'delete{again}',
            'apply' => [{ 'run' => delete_name, 'exits' => 0 }]
          }
        ]
      ).concat(@post || [])

      apply_test_environment @vars

      @vars
    end
    # rubocop:enable Metrics/MethodLength
  end

  # A template for a test that follows the readonly test workflow:
  #
  #   @pre
  #   - run
  #   - run {again} : idempotency
  #   @post
  #
  # If 'pre' is defined it is attached to the beginning of the test, in case
  # some tests require non-standard startup.
  #
  # If 'post' is defined it is attached to the end of the test, in case some
  # tests require non-standard cleanup.
  class VirtualTest < TestTemplate
    private

    def build_plan
      vars = build_vars(Google::StringUtils.underscore(@name))
      apply_test_environment vars
      vars
    end

    def build_vars(file)
      {
        'name' => @name,
        'phases' => (@pre || []).concat(
          [
            {
              'name' => 'run',
              'apply' => [{ 'run' => "#{file}.pp", 'exits' => 0 }]
            },
            {
              'name' => 'run{again}',
              'apply' => [{ 'run' => "#{file}.pp", 'exits' => 0 }]
            }
          ]
        ).concat(@post || [])
      }
    end
  end
end

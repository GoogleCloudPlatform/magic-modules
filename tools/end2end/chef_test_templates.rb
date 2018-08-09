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

module Chef
  # A template for a test that follows the standard test workflow:
  #
  #   @pre
  #   - cleanup
  #   - create
  #   - create {again} : idempotency
  #   - delete
  #   - delete {again} : idempotency
  #   @post
  #
  # If 'post' is defined it is attached to the end of the test, in case some
  # tests require non-standard cleanup.
  # rubocop:disable Metrics/ClassLength
  class StandardTest < TestTemplate
    private

    DEFAULT_RESOURCE_COUNT = 2 # 1=resource + 1=auth

    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/MethodLength
    def build_plan
      init
      mod_name = "google-#{Google::StringUtils.underscore(@module)}"
      file = Google::StringUtils.underscore(@resource_name || @name)

      # Possible statements for creating an object.
      create_change = [
        'Chef Client finished,',
        "#{@affected['create']}/#{@resources['create']}",
        'resources updated'
      ].join(' ')
      create_no_change = [
        "Chef Client finished, 0/#{@resources['create']}",
        'resources updated'
      ].join(' ')

      # Possible statement for deleting an object.
      delete_change = [
        'Chef Client finished,',
        "#{@affected['delete']}/#{@resources['delete']}",
        'resources updated'
      ].join(' ')
      delete_no_change = [
        "Chef Client finished, 0/#{@resources['delete']}",
        'resources updated'
      ].join(' ')

      @vars = {
        'name' => @name,
        'phases' => (@pre || []).concat(
          [
            {
              'name' => 'cleanup',
              'apply' => [
                {
                  'run' => "#{mod_name}::#{@recipes['delete']}",
                  'outputs' => [
                    [
                      "Deleting #{@module}_#{file}[.*]",
                      "#{@module}_#{file}[.*] action delete (up to date)"
                    ],
                    (0..@affected['delete']).map do |index|
                      ['Chef Client finished,',
                       "#{index}/#{@resources['delete']} resources",
                       'updated'].join(' ')
                    end
                  ]
                }
              ]
            },
            {
              'name' => 'create',
              'apply' => [
                {
                  'run' => "#{mod_name}::#{@recipes['create']}",
                  'outputs' => [
                    ["Creating #{@module}_#{file}[.*]"],
                    [create_change]
                  ]
                }
              ]
            },
            {
              'name' => 'create{again}',
              'apply' => [
                {
                  'run' => "#{mod_name}::#{@recipes['create']}",
                  'outputs' => [
                    ["#{@module}_#{file}[.*] action create (up to date)"],
                    [create_no_change]
                  ]
                }
              ]
            }
          ]
        )
      }

      @vars['phases'] = @vars['phases'].concat(@change || [])
      @vars['phases'] = @vars['phases'].concat(
        [
          {
            'name' => 'delete',
            'apply' => [
              {
                'run' => "#{mod_name}::#{@recipes['delete']}",
                'outputs' => [
                  ["Deleting #{@module}_#{file}[.*]"],
                  [delete_change]
                ]
              }
            ]
          },
          {
            'name' => 'delete{again}',
            'apply' => [
              {
                'run' => "#{mod_name}::#{@recipes['delete']}",
                'outputs' => [
                  ["#{@module}_#{file}[.*] action delete (up to date)"],
                  [delete_no_change]
                ]
              }
            ]
          }
        ]
      ).concat(@post || [])

      apply_test_environment @vars

      @vars
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/PerceivedComplexity

    # Setup all values with proper defaults.
    # This ensures thats all values with be hashes with custom values or
    # defaults.
    # rubocop:disable Metrics/CyclomaticComplexity
    # rubocop:disable Metrics/PerceivedComplexity
    def init
      file = Google::StringUtils.underscore(@resource_name || @name)

      @recipes = setup_hash(@recipes)
      @recipes['create'] = "tests~#{file}" unless @recipes['create']
      @recipes['delete'] = "tests~delete_#{file}" unless @recipes['delete']

      @resources = setup_hash(@resource_count)
      @resources['create'] = DEFAULT_RESOURCE_COUNT unless @resources['create']
      @resources['delete'] = DEFAULT_RESOURCE_COUNT unless @resources['delete']

      @affected = setup_hash(@affected_count)
      unless @affected['create']
        @affected['create'] =
          DEFAULT_RESOURCE_COUNT - 1
      end

      return if @affected['delete']
      @affected['delete'] = DEFAULT_RESOURCE_COUNT - 1
    end
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/PerceivedComplexity

    # This will take in either a partially/fully setup Hash for testing
    # or a default value.
    # Will return the partially / fully Hash or a fully setup Hash with default
    # value.
    def setup_hash(item)
      return item if item.is_a? ::Hash
      {
        'create' => item,
        'delete' => item
      }
    end
  end

  # A template for a test that follows the readonly test workflow:
  #
  #   @pre
  #   - create
  #   @post
  #
  # If 'post' is defined it is attached to the end of the test, in case some
  # tests require non-standard cleanup.
  class VirtualTest < TestTemplate
    private

    # rubocop:disable Metrics/MethodLength
    def build_plan
      mod_name = "google-#{Google::StringUtils.underscore(@module)}"
      file = Google::StringUtils.underscore(@resource_name || @name)

      resources = @resource_count || 2 # 1=resource + 1=auth
      affected = @affected_count || 0 # readonly resources are never updated.
      changed_resource = \
        "Chef Client finished, #{affected}/#{resources} resources updated"
      create_name = @create || "tests~#{file}"
      @vars = {
        'name' => @name,
        'phases' => (@pre || []).concat(
          [
            {
              'name' => 'create',
              'apply' => [
                {
                  'run' => "#{mod_name}::#{create_name}",
                  'outputs' => [
                    ["#{@module}_#{file}[.*] action create (up to date)"],
                    [changed_resource]
                  ]
                }
              ]
            }
          ]
        )
      }

      apply_test_environment @vars

      @vars
    end
    # rubocop:enable Metrics/MethodLength
  end
  # rubocop:enable Metrics/ClassLength
end

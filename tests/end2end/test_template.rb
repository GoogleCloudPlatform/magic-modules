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

# Basic services for test templates.
class TestTemplate
  attr_reader :env # if specified the environment is spread to all tests
  attr_reader :pre
  attr_reader :post

  def validate
    raise "'name' parameter is required" if @name.nil? || @name.empty?
    log "#{@name} standard test"

    init if respond_to? :init
    build_plan
  end

  def [](variable)
    return @verifiers if variable == 'verifiers'
    @vars[variable]
  end

  def key?(variable)
    @vars.key?(variable)
  end

  private

  def apply_test_environment(vars)
    return if @env.nil?

    vars['phases'].map { |p| p['apply'] }
                  .flatten
                  .each { |a| a['env'] = @env.clone.merge(a['env'] || {}) }
  end

  def log(message)
    Logger.instance.log 'end2end', 'parser', message
  end
end

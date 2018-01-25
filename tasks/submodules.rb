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

require 'tasks/common'

TEST_RUNNER = {
  puppet: 'rake test',
  chef: 'rake test',
  ansible: 'exit 1'
}.freeze

def test_module(provider, mod)
  output = PROVIDER_FOLDERS[provider.to_sym] % mod
  %x(cd #{output} && #{TEST_RUNNER[provider]})
end

def all_test
  provider_list.map do |x|
    all_tasks_for_provider(x, 'test_mod:')
  end
end

def all_tasks_for_provider(prov, prefix = '')
  modules_for_provider(prov).map { |mod| "#{prefix}#{prov}:#{mod}" }
end

# Compiling Tasks
desc 'Test all modules'
multitask test_mod: all_test.flatten

namespace :test_mod do
  provider_list.each do |provider|
    # Each provider should default to compiling everything.
    desc "Test all modules for #{provider.capitalize}"
    multitask provider.to_sym => all_tasks_for_provider(provider)

    namespace provider.to_sym do
      modules_for_provider(provider).each do |mod|
        # Each module should have its own task for compiling.
        task mod.to_sym do
          test_module(provider, mod)
        end
      end
    end
  end
end

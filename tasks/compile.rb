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

def compile_module(provider, mod)
  output = PROVIDER_FOLDERS[provider.to_sym] % mod
  `./compiler -p products/#{mod} -e #{provider} -o #{output}`
end

def all_compile
  provider_list.map do |x|
    all_tasks_for_provider(x, 'compile:')
  end
end

def all_tasks_for_provider(prov, prefix = '')
  modules_for_provider(prov).map { |x| "#{prefix}#{prov}:#{x}".to_sym }
end

# Compiling Tasks
desc 'Compile all modules'
task compile: all_compile

namespace :compile do
  provider_list.each do |provider|
    # Each provider should default to compiling everything.
    desc "Compile all modules for #{provider.capitalize}"
    multitask provider.to_sym => all_tasks_for_provider(provider)

    namespace provider.to_sym do
      modules_for_provider(provider).each do |mod|
        # Each module should have its own task for compiling.
        desc "Compile the #{mod} module for #{provider.capitalize}"
        task mod.to_sym do
          compile_module(provider, mod)
        end
      end
    end
  end
end

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

# Configuration
$LOAD_PATH.unshift File.dirname(__FILE__)
Dir.chdir(File.dirname(__FILE__))

class Providers
  PROVIDER_FOLDERS = {
    ansible: 'build/ansible',
    ansible_facts: 'build/ansible',
    puppet: 'build/puppet/%s',
    chef: 'build/chef/%s',
    terraform: 'build/terraform'
  }.freeze

  def self.provider_list
    PROVIDER_FOLDER.keys
  end

  def initialize(name)
    @name = name
  end

  def products
    products = File.join(File.dirname(__FILE__), '..', 'products')
    files = Dir.glob("#{products}/**/#{@name}.yaml")
    files.map do |file|
      match = file.match(%r{^.*products\/([_a-z]*)\/.*yaml.*})
      match&.captures&.at(0)
    end.compact
  end
end

# Compiling Tasks
desc 'Compile all modules'
multitask compile: Providers.provider_list.map do |x|
  Providers.new(x).products.map { |y| "compile_#{x}_#{y}" }
end

namespace :compile do
  Providers.provider_list.each do |provider|
    # Each provider should default to compiling everything.
    desc "Compile all modules for #{provider.capitalize}"
    prov = Provider.new(provider)
    multitask provider.to_sym => prov.products.map { |m| "compile_#{provider}_#{m}" }

    namespace provider.to_sym do
      prov.products.each do |mod|
        # Each module should have its own task for compiling.
        desc "Compile the #{mod} module for #{provider.capitalize}"
        task mod.to_sym { compile_module(provider, mod) }
      end
    end
  end
end

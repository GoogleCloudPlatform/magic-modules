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

PROVIDER_FOLDERS = {
  ansible: 'build/ansible',
  terraform: 'build/terraform',
  inspec: 'build/inspec'
}.freeze

# Requires
require 'rspec/core/rake_task'
require 'rubocop/rake_task'
require 'tempfile'

RSpec::Core::RakeTask.new(:spec)
RuboCop::RakeTask.new

# Handles finding the list of products for a given provider.
class Providers
  def self.provider_list
    PROVIDER_FOLDERS.keys
  end

  # All possible products that exist.
  def self.all_products
    products = File.join(File.dirname(__FILE__), 'products')
    Dir.glob("#{products}/**/api.yaml")
  end

  def initialize(name)
    @name = name
  end

  def products
    products = File.join(File.dirname(__FILE__), 'products')
    files = Dir.glob("#{products}/**/#{@name}.yaml")
    files.map do |file|
      match = file.match(%r{^.*products\/([_a-z]*)\/.*yaml.*})
      match&.captures&.at(0)
    end.compact
  end

  def compile_module(mod)
    folder = format(PROVIDER_FOLDERS[@name.to_sym], mod: mod)
    flag = "COMPILER_#{folder.gsub('build/', '').tr('/', '_').upcase}_OUTPUT"
    output = ENV[flag] || format(PROVIDER_FOLDERS[@name.to_sym], mod: mod)
    %x(bundle exec compiler -p products/#{mod} -e #{@name} -o #{output})
  end

  def compilation_targets
    products.map { |prod| "compile:#{@name}:#{prod}" }
  end
end

# API Linter Tasks
desc 'Runs the API Linter'
RSpec::Core::RakeTask.new(:lint) do |t|
  t.rspec_opts = '--pattern tools/linter/run.rb'
end

# Test Tasks
desc 'Run all of the MM tests (rubocop, rspec)'
multitask test: %w[rubocop spec]

# Compiling Tasks
compile_list = Providers.provider_list.map do |x|
  Providers.new(x).compilation_targets
end.flatten

desc 'Compile all modules'
multitask compile: compile_list

namespace :compile do
  Providers.provider_list.each do |provider|
    # Each provider should default to compiling everything.
    desc "Compile all modules for #{provider.capitalize}"
    prov = Providers.new(provider)
    multitask provider.to_sym => prov.compilation_targets

    namespace provider.to_sym do
      prov.products.each do |mod|
        # Each module should have its own task for compiling.
        desc "Compile the #{mod} module for #{provider.capitalize}"
        task mod.to_sym do
          prov.compile_module(mod)
        end
      end
    end
  end
end

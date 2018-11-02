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

# Requires for YAML linting.
require 'api/async'
require 'api/bundle'
require 'api/product'
require 'api/resource'
require 'api/type'
require 'compile/core'
require 'google/yaml_validator'

RSpec::Core::RakeTask.new(:spec)
RuboCop::RakeTask.new

# YAML Linting
# This class calls our provider code to get the printed contents of the
# compiled YAML. We run the linter on this printed version (so, no embedded
# ERB)
class YamlLinter
  include Compile::Core

  def yaml_contents(file)
    source = compile(file)
    config = Google::YamlValidator.parse(source)
    unless config.class <= Api::Product
      raise StandardError, "#{file} is #{config.class}"\
        ' instead of Api::Product' \
    end
    # Compile step #2: Now that we have the target class, compile with that
    # class features
    config.compile(file, 0)
  end
end

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

# Test Tasks
desc 'Run all of the MM tests (rubocop, rspec)'
multitask test: %w[rubocop spec]

desc 'Lints all of the compiled YAML files'
task :yamllint do
  Providers.all_products.each do |file|
    tempfile = Tempfile.new
    tempfile.write(YamlLinter.new.yaml_contents(file))
    tempfile.rewind
    puts %x(yamllint -c #{File.join(File.dirname(__FILE__), '.yamllint')} #{tempfile.path})
    tempfile.close
    tempfile.unlink
  end
end

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

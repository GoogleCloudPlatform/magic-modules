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

PROVIDER_FOLDERS = {
  ansible: 'build/ansible',
  puppet: 'build/puppet/%s',
  chef: 'build/chef/%s',
  terraform: 'build/terraform'
}.freeze

# Represents a list of tasks where each task acts upon all modules.
class MMRakeTasks
  def add_task(name, &block)
    task_layout(name, &block)
  end

  private

  def provider_list
    PROVIDER_FOLDERS.keys
  end

  # Given a prefix, return all tasks under that prefix.
  def all_tasks_for_prefix(prefix)
    provider_list.map do |x|
      all_tasks_for_provider_and_prefix(x, prefix)
    end
  end

  # Given a provider and prefix give all tasks for that provider + prefix
  def all_tasks_for_provider_and_prefix(prov, prefix = '')
    modules_for_provider(prov).map { |x| "#{prefix}:#{prov}:#{x}" }
  end

  def modules_for_provider(provider)
    products = File.join(File.dirname(__FILE__), '..', 'products')
    files = Dir.glob("#{products}/**/#{provider}.yaml")
    files.map do |file|
      match = file.match(%r{^.*products\/([_a-z]*)\/.*yaml.*})
      match&.captures&.at(0)
    end.compact
  end

  # rubocop:disable Metrics/AbcSize
  # rubocop:disable Metrics/MethodLength
  # rubocop:disable Style/EvalWithLocation
  def task_layout(name, &block)
    @self_before_instance_eval = eval 'self', block.binding
    instance_eval do
      # Compiling Tasks
      desc "#{name.capitalize} all modules"
      task :compile do
        run_multiple_tasks(all_tasks_for_prefix(name).flatten)
      end

      namespace :compile do
        provider_list.each do |provider|
          # Each provider should default to compiling everything.
          desc "Compile all modules for #{provider.capitalize}"
          task provider.to_sym do
            run_multiple_tasks(
              all_tasks_for_provider_and_prefix(provider, name)
            )
          end

          namespace provider.to_sym do
            modules_for_provider(provider).each do |mod|
              # Each module should have its own task for compiling.
              desc [
                "#{name.capitalize} the #{mod} module for",
                provider.capitalize
              ].join(' ')
              task mod.to_sym do
                yield(provider, mod)
              end
            end
          end
        end
      end
    end
  end
  # rubocop:enable Metrics/AbcSize
  # rubocop:enable Metrics/MethodLength
  # rubocop:enable Style/EvalWithLocation

  # rubocop:disable Style/MethodMissing
  def method_missing(method, *args, &block)
    @self_before_instance_eval.send method, *args, &block
  end
  # rubocop:enable Style/MethodMissing
end

# Run multiple tasks in parallel and report back failures.
def run_multiple_tasks(tasks)
  @failures = []
  tasks.map do |task|
    Thread.new do
      Rake::Task[task].invoke
    rescue StandardError
      @failures << task
    end
  end.map(&:join)
  if @failures.empty?
    puts 'Success!'
  else
    puts '##################################'
    puts '#          FAILURES              #'
    puts '##################################'
    puts @failures.join("\n").to_s
  end
end

# rubocop:disable Style/CommandLiteral
# rubocop:disable Style/SpecialGlobalVars
def run_command(command)
  `#{command}`
  raise 'Failed' if $?.exitstatus != 0
end
# rubocop:enable Style/CommandLiteral
# rubocop:enable Style/SpecialGlobalVars

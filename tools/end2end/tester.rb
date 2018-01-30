#!/usr/bin/env ruby
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

# Executes all examples against a real Google Cloud Platform project. The
# account requires to have Owner (or Editor) to all resources being tested.
#
# Usage: tools/end2end/run <project...>
#        <project> can be of the form:
#           - <provider>:<product>          e.g. puppet:dns  only this module
#           - <provider>:                   e.g. puppet:     all puppet modules
#           - :<product>                    e.g.       :dns  all DNS providers
#   to test a single object:
#           - :<product>:<object>           e.g.       :dns:ManagedZone
#           - <provider>:<product>:<object> e.g. puppet:dns:ManagedZone
#
# Environment variables:
#   - PARALLEL=true|false controls how tests are run

$LOAD_PATH.unshift File.dirname(__FILE__)
# Add root so we have google/<library> available
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '..', '..')

require 'benchmark'
require 'chef_test_templates'
require 'chef_tester'
require 'logger'
require 'puppet_test_templates'
require 'puppet_tester'
require 'singleton'
require 'yaml'

# A simple class to track execution times.
class Timers
  include Singleton

  def initialize
    @timers = {}
  end

  def measure(context, &block)
    @timers[context] = Benchmark.measure(&block)
  end

  def print
    @timers.each do |k, v|
      Logger.instance.log(
        'timer', k.join('/'),
        format('%<phase>-23s : %<elapsed>s',
               phase: k.join('/'),
               elapsed: Time.at(v.total).utc .strftime('%H:%M:%S.%L'))
      )
    end
  end
end

module Colors
  B_RED = "\e[37;41m".freeze
  B_GREEN = "\e[37;42m".freeze
  B_YELLOW = "\e[37;43m".freeze
  NC = "\e[0m".freeze
end

def log(provider, product, message)
  Logger.instance.log provider, product, message
end

def bundle_install
  Open3.popen2e(%w[bundle install]) do |_, std_out_and_err, thread|
    output = Array[*std_out_and_err]
    output.each { |line| log 'bundle', 'install', line }

    exit_code = thread.value.exitstatus
    raise 'Bundle install failed' unless exit_code.zero?
  end
end

def true?(obj)
  obj.to_s == 'true'
end

def todo?(tests)
  tests.any? do |t|
    t.tests.any? do |u|
      u['phases'].any? { |p| p['apply'].any? { |a| a.key?('todo') } }
    end
  end
end

# Placeholder for shared configuration constants.
module Config
  # The RUN_ID is available as a variable for the test environments as
  # {{run_id}}.
  RUN_ID = ENV['RUN_ID'] || (rand * 1_000_000).to_i
  log 'end2end', nil, "Run ID: #{RUN_ID}"
end

# To have the output organized at the end (as it runs in parallel it will be
# mixed between runs you can execute like this:
#
#   tools/end2end/run | tee output.log; cat output.log | sort
PARALLEL = true?(ENV['PARALLEL'] || true)
log 'end2end', nil, PARALLEL ? 'Parallel mode' : 'Serial mode'

Dir.chdir File.dirname(__FILE__)

# Load YAML files from submodules and plan.yaml in root
plans = Dir[File.join('..', '..', 'products', '**', '*-e2e.yaml')]
plans << 'plan.yaml' if File.exist? 'plan.yaml'

tests = plans.map do |file|
  YAML.safe_load(File.read(file),
                 [Puppet::Tester, Puppet::StandardTest, Puppet::VirtualTest,
                  Chef::Tester, Chef::StandardTest, Chef::VirtualTest])
end.flatten

tests.each { |t| t.test_matrix = ARGV }
     .each { |t| t.validate if t.respond_to?(:validate) }

# Ensure we are not missing any products
missing = Dir[File.join('..', '..', 'build', '*', '*')]
          .reject { |f| f.start_with?('../../build/presubmit/') }
          .reject { |f| f.end_with?('/auth') }
          .reject { |f| f.end_with?('_bundle') }
          .select do |f|
            tests.select do |t|
              f == File.join('..', '..', 'build', t.provider, t.product)
            end.empty?
          end

bad = Colors::B_RED
nc = Colors::NC
puts "#{bad} The following modules are missing: #{missing} #{nc}" unless \
  missing.empty?

Timers.instance.measure(%w[bundle install]) do
  bundle_install
end

unless ARGV.empty?
  tests = tests.reject do |t|
    ARGV.select do |a|
      parts = a.split(':')
      (parts[0].empty? || parts[0] == t.provider) \
        && (parts.length < 2 || parts[1].empty? || parts[1] == t.product)
    end.empty?
  end
end

results = nil

if PARALLEL
  # Run everything in parallel.
  test_tasks = tests.map { |t| Thread.new { t.test } }
  results = test_tasks.map(&:value)
else
  results = tests.map(&:test)
end

if results.nil? || results.empty?
  log nil, nil, 'Nothing was selected to be tested'
  success = false
else
  success = results.reject { |v| v }.empty?
end

Timers.instance.print

if success
  good = Colors::B_GREEN
  nc = Colors::NC
  log \
    nil, nil,
    [
      "#{good} _______ _     _ _______ _______ _______ _______ _______ #{nc}",
      "#{good} |______ |     | |       |       |______ |______ |______ #{nc}",
      "#{good} ______| |_____| |_____  |_____  |______ ______| ______| #{nc}",
      "#{good}                                                         #{nc}"
    ].join("\n")
  if todo?(tests)
    warn = Colors::B_YELLOW
    log \
      nil, nil,
      "#{warn} There are unaddressed TODO items. Just saying...        #{nc}"
  end
else
  bad = Colors::B_RED
  nc = Colors::NC
  log \
    nil, nil,
    [
      "#{bad} _______ _______ _____        _______ ______      #{nc}",
      "#{bad} |______ |_____|   |   |      |______ |     \\     #{nc}",
      "#{bad} |       |     | __|__ |_____ |______ |_____/     #{nc}",
      "#{bad}                                                  #{nc}"
    ].join("\n")
  if todo?(tests)
    warn = Colors::B_YELLOW
    log \
      nil, nil,
      "#{warn} There are unaddressed TODO items. Just saying... #{nc}"
  end
end

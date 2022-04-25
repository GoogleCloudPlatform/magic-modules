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

# Requires
require 'rspec/core/rake_task'
require 'rubocop/rake_task'

RSpec::Core::RakeTask.new(:spec)
RuboCop::RakeTask.new

# API Linter Tasks
desc 'Runs the API Linter'
RSpec::Core::RakeTask.new(:lint) do |t|
  t.rspec_opts = '--pattern tools/linter/run.rb'
end

# Test Tasks
desc 'Run all of the MM tests (rubocop, rspec)'
multitask test: %w[rubocop spec]

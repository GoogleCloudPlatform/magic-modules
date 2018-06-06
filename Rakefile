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

$LOAD_PATH.unshift File.dirname(__FILE__)
Dir.chdir(File.dirname(__FILE__))

# Load all tasks from tasks/
require 'parallel_tests'
require 'rubocop/rake_task'
require 'tasks/common'

@task = MMRakeTasks.new

###############################################################################
# Tasks
###############################################################################
@task.add_task('compile') { |prov, mod| compile_module(prov, mod) }

namespace 'test' do
  desc 'Run RSpec code example'
  task :spec do |_, _|
    abort unless system('parallel_rspec spec/')
  end

  RuboCop::RakeTask.new
end

###############################################################################
# Helpers
###############################################################################
def compile_module(provider, mod)
  folder = PROVIDER_FOLDERS[provider.to_sym] % mod
  flag = "COMPILER_#{folder.gsub('build/', '').tr('/', '_').upcase}_OUTPUT"
  output = ENV[flag] || (PROVIDER_FOLDERS[provider.to_sym] % mod)
  run_command(
    "bundle exec compiler -p products/#{mod} -e #{provider} -o #{output}"
  )
end

def test_module(provider, mod)
  output = PROVIDER_FOLDERS[provider.to_sym] % mod
  %x(cd #{output} && #{TEST_RUNNER[provider]})
end

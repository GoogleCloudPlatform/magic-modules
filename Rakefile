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
Dir[File.join('tasks', '*.rb')].reject { |p| File.directory? p }
                               .each do |f|
  require f
end

# Find all tasks under the test namespace
# Ignore those with multiple levels like rubocop:auto_correct
tests = Rake.application.tasks.select do |task|
  /^test:[a-z]*$/ =~ task.name
end.map(&:name)

multitask 'test' => tests

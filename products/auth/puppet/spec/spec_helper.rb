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

# Our default timezone is UTC, to avoid local time compromise test code seed
# generation.
ENV['TZ'] = 'UTC'

require 'simplecov'
SimpleCov.start

$LOAD_PATH.unshift(File.expand_path('.'))
$LOAD_PATH.unshift(File.expand_path('../puppet-google-auth/lib'))

require 'fakeweb'
require 'fake_web/registry'
FakeWeb.allow_net_connect = false

files = []
files << 'spec/copyright.rb'
files << 'spec/copyright_spec.rb'
files << File.join('lib', '**', '*.rb')

# Require all files so we can track them via code coverage
Dir[*files].reject { |p| File.directory? p }
           .each do |f|
             puts "Auto requiring #{f}" \
               if ENV['RSPEC_DEBUG']
             require f
           end

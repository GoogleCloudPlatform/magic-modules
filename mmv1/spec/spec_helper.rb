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

RSpec.configure do |config|
  config.mock_with :mocha
end

$LOAD_PATH.unshift(File.expand_path('.'))

ENV['GOOGLE_LOGGER'] = '0' unless ENV['RSPEC_DEBUG']

%w[api google provider].each do |subsystem|
  Dir[File.join(subsystem, '**', '*.rb')].reject { |p| File.directory? p }
                                         .each do |f|
                                           puts "Auto requiring #{f}" \
                                             if ENV['RSPEC_DEBUG']
                                           require f
                                         end
end

RSpec.configure do |config|
  config.mock_with :mocha
end

RSpec::Matchers.define :have_attribute_of_length do |expected|
  match do |actual|
    expected.each do |name, length|
      my_len = actual.send(name).length
      raise "#{name} should have #{length} elements, but found only #{my_len}" \
        unless my_len == length
    end
  end
end

RSpec::Matchers.define :contain_array do |expected|
  match do |actual|
    while actual.include?(expected[0])
      start = actual.index(expected[0])
      actual = actual.drop(start) unless start.nil?
      return true if actual[0, expected.size] == expected

      actual = actual.drop(1)
    end
    false
  end
end

require 'pp'

Google::LOGGER.info 'Running tests'

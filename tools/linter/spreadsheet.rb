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
################
# Discovery Doc Builder
#
# This script takes in a yaml file with a Docs object that
# describes which Discovery APIs are being taken in.
#
# The script will then build api.yaml files using
# the Discovery API

# Load everything from MM root.
$LOAD_PATH.unshift File.join(File.dirname(__FILE__), '../../')
Dir.chdir(File.join(File.dirname(__FILE__), '../../'))

require 'google/logger'
require 'tools/linter/discovery'
require 'tools/linter/fetcher'
require 'tools/linter/tests/test_runner'
require 'tools/linter/spreadsheet/csv_formatter'

require 'yaml'
require 'rspec'

# Configuration
doc_file = 'tools/linter/docs.yaml'
RSpec.configure do |c|
  c.add_formatter(CsvFormatterForMM)
  c.inclusion_filter = [:property]
end
Google::LOGGER.level = Logger::ERROR
VALID_KEYS = %w[filename url].freeze

# Running tests.
docs = YAML.safe_load(File.read(doc_file))
docs.each do |doc|
  raise "#{doc.keys} not in #{VALID_KEYS}" unless doc.keys.sort == %w[filename url]

  # Run tests on regular API
  api = ApiFetcher.api_from_file(doc['filename'])
  builder = Discovery::Builder.new(doc['url'], api.objects.map(&:name))
  run_tests(builder, api, {property: true})

  # Run tests on TF API
  api = ApiFetcher.provider_from_file(doc['filename'], 'terraform')
  run_tests(builder, api, {property: true}, {provider: 'terraform'})
end

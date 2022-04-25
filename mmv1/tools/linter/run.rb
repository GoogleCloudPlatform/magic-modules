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

require 'yaml'
require 'rspec'

Google::LOGGER.level = Logger::ERROR
REQUIRED_KEYS = %w[filename url product version].freeze
VALID_KEYS = (REQUIRED_KEYS + %w[aliases]).freeze

doc_file = 'tools/linter/docs.yaml'
docs = YAML.safe_load(File.read(doc_file))

docs.each do |doc|
  raise "Missing required keys #{REQUIRED_KEYS} for #{doc['product']}" \
    unless REQUIRED_KEYS & doc.keys == REQUIRED_KEYS

  doc.keys.each do |key|
    raise "#{key} is an invalid key." unless VALID_KEYS.include? key
  end

  api = ApiFetcher.api_from_file(doc['filename'])
  builder = Discovery::Builder.new(doc, api.objects.map(&:name))
  run_tests(builder, api, resource: true, property: true)
end

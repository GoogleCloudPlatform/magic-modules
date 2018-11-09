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

require 'tools/linter/api'
require 'tools/linter/discovery'
require 'tools/linter/test_runner'
require 'tools/linter/tests'

require 'yaml'
require 'rspec'

VALID_KEYS = %w[filename url].freeze

doc_file = 'tools/linter/docs.yaml'
docs = YAML::load(File.read(doc_file))

docs.each do |doc|
  raise "#{doc.keys} not in #{VALID_KEYS}" unless doc.keys.sort == %w[filename url]
  api = ApiFetcher.new(doc['filename']).fetch
  builder = DiscoveryBuilder.new(doc['url'], api.objects.map(&:name))
  builder.resources.each do |disc_res|
    api_obj = api.objects.select { |p| p.name == disc_res.name }.first

    # RSpec tests begin here.
    RSpec.describe disc_res.name do
      TestRunner.new(disc_res, api_obj).run do |disc_prop, api_prop, name|
        context name do
          include_examples "property_tests", disc_prop, api_prop
        end
      end
    end
  end
end

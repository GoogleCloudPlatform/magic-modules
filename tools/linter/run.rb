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
require 'tools/linter/tests'

require 'yaml'
require 'rspec'

Google::LOGGER.level = Logger::ERROR
VALID_KEYS = %w[filename url].freeze

doc_file = 'tools/linter/docs.yaml'
docs = YAML.safe_load(File.read(doc_file))

docs.each do |doc|
  raise "#{doc.keys} not in #{VALID_KEYS}" unless doc.keys.sort == %w[filename url]
  api = ApiFetcher.api_from_file(doc['filename'])
  builder = Discovery::Builder.new(doc['url'], api.objects.map(&:name))

  # First context: product name
  RSpec.describe api.prefix do
    builder.resources.each do |disc_resource|
      api_obj = api.objects.select { |p| p.name == disc_resource.name }.first
      # Second context: resource name
      describe disc_resource.name do
        # Run all resource tests on this resource
        include_examples 'resource_tests', disc_resource, api_obj
        PropertyFetcher.fetch_property_pairs(disc_resource.properties,
                                             api_obj.all_user_properties) \
                                            do |disc_prop, api_prop, name|
          # Third context: property name
          context name do
            # Run all tests on this property
            include_examples 'property_tests', disc_prop, api_prop
          end
        end
      end
    end
  end
end

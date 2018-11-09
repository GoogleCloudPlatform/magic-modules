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

require 'optparse'
require 'yaml'

doc_file = 'tools/linter/docs.yaml'
OptionParser.new do |opts|
  opts.banner = "Discovery doc runner. Usage: run.rb [docs.yaml]"
  opts.on("-f", "--file [file]") { |file| doc_file = file }
end.parse!

docs = YAML::load(File.read(doc_file))

docs.each do |doc|
  builder = DiscoveryBuilder.new(doc['url'], doc['objects'])
  api = ApiFetcher.new(doc['filename']).fetch
  builder.resources.each do |disc_res|
    TestRunner.new(disc_res, api.objects.select { |p| p.name == disc_res.name }.first).run
  end
end

#! /usr/bin/env ruby
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

require 'puppet_forge'
require 'puppet_blacksmith'
require 'json'

config = JSON.parse(STDIN.read)
unless config['source'].key? 'module_name'
  raise 'You need to define `module_name`'
end
raise 'This is being called without an output directory.' if ARGV.empty?

module_name = config['source']['module_name']
release = PuppetForge::Module.find(module_name).releases.first.version
major, minor, patch = release.split('.')

if major.nil? || minor.nil? || patch.nil?
  raise "Cowardly refusing to work with non-semver release ID #{release}"
end

if config['params']['patch_bump'] == true
  patch += 1
else
  patch = 0
  minor += 1
end

output_folder = ARGV[0]
metadata = JSON.parse(File.open(File.join(output_folder, 'metadata.json')).read)
unless metadata['name'] == module_name
  raise "Cowardly refusing to push #{metadata['name']} to #{module_name}"
end
metadata['version'] = "#{major}.#{minor}.#{patch}"
File.write(File.join(output_folder, 'metadata.json'), JSON.dump(metadata))

Dir.chdir(output_folder) { %x(puppet module build) }

Blacksmith::Forge.initialize(config['source']['username'],
                             config['source']['password'])
Blacksmith::Forge.push!(metadata['name'].split('-').last,
                        File.join(output_folder,
                                  'pkg',
                                  "#{metadata['name']}-#{metadata['version']}" \
                                  '.tar.gz'))
puts JSON.dump('version' => { 'release' => metadata['version'] })

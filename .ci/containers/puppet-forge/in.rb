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
require 'json'

config = JSON.parse(STDIN.read)
raise 'You need to define `module_name`' unless config['source'].key? 'module_name'
raise 'Surprised by lack of ARGV[0]' if ARGV[0].nil?
release_slug = "#{config['source']['module_name']}-#{config['version']['release']}"
release = PuppetForge::Release.find(release_slug)
release_tarball = release_slug + '.tar.gz'
release.download(Pathname(release_tarball))
release.verify(Pathname(release_tarball))
PuppetForge::Unpacker.unpack(release_tarball, ARGV[0], '/tmp/resource_download')

puts JSON.dump('version' => { 'release' => release.version })

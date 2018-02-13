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

PROVIDER_FOLDERS = {
  ansible: 'build/ansible',
  puppet: 'build/puppet/%s',
  chef: 'build/chef/%s',
  terraform: 'build/terraform'
}.freeze

# Give a list of all providers served by MM
def provider_list
  PROVIDER_FOLDERS.keys
end

# Give a list of all products served by a provider
def modules_for_provider(provider)
  products = File.join(File.dirname(__FILE__), '..', 'products')
  files = Dir.glob("#{products}/**/#{provider}.yaml")
  files.map do |file|
    match = file.match(%r{^.*products\/([a-z]*)\/.*yaml.*})
    match&.captures&.at(0)
  end.compact
end

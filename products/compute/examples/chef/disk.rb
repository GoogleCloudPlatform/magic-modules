<% if false # the license inside this if block assertains to this file -%>
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
<% end -%>
<% unless name == 'README.md' -%>

<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

gcompute_zone 'us-central1-a' do
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<% end -%>
gcompute_disk <%= example_resource_name('data-disk-1') -%> do
  action :create
  size_gb 50
  disk_encryption_key({
    raw_key: 'SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0='
  })
  zone 'us-central1-a'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

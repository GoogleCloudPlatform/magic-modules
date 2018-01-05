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
<% if name != 'README.md' -%>

<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

gcompute_region <%= example_resource_name('some-region') -%> do
  action :create
  r_label 'us-west1'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<% else # name == README.md -%>
# Subnetwork requires a network and a region, so define them in your recipe:
#   - gcompute_region 'some-region' do ... end
<% end # name == README.md -%>
gcompute_subnetwork <%= example_resource_name('servers') -%> do
  action :delete
  region <%= example_resource_name('some-region') %>
  project 'google.com:graphite-playground'
  credential 'mycred'
end

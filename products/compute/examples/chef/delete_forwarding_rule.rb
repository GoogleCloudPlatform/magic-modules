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

gcompute_region <%= example_resource_name('some-region') -%> do
  name 'us-west1'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_address <%= example_resource_name('some-address') -%> do
  action :create
  region <%= example_resource_name('some-region') %>
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_target_pool <%= example_resource_name('target-pool') -%> do
  action :create
  region <%= example_resource_name('some-region') %>
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<% end # name == README.md -%>
gcompute_forwarding_rule <%= example_resource_name('fwd-rule-test') -%> do
  action :delete
  ip_address gcompute_address_ref(
    <%= example_resource_name('some-address') -%>,
    'us-west1', 'google.com:graphite-playground'
  )
  ip_protocol 'TCP'
  port_range '80'
  target <%= example_resource_name('target-pool') %>
  region <%= example_resource_name('some-region') %>
  project 'google.com:graphite-playground'
  credential 'mycred'
end

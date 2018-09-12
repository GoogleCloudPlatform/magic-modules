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

<% end -%>
<% if name == "README.md" -%>
# Route requires a network so define one in your recipe:
#   - gcompute_network 'my-network' do ... end
<% else # name == README.md -%>
gcompute_network <%= example_resource_name('my-network') -%> do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

<% end # name == README.md -%>
gcompute_route <%= example_resource_name('corp-route') -%> do
  action :create
  dest_range '192.168.6.0/24'
  next_hop_gateway 'global/gateways/default-internet-gateway'
  tags %w[backends databases] # %w[] best for single words. use ['.'] w/ spaces
  network <%= example_resource_name('my-network') %>
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

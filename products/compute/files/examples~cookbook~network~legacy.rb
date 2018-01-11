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

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

raise "Missing parameter 'network_id'. Please read docs at #{__FILE__}" \
  unless ENV.key?('network_id')
<% end -%>
# The environment variable 'network_id' defines a suffix for a network name when
# using this example. If running from the command line, you can pass this suffix
# in via the command line:
#
# network_id="some_suffix" chef-client -z --runlist \
#   "recipe[gcompute::examples~network~auto]"
puts 'Creating network in Legacy mode'
<% res_name = 'mynetwork-#{ENV[\'network_id\']}' -%>
gcompute_network <%= example_resource_name(res_name) -%> do
  # On a legacy network you cannot specify the auto_create_subnetworks
  # parameter.
  # | auto_create_subnetworks => false,
  action :create
  ipv4_range '192.168.0.0/16'
  gateway_ipv4 '192.168.0.1'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<%# The license inside this block applies to this file
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
-%>
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

# TODO(nelsonjr): Create a test case to verify that errors are properly
# displayed, such as a block like the one below, which will fail because you
# cannot specify auto_create_subnetworks and ipv4Range at the same time:
# | gcompute_network { "mynetwork-${network_id}":
# |   auto_create_subnetworks => true,
# |   ipv4_range              => '192.168.0.0/16',
# |   gateway_ipv4            => '192.168.0.1',
# |   project                 => $project, # e.g. 'my-test-project'
# |   credential              => 'mycred',
# | }

notice('Creating network with automatically assigned subnetworks')
<% end # name == README.md -%>
gcompute_network { <%= example_resource_name('mynetwork-${network_id}') -%>:
  auto_create_subnetworks => true,
  project                 => $project, # e.g. 'my-test-project'
  credential              => 'mycred',
}

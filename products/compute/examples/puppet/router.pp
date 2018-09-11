<% if false # the license inside this if block pertains to this file -%>
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
<% end -%>
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<% end # name == README.md -%>
<% if name == "README.md" -%>
# Router requires a network and a region, so define them in your manifest:
#   - gcompute_network { 'my-network': ensure => present }
#   - gcompute_region { 'some-region': ... }
<% else # name == README.md -%>
gcompute_network { <%= example_resource_name('my-network') -%>:
  ensure     => present,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

gcompute_region { <%= example_resource_name('some-region') -%>:
  name       => 'us-west1',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

<% end # name == README.md -%>
gcompute_router { <%= example_resource_name('my-router') -%>:
  ensure     => present,
  network    => <%= example_resource_name('my-network') -%>,
  bgp        => {
    asn                  => 64514,
    advertise_mode       => 'CUSTOM',
    advertised_groups    => ['ALL_SUBNETS'],
    advertised_ip_ranges => [
      {
        range => '1.2.3.4',
      },
      {
        range => '6.7.0.0/16',
      }
    ]
  },
  region     => <%= example_resource_name('some-region') -%>,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

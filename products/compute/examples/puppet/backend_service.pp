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
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<% end # name == README.md -%>
<% if name == "README.md" -%>
# Backend Service requires various other services to be setup beforehand. Please
# make sure they are defined as well:
#   - gcompute_instance_group { ... }
#   - Health check
<% else # name == README.md -%>
gcompute_zone { 'us-central1-a':
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

gcompute_instance_group { <%= example_resource_name('my-puppet-masters') -%>:
  ensure     => present,
  zone       => 'us-central1-a',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

gcompute_health_check { <%= example_resource_name('example-hc') -%>:
  ensure              => present,
  type                => 'TCP',
  tcp_health_check    => {
    port_name => 'service-health',
    request   => 'ping',
    response  => 'pong',
  },
  healthy_threshold   => 10,
  timeout_sec         => 2,
  unhealthy_threshold => 5,
  project             => $project, # e.g. 'my-test-project'
  credential          => 'mycred',
}

<% end # name == README.md -%>
gcompute_backend_service { <%= example_resource_name('my-app-backend') -%>:
  ensure        => present,
  backends      => [
    { group => <%= example_resource_name('my-puppet-masters') -%> },
  ],
  enable_cdn    => true,
  health_checks => [
    'example-hc'
  ],
  project       => $project, # e.g. 'my-test-project'
  credential    => 'mycred',
}

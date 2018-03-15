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
gdns_managed_zone { <%= example_resource_name('some-managed-zone') -%>:
  ensure      => present,
  name        => 'testzone-4-com',
  dns_name    => 'testzone-4.com.',
  description => 'Test Example Zone',
  project     => 'google.com:graphite-playground',
  credential  => 'mycred',
}

<% res_name = 'www.testzone-4.com.' -%>
gdns_resource_record_set { <%= example_resource_name(res_name) -%>:
  ensure       => present,
  managed_zone => <%= example_resource_name('some-managed-zone') -%>,
  type         => 'A',
  ttl          => 600,
  target       => [
    '10.1.2.3',
    '40.5.6.7',
    '80.9.10.11'
  ],
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

<% res_name = 'sites.testzone-4.com.' -%>
gdns_resource_record_set { <%= example_resource_name(res_name) -%>:
  ensure       => present,
  managed_zone => <%= example_resource_name('some-managed-zone') -%>,
  type         => 'CNAME',
  target       => 'www.testzone-4.com.',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

<% res_name = 'deleteme.testzone-4.com.' -%>
gdns_resource_record_set { <%= example_resource_name(res_name) -%>:
  ensure       => absent,
  managed_zone => <%= example_resource_name('some-managed-zone') -%>,
  type         => 'A',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

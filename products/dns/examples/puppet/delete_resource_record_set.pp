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

# An example Puppet manifest that ensures Google Cloud Computing DNS Resource
# Record Set in a project do not exist.

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gdns_managed_zone { <%= example_resource_name('testzone-4-com') -%>:
  ensure      => present,
  name        => 'testzone-4-com',
  dns_name    => 'testzone-4.com.',
  description => 'Test Example Zone',
  project     => 'google.com:graphite-playground',
  credential  => 'mycred',
}

gdns_resource_record_set { <%= example_resource_name('www.testzone-4.com.') -%>:
  ensure       => absent,
  managed_zone => <%= example_resource_name('testzone-4-com') -%>,
  type         => 'A',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

<% res_name = 'sites.testzone-4.com.' -%>
gdns_resource_record_set { <%= example_resource_name(res_name) -%>:
  ensure       => absent,
  managed_zone => <%= example_resource_name('testzone-4-com') -%>,
  type         => 'CNAME',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

<% end # name == README.md -%>
<% res_name = 'deleteme.testzone-4.com.' -%>
gdns_resource_record_set { <%= example_resource_name(res_name) -%>:
  ensure       => absent,
  managed_zone => <%= example_resource_name('testzone-4-com') -%>,
  type         => 'A',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

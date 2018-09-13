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

# An example Puppet manifest that creates a Google Cloud Computing DNS Managed
# Zone in a project.

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

# Ensures a managed zone exists and has the correct values.
#
gdns_managed_zone { <%= example_resource_name('test-example-zone') -%>:
  ensure      => present,
  dns_name    => 'test.somewild-example.com.',
  description => 'Test Example Zone',

  # You can also set output-only values as well. Puppet will ignore the values
  # when creating the resource, but will assert that its value matches what you
  # specified.
  #
  # This important to ensure that, for example, the top-level registrar is using
  # the correct DNS server names. Although this can cause failures in a run from
  # a clean project, it is useful to ensure that there are no mismatches in the
  # different services.
  #
  # id            => 8550163345207615620,
  # name_servers  => [
  #   'ns-cloud-a1.googledomains.com.',
  #   'ns-cloud-a2.googledomains.com.',
  #   'ns-cloud-a3.googledomains.com.',
  #   'ns-cloud-a4.googledomains.com.',
  # ],
  # creation_time => '2016-12-02T04:59:24.333Z',

  project     => $project, # e.g. 'my-test-project'
  credential  => 'mycred',
}

# Ensures a managed zone exists and has the correct values.
gdns_managed_zone { <%= example_resource_name('testzone-2-com') -%>:
  ensure     => absent,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

# Ensures a managed zone exists and has the correct values.
<% end # name == README.md -%>
<% res_name = 'id-for-testzone-3-com' -%>
gdns_managed_zone { <%= example_resource_name(res_name) -%>:
  ensure      => present,
  name        => 'testzone-3-com',
  dns_name    => 'test.somewild-example.com.',
  description => 'Test Example Zone',
  project     => $project, # e.g. 'my-test-project'
  credential  => 'mycred',
}

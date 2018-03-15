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

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<%
  account = example_resource_name(
    'test-account@graphite-playground.google.com.iam.gserviceaccount.com'
  )
-%>
giam_service_account { 'myaccount':
  ensure       => present,
  name         =>
    <%= account -%>,
  display_name => 'My Puppet test key',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

<% end # name == README.md -%>
giam_service_account_key { 'test-name':
  ensure           => present,
  service_account  => 'myaccount',
  path             => '/home/nelsona/test.json',
  key_algorithm    => 'KEY_ALG_RSA_2048',
  private_key_type => 'TYPE_GOOGLE_CREDENTIALS_FILE',
  project          => 'google.com:graphite-playground',
  credential       => 'mycred',
}

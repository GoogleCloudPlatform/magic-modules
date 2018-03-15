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

# To delete a key you need to either provide a key file or a key ID.
#
# $ FACTER_key_id puppet apply <%= name %>
if !defined('$key_id') {
  fail('For this example to run you need to define a fact named "key_id".
        Please refer to the documentation inside the example file
        "<%= name -%>"')
}

<% end # name == README.md -%>
giam_service_account_key { 'mykey':
  ensure           => absent,
  key_id           => $key_id,
  service_account  => 'myaccount',
  project          => 'google.com:graphite-playground',
  credential       => 'mycred',
}

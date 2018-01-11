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

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<% end # name == README.md -%>
giam_service_account { 'myaccount':
  ensure       => present,
  name         =>
    'test-23489723@graphite-playground.google.com.iam.gserviceaccount.com',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}

giam_service_account_key { 'test-name':
  ensure               => present,
  #key_id              => '9669de7d22f7be4d630783e4560f37df78b98297',
  key_file             => '/home/alexstephen/test.json',
  overwrite_if_missing => true,
  service_account      => 'myaccount',
  key_algorithm        => 'KEY_ALG_UNSPECIFIED',
  private_key_type     => 'TYPE_GOOGLE_CREDENTIALS_FILE',
  project              => 'google.com:graphite-playground',
  credential           => 'mycred',
}

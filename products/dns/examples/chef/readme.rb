<% if false # the license inside this if block pertains to this file -%>
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
<%= compile 'templates/chef/example~auth.rb.erb' -%>

gdns_managed_zone 'testzone-3-com' do
  action :create
  dns_name 'test.somewild-example.com.'
  description 'Test Example Zone'
  credential 'mycred'
  project ENV['PROJECT'] # ex: 'my-test-project'
end

gdns_resource_record_set 'www.testzone-4.com.' do
  action :create
  managed_zone 'testzone-3-com'
  type 'A'
  ttl 600
  target [
    '10.1.2.3',
    '40.5.6.7',
    '80.9.10.11'
  ]
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

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

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

<% end -%>
gdns_resource_record_set <%= example_resource_name('www.testzone-4.com.') -%> do
  action :delete
  managed_zone <%= example_resource_name('testzone-4-com') %>
  type 'A'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

<% res_name = 'sites.testzone-4.com.' -%>
gdns_resource_record_set <%= example_resource_name(res_name) -%> do
  action :delete
  managed_zone <%= example_resource_name('testzone-4-com') %>
  type 'CNAME'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end
<% unless name == 'README.md' -%>

gdns_managed_zone <%= example_resource_name('testzone-4-com') -%> do
  action :delete
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end
<% end -%>

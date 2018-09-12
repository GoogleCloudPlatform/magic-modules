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
<% res_name = 'sql-test-#{ENV[\'sql_instance_suffix\']}' -%>
<% if name != 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<% end -%>
<%= compile 'templates/chef/example~auth.rb.erb' -%>

gsql_instance  <%= example_resource_name(res_name) -%> do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gsql_database <%= example_resource_name('webstore') -%> do
  action :create
  charset 'utf8'
  instance <%= example_resource_name(res_name) %>
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

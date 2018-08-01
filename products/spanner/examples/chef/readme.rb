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
gspanner_instance <%= example_resource_name('my-spanner') -%> do
  action :create
  display_name 'My Spanner Instance'
  node_count 2
  labels [
    {
      'cost-center' => 'ti-1700004'
    }
  ]
  config 'regional-us-central1'
  project ENV['PROJECT']
  credential 'mycred'
end

gspanner_database <%= example_resource_name('webstore') -%> do
  action :create
  instance <%= example_resource_name('my-spanner') %>
  extra_statements [
    'CREATE TABLE customers (
       customer_id INT64 NOT NULL,
       last_name STRING(MAX)
     ) PRIMARY KEY (customer_id)'
  ]
  project ENV['PROJECT']
  credential 'mycred'
end

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
<% if name != "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

raise "Missing parameter 'cluster_id'. Please read docs at #{__FILE__}" \
  unless ENV.key?('cluster_id')

<% end # name == README.md -%>
<% res_name = 'mycluster-#{ENV[\'cluster_id\']}' -%>
gcontainer_cluster <%= example_resource_name(res_name) -%> do
  action :create
  initial_node_count 2
  master_auth(
    username: 'cluster_admin',
    password: 'my-secret-password'
  )
  node_config(
    machine_type: 'n1-standard-4', # we want 4-cores for our cluster
    disk_size_gb: 500              # ... and a lot of disk space
  )
  zone 'us-central1-a'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

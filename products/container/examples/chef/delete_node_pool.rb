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
<% cluster_name = 'mycluster-#{ENV[\'cluster_id\']}' -%>
<% if name != "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

raise "Missing parameter 'cluster_id'. Please read docs at #{__FILE__}" \
  unless ENV.key?('cluster_id')

gcontainer_cluster <%= example_resource_name(cluster_name) -%> do
  action :create
  zone 'us-central1-a'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

<% else -%>
# Tip: Have environment variable 'cluster-id' set.
# Tip: Insert a gcontainer cluster with name mycluster-${cluster_id}
<% end # name == README -%>
gcontainer_node_pool <%= example_resource_name('web-servers') -%> do
  action :delete
  cluster <%= example_resource_name(cluster_name) %>
  zone 'us-central1-a'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

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
<% unless name == 'README.md' -%>

<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

gcompute_global_address <%= example_resource_name('my-app-lb-address') -%> do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_instance_group <%= example_resource_name('my-chef-servers') -%> do
  action :create
  zone 'us-central1-a'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_backend_service <%= example_resource_name('my-app-backend') -%> do
  action :create
  backends [
    { group: <%= example_resource_name('my-chef-servers') -%> }
  ]
  enable_cdn true
  health_checks [
    gcompute_health_check_ref('another-hc', ENV['PROJECT']) # ex: 'my-test-project'
  ]
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_url_map <%= example_resource_name('my-url-map') -%> do
  action :create
  default_service <%= example_resource_name('my-app-backend') %>
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_target_http_proxy <%= example_resource_name('my-http-proxy') -%> do
  action :create
  url_map <%= example_resource_name('my-url-map') %>
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

<% end # name == README.md -%>
gcompute_global_forwarding_rule <%= example_resource_name('test1') -%> do
  action :create
  ip_address gcompute_global_address_ref(
    <%= example_resource_name('my-app-lb-address') -%>,
    ENV['PROJECT'] # ex: 'my-test-project'
  )
  ip_protocol 'TCP'
  port_range '80'
  target gcompute_target_http_proxy_ref(
    <%= example_resource_name('my-http-proxy') -%>,
    ENV['PROJECT'] # ex: 'my-test-project'
  )
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

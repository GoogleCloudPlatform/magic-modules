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

gcompute_zone 'us-central1-a' do
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_instance_group <%= example_resource_name('my-chef-servers') -%> do
  action :create
  zone 'us-central1-a'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

# Google::Functions must be included at runtime to ensure that the
# gcompute_health_check_ref function can be used in health_check blocks.
::Chef::Resource.send(:include, Google::Functions)

gcompute_backend_service <%= example_resource_name('my-app-backend') -%> do
  action :create
  backends [
    { group: <%= example_resource_name('my-chef-servers') -%> }
  ]
  enable_cdn true
  health_checks [
    gcompute_health_check_ref('another-hc', 'google.com:graphite-playground')
  ]
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_url_map <%= example_resource_name('my-url-map') -%> do
  action :create
  default_service <%= example_resource_name('my-app-backend') %>
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<% end # name == README.md -%>
gcompute_target_http_proxy <%= example_resource_name('my-http-proxy') -%> do
  action :create
  url_map <%= example_resource_name('my-url-map') %>
  project 'google.com:graphite-playground'
  credential 'mycred'
end

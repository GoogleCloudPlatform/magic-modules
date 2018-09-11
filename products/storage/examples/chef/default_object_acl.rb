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
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

gstorage_bucket <%= example_resource_name('storage-module-test') -%> do
  action :create
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<% else # name == README.md -%>
# Default Object ACL requires a bucket. Please ensure its existence with
# the gstorage_bucket { ... } resource
<% end # name == README.md -%>
<% res_name = 'user-nelsona@google.com' -%>
gstorage_default_object_acl <%= example_resource_name(res_name) -%> do
  action :create
  bucket <%= example_resource_name('storage-module-test') %>
  entity 'user-nelsona@google.com'
  role 'WRITER'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

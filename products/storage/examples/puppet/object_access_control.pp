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
<% bucket_name = 'puppet-storage-module-test' -%>
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gstorage_bucket { <%= example_resource_name(bucket_name) -%>:
  ensure     => present,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

<% else # name == README.md -%>
# Object Access Control requires a bucket. Please ensure its existence with
# the gstorage_bucket { ... } resource.  The object does not have to exist.
<% end # name == README.md -%>
<% res_name = 'user-nelsona@google.com' -%>
gstorage_object_access_control { <%= example_resource_name(res_name) -%>:
  bucket     => <%= example_resource_name(bucket_name) -%>,
  object     => 'acl-controlled-file.txt'
  entity     => 'user-nelsona@google.com',
  role       => 'WRITER',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

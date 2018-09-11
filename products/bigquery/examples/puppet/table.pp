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
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<% end # name == README.md -%>
gbigquery_dataset { <%= example_resource_name('example_dataset') -%>:
  ensure            => present,
  dataset_reference => {
    dataset_id => 'example_dataset'
  },
  project           => $project, # e.g. 'my-test-project'
  credential        => 'mycred',
}

gbigquery_table { <%= example_resource_name('example_table') -%>:
  ensure          => present,
  dataset         => <%= example_resource_name('example_dataset') -%>,
  table_reference => {
    dataset_id => <%= example_resource_name('example_dataset') -%>,
    project_id => $project,
    table_id   => <%= example_resource_name('example_table') %>
  },
  project         => $project, # e.g. 'my-test-project'
  credential      => 'mycred',
}

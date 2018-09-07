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

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gcompute_region { <%= example_resource_name('some-region') -%>:
  name       => 'us-west1',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

<% end # name == README.md -%>
gcompute_forwarding_rule { <%= example_resource_name('test1') -%>:
  ensure     => absent,
  region     => <%= example_resource_name('some-region') -%>,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

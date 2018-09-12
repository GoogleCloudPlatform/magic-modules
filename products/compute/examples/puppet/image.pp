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
<% if name != "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gcompute_zone { 'us-central1-a':
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

gcompute_disk { <%= example_resource_name('data-disk-1') -%>:
  ensure     => present,
  zone       => 'us-central1-a',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

<% else -%>
# Tip: Be sure to include a valid gcompute_disk object
<% end # name == README.md -%>
gcompute_image { <%= example_resource_name('test-image') -%>:
  ensure      => present,
  source_disk => <%= example_resource_name('data-disk-1') -%>,
  project     => $project, # e.g. 'my-test-project'
  credential  => 'mycred'
}

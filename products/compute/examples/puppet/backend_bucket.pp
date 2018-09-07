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

<%
  # TODO(nelsonjr): http://b/63088154 Google Cloud Platform API is returning
  # access denied if we use a more restricted scope such as
  # https://www.googleapis.com/auth/compute. For the time being use an all
  # mighty scope instead: https://www.googleapis.com/auth/cloud-platform.
  original_scopes = data[:scopes]
  data[:scopes] = ['https://www.googleapis.com/auth/cloud-platform']
-%>
<%= compile 'templates/puppet/examples~credential.pp.erb' -%>
<% data[:scopes] = original_scopes # restore the scopes -%>

<% end # name == README.md -%>
gcompute_backend_bucket { <%= example_resource_name('be-bucket-connection') -%>:
  ensure      => present,
  bucket_name => <%= example_resource_name('backend-bucket-test') -%>,
  description => 'A BackendBucket to connect LNB w/ Storage Bucket',
  enable_cdn  => true,
  project     => $project, # e.g. 'my-test-project'
  credential  => 'mycred',
}

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

<% end -%>
# Ensures a project exists and has the correct values.
#
# All project settings are read-only, yet we are setting them anyway. Chef will
# use these values to check if they match, and fail the run otherwise.
#
# This important to ensure that your project quotas are set properly and avoid
# discrepancies from it to fail in production.
gdns_project <%= example_resource_name('google.com:graphite-playground') -%> do
  quota_managed_zones 10_000
  quota_total_rrdata_size_per_change 100_000
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

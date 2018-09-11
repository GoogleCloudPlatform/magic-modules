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
<% unless name == 'README.md' -%>

<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

<% end -%>
gdns_managed_zone <%= example_resource_name('testzone-3-com') -%> do
  action :create
  dns_name 'test.somewild-example.com.'
  description 'Test Example Zone'

  # You can also set output-only values as well. Chef will ignore the values
  # when creating the resource, but will assert that its value matches what you
  # specified.
  #
  # This important to ensure that, for example, the top-level registrar is using
  # the correct DNS server names. Although this can cause failures in a run from
  # a clean project, it is useful to ensure that there are no mismatches in the
  # different services.
  #
  # id 579_667_184_320_567_887
  # name_servers [
  #   'ns-cloud-b1.googledomains.com.',
  #   'ns-cloud-b2.googledomains.com.',
  #   'ns-cloud-b3.googledomains.com.',
  #   'ns-cloud-b4.googledomains.com.'
  # ]
  # creation_time '2016-12-02T04:59:24.333Z'

  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

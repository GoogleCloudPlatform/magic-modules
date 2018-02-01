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
<% if name != 'README.md' -%>

<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

gcompute_zone 'us-west1-a' do
  action :create
  project 'google.com:graphite-playground'
  credential 'mycred'
end

# Google::Functions must be included at runtime to ensure that the
# gcompute_image_family function can be used in gcompute_disk blocks.
::Chef::Resource.send(:include, Google::Functions)

gcompute_disk <%= example_resource_name('instance-test-os-1') -%> do
  action :create
  source_image gcompute_image_family('ubuntu-1604-lts', 'ubuntu-os-cloud')
  zone 'us-west1-a'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_network <%= example_resource_name('mynetwork-test') -%> do
  action :create
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_region 'us-west1' do
  action :create
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_address <%= example_resource_name('instance-test-ip') -%> do
  action :create
  region 'us-west1'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

gcompute_machine_type 'n1-standard-1' do
  action :create
  zone 'us-west1-a'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

<% else -%>
# Power Tips:
#   1) Remember to define the resources needed to allocate the VM:
#      a) gcompute_disk (to be used in 'disks' property)
#      b) gcompute_network (to be used in 'network' property)
#      c) gcompute_address (to be used in 'access_configs', if your machine
#         needs external ingress access)
#      d) gcompute_zone (to determine where the VM will be allocated)
#      e) gcompute_machine_type (to determine the kind of machine to be created)
#   2) Don't forget to define a source_image for the OS of the boot disk
#      a) You can use the provided gcompute_image_family function to specify the
#         latest version of an operating system of a given family
#         e.g. Ubuntu 16.04
<% end -%>
gcompute_instance <%= example_resource_name('instance-test') -%> do
  action :create
  machine_type 'n1-standard-1'
  disks [
    {
      boot: true,
      auto_delete: true,
      source: <%= example_resource_name('instance-test-os-1') %>
    }
  ]
  metadata (
    'startup-script-url' => 'gs://graphite-playground/bootstrap.sh',
    'cost-center' => '12345'
  )
  network_interfaces [
    {
      network: <%= example_resource_name('mynetwork-test') %>,
      access_configs: [
        {
          name: 'External NAT',
          nat_ip: <%= example_resource_name('instance-test-ip') -%>,
          type: 'ONE_TO_ONE_NAT'
        }
      ]
    }
  ]
  zone 'us-west1-a'
  project 'google.com:graphite-playground'
  credential 'mycred'
end

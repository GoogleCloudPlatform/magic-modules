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
<% if name != 'README.md' -%>

<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :chef) -%>

<%= compile 'templates/chef/example~auth.rb.erb' -%>

# TODO(nelsonjr): Reactiveate example based on disk once http://b/66871792 is
# resolved.
#gcompute_disk <%= example_resource_name('os-disk-1') -%> do
#  action :create
#  zone 'us-west1-a'
#  source_image 'projects/ubuntu-os-cloud/global/images/family/ubuntu-1604-lts'
#  project ENV['PROJECT'] # ex: 'my-test-project'
#  credential 'mycred'
#end

# Google::Functions must be included at runtime to ensure that the
# gcompute_image_family function can be used in gcompute_disk blocks.
::Chef::Resource.send(:include, Google::Functions)

gcompute_network <%= example_resource_name('mynetwork-test') -%> do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

<% else -%>
# Power Tips:
#   1) Remember to define the resources needed to allocate the VM:
#      a) gcompute_disk_type (to be used in 'diskType' property)
#      b) gcompute_network (to be used in 'network_interfaces' property)
#      c) gcompute_subnetwork (to be used in the 'subnetwork' property)
#      d) gcompute_disk (to be used in the 'sourceDisk' property)
#   2) Don't forget to define a source_image for the OS of the boot disk
<% end -%>
<% res_name = example_resource_name('instance-template-test') -%>
gcompute_instance_template <%= res_name -%> do
  action :create
  properties(
    machine_type: 'n1-standard-1',
    disks: [
      {
        # Tip: Auto delete will prevent disks from being left behind on
        # deletion.
        auto_delete: true,
        boot: true,
        initialize_params: {
          disk_size_gb: 100,
          source_image:
            gcompute_image_family('ubuntu-1604-lts', 'ubuntu-os-cloud')
        }
      }
    ],
    metadata: {
      'startup-script-url' => 'gs://graphite-playground/bootstrap.sh',
      'cost-center' => '12345'
    },
    network_interfaces: [
      {
        access_configs: {
          name: 'test-config',
          type: 'ONE_TO_ONE_NAT',
        },
        network: <%= example_resource_name('mynetwork-test') %>
      }
    ]
  )
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

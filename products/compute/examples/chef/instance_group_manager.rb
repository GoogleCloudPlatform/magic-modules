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

gcompute_zone 'us-west1-a' do
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gcompute_machine_type 'n1-standard-1' do
  project ENV['PROJECT'] # ex: 'my-test-project'
  zone 'us-west1-a'
  credential 'mycred'
end

gcompute_network 'mynetwork-test' do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

# Google::Functions must be included at runtime to ensure that the
# gcompute_image_family function can be used in gcompute_disk blocks.
::Chef::Resource.send(:include, Google::Functions)

gcompute_instance_template 'instance-template' do
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
          source_image:
            gcompute_image_family('ubuntu-1604-lts', 'ubuntu-os-cloud')
        }
      }
    ],
    network_interfaces: [
      {
        network: 'mynetwork-test',
        access_configs: [
          {
            name: 'External NAT',
            type: 'ONE_TO_ONE_NAT'
          }
        ]
      }
    ]
  )
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

<% end # name == README.md -%>
gcompute_instance_group_manager 'test1' do
  action :create
  base_instance_name 'test1-child'
  instance_template 'instance-template'
  target_size 3
  zone 'us-west1-a'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

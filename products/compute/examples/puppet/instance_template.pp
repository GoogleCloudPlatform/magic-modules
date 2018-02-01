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
<% if name != "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gcompute_zone { 'us-central1-a':
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

gcompute_machine_type { 'n1-standard-1':
  project    => 'google.com:graphite-playground',
  zone       => 'us-central1-a',
  credential => 'mycred',
}

# TODO(nelsonjr): Reactivate example based on disk once http://b/66871792 is
# resolved.
# | gcompute_disk { <%= example_resource_name('os-disk-1') -%>:
# |   ensure       => present,
# |   zone         => 'us-central1-a',
# |   source_image =>
# |     gcompute_image_family('ubuntu-1604-lts', 'ubuntu-os-cloud'),
# |   project      => 'google.com:graphite-playground',
# |   credential   => 'mycred',
# | }

gcompute_network { <%= example_resource_name('mynetwork-test') -%>:
  ensure     => present,
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

<% else -%>
# Power Tips:
#   1) Remember to define the resources needed to allocate the VM:
#      a) gcompute_disk_type (to be used in 'diskType' property)
#      b) gcompute_machine_type (to be used in 'machine_type' property)
#      c) gcompute_network (to be used in 'network_interfaces' property)
#      d) gcompute_subnetwork (to be used in the 'subnetwork' property)
#      e) gcompute_disk (to be used in the 'sourceDisk' property)
#   2) Don't forget to define a source_image for the OS of the boot disk
#      a) You can use the provided gcompute_image_family function to specify the
#         latest version of an operating system of a given family
#         e.g. Ubuntu 16.04
<% end # name == README.md -%>
gcompute_instance_template { <%= example_resource_name('instance-template') -%>:
  ensure     => present,
  properties => {
    machine_type       => 'n1-standard-1',
    disks              => [
      {
        # Tip: Auto delete will prevent disks from being left behind on
        # deletion.
        auto_delete       => true,
        boot              => true,
        initialize_params => {
          source_image =>
            gcompute_image_family('ubuntu-1604-lts', 'ubuntu-os-cloud'),
        }
      }
    ],
    metadata           => {
      'startup-script-url'   => 'gs://graphite-playground/bootstrap.sh',
      'cost-center'          => '12345',
    },
    network_interfaces => [
      {
        access_configs => {
          name => 'test-config',
          type => 'ONE_TO_ONE_NAT',
        },
        network        => <%= example_resource_name('mynetwork-test') %>,
      }
    ]
  },
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

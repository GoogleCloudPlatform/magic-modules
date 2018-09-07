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
<% if name != "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gcompute_zone { 'us-central1-a':
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

gcompute_disk { <%= example_resource_name('instance-test-os-1') -%>:
  ensure       => present,
  size_gb      => 50,
  source_image =>
    'projects/ubuntu-os-cloud/global/images/family/ubuntu-1604-lts',
  zone         => 'us-central1-a',
  project      => $project, # e.g. 'my-test-project'
  credential   => 'mycred',
}

# Tips
#   1) You can use network 'default' if do not use VLAN or other traffic
#      seggregation on your project.
#   2) Don't forget to define the firewall rules if you specify a custom
#      network to ensure the traffic can reach your machine
gcompute_network { <%= example_resource_name('default') -%>:
  ensure     => present,
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

gcompute_region { 'us-central1':
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

# Defines the machine type to be used by the VM. This definition is required
# only once per catalog as it is shared to any objects that use the
# 'n1-standard-1' defined below.
gcompute_machine_type { 'n1-standard-1':
  zone       => 'us-central1-a',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

# Ensures the 'instance-test-ip' external IP address exists. If it does not
# exist it will allocate an ephemeral one.
gcompute_address { <%= example_resource_name('instance-test-ip') -%>:
  region     => 'us-central1',
  project    => $project, # e.g. 'my-test-project'
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
#      f) gcompute_address (to be used in 'access_configs', if your machine
#         needs external ingress access)
#   2) Don't forget to define a source_image for the OS of the boot disk
#      a) You can use the provided gcompute_image_family function to specify the
#         latest version of an operating system of a given family
#         e.g. Ubuntu 16.04
<% end # name == README.md -%>
gcompute_instance { <%= example_resource_name('instance-test') -%>:
  ensure             => present,
  machine_type       => 'n1-standard-1',
  disks              => [
    {
      auto_delete => true,
      boot        => true,
      source      => <%= example_resource_name('instance-test-os-1') %>
    }
  ],
  metadata           => {
    startup-script-url   => 'gs://graphite-playground/bootstrap.sh',
    cost-center          => '12345',
  },
  network_interfaces => [
    {
      network        => <%= example_resource_name('default') %>,
      access_configs => [
        {
          name   => 'External NAT',
          nat_ip => <%= example_resource_name('instance-test-ip') -%>,
          type   => 'ONE_TO_ONE_NAT',
        },
      ],
    }
  ],
  zone               => 'us-central1-a',
  project            => $project, # e.g. 'my-test-project'
  credential         => 'mycred',
}

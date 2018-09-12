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
$vm = $_name
$image = split($image_family, ':') # input: <family-name>:<image-name>

# Convert provided zone to a region
$region_parts = split($zone, '-')
$region       = "${region_parts[0]}-${region_parts[1]}"

# Convenience because title is only used by Puppet, not GCP
$cred = 'bolt-credential'

gauth_credential { $cred:
  provider => serviceaccount,
  path     => $credential,
  scopes   => ['https://www.googleapis.com/auth/cloud-platform'],
}

gcompute_network { $network_name:
  project    => $project,
  credential => $cred,
}

gcompute_region { $region:
  project    => $project,
  credential => $cred,
}

gcompute_zone { $zone:
  project    => $project,
  credential => $cred,
}

gcompute_machine_type { $machine_type:
  zone       => $zone,
  project    => $project,
  credential => $cred,
}

if $ensure == absent {
  Gcompute_instance[$vm] -> Gcompute_disk[$vm]
  if $allocate_static_ip {
    Gcompute_instance[$vm] -> Gcompute_address[$vm]
  }
}

gcompute_disk { $vm:
  ensure       => $ensure,
  size_gb      => $size_gb,
  source_image => gcompute_image_family($image[0], $image[1]),
  zone         => $zone,
  project      => $project,
  credential   => $cred,
}

if $allocate_static_ip {
  gcompute_address { $vm:
    ensure     => $ensure,
    region     => $region,
    project    => $project,
    credential => $cred,
  }

  gcompute_instance { $vm:
    ensure             => $ensure,
    machine_type       => $machine_type,
    disks              => [
      {
        boot        => true,
        source      => $vm,
        auto_delete => true,
      }
    ],
    network_interfaces => [
      {
        network        => $network_name,
        access_configs => [
          {
            name   => 'External NAT',
            nat_ip => $vm,
            type   => 'ONE_TO_ONE_NAT',
          }
        ],
      }
    ],
    zone               => $zone,
    project            => $project,
    credential         => $cred,
  }
} else {
  gcompute_instance { $vm:
    ensure             => $ensure,
    machine_type       => $machine_type,
    disks              => [
      {
        boot        => true,
        source      => $vm,
        auto_delete => true,
      }
    ],
    network_interfaces => [
      {
        network        => $network_name,
        access_configs => [
          {
            name   => 'External NAT',
            type   => 'ONE_TO_ONE_NAT',
          }
        ],
      }
    ],
    zone               => $zone,
    project            => $project,
    credential         => $cred,
  }
}

notify { 'task:success':
  message => "$vm",
  require => Gcompute_instance[$vm],
}

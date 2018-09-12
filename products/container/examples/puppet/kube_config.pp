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

if $cluster_id == undef {
  fail('Please specify $cluster_id variable to run this example.')
}

gcontainer_cluster { "mycluster-${cluster_id}":
  ensure             => present,
  initial_node_count => 2,
  master_auth        => {
    username => 'cluster_admin',
    password => 'my-secret-password',
  },
  node_config        => {
    machine_type => 'n1-standard-4', # we want a 4-core machine for our cluster
    disk_size_gb => 500,             # ... and a lot of disk space
  },
  zone               => 'us-central1-a',
  project            => $project, # e.g. 'my-test-project'
  credential         => 'mycred',
}

file { '/home/nelsona/.kube':
  ensure => directory,
}

file { '/home/nelsona/.puppetlabs/etc/puppet':
  ensure => directory,
}

<% end # name == README.md -%>
# ~/.kube/config is used by Kubernetes client (kubectl)
gcontainer_kube_config { '/home/nelsona/.kube/config':
  ensure     => present,
  context    => "gke-mycluster-${cluster_id}",
  cluster    => "mycluster-${cluster_id}",
  zone       => 'us-central1-a',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

# A file named ~/.puppetlabs/etc/puppet/kubernetes is used by the
# garethr-kubernetes module.
gcontainer_kube_config { '/home/nelsona/.puppetlabs/etc/puppet/kubernetes.conf':
  ensure     => present,
  cluster    => "mycluster-${cluster_id}",
  zone       => 'us-central1-a',
  project    => $project, # e.g. 'my-test-project'
  credential => 'mycred',
}

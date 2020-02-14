### Test that a GCP container cluster is in a particular state e.g. "RUNNING"

    describe google_container_cluster(project: 'chef-inspec-gcp', location: 'europe-west2-a', name: 'inspec-gcp-kube-cluster') do
      its('status') { should eq 'RUNNING' }
    end

### Test that a GCP container cluster has the expected kube master user/password

    describe google_container_cluster(project: 'chef-inspec-gcp', location: 'europe-west2-a', name: 'inspec-gcp-kube-cluster') do
      its('master_auth.username'){ should eq "user_name"}
      its('master_auth.password'){ should eq "choose_something_strong"}
    end

### Test that the locations where the GCP container cluster is running match those expected

    describe google_container_cluster(project: 'chef-inspec-gcp', location: 'europe-west2-a', name: 'inspec-gcp-kube-cluster') do
      its('locations.sort'){should cmp ["europe-west2-a", "europe-west2-b", "europe-west2-c"].sort}
    end

### Test GCP container cluster network and subnetwork settings

    describe google_container_cluster(project: 'chef-inspec-gcp', location: 'europe-west2-a', name: 'inspec-gcp-kube-cluster') do
      its('network'){should eq "default"}
      its('subnetwork'){should eq "default"}
    end

### Test GCP container cluster node pool configuration settings

    describe google_container_cluster(project: 'chef-inspec-gcp', location: 'europe-west2-a', name: 'inspec-gcp-kube-cluster') do
      its('node_config.disk_size_gb'){should eq 100}
      its('node_config.image_type'){should eq "COS"}
      its('node_config.machine_type'){should eq "n1-standard-1"}
      its('node_ipv4_cidr_size'){should eq 24}
      its('node_pools.count'){should eq 1}
    end
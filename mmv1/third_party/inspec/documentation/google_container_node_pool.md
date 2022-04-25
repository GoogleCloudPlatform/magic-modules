### Test that a GCP container node pool is in a particular state e.g. "RUNNING"

    describe google_container_node_pool(project: 'chef-inspec-gcp', locations: 'europe-west2-a', cluster_name: 'inspec-gcp-kube-cluster', nodepool_name: 'inspec-gcp-kube-node-pool') do
      its('status') { should eq 'RUNNING' }
    end

### Test GCP container node pool disk size in GB is as expected

    describe google_container_node_pool(project: 'chef-inspec-gcp', locations: 'europe-west2-a', cluster_name: 'inspec-gcp-kube-cluster', nodepool_name: 'inspec-gcp-kube-node-pool') do
      its('node_config.disk_size_gb'){should eq 100}
    end

### Test GCP container node pool machine type is as expected

    describe google_container_node_pool(project: 'chef-inspec-gcp', locations: 'europe-west2-a', cluster_name: 'inspec-gcp-kube-cluster', nodepool_name: 'inspec-gcp-kube-node-pool') do
      its('node_config.machine_type'){should eq "n1-standard-1"}
    end

### Test GCP container node pool node image type is as expected

    describe google_container_node_pool(project: 'chef-inspec-gcp', locations: 'europe-west2-a', cluster_name: 'inspec-gcp-kube-cluster', nodepool_name: 'inspec-gcp-kube-node-pool') do
      its('node_config.image_type'){should eq "COS"}
    end

### Test GCP container node pool initial node count is as expected

    describe google_container_node_pool(project: 'chef-inspec-gcp', locations: 'europe-west2-a', cluster_name: 'inspec-gcp-kube-cluster', nodepool_name: 'inspec-gcp-kube-node-pool') do
      its('initial_node_count'){should eq 3}
    end
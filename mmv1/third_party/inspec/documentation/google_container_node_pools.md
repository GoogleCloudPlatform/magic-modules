### Test that there are no more than a specified number of node pools available for the project

    describe google_container_node_pools(project: 'chef-inspec-gcp') do
      its('count') { should be <= 10}
    end

### Test that an expected node pool is available for the project

    describe google_container_node_pools(project: 'chef-inspec-gcp') do
      its('node_pool_names') { should include "us-east1-b" }
    end

### Test that a subset of all node pools matching "mypool*" are "UP"

    google_container_node_pools(project: 'chef-inspec-gcp', location: 'europe-west2-a', cluster_name: 'inspec-gcp-cluster').where(node_pool_name: /^mypool/).node_pool_names.each do |node_pool_name|
      describe google_container_node_pool(project: 'chef-inspec-gcp', location: 'europe-west2-a', cluster_name: 'inspec-gcp-cluster', nodepool_name: node_pool_name) do
        it { should exist }
        its('status') { should eq 'RUNNING' }
      end
    end
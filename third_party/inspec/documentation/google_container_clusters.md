### Test that there are no more than a specified number of clusters available for the project in a particular zone

    describe google_container_clusters(project: 'chef-inspec-gcp', location: 'europe-west2-a') do
      its('count') { should be <= 5}
    end

### Test that an expected cluster is available for the project

    describe google_container_clusters(project: 'chef-inspec-gcp', location: 'europe-west2-a') do
      its('cluster_names') { should include "my-cluster" }
    end

### Test whether any clusters are in status "STOPPING"

    describe google_container_clusters(project: 'chef-inspec-gcp', location: 'europe-west2-a') do
      its('cluster_statuses') { should_not include "STOPPING" }
    end

### Test that a subset of all clusters matching "kube*" are "RUNNING"

    google_container_clusters(project: gcp_project_id).where(cluster_name: /^kube/).cluster_names.each do |cluster_name|
      describe google_container_cluster(project: 'chef-inspec-gcp', location: 'europe-west2-a', name: cluster_name) do
        it { should exist }
        its('status') { should eq 'RUNNING' }
      end
    end
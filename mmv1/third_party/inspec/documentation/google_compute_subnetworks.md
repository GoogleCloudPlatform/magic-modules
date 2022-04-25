### Test that there are no more than a specified number of subnetworks available for the project and region

    describe google_compute_subnetworks(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('count') { should be <= 100}
    end

### Test that an expected subnetwork identifier is present in the project and region

    describe google_compute_subnetworks(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('subnetwork_ids') { should include 12345678975432 }
    end


### Test that an expected subnetwork name is available for the project and region

    describe google_compute_subnetworks(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('subnetwork_names') { should include "subnetwork-name" }
    end

### Test that an expected subnetwork network name is not present for the project and region

    describe google_compute_subnetworks(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('subnetwork_networks') { should not include "network-name" }
    end
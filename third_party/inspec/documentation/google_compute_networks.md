### Test that there are no more than a specified number of networks available for the project

    describe google_compute_networks(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected network identifier is present in the project 

    describe google_compute_networks(project: 'chef-inspec-gcp') do
      its('network_ids') { should include 12345678975432 }
    end

### Test that an expected network name is available for the project

    describe google_compute_networks(project: 'chef-inspec-gcp') do
      its('network_names') { should include "network-name" }
    end

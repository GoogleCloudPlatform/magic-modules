### Test that a GCP compute network exists

    describe google_compute_network(project: 'chef-inspec-gcp',  name: 'gcp-inspec-network') do
      it { should exist }
    end

### Test when a GCP compute network was created

    describe google_compute_network(project: 'chef-inspec-gcp',  name: 'gcp-inspec-network') do
      its('creation_timestamp_date') { should be > Time.now - 365*60*60*24*10 }
    end    
    
### Test for an expected network identifier 

    describe google_compute_network(project: 'chef-inspec-gcp',  name: 'gcp-inspec-network') do
      its('id') { should eq 12345567789 }
    end    


### Test whether a single attached subnetwork name is correct 

    describe google_compute_network(project: 'chef-inspec-gcp',  name: 'gcp-inspec-network') do
      its ('subnetworks.count') { should eq 1 }
      its ('subnetworks.first') { should match "subnetwork-name"}
    end    
    
### Test whether the network is configured to automatically create subnetworks or not

    describe google_compute_network(project: 'chef-inspec-gcp',  name: 'gcp-inspec-network') do
      its ('auto_create_subnetworks'){ should be false }
    end    


### Check the network routing configuration routing mode 

    describe google_compute_network(project: 'chef-inspec-gcp',  name: 'gcp-inspec-network') do
      its ('routing_config.routing_mode') { should eq "REGIONAL" }
    end
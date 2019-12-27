### Test that a GCP compute subnetwork exists

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      it { should exist }
    end

### Test when a GCP compute subnetwork was created

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      its('creation_timestamp') { should be > Time.now - 365*60*60*24*10 }
    end

### Test for an expected subnetwork identifier 

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      its('id') { should eq 12345567789 }
    end    

### Test that a subnetwork gateway address is as expected

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      its('gateway_address') { should eq "10.2.0.1" }
    end  

### Test that a subnetwork IP CIDR range is as expected

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      its('ip_cidr_range') { should eq "10.2.0.0/29" }
    end  

### Test that a subnetwork is associated with the expected network

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      its('network') { should match "gcp_network_name" }
    end  

### Test whether VMs in this subnet can access Google services without assigning external IP addresses through Private Google Access

    describe google_compute_subnetwork(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-subnetwork') do
      its('private_ip_google_access') { should be false }
    end
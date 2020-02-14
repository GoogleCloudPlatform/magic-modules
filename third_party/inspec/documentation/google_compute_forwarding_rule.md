### Test that a GCP compute forwarding_rule exists

    describe google_compute_forwarding_rule(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-forwarding_rule') do
      it { should exist }
    end

### Test when a GCP compute forwarding_rule was created

    describe google_compute_forwarding_rule(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-forwarding_rule') do
      its('creation_timestamp_date') { should be > Time.now - 365*60*60*24*10 }
    end

### Test for an expected forwarding_rule identifier 

    describe google_compute_forwarding_rule(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-forwarding_rule') do
      its('id') { should eq 12345567789 }
    end    

### Test that a forwarding_rule load_balancing_scheme is as expected

    describe google_compute_forwarding_rule(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-forwarding_rule') do
      its('load_balancing_scheme') { should eq "INTERNAL" }
    end  

### Test that a forwarding_rule IP address is as expected

    describe google_compute_forwarding_rule(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-forwarding_rule') do
      its('ip_address') { should eq "10.0.0.1" }
    end  

### Test that a forwarding_rule is associated with the expected network

    describe google_compute_forwarding_rule(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-forwarding_rule') do
      its('network') { should match "gcp_network_name" }
    end  
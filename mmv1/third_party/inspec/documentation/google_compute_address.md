### Test that a GCP compute address IP exists

    describe google_compute_address(project: 'chef-inspec-gcp', location: 'europe-west2', name: 'compute-address') do
      its('address_ip_exists')  { should be true }
    end

### Test that a GCP compute address is in a particular status

    describe google_compute_address(project: 'chef-inspec-gcp', location: 'europe-west2', name: 'compute-address') do
      its('status') { should eq "IN_USE" }
    end

### Test that a GCP compute address IP has the expected number of users

    describe google_compute_address(project: 'chef-inspec-gcp', location: 'europe-west2', name: 'compute-address') do
      its('user_count') { should eq 1 }
    end

### Test that the first user of a GCP compute address has the expected resource name

    describe google_compute_address(project: 'chef-inspec-gcp', location: 'europe-west2', name: 'compute-address') do
      its('user_resource_name') { should eq "gcp_ext_vm_name" }
    end

### Test that a GCP compute region instance group manager has the expected size

    describe google_compute_region_instance_group_manager(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-test') do
      its('target_size') { should eq 2 }
    end

### Test that a GCP compute region instance group manager has a port with supplied name and value

    describe google_compute_region_instance_group_manager(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-test') do
      its('named_ports') { should include "http" }
    end